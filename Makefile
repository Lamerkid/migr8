BIN := "./bin/"
LDFLAGS := "-X main.version=develop"

build:
	go build -ldflags=$(LDFLAGS) -v -o $(BIN) ./cmd/migr8

test:
	go test -race -count 100 ./...

lint:
	golangci-lint run ./...

integration-tests:
	@echo "Starting postgres..."
	@docker-compose -f test/integration/docker-compose.yaml up -d --build postgres
	@echo "Running integration tests..."
	@go test -v -tags=integration ./test/integration/...
	@docker-compose -f test/integration/docker-compose.yaml down -v --remove-orphans

.PHONY: build test lint integration-tests
