# Makefile for obs-tools-usage

# Variables
PROTO_DIR = api/proto
PROTO_FILES = $(shell find $(PROTO_DIR) -name "*.proto")
GO_OUT_DIR = .

# Go modules
.PHONY: mod-tidy
mod-tidy:
	go mod tidy

# Generate protobuf files
.PHONY: proto
proto: $(PROTO_FILES)
	@echo "Generating protobuf files..."
	protoc --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "Protobuf files generated successfully!"

# Install protobuf dependencies
.PHONY: install-proto-deps
install-proto-deps:
	@echo "Installing protobuf dependencies..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Protobuf dependencies installed!"

# Build the application
.PHONY: build
build:
	@echo "Building application..."
	go build -o bin/product-service cmd/product/main.go
	@echo "Application built successfully!"

# Run the application
.PHONY: run
run: build
	@echo "Running product service..."
	./bin/product-service

# Run with hot reload (requires air)
.PHONY: dev
dev:
	@echo "Running in development mode..."
	air

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	@echo "Clean completed!"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Lint the code
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# Format the code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Generate all (proto + build)
.PHONY: generate
generate: proto build

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  mod-tidy          - Tidy go modules"
	@echo "  proto             - Generate protobuf files"
	@echo "  install-proto-deps- Install protobuf dependencies"
	@echo "  build             - Build the application"
	@echo "  run               - Run the application"
	@echo "  dev               - Run in development mode (requires air)"
	@echo "  clean             - Clean build artifacts"
	@echo "  test              - Run tests"
	@echo "  test-coverage     - Run tests with coverage"
	@echo "  lint              - Run linter"
	@echo "  fmt               - Format code"
	@echo "  generate          - Generate proto files and build"
	@echo "  help              - Show this help message"
