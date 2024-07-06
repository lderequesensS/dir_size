.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	staticcheck ./...
.PHONY:fmt

vet: fmt lint
	go vet ./...
.PHONY:vet

build:
	go build -o dir_size main.go
.PHONY:build
