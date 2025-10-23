#!/bin/bash

# Development server startup script
# This script starts the product service with all dependencies

set -e

echo "🚀 Starting Microservices in Development Mode..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go first."
    exit 1
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Navigate to project root
cd "$(dirname "$0")/.."

print_status "Building the applications..."
go build -o bin/product-service cmd/product/main.go
if [ $? -eq 0 ]; then
    print_success "Product service built successfully!"
else
    print_error "Failed to build product service"
    exit 1
fi

go build -o bin/basket-service cmd/basket/main.go
if [ $? -eq 0 ]; then
    print_success "Basket service built successfully!"
else
    print_error "Failed to build basket service"
    exit 1
fi

print_status "Starting dependencies with Docker Compose..."
docker-compose up -d postgres redis

# Wait for PostgreSQL to be ready
print_status "Waiting for PostgreSQL to be ready..."
sleep 5

# Check if PostgreSQL is ready
until docker-compose exec postgres pg_isready -U postgres; do
    print_status "Waiting for PostgreSQL..."
    sleep 2
done

print_success "PostgreSQL is ready!"

# Wait for Redis to be ready
print_status "Waiting for Redis to be ready..."
sleep 3

# Check if Redis is ready
until docker-compose exec redis redis-cli ping; do
    print_status "Waiting for Redis..."
    sleep 2
done

print_success "Redis is ready!"

# Set environment variables
export ENVIRONMENT=development
export PORT=8080
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=product_service
export DB_SSL_MODE=disable
export LOG_LEVEL=debug

print_status "Starting Services..."
print_status ""
print_status "=== PRODUCT SERVICE ==="
print_status "HTTP API: http://localhost:8080"
print_status "gRPC API: localhost:50050"
print_status "Health Check: http://localhost:8080/health"
print_status "Metrics: http://localhost:8080/metrics"
print_status ""
print_status "=== BASKET SERVICE ==="
print_status "HTTP API: http://localhost:8081"
print_status "gRPC API: localhost:50051"
print_status "Health Check: http://localhost:8081/health"
print_status "Metrics: http://localhost:8081/metrics"
print_status ""
print_status "Press Ctrl+C to stop all services"

# Start services in background
./bin/product-service &
PRODUCT_PID=$!

./bin/basket-service &
BASKET_PID=$!

# Function to cleanup background processes
cleanup() {
    print_status "Stopping services..."
    kill $PRODUCT_PID 2>/dev/null || true
    kill $BASKET_PID 2>/dev/null || true
    print_success "Services stopped!"
    exit 0
}

# Set trap for cleanup
trap cleanup SIGINT SIGTERM

# Wait for processes
wait
