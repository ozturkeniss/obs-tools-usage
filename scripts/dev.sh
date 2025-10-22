#!/bin/bash

# Development server startup script
# This script starts the product service with all dependencies

set -e

echo "🚀 Starting Product Service in Development Mode..."

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

print_status "Building the application..."
go build -o bin/product-service cmd/product/main.go
if [ $? -eq 0 ]; then
    print_success "Application built successfully!"
else
    print_error "Failed to build application"
    exit 1
fi

print_status "Starting PostgreSQL with Docker Compose..."
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
print_status "Waiting for PostgreSQL to be ready..."
sleep 5

# Check if PostgreSQL is ready
until docker-compose exec postgres pg_isready -U postgres; do
    print_status "Waiting for PostgreSQL..."
    sleep 2
done

print_success "PostgreSQL is ready!"

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

print_status "Starting Product Service..."
print_status "HTTP API: http://localhost:8080"
print_status "gRPC API: localhost:50050"
print_status "Health Check: http://localhost:8080/health"
print_status "Metrics: http://localhost:8080/metrics"
print_status "Press Ctrl+C to stop the service"

# Start the service
./bin/product-service
