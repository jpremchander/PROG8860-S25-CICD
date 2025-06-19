.PHONY: help build run test clean docker-build docker-run lint fmt dev deps

# Variables
APP_NAME=yelpcamp-app
DOCKER_IMAGE=yelpcamp-go
GO_VERSION=1.21

# Default target
help: ## Show this help message
	@echo 'YelpCamp Go - Development Commands'
	@echo ''
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
dev: ## Run with hot reload (requires air)
	air

run: ## Run the application
	go run .

build: ## Build the application
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/$(APP_NAME) .

# Testing
test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests and show coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Code quality
lint: ## Run linter
	golangci-lint run --timeout=5m

fmt: ## Format code
	go fmt ./...
	goimports -w .

vet: ## Run go vet
	go vet ./...

# Dependencies
deps: ## Install dependencies
	go mod download
	go mod tidy

deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy

# Docker Development
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):latest .

docker-build-prod: ## Build production Docker image
	docker build -f Dockerfile.multi-stage -t $(DOCKER_IMAGE):prod .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):latest

# Docker Compose
compose-dev: ## Run development environment
	docker-compose -f docker-compose.dev.yml up --build

compose-prod: ## Run production environment
	docker-compose -f docker-compose.prod.yml up --build -d

compose-down: ## Stop all containers
	docker-compose -f docker-compose.dev.yml down
	docker-compose -f docker-compose.prod.yml down

# Database
mongo-shell: ## Connect to MongoDB shell
	docker exec -it $$(docker-compose ps -q mongo) mongosh -u yelpcamp_dev -p dev_password_123 yelpcamp_dev

# Jenkins
jenkins-build: ## Simulate Jenkins build locally
	@echo "🔄 Simulating Jenkins build..."
	make lint
	make test
	make build
	make docker-build
	@echo "✅ Local Jenkins simulation completed"

# CI/CD
ci-test: lint test ## Run CI tests locally

# Cleanup
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	docker system prune -f

# Security
security-scan: ## Run security scan
	gosec ./...

# Release
release: clean lint test build ## Prepare release build
	@echo "Release build completed successfully"

# Monitoring
logs: ## Show application logs
	docker-compose logs -f app

logs-mongo: ## Show MongoDB logs
	docker-compose logs -f mongo

status: ## Show container status
	docker-compose ps
