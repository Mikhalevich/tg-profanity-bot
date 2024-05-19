MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
ROOT := $(dir $(MKFILE_PATH))
GOBIN ?= $(ROOT)/tools/bin
BIN_PATH ?= $(ROOT)/bin
LINTER_NAME := golangci-lint
LINTER_VERSION := v1.57.2

all: build

.PHONY: build
build:
	go build -mod=vendor -o $(BIN_PATH)/bot ./cmd/bot/main.go
	go build -mod=vendor -o $(BIN_PATH)/consumer ./cmd/consumer/main.go

.PHONY: test
test:
	go test ./...

.PHONY: bench
bench:
	go test -bench=. ./...

.PHONY: run
run:
	$(BIN_PATH)/bot -config=./config/config.yaml

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: install-linter
install-linter:
	if [ ! -f $(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) ]; then \
		echo INSTALLING $(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) $(LINTER_VERSION) ; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN)/$(LINTER_VERSION) $(LINTER_VERSION) ; \
		echo DONE ; \
	fi

.PHONY: lint
lint: install-linter
	$(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) run --config .golangci.yml
