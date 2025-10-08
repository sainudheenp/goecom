.PHONY: help run build test lint migrate-up migrate-down seed docker-up docker-down clean

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	go run ./cmd/server

build: ## Build the application
	go build -o bin/server ./cmd/server

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out

lint: ## Run linter
	golangci-lint run

migrate-up: ## Run database migrations up
	@bash scripts/migrate.sh up

migrate-down: ## Run database migrations down
	@bash scripts/migrate.sh down

seed: ## Seed database with sample data
	@bash scripts/seed.sh

docker-up: ## Start docker-compose services
	docker-compose up --build -d

docker-down: ## Stop docker-compose services
	docker-compose down

docker-logs: ## View docker-compose logs
	docker-compose logs -f

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out

deps: ## Download dependencies
	go mod download
	go mod tidy

fmt: ## Format code
	go fmt ./...
	goimports -w .

integration-test: ## Run integration tests
	go test -v -tags=integration ./test/...
