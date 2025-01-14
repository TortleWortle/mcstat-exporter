EXECUTABLE=mcstat-exporter
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64

VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"


GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
INTERNAL := $(wildcard *.go)

all: clean build

build: windows linux darwin

install:
	go install ./cmd/$(EXECUTABLE)

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o build/$(EXECUTABLE)/$(WINDOWS) ./cmd/$(EXECUTABLE)

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/$(EXECUTABLE)/$(LINUX) ./cmd/$(EXECUTABLE)

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o build/$(EXECUTABLE)/$(DARWIN) ./cmd/$(EXECUTABLE)

run:
	go run ./cmd/$(EXECUTABLE)

.PHONY: test
test:
	go test ./internal

clean:
	rm -rf build