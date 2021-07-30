GOROOT =
PATH := ${PWD}/cache/go/bin:${PWD}/cache/go/misc/wasm:${PATH}
export
LINT_VERSION = 1.41.1
SHELL := /usr/bin/env bash
GO_VERSION = 1.16.6

.PHONY: all
all: lint test

.PHONY: lint-deps
lint-deps:
	@if ! which golangci-lint >/dev/null || [[ "$$(golangci-lint version 2>&1)" != *${LINT_VERSION}* ]]; then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v${LINT_VERSION}; \
	fi

.PHONY: lint
lint: lint-deps go
	golangci-lint run
	GOOS=js GOARCH=wasm golangci-lint run --build-tags js,wasm

.PHONY: test-deps
test-deps:
	@go install github.com/mattn/goveralls@v0.0.9

.PHONY: test
test: test-deps go
	go test -race -coverprofile=cover.out ./...
	GOOS=js GOARCH=wasm go test -cover ./...
	@if [[ "$$CI" == true ]]; then \
		set -ex; \
		goveralls -coverprofile=cover.out -service=github; \
	fi

cache:
	mkdir -p cache

.PHONY: go
go: cache/go${GO_VERSION}

cache/go${GO_VERSION}: cache
	if [[ ! -e cache/go${GO_VERSION} ]]; then \
		set -ex; \
		TMP=$$(mktemp -d); trap 'rm -rf "$$TMP"' EXIT; \
		git clone \
			--depth 1 \
			--single-branch \
			--branch hackpad-go${GO_VERSION} \
			https://github.com/hack-pad/go.git \
			"$$TMP"; \
		pushd "$$TMP/src"; \
		./make.bash; \
		export PATH="$$TMP/bin:$$PATH"; \
		go version; \
		mkdir -p ../bin/js_wasm; \
		go build -o ../bin/js_wasm/ std cmd/go cmd/gofmt; \
		go tool dist test -rebuild -list; \
		go build -o ../pkg/tool/js_wasm/ std cmd/buildid cmd/pack cmd/cover cmd/vet; \
		go install ./...; \
		popd; \
		mv "$$TMP" cache/go${GO_VERSION}; \
		ln -sfn go${GO_VERSION} cache/go; \
	fi
	touch cache/go${GO_VERSION}
	touch cache/go.mod  # Makes it so linters will ignore this dir

