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

# Development
.PHONY: dev
dev:
	@echo "Starting development server..."
	@chmod +x scripts/dev.sh
	@./scripts/dev.sh

# Build the application
.PHONY: build
build:
	@echo "Building application..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh

# Build notification service
.PHONY: build-notification
build-notification:
	@echo "Building notification service..."
	go build -o bin/notification-service cmd/notification/main.go

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@chmod +x scripts/test.sh
	@./scripts/test.sh

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Lint and format code
.PHONY: lint
lint:
	@echo "Running linter..."
	@chmod +x scripts/lint.sh
	@./scripts/lint.sh

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@chmod +x scripts/lint.sh
	@./scripts/lint.sh format --fix

# Database operations
.PHONY: db-migrate
db-migrate:
	@echo "Running database migrations..."
	@chmod +x scripts/db.sh
	@./scripts/db.sh migrate

.PHONY: db-seed
db-seed:
	@echo "Seeding database..."
	@chmod +x scripts/db.sh
	@./scripts/db.sh seed

.PHONY: db-backup
db-backup:
	@echo "Creating database backup..."
	@chmod +x scripts/db.sh
	@./scripts/db.sh backup

.PHONY: db-status
db-status:
	@echo "Checking database status..."
	@chmod +x scripts/db.sh
	@./scripts/db.sh status

# Docker operations
.PHONY: docker-build
docker-build:
	@echo "Building Docker images..."
	@chmod +x scripts/docker.sh
	@./scripts/docker.sh build

.PHONY: docker-run
docker-run:
	@echo "Running Docker containers..."
	@chmod +x scripts/docker.sh
	@./scripts/docker.sh run

.PHONY: docker-stop
docker-stop:
	@echo "Stopping Docker containers..."
	@chmod +x scripts/docker.sh
	@./scripts/docker.sh stop

.PHONY: docker-clean
docker-clean:
	@echo "Cleaning Docker resources..."
	@chmod +x scripts/docker.sh
	@./scripts/docker.sh clean

# Cleanup
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@chmod +x scripts/clean.sh
	@./scripts/clean.sh build

.PHONY: clean-all
clean-all:
	@echo "Cleaning everything..."
	@chmod +x scripts/clean.sh
	@./scripts/clean.sh all --force

# Generate all (proto + build)
.PHONY: generate
generate: proto build

# Install development dependencies
.PHONY: install-deps
install-deps: install-proto-deps
	@echo "Installing development dependencies..."
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "Development dependencies installed!"

# Setup project
.PHONY: setup
setup: install-deps mod-tidy proto
	@echo "Project setup completed!"

# Run the application
.PHONY: run
run: build
	@echo "Running microservices..."
	@chmod +x scripts/services.sh
	@./scripts/services.sh start --service all

# Services management
.PHONY: services-start
services-start:
	@chmod +x scripts/services.sh
	@./scripts/services.sh start

.PHONY: services-stop
services-stop:
	@chmod +x scripts/services.sh
	@./scripts/services.sh stop

.PHONY: services-restart
services-restart:
	@chmod +x scripts/services.sh
	@./scripts/services.sh restart

.PHONY: services-status
services-status:
	@chmod +x scripts/services.sh
	@./scripts/services.sh status

.PHONY: services-logs
services-logs:
	@chmod +x scripts/services.sh
	@./scripts/services.sh logs

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  setup          - Setup project (install deps, generate proto)"
	@echo "  dev            - Start development server"
	@echo "  build          - Build the application"
	@echo "  run            - Run microservices"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  lint           - Run linter and format checks"
	@echo "  fmt            - Format code"
	@echo "  proto          - Generate protobuf files"
	@echo "  mod-tidy       - Tidy go modules"
	@echo "  install-deps   - Install development dependencies"
	@echo ""
	@echo "Services management:"
	@echo "  services-start    - Start services"
	@echo "  services-stop     - Stop services"
	@echo "  services-restart  - Restart services"
	@echo "  services-status   - Show service status"
	@echo "  services-logs     - Show service logs"
	@echo ""
	@echo "Database targets:"
	@echo "  db-migrate     - Run database migrations"
	@echo "  db-seed        - Seed database with initial data"
	@echo "  db-backup      - Create database backup"
	@echo "  db-status      - Show database status"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-build   - Build Docker images"
	@echo "  docker-run     - Run Docker containers"
	@echo "  docker-stop    - Stop Docker containers"
	@echo "  docker-clean   - Clean Docker resources"
	@echo ""
	@echo "Cleanup targets:"
	@echo "  clean          - Clean build artifacts"
	@echo "  clean-all      - Clean everything"
	@echo ""
	@echo "  help           - Show this help message"