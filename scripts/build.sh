#!/bin/bash

# Production build script
# This script builds the application for production deployment

set -e

echo "🏗️  Building Microservices for Production..."

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

# Clean previous builds
print_status "Cleaning previous builds..."
rm -rf bin/
mkdir -p bin/
mkdir -p fiberv2-gateway/bin/

# Download dependencies
print_status "Downloading dependencies..."
go mod download

# Run tests
print_status "Running tests..."
go test ./...

if [ $? -eq 0 ]; then
    print_success "All tests passed!"
else
    print_error "Tests failed!"
    exit 1
fi

# Format code
print_status "Formatting code..."
go fmt ./...

# Run linter (if available)
if command -v golangci-lint &> /dev/null; then
    print_status "Running linter..."
    golangci-lint run
else
    print_warning "golangci-lint not found, skipping linting"
fi

# Build for different platforms
print_status "Building Product Service for Linux AMD64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/product-service-linux-amd64 cmd/product/main.go

print_status "Building Product Service for Linux ARM64..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/product-service-linux-arm64 cmd/product/main.go

print_status "Building Product Service for Windows AMD64..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/product-service-windows-amd64.exe cmd/product/main.go

print_status "Building Product Service for macOS AMD64..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/product-service-darwin-amd64 cmd/product/main.go

print_status "Building Product Service for macOS ARM64..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/product-service-darwin-arm64 cmd/product/main.go

print_status "Building Basket Service for Linux AMD64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/basket-service-linux-amd64 cmd/basket/main.go

print_status "Building Basket Service for Linux ARM64..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/basket-service-linux-arm64 cmd/basket/main.go

print_status "Building Basket Service for Windows AMD64..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/basket-service-windows-amd64.exe cmd/basket/main.go

print_status "Building Basket Service for macOS AMD64..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/basket-service-darwin-amd64 cmd/basket/main.go

print_status "Building Basket Service for macOS ARM64..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/basket-service-darwin-arm64 cmd/basket/main.go

print_status "Building Payment Service for Linux AMD64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/payment-service-linux-amd64 cmd/payment/main.go

print_status "Building Payment Service for Linux ARM64..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/payment-service-linux-arm64 cmd/payment/main.go

print_status "Building Payment Service for Windows AMD64..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/payment-service-windows-amd64.exe cmd/payment/main.go

print_status "Building Payment Service for macOS AMD64..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/payment-service-darwin-amd64 cmd/payment/main.go

print_status "Building Payment Service for macOS ARM64..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/payment-service-darwin-arm64 cmd/payment/main.go

print_status "Building Gateway Service for Linux AMD64..."
cd fiberv2-gateway
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/gateway-linux-amd64 cmd/main.go

print_status "Building Gateway Service for Linux ARM64..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/gateway-linux-arm64 cmd/main.go

print_status "Building Gateway Service for Windows AMD64..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/gateway-windows-amd64.exe cmd/main.go

print_status "Building Gateway Service for macOS AMD64..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/gateway-darwin-amd64 cmd/main.go

print_status "Building Gateway Service for macOS ARM64..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/gateway-darwin-arm64 cmd/main.go
cd ..

# Create symlinks for the current platform
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    if [[ $(uname -m) == "x86_64" ]]; then
        ln -sf product-service-linux-amd64 bin/product-service
        ln -sf basket-service-linux-amd64 bin/basket-service
        ln -sf payment-service-linux-amd64 bin/payment-service
        ln -sf gateway-linux-amd64 fiberv2-gateway/bin/gateway
    elif [[ $(uname -m) == "aarch64" ]]; then
        ln -sf product-service-linux-arm64 bin/product-service
        ln -sf basket-service-linux-arm64 bin/basket-service
        ln -sf payment-service-linux-arm64 bin/payment-service
        ln -sf gateway-linux-arm64 fiberv2-gateway/bin/gateway
    fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
    if [[ $(uname -m) == "x86_64" ]]; then
        ln -sf product-service-darwin-amd64 bin/product-service
        ln -sf basket-service-darwin-amd64 bin/basket-service
        ln -sf payment-service-darwin-amd64 bin/payment-service
        ln -sf gateway-darwin-amd64 fiberv2-gateway/bin/gateway
    elif [[ $(uname -m) == "arm64" ]]; then
        ln -sf product-service-darwin-arm64 bin/product-service
        ln -sf basket-service-darwin-arm64 bin/basket-service
        ln -sf payment-service-darwin-arm64 bin/payment-service
        ln -sf gateway-darwin-arm64 fiberv2-gateway/bin/gateway
    fi
fi

# Show build results
print_success "Build completed successfully!"
echo ""
print_status "Build artifacts:"
ls -la bin/

# Show file sizes
print_status "Binary sizes:"
for file in bin/*-service-*; do
    if [ -f "$file" ]; then
        size=$(du -h "$file" | cut -f1)
        echo "  $(basename "$file"): $size"
    fi
done
