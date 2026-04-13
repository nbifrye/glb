VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build clean test

build:
	go build $(LDFLAGS) -o bin/glb ./cmd/glb/

clean:
	rm -rf bin/

test:
	go test ./...
