OWNER    = nicholasdille
PACKAGE  = insulatr
IMAGE    = $(OWNER)/$(PACKAGE)
STATIC   = insulatr-$(shell uname -m)
SOURCE   = $(shell echo *.go)
PWD      = $(shell pwd)
BIN      = $(PWD)/bin
TOOLS    = $(PWD)/tools
GOMOD    = $(PWD)/go.mod
GOFMT    = gofmt
SEMVER   = $(TOOLS)/semver
BUILDDEF = insulatr.yaml

GIT_COMMIT = $(shell git rev-list -1 HEAD)
BUILD_TIME = $(shell date +%Y%m%d-%H%M%S)
GIT_TAG = $(shell git describe --tags 2>/dev/null)

M = $(shell printf "\033[34;1m▶\033[0m")

.DEFAULT_GOAL := $(PACKAGE)

.PHONY: clean prepare deps deppatch depupdate deptidy format linter check static binary check-docker docker test run check-changes $(PACKAGE) bump-% build-% release-% tag-% changelog changelog-% release $(IMAGE)-% check-tag push-% latest-%

.SECONDARY:

clean: clean-docker; $(info $(M) Cleaning...)
	@rm -rf $(BIN)
	@rm -rf $(TOOLS)

prepare: | $(BIN) $(TOOLS) $(SEMVER)

$(BIN): ; $(info $(M) Preparing binary...)
	@mkdir -p $(BIN)

$(TOOLS): ; $(info $(M) Preparing tools...)
	@mkdir -p $(TOOLS)

##################################################
# TOOLS
##################################################

deps: $(GOMOD)

deppatch: ; $(info $(M) Updating dependencies to the latest patch...)
	@go get -u=patch

depupdate: ; $(info $(M) Updating dependencies to the latest version...)
	@go get -u

deptidy: ; $(info $(M) Updating dependencies to the latest version...)
	@go mod tidy

$(GOMOD): ; $(info $(M) Initializing dependencies...)
	@test -f go.mod || go mod init

format: ; $(info $(M) Running formatter...)
	@gofmt -l -w $(SOURCE)

# go get github.com/golang/lint/golint
lint: ; $(info $(M) Running linter...)
	@golint $(PACKAGE)

# go get github.com/KyleBanks/depth/cmd/depth
deptree: ; $(info $(M) Creating dependency tree...)
	@depth .

semver: $(SEMVER)

$(SEMVER): $(TOOLS) ; $(info $(M) Installing semver...)
	@test -f $@ && test -x $@ || ( \
		curl -sLf https://github.com/fsaintjacques/semver-tool/raw/2.1.0/src/semver > $@; \
		chmod +x $@; \
	)

##################################################
# BUILD
##################################################

check: format lint

%.sha256: % ; $(info $(M) Creating SHA256 for $*...)
	@echo sha256sum $* > $@

%.asc: % ; $(info $(M) Creating signature for $*...)
	@gpg --local-user $$(git config --get user.signingKey) --sign --armor --detach-sig --yes $*

binary $(PACKAGE): $(BIN)/$(PACKAGE)

$(BIN)/$(PACKAGE): $(SOURCE) | prepare ; $(info $(M) Building $(PACKAGE)...)
	@go build -ldflags "-s -w -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME) -X main.Version=$(GIT_TAG)" -o $@ $(SOURCE)

static: $(BIN)/$(STATIC)

$(BIN)/$(STATIC): $(SOURCE) | prepare ; $(info $(M) Building static $(PACKAGE)...)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags "-s -w -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME) -X main.Version=$(GIT_TAG)" -o $@ $(SOURCE)

##################################################
# TEST
##################################################

scp-%: $(BIN)/$(PACKAGE) ; $(info $(M) Copying to $*)
	@tar -cz bin/$(PACKAGE) $(BUILDDEF) | ssh $* tar -xvz

ssh-%: scp-% ; $(info $(M) Running remotely on $*)
	@ssh $* ./bin/$(PACKAGE) --file $(BUILDDEF) $(PARAMS)
