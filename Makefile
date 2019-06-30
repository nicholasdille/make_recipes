PWD=$(shell pwd)
TOOLS=$(PWD)/tools
SEMVER=$(shell test -f $(TOOLS)/semver && test -x $(TOOLS)/semver || ( mkdir -p $(TOOLS) && curl -sLf https://github.com/fsaintjacques/semver-tool/raw/2.1.0/src/semver > $(TOOLS)/semver && chmod +x $(TOOLS)/semver && echo $(TOOLS)/semver))

GITHUB_OWNER=nicholasdille
GITHUB_REPOSITORY=insulatr
RELEASE_VERSION=1.0.4
MAJOR_VERSION = $(shell $(SEMVER) get major $(RELEASE_VERSION))
MINOR_VERSION = $(shell $(SEMVER) get minor $(RELEASE_VERSION))

RELEASE_ASSETS=$(wildcard Makefile.*)

DOCKER_OWNER=nicholasdille
DOCKER_REPOSITORY=test
DOCKER_TAG=1.0

include Makefile.docker
include Makefile.release
#include Makefile.go

.PHONY: docker

.SECONDARY:

docker: .docker/$(DOCKER_OWNER)/$(DOCKER_REPOSITORY)/$(DOCKER_TAG).image .docker/$(DOCKER_OWNER)/$(DOCKER_REPOSITORY)/1.tag
