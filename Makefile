.PHONY: build test lint ci

build:
	go build -v ./...

test:
	go test -v ./...

lint:
	go vet ./...

ci: build test lint
