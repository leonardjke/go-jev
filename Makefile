.PHONY: all build test test-e2e cover fmt vet lint tidy clean

all: fmt vet test

build:
	go build ./...

test:
	go test -race ./...

# Talks to the live API. Skips unless API_KEY is set.
test-e2e:
	go test -tags e2e -v ./...

cover:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

fmt:
	gofmt -l -w .

vet:
	go vet ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

clean:
	rm -f coverage.out
