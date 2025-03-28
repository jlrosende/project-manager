.PHONY: build run run-build lint docker-build

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

build:
	go build \
		-v \
		-ldflags "\
			-X 'github.com/jlrosende/project-manager/internal.mayor=$(MAYOR)' \
			-X 'github.com/jlrosende/project-manager/internal.minor=$(MINOR)' \
			-X 'github.com/jlrosende/project-manager/internal.patch=$(PATCH)' \
			-X 'github.com/jlrosende/project-manager/internal.build=$(BUILD)' \
		" \
		-o dist/pm \
		./cmd/main.go

run:
	go run cmd/main.go $(args)

run-build: build
	./dist/pm $(args)

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${lint_v} run -v

docker-build:
	docker buildx build \
			--progress=plain \
			--platform linux/amd64 \
			--build-arg GOCACHE=$(GOCACHE) \
			--build-arg GOMODCACHE=$(GOMODCACHE) \
			--output "type=docker" \
			-t sisusfox:latest \
			.
