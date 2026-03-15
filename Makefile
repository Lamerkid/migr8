include .env
export

BIN := "./bin/"
LDFLAGS := "-X main.version=develop"

lint:
	golangci-lint run ./...

test:
	go test -race -count 100 ./...

integration-tests:
	@echo "Starting postgres..."
	@docker-compose -f test/integration/docker-compose.yaml up -d --build postgres
	@echo "Running integration tests..."
	@echo "DSN: $$M8_DSN"
	@echo "DIR: $$M8_DIR"
	@go test -v -tags=integration ./test/integration/...
	@docker-compose -f test/integration/docker-compose.yaml down -v --remove-orphans

build:
	go build -ldflags=$(LDFLAGS) -v -o $(BIN) ./cmd/migr8

up:
	docker-compose -f deployments/docker-compose.yaml up --build

.PHONY: lint test integration-tests build up
