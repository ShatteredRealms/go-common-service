# Overview [![codecov](https://codecov.io/gh/ShatteredRealms/go-common-service/graph/badge.svg?token=8nwCfHwBSh)](https://codecov.io/gh/ShatteredRealms/go-common-service) [![Unit and Integration Tests](https://github.com/ShatteredRealms/go-common-service/actions/workflows/test.yml/badge.svg)](https://github.com/ShatteredRealms/go-common-service/actions/workflows/test.yml)
Common utilities and packages for Shattered Realms golang microservices.

## Packages
* `api`
    * Common proto file definitions and gRPC API definition for common services
* `cmd`
    * Short lived tests for testing all packages
* `pkg/bus`
    * Event bus implementation for pub/sub
    * Event bus message definitions and basic repository and service for microservices to use as a base
* `pkg/config`
    * Base configuration layout used by microservices
* `pkg/log`
    * Logging package common used by microservices
* `pkg/mocks`
    * Generated mocks for testing
* `pkg/model`
    * Base database model definitions
    * Common database models
* `pkg/pb`
    * Generated protobuf files for gRPC services and proto definitions
* `pkg/repository`
    * Setup connection to postgres database with redis cache and opentelemetry tracing
    * Setup connection to kafka
* `pkg/srospan`
    * OpenTeleemtry Span attributes for microservices
* `pkg/srv`
    * Basic health gRPC service implementation
    * Common errors
    * Base service context for microservices

# Development
## Requirements
**Required:**
* Golang ~> 1.23

**Preferred:**
* Make
* Docker
* Protoc
* inotifywait

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

## Testing
* `make test` runs all tests and builds a coverage report
* `make report` views the coverage report in the browser
* `make test-watch` runs the tests with coverage on all `.go` file changes.
* `make report-watch` updates coverage results as they update.
* `make dev-watch -j` simply runs both `test-watch` and `report-watch` in parallel.

## Versioning
Run `make git` to generate a new version tag by incrementing the patch number. This is the same as `make git-patch`. To increment the minor or major version, run `make git-minor` or `make git-major` respectively.
