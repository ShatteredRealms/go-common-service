# Overview


# Development
## Requirements
**Required:**
* Golang ~> 1.23
**Preferred:**
* Make
* Docker
* Protoc

## Getting Started
The following commands download all dependencies and installs any required tools for development.
```bash
go mod download
make build-tools
```

Anytime changes to the gRPC API are made, the following command should be run to generate protobuf files and update the mocks used for testing.
```bash
make protos
```

## Building
The application version is automatically taken as the latest version tag.
**Binary:** `make build` the output result will be placed in the `bin` folder in the project root directory.\ 
**Docker:** \
    `make docker` builds a docker image called `sro-<app-name>` with 3 tags variants `latest`, `${version}`, and `${version}-${commit hash}`.\
    `make push` push the image to the docker repository.\
    `make docker push` runs docker then push.\

Note: A helper command `make aws-docker-login` is available to authenticate with the default aws credential context if pushing to AWS ECR.

## Testing
`make test` runs all tests and builds a coverage report, and `test-watch` runs the tests with coverage on all `.go` file changes.
`make report` views the coverage report in the browser, and `test-watch` updates coverage results as they update.

## Building Container
Deployment is done using docker. If using an AWS docker repository, running `make aws-docker-login` will authenticate with the default aws credential context. To push the images, run `make push`. This will build the image and push them to the docker repository.
