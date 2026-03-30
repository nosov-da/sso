.PHONY: help build run test coverage clean docker-build docker-run deps fmt lint

# Variables
BINARY_NAME=sso
GO_VERSION=1.24.4
DOCKER_IMAGE=sso:latest

help:
	@echo "Available targets:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run tests"
	@echo "  make coverage    - Run tests with coverage"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make deps        - Download dependencies"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Run linter"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run  - Run Docker container"

build:
	@echo "Building $(BINARY_NAME)..."
	CGO_ENABLED=0 go build -o bin/$(BINARY_NAME) ./cmd/sso

run:
	@echo "Running $(BINARY_NAME)..."
	go run ./cmd/sso

test:
	@echo "Running tests..."
	go test -v ./...

coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	@echo "Cleaning..."
	rm -f bin/$(BINARY_NAME)
	rm -f coverage.out coverage.html
	go clean -testcache

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Running linter..."
	golangci-lint run ./...

docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 \
		-e SSO_JWT_SECRET="your-secret-key" \
		-e ENVIRONMENT="development" \
		$(DOCKER_IMAGE)
