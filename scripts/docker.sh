#!/bin/bash

# Docker build and deploy script
# This script builds Docker images and manages containers

set -e

echo "🐳 Docker Build and Deploy Script..."

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

# Default values
ACTION="build"
TAG="latest"
PUSH=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        build)
            ACTION="build"
            shift
            ;;
        run)
            ACTION="run"
            shift
            ;;
        stop)
            ACTION="stop"
            shift
            ;;
        clean)
            ACTION="clean"
            shift
            ;;
        push)
            ACTION="push"
            shift
            ;;
        --tag)
            TAG="$2"
            shift 2
            ;;
        --push)
            PUSH=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [ACTION] [OPTIONS]"
            echo ""
            echo "ACTIONS:"
            echo "  build    Build Docker images (default)"
            echo "  run      Run containers with docker-compose"
            echo "  stop     Stop all containers"
            echo "  clean    Clean up containers and images"
            echo "  push     Push images to registry"
            echo ""
            echo "OPTIONS:"
            echo "  --tag TAG    Set image tag (default: latest)"
            echo "  --push       Push images after building"
            echo "  -h, --help   Show this help message"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

case $ACTION in
    build)
        print_status "Building Docker images with tag: $TAG"
        
        # Build product service image
        print_status "Building product service image..."
        docker build -f dockerfiles/product.dockerfile -t product-service:$TAG .
        
        if [ $? -eq 0 ]; then
            print_success "Product service image built successfully!"
        else
            print_error "Failed to build product service image!"
            exit 1
        fi
        
        # Build basket service image
        print_status "Building basket service image..."
        docker build -f dockerfiles/basket.dockerfile -t basket-service:$TAG .
        
        if [ $? -eq 0 ]; then
            print_success "Basket service image built successfully!"
        else
            print_error "Failed to build basket service image!"
            exit 1
        fi
        
        # Build payment service image
        print_status "Building payment service image..."
        docker build -f dockerfiles/payment.dockerfile -t payment-service:$TAG .
        
        if [ $? -eq 0 ]; then
            print_success "Payment service image built successfully!"
        else
            print_error "Failed to build payment service image!"
            exit 1
        fi
        
        # Build gateway service image
        print_status "Building gateway service image..."
        docker build -f dockerfiles/gateway.dockerfile -t gateway:$TAG .
        
        if [ $? -eq 0 ]; then
            print_success "Gateway service image built successfully!"
        else
            print_error "Failed to build gateway service image!"
            exit 1
        fi
        
        # Build all services with docker-compose
        print_status "Building all services with docker-compose..."
        docker-compose build
        
        if [ $? -eq 0 ]; then
            print_success "All services built successfully!"
        else
            print_error "Failed to build services with docker-compose!"
            exit 1
        fi
        
        # Show built images
        print_status "Built images:"
        docker images | grep -E "(product-service|basket-service|payment-service|gateway|postgres|redis|mariadb|kafka|zookeeper)"
        
        if [ "$PUSH" = true ]; then
            print_status "Pushing images to registry..."
            docker push product-service:$TAG
            docker push basket-service:$TAG
            docker push payment-service:$TAG
            docker push gateway:$TAG
        fi
        ;;
        
    run)
        print_status "Starting services with docker-compose..."
        docker-compose up -d
        
        if [ $? -eq 0 ]; then
            print_success "Services started successfully!"
            print_status "Services running:"
            docker-compose ps
            echo ""
            print_status "Access points:"
            print_status "  🌐 FiberV2 Gateway: http://localhost:8083"
            print_status "  📊 Gateway Admin: http://localhost:8083/admin"
            print_status "  ❤️  Gateway Health: http://localhost:8083/health"
            print_status "  📈 Gateway Metrics: http://localhost:8083/metrics"
            print_status ""
            print_status "  🛍️  Product Service HTTP API: http://localhost:8080"
            print_status "  🛒 Basket Service HTTP API: http://localhost:8081"
            print_status "  💳 Payment Service HTTP API: http://localhost:8082"
            print_status ""
            print_status "  🛍️  Product Service gRPC API: localhost:50050"
            print_status "  🛒 Basket Service gRPC API: localhost:50051"
            print_status "  💳 Payment Service gRPC API: localhost:50052"
            print_status ""
            print_status "  🗄️  PostgreSQL: localhost:5432"
            print_status "  🔴 Redis: localhost:6379"
            print_status "  🗄️  MariaDB: localhost:3306"
            print_status "  📨 Kafka: localhost:9092"
            print_status ""
            print_status "  ❤️  Health Checks:"
            print_status "    Product Health: http://localhost:8080/health"
            print_status "    Basket Health: http://localhost:8081/health"
            print_status "    Payment Health: http://localhost:8082/health"
            print_status "    Gateway Health: http://localhost:8083/health"
            print_status ""
            print_status "  📈 Metrics:"
            print_status "    Product Metrics: http://localhost:8080/metrics"
            print_status "    Basket Metrics: http://localhost:8081/metrics"
            print_status "    Payment Metrics: http://localhost:8082/metrics"
            print_status "    Gateway Metrics: http://localhost:8083/metrics"
        else
            print_error "Failed to start services!"
            exit 1
        fi
        ;;
        
    stop)
        print_status "Stopping all services..."
        docker-compose down
        
        if [ $? -eq 0 ]; then
            print_success "Services stopped successfully!"
        else
            print_error "Failed to stop services!"
            exit 1
        fi
        ;;
        
    clean)
        print_status "Cleaning up containers and images..."
        
        # Stop and remove containers
        docker-compose down --volumes --rmi all
        
        # Remove service images
        docker rmi $(docker images "product-service*" -q) 2>/dev/null || true
        docker rmi $(docker images "basket-service*" -q) 2>/dev/null || true
        docker rmi $(docker images "payment-service*" -q) 2>/dev/null || true
        docker rmi $(docker images "gateway*" -q) 2>/dev/null || true
        
        # Clean up unused resources
        docker system prune -f
        
        print_success "Cleanup completed!"
        ;;
        
    push)
        print_status "Pushing images to registry..."
        
        # Tag and push services
        docker tag product-service:$TAG product-service:latest
        docker tag basket-service:$TAG basket-service:latest
        docker tag payment-service:$TAG payment-service:latest
        docker tag gateway:$TAG gateway:latest
        
        docker push product-service:$TAG
        docker push product-service:latest
        docker push basket-service:$TAG
        docker push basket-service:latest
        docker push payment-service:$TAG
        docker push payment-service:latest
        docker push gateway:$TAG
        docker push gateway:latest
        
        print_success "Images pushed successfully!"
        ;;
esac
