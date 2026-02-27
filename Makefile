.PHONY: build test lint generate generate-check install

BINARY := yougile
CMD_PATH := ./cmd/yougile
PREFIX ?= /usr/local

build:
	go build -o bin/$(BINARY) $(CMD_PATH)

install: build
	install -d $(PREFIX)/bin
	install -m 755 bin/$(BINARY) $(PREFIX)/bin/$(BINARY)

test:
	go test ./...

lint:
	golangci-lint run ./...

generate:
	oapi-codegen -package client -generate types,client -o pkg/client/api.gen.go docs/api.json

# CI: ensure generated code is up to date (run after make generate, then git diff)
generate-check: generate
	@git diff --exit-code pkg/client/ docs/api.json || (echo "run 'make generate' and commit pkg/client/ and docs/api.json"; exit 1)
