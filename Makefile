# Makefile for common development tasks.
#
# Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
# SPDX-License-Identifier: BSD-3-Clause
#
# This file is part of neon <https://github.com/bow/neon>.

APP_NAME   := neon
REPO_NAME  := $(shell (head -n 1 go.mod | sed 's/^module //g'))

BIN_DIR  ?= $(CURDIR)/bin
BIN_NAME ?= $(APP_NAME)
BIN_PATH := $(BIN_DIR)/$(BIN_NAME)
ifeq ($(shell go env GOOS 2> /dev/null),windows)
BIN_PATH := $(addsuffix .exe,$(BIN_PATH))
endif

# Linker flags for go-build
# BASE_LD_FLAGS are linker flags that can not be overwritten.
GIT_TAG       := $(shell git describe --tags --always --dirty 2> /dev/null || echo "untagged")
GIT_COMMIT    := $(shell git rev-parse --quiet --verify HEAD || echo "?")
GIT_DIRTY     := $(shell test -n "$(shell git status --porcelain)" && echo "-dirty" || true)
BASE_LD_FLAGS := -X ${REPO_NAME}/internal.version=$(GIT_TAG)
BASE_LD_FLAGS += -X ${REPO_NAME}/internal.gitCommit=$(GIT_COMMIT)$(GIT_DIRTY)

# Allow for optional LD flags from env, appended to base flags, stripping trailing whitespaces.
LD_FLAGS := $(strip $(BASE_LD_FLAGS) $(LD_FLAGS))

CGO_ENABLED ?= 0

# Protobuf settings.
PROTO_DIR := $(CURDIR)/api
PROTO_FILES := $(shell find $(PROTO_DIR) -type f -name "*.proto" -print)

# DB settings.
DEV_DB_FILE := dev.db

# Because these tools' counterparts in nixpkgs do not work or are unavailable.
GOPATH := $(shell go env GOPATH)
GOPATH_BIN = $(GOPATH)/bin
MOCKGEN_EXE = $(GOPATH_BIN)/mockgen
NANCY_EXE = $(GOPATH_BIN)/nancy

all: help


.PHONY: bin
bin: $(BIN_PATH)  ## Compile an executable binary.

$(BIN_PATH): $(shell find . -type f -name '*.go' -print | grep -v '_mock_test') go.mod
	go mod tidy && CGO_ENABLED=$(CGO_ENABLED) go build -ldflags '$(LD_FLAGS)' -o $@


.PHONY: clean
clean:  ## Remove all build artifacts.
	rm -f bin/* coverage.html .coverage.out .junit.xml $(DEV_DB_FILE) result && (docker rmi ghcr.io/bow/$(APP_NAME) 2> /dev/null || true)


.PHONY: dev
dev:  ## Install dependencies for local development.
	@if command -v nix-env > /dev/null && command -v direnv > /dev/null; then \
		printf "Configuring a local dev environment and setting up git pre-commit hooks...\n" >&2 \
			&& direnv allow . > /dev/null \
			&& DIRENV_LOG_FORMAT="" direnv exec $(CURDIR) go install go.uber.org/mock/mockgen@v0.4.0 \
			&& DIRENV_LOG_FORMAT="" direnv exec $(CURDIR) go install github.com/sonatype-nexus-community/nancy@latest \
			&& printf "Done.\n" >&2; \
	elif command -v nix-env > /dev/null; then \
		printf "Error: direnv seems to be unconfigured or missing\n" >&2 && exit 1; \
	elif command -v direnv > /dev/null; then \
		printf "Error: nix seems to be unconfigured or missing\n" >&2 && exit 1; \
	else \
		printf "Error: both direnv and nix seem to be unconfigured and/or missing" >&2 && exit 1; \
	fi


.PHONY: fmt
fmt:  ## Apply gofmt.
	go fmt ./...


.PHONY: help
help:  ## Show this help.
	$(eval PADLEN=$(shell grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| cut -d':' -f1 \
		| awk '{cur = length($$0); lengths[cur] = lengths[cur] $$0 ORS; max=(cur > max ? cur : max)} END {printf "%s", max}' \
		|| (true && echo 0)))
	@(grep --version > /dev/null 2>&1 || (>&2 "error: GNU grep not installed"; exit 1)) \
		&& printf "\033[36m◉ %s dev console\033[0m\n" "$(APP_NAME)" >&2 \
		&& grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
			| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m» \033[33m%-*s\033[0m \033[36m· \033[0m%s\n", $(PADLEN), $$1, $$2}' \
			| sort


.PHONY: img
img:  ## Build a docker image and load it into a running daemon.
	nix build .#dockerArchiveStreamer && ./result | docker image load


.PHONY: lint
lint:  ## Lint the code.
	golangci-lint run


.PHONY: mocks
mocks: internal/datastore/parser_mock_test.go  ## Generate mocks from interfaces.
mocks: internal/server/datastore_mock_test.go
mocks: $(addprefix internal/reader/,operator_mock_test.go backend_mock_test.go state_mock_test.go)
mocks: internal/reader/backend/client_mock_test.go

internal/datastore/parser_mock_test.go: internal/datastore/parser.go
	$(MOCKGEN_EXE) -source=$< -package=datastore Parser > $@

internal/server/datastore_mock_test.go: internal/datastore/datastore.go
	$(MOCKGEN_EXE) -source=$< -package=server Datastore > $@

internal/reader/operator_mock_test.go: internal/reader/ui/operator.go
	$(MOCKGEN_EXE) -source=$< -package=reader Operator > $@

internal/reader/backend_mock_test.go: internal/reader/backend/backend.go
	$(MOCKGEN_EXE) -source=$< -package=reader Backend > $@

internal/reader/state_mock_test.go: internal/reader/state/state.go
	$(MOCKGEN_EXE) -source=$< -package=reader State > $@

internal/reader/backend/client_mock_test.go: api/neon_grpc.pb.go
	$(MOCKGEN_EXE) -source=$< -package=backend NeonClient > $@


.PHONY: protos
protos: $(PROTO_FILES)  ## Generate code from protobuf.
	@protoc \
		-I=$(PROTO_DIR) \
		--go_opt=Mneon.proto="$(REPO_NAME)/api;api" \
		--go-grpc_opt=Mneon.proto="$(REPO_NAME)/api;api" \
		--go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)


.PHONY: scan-sec
scan-sec:  scan-sec-ast scan-sec-deps  ## Perform all security analyses.


.PHONY: scan-sec-deps
scan-sec-deps:  ## Scan dependencies for reported vulnerabilities.
	go list -json -deps ./... | $(NANCY_EXE) sleuth


.PHONY: scan-sec-ast
scan-sec-ast:  ## Perform static security analysis on the AST.
	gosec -exclude=G304 ./...


.PHONY: serve
serve: bin  ## Compile the binary and run the server in development mode.
	$(BIN_PATH) server --db-path $(DEV_DB_FILE)


.PHONY: test .coverage.out
test: mocks .coverage.out  ## Run the test suite.

.coverage.out:
	gotestsum --format dots-v2 --junitfile .junit.xml -- ./... -coverprofile=$@.all -covermode=atomic -coverpkg ./internal/...,./cmd/...,./api,./. \
		&& grep -v "_mock_test.go" $@.all | grep -v "/api/" > $@ \
		&& go tool cover -func=$@


.PHONY: test-cov-xml
test-cov-xml: .coverage.out  ## Run the test suite and output coverage to XML.
	gocover-cobertura < $< > .coverage.xml


.PHONY: test-cov-html
test-cov-html: .coverage.out  ## Run the test suite and output coverage to HTML.
	go tool cover -html=.coverage.out -o coverage.html
