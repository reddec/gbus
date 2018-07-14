BUILD_ENVPARMS:=CGO_ENABLED=0 CC=gcc
PROJECT_DIR?=$(shell pwd)

.PHONY: test
test:
	$(BUILD_ENVPARMS) vgo test -v ./...

.PNONY: build
build:
	$(BUILD_ENVPARMS) vgo build -o bin/events-bus-gen ./cmd/events-bus-gen

.PHONY: gen
gen: build
	vgo generate ./...