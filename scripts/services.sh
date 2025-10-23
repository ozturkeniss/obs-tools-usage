#!/bin/bash

# Services management script
# This script helps manage microservices (start, stop, restart, status)

set -e

echo "🔧 Microservices Management Script..."

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
ACTION="status"
SERVICE="all"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        start)
            ACTION="start"
            shift
            ;;
        stop)
            ACTION="stop"
            shift
            ;;
        restart)
            ACTION="restart"
            shift
            ;;
        status)
            ACTION="status"
            shift
            ;;
        logs)
            ACTION="logs"
            shift
            ;;
        --service)
            SERVICE="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [ACTION] [OPTIONS]"
            echo ""
            echo "ACTIONS:"
            echo "  start     Start services"
            echo "  stop      Stop services"
            echo "  restart   Restart services"
            echo "  status    Show service status (default)"
            echo "  logs      Show service logs"
            echo ""
            echo "OPTIONS:"
            echo "  --service SERVICE    Target specific service (product|basket|payment|gateway|all)"
            echo "  -h, --help          Show this help message"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Function to check if service is running
check_service() {
    local service=$1
    local port=$2
    
    if nc -z localhost $port 2>/dev/null; then
        return 0  # Service is running
    else
        return 1  # Service is not running
    fi
}

# Function to show service status
show_status() {
    print_status "Service Status:"
    echo ""
    
    # Check Product Service
    if check_service "product-service" 8080; then
        print_success "Product Service HTTP (8080): ✅ Running"
    else
        print_error "Product Service HTTP (8080): ❌ Not running"
    fi
    
    if check_service "product-service" 50050; then
        print_success "Product Service gRPC (50050): ✅ Running"
    else
        print_error "Product Service gRPC (50050): ❌ Not running"
    fi
    
    # Check Basket Service
    if check_service "basket-service" 8081; then
        print_success "Basket Service HTTP (8081): ✅ Running"
    else
        print_error "Basket Service HTTP (8081): ❌ Not running"
    fi
    
    if check_service "basket-service" 50051; then
        print_success "Basket Service gRPC (50051): ✅ Running"
    else
        print_error "Basket Service gRPC (50051): ❌ Not running"
    fi
    
    # Check dependencies
    if check_service "postgres" 5432; then
        print_success "PostgreSQL (5432): ✅ Running"
    else
        print_error "PostgreSQL (5432): ❌ Not running"
    fi
    
    if check_service "redis" 6379; then
        print_success "Redis (6379): ✅ Running"
    else
        print_error "Redis (6379): ❌ Not running"
    fi
    
    echo ""
    print_status "Access URLs:"
    echo "  Product HTTP API: http://localhost:8080"
    echo "  Product gRPC API: localhost:50050"
    echo "  Basket HTTP API: http://localhost:8081"
    echo "  Basket gRPC API: localhost:50051"
    echo "  Product Health: http://localhost:8080/health"
    echo "  Basket Health: http://localhost:8081/health"
    echo "  Product Metrics: http://localhost:8080/metrics"
    echo "  Basket Metrics: http://localhost:8081/metrics"
}

# Function to start services
start_services() {
    print_status "Starting services..."
    
    # Start dependencies first
    print_status "Starting dependencies (PostgreSQL, Redis)..."
    docker-compose up -d postgres redis
    
    # Wait for dependencies
    print_status "Waiting for dependencies to be ready..."
    sleep 5
    
    # Check if binaries exist
    if [ ! -f "bin/product-service" ]; then
        print_status "Building product service..."
        go build -o bin/product-service cmd/product/main.go
    fi
    
    if [ ! -f "bin/basket-service" ]; then
        print_status "Building basket service..."
        go build -o bin/basket-service cmd/basket/main.go
    fi
    
    # Start services based on selection
    case $SERVICE in
        product|all)
            print_status "Starting product service..."
            nohup ./bin/product-service > logs/product-service.log 2>&1 &
            echo $! > logs/product-service.pid
            print_s_started "Product service started (PID: $(cat logs/product-service.pid))"
            ;;
        basket|all)
            print_status "Starting basket service..."
            nohup ./bin/basket-service > logs/basket-service.log 2>&1 &
            echo $! > logs/basket-service.pid
            print_success "Basket service started (PID: $(cat logs/basket-service.pid))"
            ;;
    esac
    
    # Create logs directory if it doesn't exist
    mkdir -p logs
    
    print_success "Services started!"
    show_status
}

# Function to stop services
stop_services() {
    print_status "Stopping services..."
    
    case $SERVICE in
        product|all)
            if [ -f "logs/product-service.pid" ]; then
                PID=$(cat logs/product-service.pid)
                if kill -0 $PID 2>/dev/null; then
                    kill $PID
                    print_success "Product service stopped (PID: $PID)"
                else
                    print_warning "Product service was not running"
                fi
                rm -f logs/product-service.pid
            fi
            ;;
        basket|all)
            if [ -f "logs/basket-service.pid" ]; then
                PID=$(cat logs/basket-service.pid)
                if kill -0 $PID 2>/dev/null; then
                    kill $PID
                    print_success "Basket service stopped (PID: $PID)"
                else
                    print_warning "Basket service was not running"
                fi
                rm -f logs/basket-service.pid
            fi
            ;;
    esac
    
    if [ "$SERVICE" = "all" ]; then
        print_status "Stopping dependencies..."
        docker-compose down
        print_success "Dependencies stopped!"
    fi
    
    print_success "Services stopped!"
}

# Function to restart services
restart_services() {
    print_status "Restarting services..."
    stop_services
    sleep 2
    start_services
}

# Function to show logs
show_logs() {
    case $SERVICE in
        product)
            if [ -f "logs/product-service.log" ]; then
                tail -f logs/product-service.log
            else
                print_error "Product service log file not found"
            fi
            ;;
        basket)
            if [ -f "logs/basket-service.log" ]; then
                tail -f logs/basket-service.log
            else
                print_error "Basket service log file not found"
            fi
            ;;
        all)
            print_status "Showing logs for all services..."
            if [ -f "logs/product-service.log" ]; then
                print_status "=== PRODUCT SERVICE LOGS ==="
                tail -20 logs/product-service.log
                echo ""
            fi
            if [ -f "logs/basket-service.log" ]; then
                print_status "=== BASKET SERVICE LOGS ==="
                tail -20 logs/basket-service.log
                echo ""
            fi
            ;;
    esac
}

# Main execution
case $ACTION in
    start)
        start_services
        ;;
    stop)
        stop_services
        ;;
    restart)
        restart_services
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs
        ;;
esac
