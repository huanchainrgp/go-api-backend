# Makefile for Go API Test1

.PHONY: help build run test clean docker-build docker-run docker-compose-up docker-compose-down dev-setup

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Go commands
build: ## Build the application
	go build -o bin/main main.go

run: ## Run the application
	go run main.go

test: ## Run tests
	go test ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

# Docker commands
docker-build: ## Build Docker image
	docker build -t go-api-test1 .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env go-api-test1

docker-compose-up: ## Start services with Docker Compose
	docker-compose up -d

docker-compose-down: ## Stop services with Docker Compose
	docker-compose down

# Development setup
dev-setup: ## Set up development environment
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp env.example .env; echo "Created .env file from template"; fi
	@echo "Installing dependencies..."
	go mod tidy
	@echo "Development environment ready!"
	@echo "Run 'make run' to start the development server"

# Database commands
db-migrate: ## Run database migrations
	go run main.go migrate

# Swagger documentation
swagger-generate: ## Generate Swagger documentation
	swag init -g main.go

# Production deployment
deploy: ## Deploy to production
	docker-compose -f docker-compose.yml up -d

# Development with hot reload (requires air)
dev: ## Run development server with hot reload
	air

# Install development tools
install-tools: ## Install development tools
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/cosmtrek/air@latest
