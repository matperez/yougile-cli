.PHONY: build test lint generate generate-check

BINARY := yougile
CMD_PATH := ./cmd/yougile

build:
	go build -o bin/$(BINARY) $(CMD_PATH)

test:
	go test ./...

lint:
	golangci-lint run ./...

generate:
	oapi-codegen -package client -generate types,client -o pkg/client/api.gen.go docs/api.json

# CI: ensure generated code is up to date (run after make generate, then git diff)
generate-check: generate
	@git diff --exit-code pkg/client/ docs/api.json || (echo "run 'make generate' and commit pkg/client/ and docs/api.json"; exit 1)
