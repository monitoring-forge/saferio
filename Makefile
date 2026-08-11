VERSION=0.0.1

check:
	go test -v ./...
	go test -race

lint:
	golangci-lint run ./...