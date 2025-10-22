#!/bin/bash

# Test script
# This script runs all tests and generates coverage reports

set -e

echo "🧪 Running Product Service Tests..."

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

# Navigate to project root
cd "$(dirname "$0")/.."

# Clean previous test results
print_status "Cleaning previous test results..."
rm -rf coverage.out coverage.html

# Run unit tests
print_status "Running unit tests..."
go test -v ./...

if [ $? -eq 0 ]; then
    print_success "Unit tests passed!"
else
    print_error "Unit tests failed!"
    exit 1
fi

# Run tests with coverage
print_status "Running tests with coverage..."
go test -coverprofile=coverage.out -covermode=atomic ./...

if [ $? -eq 0 ]; then
    print_success "Coverage tests completed!"
else
    print_error "Coverage tests failed!"
    exit 1
fi

# Generate coverage report
print_status "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

if [ $? -eq 0 ]; then
    print_success "Coverage report generated: coverage.html"
else
    print_error "Failed to generate coverage report!"
    exit 1
fi

# Show coverage summary
print_status "Coverage summary:"
go tool cover -func=coverage.out | tail -1

# Run benchmarks
print_status "Running benchmarks..."
go test -bench=. -benchmem ./...

# Run race detection
print_status "Running race detection..."
go test -race ./...

if [ $? -eq 0 ]; then
    print_success "Race detection passed!"
else
    print_warning "Race detection found issues!"
fi

print_success "All tests completed successfully!"
