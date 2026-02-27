.PHONY: build test lint generate

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
