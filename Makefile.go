########## Required variables
# GO_PACKAGE - container the name of the binary

########## Variables
STATIC   = $(GO_PACKAGE)-$(shell uname -m)
SOURCE   = $(shell echo *.go)
PWD      = $(shell pwd)
BIN      = $(PWD)/bin
GOMOD    = $(PWD)/go.mod
GOFMT    = gofmt

GIT_COMMIT = $(shell git rev-list -1 HEAD)
BUILD_TIME = $(shell date +%Y%m%d-%H%M%S)
GIT_TAG = $(shell git describe --tags 2>/dev/null)

M = $(shell printf "\033[34;1m▶\033[0m")

########## Targets: Global

.PHONY: clean prepare deps deppatch depupdate deptidy format linter check

.SECONDARY:

clean: clean-docker; $(info $(M) Cleaning...)
	@rm -rf $(BIN)

prepare: | $(BIN)

$(BIN): ; $(info $(M) Preparing binary...)
	@mkdir -p $(BIN)

########## Targets: Tools

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

########## Targets: Build

check: format lint

%.sha256: % ; $(info $(M) Creating SHA256 for $*...)
	@echo sha256sum $* > $@

%.asc: % ; $(info $(M) Creating signature for $*...)
	@gpg --local-user $$(git config --get user.signingKey) --sign --armor --detach-sig --yes $*

$(BIN)/$(PACKAGE): $(SOURCE) | prepare ; $(info $(M) Building $(PACKAGE)...)
	@go build -ldflags "-s -w -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME) -X main.Version=$(GIT_TAG)" -o $@ $(SOURCE)

$(BIN)/$(STATIC): $(SOURCE) | prepare ; $(info $(M) Building static $(PACKAGE)...)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags "-s -w -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME) -X main.Version=$(GIT_TAG)" -o $@ $(SOURCE)
