package components

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/src-d/engine/docker"
)

var srcdNamespaces = []string{
	"srcd",
	"bblfsh",
}

type Component struct {
	Name    string
	Image   string
	Version string // only if there's a required version
}

func (c *Component) ImageWithVersion() string {
	return fmt.Sprintf("%s:%s", c.Image, c.Version)
}

// Kill removes the Component container. If it is not running it returns nil
func (c *Component) Kill() error {
	err := docker.RemoveContainer(c.Name)
	if err != nil && err != docker.ErrNotFound {
		return err
	}

	return nil
}

// IsInstalled returns true if the Component image is installed with the
// exact version
func (c *Component) IsInstalled(ctx context.Context) (bool, error) {
	return IsInstalled(ctx, c.ImageWithVersion())
}

// IsRunning returns true if the Component container is running using the
// exact image version
func (c *Component) IsRunning() (bool, error) {
	return docker.IsRunning(c.Name, c.ImageWithVersion())
}

const (
	BblfshVolume = "srcd-cli-bblfsh-storage"
)

var (
	Gitbase = Component{
		Name:    "srcd-cli-gitbase",
		Image:   "srcd/gitbase",
		Version: "v0.17.1",
	}

	GitbaseWeb = Component{
		Name:    "srcd-cli-gitbase-web",
		Image:   "srcd/gitbase-web",
		Version: "v0.3.0",
	}

	Bblfshd = Component{
		Name:    "srcd-cli-bblfshd",
		Image:   "bblfsh/bblfshd",
		Version: "v2.9.2-drivers",
	}

	BblfshWeb = Component{
		Name:    "srcd-cli-bblfsh-web",
		Image:   "bblfsh/web",
		Version: "v0.7.0",
	}

	workDirDependants = []Component{
		Gitbase,
		Bblfshd, // does not depend on workdir but it does depend on user dir
	}
)

// FilterFunc is a filtering function for List.
type FilterFunc func(Component) bool

func filter(cmps []Component, filters []FilterFunc) []Component {
	var result []Component
	for _, cmp := range cmps {
		var add = true
		for _, f := range filters {
			if !f(cmp) {
				add = false
				break
			}
		}

		if add {
			result = append(result, cmp)
		}
	}
	return result
}

// IsWorkingDirDependant filters Components that depend on the working directory.
var IsWorkingDirDependant FilterFunc = func(cmp Component) bool {
	for _, c := range workDirDependants {
		if c.Image == cmp.Image {
			return true
		}
	}
	return false
}

// IsInstalledFilter filters Components that have its image installed, with
// the exact version
var IsInstalledFilter FilterFunc = func(cmp Component) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	installed, err := docker.IsInstalled(ctx, cmp.Image, cmp.Version)
	if err != nil {
		//TODO log, or add err to FilterFunc
		return false
	}

	return installed
}

// IsRunningFilter filters Components that have a container running, using
// its image with the exact version
var IsRunningFilter FilterFunc = func(cmp Component) bool {
	r, err := cmp.IsRunning()
	if err != nil {
		return false
	}

	return r
}

// List returns the list of known Components, which may or may not be installed.
// If allVersions is true other Components with image versions different from
// the current ones will be included.
func List(ctx context.Context, allVersions bool, filters ...FilterFunc) ([]Component, error) {
	componentsList := []Component{
		Gitbase,
		GitbaseWeb,
		Bblfshd,
		BblfshWeb,
	}

	if allVersions {
		otherComponents := make([]Component, 0)

		for _, cmp := range componentsList {
			newCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			// Look for any other image version that might be installed
			versions, err := docker.VersionsInstalled(newCtx, cmp.Image)
			if err != nil {
				return nil, err
			}

			for _, v := range versions {
				if v == cmp.Version {
					// Already added before
					continue
				}

				otherComponents = append(otherComponents, Component{
					Name:    cmp.Name,
					Image:   cmp.Image,
					Version: v,
				})
			}
		}

		componentsList = append(componentsList, otherComponents...)
	}

	if len(filters) > 0 {
		return filter(componentsList, filters), nil
	}

	return componentsList, nil
}

var ErrNotSrcd = fmt.Errorf("not srcd component")

// Install installs a new component.
func Install(ctx context.Context, id string) error {
	if !isSrcdComponent(id) {
		return ErrNotSrcd
	}

	image, version := docker.SplitImageID(id)
	return docker.Pull(ctx, image, version)
}

func IsInstalled(ctx context.Context, id string) (bool, error) {
	if !isSrcdComponent(id) {
		return false, ErrNotSrcd
	}

	image, version := docker.SplitImageID(id)
	return docker.IsInstalled(ctx, image, version)
}

func Stop() error {
	logrus.Info("stopping containers...")

	// we actually not just stop but remove containers here
	// it's needed to make sure configuration of the containers is correct
	// without over-complicated logic for it
	if err := removeContainers(); err != nil {
		return errors.Wrap(err, "unable to stop all containers")
	}

	return nil
}

func Prune(images bool) error {
	logrus.Info("removing containers...")
	if err := removeContainers(); err != nil {
		return errors.Wrap(err, "unable to remove all containers")
	}

	logrus.Info("removing volumes...")

	if err := removeVolumes(); err != nil {
		return errors.Wrap(err, "unable to remove volumes")
	}

	logrus.Info("removing network...")

	if err := docker.RemoveNetwork(context.Background()); err != nil {
		return errors.Wrap(err, "unable to remove network")
	}

	if images {
		logrus.Info("removing images...")

		if err := removeImages(); err != nil {
			return errors.Wrap(err, "unable to remove all images")
		}
	}

	return nil
}

func removeContainers() error {
	cs, err := docker.List()
	if err != nil {
		return err
	}

	for _, c := range cs {
		if len(c.Names) == 0 {
			continue
		}

		name := strings.TrimLeft(c.Names[0], "/")
		if isFromEngine(name) {
			logrus.Infof("removing container %s", name)

			if err := docker.RemoveContainer(name); err != nil {
				return err
			}
		}
	}

	return nil
}

func removeVolumes() error {
	vols, err := docker.ListVolumes(context.Background())
	if err != nil {
		return err
	}

	for _, vol := range vols {
		if isFromEngine(vol.Name) {
			logrus.Infof("removing volume %s", vol.Name)

			if err := docker.RemoveVolume(context.Background(), vol.Name); err != nil {
				return err
			}
		}
	}

	return nil
}

func removeImages() error {
	cmps, err := List(context.Background(), true, IsInstalledFilter)
	if err != nil {
		return errors.Wrap(err, "unable to list images")
	}

	for _, cmp := range cmps {
		logrus.Infof("removing image %s", cmp.ImageWithVersion())

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		if err := docker.RemoveImage(ctx, cmp.ImageWithVersion()); err != nil {
			return err
		}
	}

	return nil
}

func stringInSlice(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// isSrcdComponent returns true if the Image repository (id) belongs to src-d
func isSrcdComponent(id string) bool {
	namespace := strings.Split(id, "/")[0]
	return stringInSlice(srcdNamespaces, namespace)
}

func isFromEngine(name string) bool {
	return strings.HasPrefix(name, "srcd-cli-")
}
