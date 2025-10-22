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
        docker images | grep -E "(product-service|postgres)"
        
        if [ "$PUSH" = true ]; then
            print_status "Pushing images to registry..."
            docker push product-service:$TAG
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
            print_status "  HTTP API: http://localhost:8080"
            print_status "  gRPC API: localhost:50050"
            print_status "  PostgreSQL: localhost:5432"
            print_status "  Health Check: http://localhost:8080/health"
            print_status "  Metrics: http://localhost:8080/metrics"
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
        
        # Remove product service images
        docker rmi $(docker images "product-service*" -q) 2>/dev/null || true
        
        # Clean up unused resources
        docker system prune -f
        
        print_success "Cleanup completed!"
        ;;
        
    push)
        print_status "Pushing images to registry..."
        
        # Tag and push product service
        docker tag product-service:$TAG product-service:latest
        docker push product-service:$TAG
        docker push product-service:latest
        
        print_success "Images pushed successfully!"
        ;;
esac
