REVISION := $(shell git describe --always)
DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
LDFLAGS	:= -ldflags="-X \"main.Revision=$(REVISION)\" -X \"main.BuildDate=${DATE}\""

.PHONY: build deps clean

build:
	go build -o bin/galaxy $(LDFLAGS) cmd/galaxy/*.go
deps:
	dep ensure
clean:
	rm -f bin/galaxy
