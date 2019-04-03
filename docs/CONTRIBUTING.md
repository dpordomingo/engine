# Contribution Guidelines

As all source{d} projects, this project follows the
[source{d} Contributing Guidelines](https://github.com/src-d/guide/blob/master/engineering/documents/CONTRIBUTING.md).


# Additional Contribution Guidelines

In addition to the [source{d} Contributing Guidelines](https://github.com/src-d/guide/blob/master/engineering/documents/CONTRIBUTING.md),
this project follows the following guidelines.


## Generated Code

Before submitting a pull request make sure all the generated code changes are also committed.

### Dependencies

Go dependencies are managed with [go Modules](https://github.com/golang/go/wiki/Modules). You can run these commands to ensure all project dependencies are up to date:

```shell
$ GO111MODULE=on go mod tidy
$ GO111MODULE=on go mod vendor
$ make no-changes-in-commit
```

### gRPC services

To generate the code for the gRPC services run:

```shell
$ make -C api proto
$ make no-changes-in-commit
```


## Testing

For unit and integration tests run:

```bash
$ make test
$ make test-integration
```

## Changelog

This project lists the important changes between releases in the [`CHANGELOG.md`](../CHANGELOG.md) file.

If you open a PR, you should also add a brief summary in the `CHANGELOG.md` mentioning the new feature, change or bugfix that you proposed.
