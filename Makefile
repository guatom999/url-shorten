.PHONY: help dev build test clean run fmt lint tidy

# Default target
help: ## Show available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

dev: ## Run in development mode
	@echo "Starting development server..."
	@cp .env.example .env 2>/dev/null || true
	@go run main.go

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/app main.go

run: ## Run the built application
	@echo "Running application..."
	@./bin/app

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/ coverage.out coverage.html
	@go clean -cache

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

setup: ## Setup development environment
	@echo "Setting up development environment..."
	@cp .env.example .env 2>/dev/null || true
	@go mod tidy
	@echo "Setup complete!"