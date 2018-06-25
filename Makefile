
BINARY_NAME := $(notdir $(shell pwd))
default: build

.PHONY: build
build:
	go build -a -o $(BINARY_NAME) .
