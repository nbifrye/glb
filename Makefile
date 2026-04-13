VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build clean test cover vet

build:
	go build $(LDFLAGS) -o bin/glb ./cmd/glb/

clean:
	rm -rf bin/ coverage.out

test:
	go test -race -count=1 ./...

cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

vet:
	go vet ./...
