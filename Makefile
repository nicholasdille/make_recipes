########## Global
OWNER=nicholasdille
PROJECT=test
VERSION=1.0.0
PWD=$(shell pwd)
TOOLS=$(PWD)/tools
SEMVER=$(shell test -f $(TOOLS)/semver && test -x $(TOOLS)/semver || ( mkdir -p $(TOOLS) && curl -sLf https://github.com/fsaintjacques/semver-tool/raw/2.1.0/src/semver > $(TOOLS)/semver && chmod +x $(TOOLS)/semver && echo $(TOOLS)/semver))

########## Docker
DOCKER_OWNER=$(OWNER)
DOCKER_REPOSITORY=$(PROJECT)
DOCKER_TAG=$(VERSION)
include Makefile.docker

########## Release
GITHUB_OWNER=$(OWNER)
GITHUB_REPOSITORY=$(PROJECT)
RELEASE_VERSION=$(VERSION)
MAJOR_VERSION = $(shell $(SEMVER) get major $(RELEASE_VERSION))
MINOR_VERSION = $(shell $(SEMVER) get minor $(RELEASE_VERSION))
RELEASE_ASSETS=$(wildcard Makefile.*)
include Makefile.release

########## Go
GO_PACKAGE=$(PROJECT)
include Makefile.golang

########## Custom
.PHONY: docker

.SECONDARY:

docker: .docker/$(DOCKER_OWNER)/$(DOCKER_REPOSITORY)/$(DOCKER_TAG).image .docker/$(DOCKER_OWNER)/$(DOCKER_REPOSITORY)/1.tag
