.PHONY: test lint tidy

tidy:
	go mod tidy

test: tidy
	go test -v -race ./...

lint:
	go vet ./...
