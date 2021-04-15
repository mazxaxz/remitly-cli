.PHONY: build
build:
	go mod download
	go build -o remitly ./cmd/main.go
