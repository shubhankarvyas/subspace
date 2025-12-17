# Makefile for Subspace Automation PoC

.PHONY: help build run demo stats clean test deps fmt lint

# Default target
help:
	@echo "Subspace Automation PoC - Makefile Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make build       Build the application binary"
	@echo "  make run         Run in normal automation mode"
	@echo "  make demo        Run in demonstration mode"
	@echo "  make stats       Show current statistics"
	@echo "  make clean       Remove build artifacts"
	@echo "  make test        Run tests"
	@echo "  make deps        Download dependencies"
	@echo "  make fmt         Format code"
	@echo "  make lint        Run linter"

# Build the binary
build:
	@echo "Building Subspace..."
	@go build -o subspace cmd/app/main.go
	@echo "Build complete: ./subspace"

# Run normal mode
run: build
	@./subspace

# Run demo mode
demo: build
	@./subspace -demo

# Show statistics
stats: build
	@./subspace -stats

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f subspace
	@rm -rf data/*.json
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies ready"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# Setup (first time)
setup: deps
	@echo "Setting up project..."
	@cp -n .env.example .env || true
	@mkdir -p data
	@echo "Setup complete. Edit .env if needed."

# Run with custom config
run-custom:
	@./subspace -config=custom-config.yaml
