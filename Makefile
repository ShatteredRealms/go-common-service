#)####################################################################################
#   _____ _           _   _                    _   _____            _               #
#  / ____| |         | | | |                  | | |  __ \          | |              #
# | (___ | |__   __ _| |_| |_ ___ _ __ ___  __| | | |__) |___  __ _| |_ __ ___  ___ #
#  \___ \| '_ \ / _` | __| __/ _ \ '__/ _ \/ _` | |  _  // _ \/ _` | | '_ ` _ \/ __|#
#  ____) | | | | (_| | |_| ||  __/ | |  __/ (_| | | | \ \  __/ (_| | | | | | | \__ \#
# |_____/|_| |_|\__,_|\__|\__\___|_|  \___|\__,_| |_|  \_\___|\__,_|_|_| |_| |_|___/#
#####################################################################################

#
# Makefile for building, running, and testing
#

APP_NAME = go-common-service

# Import dotenv
ifneq (,$(wildcard ../.env))
	include ../.env
	export
endif

# Application versions
BASE_VERSION = $(shell git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')
COMMIT_HASH = $(shell git rev-parse --short HEAD)

COVERAGE_FILE=coverage.out

# Gets the directory containing the Makefile
ROOT_DIR = $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# Base container registry
SRO_BASE_REGISTRY ?= docker.io
SRO_REGISTRY ?= $(SRO_BASE_REGISTRY)/sro

# The registry for this service
REGISTRY = $(SRO_REGISTRY)/$(APP_NAME)
time=$(shell date +%s)

PROTO_DIR=$(ROOT_DIR)/api

PROTO_FILES = $(shell find "$(PROTO_DIR)/sro" -name '*.proto')

MOCK_INTERFACES = $(shell egrep -rl --include="*.go" "type (\w*) interface {" $(ROOT_DIR)/pkg | sed "s/.go$$//")

# Versioning
VERSION=$(BASE_VERSION)
ifeq ($(VERSION),)
	VERSION := 0.0.0
endif

VERSION_PARTS=$(subst ., ,$(VERSION))
MAJOR_VERSION=$(word 1,$(VERSION_PARTS))
MINOR_VERSION=$(word 2,$(VERSION_PARTS))
PATCH_VERSION=$(word 3,$(VERSION_PARTS))


#   _____                    _
#  |_   _|                  | |
#    | | __ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    \_/\__,_|_|  \__, |\___|\__|___/
#                  __/ |
#                 |___/

.PHONY: test report mocks clean-mocks report-watch $(APP_NAME)
test:
	ginkgo --randomize-all -p --cover -covermode atomic -coverprofile=$(COVERAGE_FILE) --output-dir $(ROOT_DIR)/ --output-interceptor-mode=none $(ROOT_DIR)/pkg/...

test-watch:
	ginkgo watch --randomize-all -p --cover -covermode atomic -coverprofile=$(COVERAGE_FILE) -output-dir=$(ROOT_DIR) $(ROOT_DIR)/...

report: test
	go tool cover -func=$(ROOT_DIR)/$(COVERAGE_FILE) -o $(ROOT_DIR)/coverage.txt
	go tool cover -html=$(ROOT_DIR)/$(COVERAGE_FILE) -o $(ROOT_DIR)/coverage.html

report-watch:
	while inotifywait -e close_write $(ROOT_DIR)/$(COVERAGE_FILE); do \
		go tool cover -func=$(ROOT_DIR)/$(COVERAGE_FILE) -o $(ROOT_DIR)/coverage.txt; \
		go tool cover -html=$(ROOT_DIR)/$(COVERAGE_FILE) -o $(ROOT_DIR)/coverage.html; \
	done

dev-watch: test-watch report-watch

mocks: $(MOCK_INTERFACES)
$(MOCK_INTERFACES):
	mockgen \
		-source="$@.go" \
		-destination="$(@D)/mocks/$(@F).go"

run:
	go run $(ROOT_DIR)/cmd/$(APP_NAME)

run-watch:
	gow run $(ROOT_DIR)/cmd/$(APP_NAME)

.PHONY: clean-protos protos $(PROTO_FILES)

clean-protos:
	rm -rf "$(ROOT_DIR)/pkg/pb"

protos: clean-protos $(PROTO_FILES) move-protos mocks

$(PROTO_FILES):
	protoc "$@" \
		-I "$(PROTO_DIR)" \
		--go_out="$(ROOT_DIR)" \
		--go-grpc_out="$(ROOT_DIR)" \
		--grpc-gateway_out="$(ROOT_DIR)" \
		--grpc-gateway_opt "logtostderr=true"

move-protos:
	mv -v "$(ROOT_DIR)/github.com/ShatteredRealms/$(APP_NAME)/pkg/pb" "$(ROOT_DIR)/pkg/"
	rm -r "$(ROOT_DIR)/github.com"

install-tools:
	  cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %@latest

git: git-patch
git-major:
	git tag v$(shell echo $(MAJOR_VERSION)+1 | bc).0.0
	git push
	git push --tags
git-minor:
	git tag v$(MAJOR_VERSION).$(shell echo $(MINOR_VERSION)+1 | bc).0 
	git push
	git push --tags
git-patch:
	git tag v$(MAJOR_VERSION).$(MINOR_VERSION).$(shell echo $(PATCH_VERSION)+1 | bc)
	git push
	git push --tags
