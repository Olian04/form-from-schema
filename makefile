.PHONY: all
all: format test lint

.PHONY: generate
generate:
	go tool templ generate

.PHONY: test
test: generate
	go test -cover -race ./...

.PHONY: lint
lint:
	go vet ./...
	golangci-lint run ./...

.PHONY: format
format:
	go fmt ./...
	go tool templ fmt ./lib/targets/html
	go mod tidy