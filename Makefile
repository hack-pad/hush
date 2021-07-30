LINT_VERSION = 1.41.1
SHELL := /usr/bin/env bash

.PHONY: all
all: lint test

.PHONY: lint-deps
lint-deps:
	@if ! which golangci-lint >/dev/null || [[ "$$(golangci-lint version 2>&1)" != *${LINT_VERSION}* ]]; then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v${LINT_VERSION}; \
	fi

.PHONY: lint
lint: lint-deps
	golangci-lint run
	GOOS=js GOARCH=wasm golangci-lint run --build-tags js,wasm

.PHONY: test
test:
	go test -race -cover ./...
	GOOS=js GOARCH=wasm go test -cover ./...

