.PHONY: all build test race fmt lint bench

all: build

build:
	go build ./...

test:
	go test ./...

race:
	go test -race ./...

fmt:
	go fmt ./...

bench:
	go test -bench=. -benchmem ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	golangci-lint run
