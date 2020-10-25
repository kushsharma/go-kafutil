.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --no-builtin-rules
VERSION=`cat version`
BUILD=`date +%FT%T%z`
#COMMIT=`git rev-parse HEAD`
COMMIT=`date +%FT%T%z`
EXECUTABLE="kafutil"

# first command used as the default one if only `make` is used
all: build

PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC := $(shell which protoc)
# If protoc isn't on the path, set it to a target that's never up to date, so
# the install command always runs.
ifeq ($(PROTOC),)
    PROTOC = must-rebuild
endif

# Figure out which machine we're running on.
UNAME := $(shell uname)

$(PROTOC):
# Run the right installation command for the operating system.
ifeq ($(UNAME), Darwin)
	brew install protobuf
endif
ifeq ($(UNAME), Linux)
	sudo apt-get install protobuf-compiler
endif
# You can add instructions for other operating systems here, or use different
# branching logic as appropriate.

# If $GOPATH/bin/protoc-gen-go does not exist, we'll run this command to install
# it.
$(PROTOC_GEN_GO):
	@echo "> checking dependencies"
	go get -u github.com/golang/protobuf/protoc-gen-go

.PHONY: build test clean generate dist init build_linux build_mac

app.pb.go: ./protos/app.proto | $(PROTOC_GEN_GO) $(PROTOC)
	@echo "> building protos"
	protoc --include_imports --go_out=. --go_opt=paths=source_relative ./protos/app.proto --descriptor_set_out=./protos/desc.set

build: app.pb.go
	@go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}' -X 'main.AppName=${EXECUTABLE}'" -o ${EXECUTABLE} ./main.go
	@echo "> build complete"

run: build
	@./${EXECUTABLE}

clean:
	@rm -rf ${EXECUTABLE} dist/

build_nix:
	@env GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}' -X 'main.AppName=${EXECUTABLE}'" -o ${EXECUTABLE} ./main.go

build_mac:
	@env GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}' -X 'main.AppName=${EXECUTABLE}'" -o ${EXECUTABLE} ./main.go
