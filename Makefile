# Common development tasks.

# Cross-platform adjustments.
SYS := $(shell uname 2> /dev/null)
ifeq ($(SYS),Linux)
DATE_EXE := date
GREP_EXE := grep
SED_EXE  := sed
else ifeq ($(SYS),Darwin)
DATE_EXE := gdate
GREP_EXE := ggrep
SED_EXE  := gsed
else
$(error error: unsupported development platform)
endif

APP_NAME   := courier
GO_VERSION := $(shell (head -n 3 go.mod | $(SED_EXE) 's/^go//g' | tail -n 1))
REPO_NAME  := $(shell (head -n 1 go.mod | $(SED_EXE) 's/^module //g'))

GIT_TAG    := $(shell git describe --tags --always --dirty 2> /dev/null || echo "untagged")
GIT_COMMIT := $(shell git rev-parse --quiet --verify HEAD || echo "?")
GIT_DIRTY  := $(shell test -n "$(shell git status --porcelain)" && echo "-dirty" || true)
BUILD_TIME := $(shell $(DATE_EXE) -u '+%Y-%m-%dT%H:%M:%SZ')
IS_RELEASE := $(shell ((echo "${GIT_TAG}" | $(GREP_EXE) -qE "^v?[0-9]+\.[0-9]+\.[0-9]+$$") && echo '1') || true)

IMG_NAME   := ghcr.io/bow/$(APP_NAME)
ifeq ($(IS_RELEASE),1)
IMG_TAG    := $(GIT_TAG)
else
IMG_TAG    := latest
endif

BIN_DIR  ?= $(CURDIR)/bin
BIN_NAME ?= $(APP_NAME)
BIN_PATH := $(BIN_DIR)/$(BIN_NAME)
ifeq ($(shell go env GOOS 2> /dev/null),windows)
BIN_PATH := $(addsuffix .exe,$(BIN_PATH))
endif

# Linker flags for go-build
# BASE_LD_FLAGS are linker flags that can not be overwritten.
BASE_LD_FLAGS := -X ${REPO_NAME}/internal.version=$(GIT_TAG)
BASE_LD_FLAGS += -X ${REPO_NAME}/internal.buildTime=$(BUILD_TIME)
BASE_LD_FLAGS += -X ${REPO_NAME}/internal.gitCommit=$(GIT_COMMIT)$(GIT_DIRTY)

# Allow for optional LD flags from env, appended to base flags, stripping trailing whitespaces.
LD_FLAGS := $(strip $(BASE_LD_FLAGS) $(LD_FLAGS))

# Protobuf settings.
PROTOC_VERSION := 3.20.1
PROTOC_GEN_GO_GRPC_VERSION := v1.2.0
GO_PROTOBUF_VERSION := $(shell (cat go.mod | $(GREP_EXE) google.golang.org/protobuf | $(SED_EXE) -r 's/[[:space:]]+google.golang.org\/protobuf //g'))
PROTO_DIR := $(CURDIR)/api
PROTO_FILES := $(shell find $(PROTO_DIR) -type f -name "*.proto" -print)

# DB settings.
GOLANG_MIGRATE_VERSION := $(shell (cat go.mod | $(GREP_EXE) github.com/golang-migrate/migrate/v4 | $(SED_EXE) -r 's/[[:space:]]+github.com\/golang-migrate\/migrate\/v4 //g'))
DEV_DB_FILE := dev.db

all: help


.PHONY: bin
bin: $(BIN_PATH)  ## Compile an executable binary.

$(BIN_PATH): $(shell find . -type f -name '*.go' -print) go.mod
	go mod tidy && go build -trimpath -ldflags '$(LD_FLAGS)' -o $@


.PHONY: clean
clean:  ## Remove all build artifacts.
	rm -f bin/* coverage.html .coverage.out .junit.xml $(DEV_DB_FILE) && (docker rmi $(IMG_NAME) 2> /dev/null || true)


.PHONY: img
img:  ## Build and tag the Docker container.
	docker build --build-arg REVISION=$(GIT_COMMIT)$(GIT_DIRTY) --build-arg BUILD_TIME=$(BUILD_TIME) --tag $(IMG_NAME):$(IMG_TAG) .


.PHONY: install-dev
install-dev:  ## Install dependencies for local development.
	@if command -v asdf 1>/dev/null 2>&1; then \
		if [ ! -f .tool-versions ]; then \
			(asdf plugin add golang 2> /dev/null || true) \
				&& asdf install golang $(GO_VERSION) \
				&& asdf local golang $(GO_VERSION) > .tool-versions; \
			(asdf plugin add protoc 2> /dev/null || true) \
				&& asdf install protoc $(PROTOC_VERSION) \
				&& asdf local protoc $(PROTOC_VERSION) >> .tool-versions; \
		fi; \
		asdf reshim; \
	fi
	go install gotest.tools/gotestsum@v1.8.0 \
		&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION) \
		&& go install google.golang.org/protobuf/cmd/protoc-gen-go@$(GO_PROTOBUF_VERSION) \
		&& go install github.com/golang/mock/mockgen@v1.6.0 \
		&& go install github.com/securego/gosec/v2/cmd/gosec@latest \
		&& go install github.com/sonatype-nexus-community/nancy@latest \
		&& go install github.com/boumenot/gocover-cobertura@latest \
		&& go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.0 \
		&& go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@$(GOLANG_MIGRATE_VERSION)
	@if command -v asdf 1>/dev/null 2>&1; then \
		asdf reshim; \
	fi


.PHONY: help
help:  ## Show this help.
	$(eval PADLEN=$(shell $(GREP_EXE) -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| cut -d':' -f1 \
		| awk '{cur = length($$0); lengths[cur] = lengths[cur] $$0 ORS; max=(cur > max ? cur : max)} END {printf "%s", max}' \
		|| (true && echo 0)))
	@($(GREP_EXE) --version > /dev/null 2>&1 || (>&2 "error: GNU grep not installed"; exit 1)) \
		&& printf "\033[36m◉ %s dev console\033[0m\n" "$(APP_NAME)" >&2 \
		&& $(GREP_EXE) -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
			| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m» \033[33m%*-s\033[0m \033[36m· \033[0m%s\n", $(PADLEN), $$1, $$2}' \
			| sort


.PHONY: lint
lint:  ## Lint the code.
	golangci-lint run


.PHONY: proto
proto: $(PROTO_FILES) ## Generate code from protobuf.
	@protoc \
		-I=$(PROTO_DIR) \
		--go_opt=Mcourier.proto="$(REPO_NAME)/api;api" \
		--go-grpc_opt=Mcourier.proto="$(REPO_NAME)/api;api" \
		--go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)


.PHONY: scan-security
scan-security:  scan-security-ast scan-security-deps  ## Perform all security analyses.


.PHONY: scan-security-deps
scan-security-deps:  ## Scan dependencies for reported vulnerabilities.
	go list -json -deps ./... | nancy sleuth


.PHONY: scan-security-ast
scan-security-ast:  ## Perform static security analysis on the AST.
	gosec ./...


.PHONY: serve
serve: bin  ## Compile the binary and run the server in development mode.
	$(BIN_PATH) serve --db-path $(DEV_DB_FILE)


.PHONY: test .coverage.out
test: internal/internal_mock_test.go .coverage.out  ## Run the test suite.

.coverage.out:
	gotestsum --format dots-v2 --junitfile .junit.xml -- ./... -parallel=$(shell nproc) -coverprofile=$@ -covermode=atomic \
		&& go tool cover -func=$@

internal/internal_mock_test.go: internal/internal.go
	mockgen -source=$< -package=internal FeedParser,FeedStore > $@


.PHONY: test-cov-xml
test-cov-xml: .coverage.out  ## Run the test suite and output coverage to XML.
	gocover-cobertura < $< > .coverage.xml


.PHONY: test-cov-html
test-cov-html: .coverage.out  ## Run the test suite and output coverage to HTML.
	go tool cover -html=.coverage.out -o coverage.html
