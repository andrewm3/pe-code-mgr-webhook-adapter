IMAGE   ?= pe-code-mgr-webhook-adapter
VERSION := $(shell git describe --always --dirty)

all: test build

test:
	go test -v ./...

build:
	docker build -t $(IMAGE):$(VERSION) .

.PHONY: test build
