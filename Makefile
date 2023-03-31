BUILD_DIR ?= .build
PROJECT_NAME ?= $(shell basename $(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

include scripts/makefiles/third_party/pasdam/makefiles/go.mk
include scripts/makefiles/third_party/pasdam/makefiles/go.mod.mk
include scripts/makefiles/third_party/pasdam/makefiles/help.mk

.DEFAULT_GOAL := help

## build: Build all artifacts (binary and docker image)
.PHONY: build
build: | go-build

## clean: Remove all artifacts (binary and docker image)
.PHONY: clean
clean: | go-clean
