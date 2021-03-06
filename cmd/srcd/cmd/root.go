// Copyright © 2018 Francesc Campoy <francesc@sourced.tech>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"context"
	"regexp"
	"time"

	"github.com/src-d/engine/cmd/srcd/daemon"

	"gopkg.in/src-d/go-cli.v0"
	"gopkg.in/src-d/go-log.v1"
)

const name = "srcd"

var rootCmd = cli.NewNoDefaults(name, "The Code as Data solution by source{d}")

// Command implements the default group flags. It is meant to be embedded into
// other application commands to provide default behavior for logging, config
type Command struct {
	cli.PlainCommand
	cli.LogOptions `group:"Log Options"`

	Config string `long:"config" description:"config file (default: $HOME/.srcd/config.yml)"`
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(&cli.CompletionCommand{Name: name})
	rootCmd.RunMain()
}

var logMsgRegex = regexp.MustCompile(`.*msg="(.+)"`)

func logAfterTimeout(msg string, timeout time.Duration) func() {
	d := newDefered(timeout, msg, nil, false, 0)
	return d.Print()
}

func logAfterTimeoutWithSpinner(msg string, timeout time.Duration, spinnerInterval time.Duration) func() {
	d := newDefered(timeout, msg, nil, true, spinnerInterval)
	return d.Print()
}

func logAfterTimeoutWithServerLogs(msg string, timeout time.Duration) func() {
	d := newDefered(timeout, msg, readDaemonLogs, false, 0)
	return d.Print()
}

func readDaemonLogs(stop <-chan bool) <-chan string {
	logs, err := daemon.GetLogs()
	if err != nil {
		log.Errorf(err, "could not get logs from server container")
		return nil
	}

	ch := make(chan string)
	go func() {
		defer logs.Close()
		scanner := bufio.NewScanner(logs)

		c := make(chan bool)
		scan := func() {
			c <- scanner.Scan()
		}

		go scan()
		for {
			select {
			case <-stop:
				close(ch)
				return
			case more := <-c:
				if !more {
					close(ch)
					if err := scanner.Err(); err != nil && err != context.Canceled {
						log.Errorf(err, "can't read logs from server")
					}

					return
				}

				match := logMsgRegex.FindStringSubmatch(scanner.Text())
				if len(match) == 2 {
					ch <- match[1]
				}

				go scan()
			}
		}

	}()

	return ch
}
