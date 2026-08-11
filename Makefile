VERSION=0.0.1

.PHONY: check lint

check:
	go test -v ./...
	go test -race

lint:
	golangci-lint run ./...