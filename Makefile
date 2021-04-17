.PHONY: build
build:
	go mod download
	go build -o remitly ./cmd/main.go

.PHONY: test
test:
	go test -count=1 ./...

.PHONY: mocks
mocks:
	mockgen -destination=./pkg/remitly/mocks/client.go github.com/mazxaxz/remitly-cli/pkg/remitly Clienter
