BUILD_ENVPARMS:=CGO_ENABLED=0 CC=gcc

.PHONY: test
test:
	$(BUILD_ENVPARMS) vgo test -v ./...