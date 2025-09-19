.PHONY: build run run-build lint docker-build generate tools

MAYOR ?= 0
MINOR ?= 0
PATCH ?= 0-develop
BUILD = $(shell git rev-parse --short HEAD)

# VERSIONS
lint_v = v2.0.2

# GLOBAL ENVS
GOCACHE = $(shell go env GOCACHE)
GOMODCACHE = $(shell go env GOMODCACHE)
GOOS ?= $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

.PHONY: build
build:
	go tool goreleaser build --clean --snapshot

.PHONY: run
run:
	go run cmd/main.go $(args)

.PHONY: run-build
run-build: build
	./dist/pm $(args)

.PHONY: lint
lint:
	go tool golangci-lint run -v

.PHONY: lint-fix
lint-fix:
	go tool golangci-lint run -v --fix

.PHONY: bdocker-build
docker-build:
	docker buildx build \
			--progress=plain \
			--platform linux/amd64 \
			--build-arg GOCACHE=$(GOCACHE) \
			--build-arg GOMODCACHE=$(GOMODCACHE) \
			--output "type=docker" \
			-t jlrosende/pm:latest \
			-f build/docker/Dockerfile \
			.

.PHONY: gendocs
gendocs: gendocs-md gendocs-man gendocs-rest

gendocs-md:
	go run ./internal/tools/docgen -out ./docs/cli -format markdown
gendocs-man:
	go run ./internal/tools/docgen -out ./man -format man
gendocs-rest:
	go run ./internal/tools/docgen -out ./docs/rest -format rest

.PHONY: release
release:
	go tool run

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test: unit integration

.PHONY: unit
unit:
	go  test ./tests/... -v -tags unit

.PHONY: integration
integration:
	go  test ./tests/... -v -tags integration
