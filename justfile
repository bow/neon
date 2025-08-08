# Justfile for common development tasks.
#
# Copyright (c) 2022-2025 Wibowo Arindrarto <contact@arindrarto.dev>
# SPDX-License-Identifier: BSD-3-Clause
#
# This file is part of neon <https://github.com/bow/neon>.

app-id    := 'neon'
repo-name := shell('head -n 1 go.mod | sed "s/^module //g"')

bin-dir  := env("BIN_DIR", justfile_directory()  / "bin")
bin-name := env("BIN_NAME", app-id)
bin-path := bin-dir / bin-name + (if os() == "windows" { ".exe" } else { "" })

go-exe  := require("go")
git-exe := require("git")

git-tag    := shell("git -C $1 describe --tags --always --dirty 2> /dev/null || echo 'untagged'", justfile_directory())
git-commit := shell("git -C $1 rev-parse --quiet --verify  HEAD || echo '?'", justfile_directory())
git-dirty  := if shell("git status --porcelain") == "" { "" } else { "-dirty" }

cgo-enabled   := env("CGO_ENABLED", "0")
base-ld-flags := "-X " + repo-name + "/internal.version=" + git-tag + " -X " + repo-name + "/internal.gitCommit=" + git-commit + git-dirty
ld-flags      := trim(base-ld-flags + " " + env("LD_FLAGS", ""))

proto-dir   := justfile_directory() / "api"
proto-files := replace(shell("find $1 -type f -name '*.proto' -print", proto-dir), "\n", " ")

dev-db-file := "dev.db"

# Tool paths
gopath := shell("go env GOPATH")
nancy-exe := gopath / "bin" / "nancy"


[private]
default: help

# Compile an executable binary
build-bin:
    @mkdir -p {{bin-dir}}
    go mod tidy
    CGO_ENABLED={{cgo-enabled}} go build -ldflags '{{ld-flags}}' -o {{bin-path}}

# Build a docker image and load it into a running daemon
build-img:
    nix build .#dockerArchiveStreamer && ./result | docker image load

# Remove all build artifacts
clean:
    @rm -f bin/* coverage.html .coverage.out .junit.xml {{dev-db-file}} result
    -@docker rmi ghcr.io/bow/{{app-id}} 2> /dev/null

# Apply gofmt
fmt:
    go fmt ./...

# Show this help and exit.
help:
    @just --list --list-prefix $'{{BOLD}}{{BLUE}}→{{NORMAL}} ' --justfile {{justfile()}} --list-heading $'{{BOLD}}{{CYAN}}◉ {{YELLOW}}{{app-id}}{{CYAN}} dev console{{NORMAL}}\n'

# Install dependencies for local development not yet in Nix
install-dev:
    #!/usr/bin/env bash
    if command -v nix-env > /dev/null && command -v direnv > /dev/null; then
        printf "Configuring a local dev environment...\n" >&2 \
            && direnv allow . > /dev/null \
            && DIRENV_LOG_FORMAT="" direnv exec {{justfile_directory()}} go install github.com/sonatype-nexus-community/nancy@latest \
            && printf "Done.\n" >&2
    elif command -v nix-env > /dev/null; then
        printf "Error: direnv seems to be unconfigured or missing\n" >&2 && exit 1
    elif command -v direnv > /dev/null; then
        printf "Error: nix seems to be unconfigured or missing\n" >&2 && exit 1
    else
        printf "Error: both direnv and nix seem to be unconfigured and/or missing" >&2 && exit 1
    fi

# Lint the code
lint:
    golangci-lint run

# Generate mocks from interfaces
gen-mocks:
    #!/usr/bin/env -S parallel --shebang --ungroup --jobs {{ num_cpus() }}
    mockgen -source=internal/datastore/parser.go -package=datastore Parser > internal/datastore/parser_mock_test.go
    mockgen -source=internal/datastore/datastore.go -package=server Datastore > internal/server/datastore_mock_test.go
    mockgen -source=internal/reader/ui/operator.go -package=reader Operator > internal/reader/operator_mock_test.go
    mockgen -source=internal/reader/backend/backend.go -package=reader Backend > internal/reader/backend_mock_test.go
    mockgen -source=internal/reader/state/state.go -package=reader State > internal/reader/state_mock_test.go
    mockgen -package=backend google.golang.org/grpc ServerStreamingClient > internal/reader/backend/client_grpc_mock_test.go
    mockgen -source=api/neon_grpc.pb.go -package=backend NeonClient > internal/reader/backend/client_mock_test.go

# Generate code from protobuf
gen-protos:
    protoc \
        -I={{proto-dir}} \
        --go_opt=Mneon.proto="{{repo-name}}/api;api" \
        --go-grpc_opt=Mneon.proto="{{repo-name}}/api;api" \
        --go_out={{proto-dir}} --go_opt=paths=source_relative \
        --go-grpc_out={{proto-dir}} --go-grpc_opt=paths=source_relative \
        {{proto-files}}

# Perform all security analyses
scan-sec: scan-sec-ast scan-sec-deps

# Scan dependencies for reported vulnerabilities
scan-sec-deps:
    go list -json -deps ./... | {{nancy-exe}} sleuth

# Perform static security analysis on the AST
scan-sec-ast:
    gosec -exclude=G304 ./...

# Compile the binary and run the server in development mode
serve db=(dev-db-file): build-bin
    {{bin-path}} server --db-path {{db}}

# Run the test suite
test: gen-mocks test-cov

[private]
test-cov:
    gotestsum --format dots-v2 --junitfile .junit.xml -- ./... -coverprofile=.coverage.out.all -covermode=atomic -coverpkg ./internal/...,./cmd/...,./api,./.
    @grep -v "_mock_test.go" .coverage.out.all | grep -v "/api/" > .coverage.out
    @go tool cover -func=.coverage.out

# Run the test suite and output coverage to XML
test-cov-xml: test-cov
    gocover-cobertura < .coverage.out > .coverage.xml

# Run the test suite and output coverage to HTML
test-cov-html: test-cov
    go tool cover -html=.coverage.out -o coverage.html

# Update dependencies and nix flake
update:
    nix flake update
    go get -u ./...
    go mod tidy
    gomod2nix
