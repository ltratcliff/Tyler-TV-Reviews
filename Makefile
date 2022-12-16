.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golangci-lint run ./...
.PHONY:lint

vet: lint
	go vet ./...
.PHONY:vet

vuln: vet
	govulncheck  ./...
.PHONY:vuln

build: vuln
	go build main.go
.PHONY:build
