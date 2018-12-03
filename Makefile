BIN      := pe-code-mgr-webhook-adapter
VERSION  := $(shell git describe --always --dirty --tags)
REGISTRY ?= andrewm3
IMAGE    := $(REGISTRY)/$(BIN)

all: test build

test:
	go test -v ./...

build:
	docker build -t $(IMAGE):$(VERSION) .

push:
	docker push $(IMAGE):$(VERSION)

.PHONY: test build
