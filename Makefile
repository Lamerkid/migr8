BIN := "./bin/"
LDFLAGS := "-X main.version=develop"

build:
	go build -ldflags=$(LDFLAGS) -v -o $(BIN) ./cmd/migr8

.PHONY: build
