MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
ROOT := $(dir $(MKFILE_PATH))
GOBIN ?= $(ROOT)/tools/bin
ENV_PATH = PATH=$(GOBIN):$(PATH)
BIN_PATH ?= $(ROOT)/bin
LINTER_NAME := golangci-lint
LINTER_VERSION := v1.60.2

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

.PHONY: run-bot
run-bot:
	$(BIN_PATH)/bot -config=./config/config-bot.yaml

.PHONY: compose-sync-up
compose-sync-up:
	docker-compose -f ./script/docker/docker-compose-sync.yml up --build

.PHONY: compose-db-sync-up
compose-db-sync-up:
	docker-compose -f ./script/docker/docker-compose-db-sync.yml up --build

.PHONY: compose-async-up
compose-async-up:
	docker-compose -f ./script/docker/docker-compose-async.yml up --build

.PHONY: compose-db-async-up
compose-db-async-up:
	docker-compose -f ./script/docker/docker-compose-db-async.yml up --build

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

.PHONY: tools
tools: install-linter
	@if [ ! -f $(GOBIN)/mockgen ]; then\
		echo "Installing mockgen";\
		GOBIN=$(GOBIN) go install go.uber.org/mock/mockgen@v0.5.0;\
	fi

.PHONY: generate
generate: tools
	$(ENV_PATH) go generate ./...
