.PHONY: build

MAYOR ?= 0
MINOR ?= 0
PATCH ?= 0-develop
BUILD = $(shell git rev-parse --short HEAD)

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