#!/bin/bash

# Cleanup script
# This script cleans up build artifacts, logs, and temporary files

set -e

echo "🧹 Cleanup Script..."

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
ACTION="all"
FORCE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        build)
            ACTION="build"
            shift
            ;;
        logs)
            ACTION="logs"
            shift
            ;;
        docker)
            ACTION="docker"
            shift
            ;;
        all)
            ACTION="all"
            shift
            ;;
        --force)
            FORCE=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [ACTION] [OPTIONS]"
            echo ""
            echo "ACTIONS:"
            echo "  build     Clean build artifacts"
            echo "  logs      Clean log files"
            echo "  docker    Clean Docker containers and images"
            echo "  all       Clean everything (default)"
            echo ""
            echo "OPTIONS:"
            echo "  --force   Force cleanup without confirmation"
            echo "  -h, --help Show this help message"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Function to clean build artifacts
clean_build() {
    print_status "Cleaning build artifacts..."
    
    # Remove bin directory
    if [ -d "bin" ]; then
        rm -rf bin/
        print_success "Removed bin/ directory"
    else
        print_status "bin/ directory not found"
    fi
    
    # Remove coverage files
    if [ -f "coverage.out" ]; then
        rm -f coverage.out
        print_success "Removed coverage.out"
    fi
    
    if [ -f "coverage.html" ]; then
        rm -f coverage.html
        print_success "Removed coverage.html"
    fi
    
    # Remove test cache
    go clean -testcache
    
    # Remove module cache (optional)
    if [ "$FORCE" = true ]; then
        print_status "Cleaning module cache..."
        go clean -modcache
    fi
}

# Function to clean log files
clean_logs() {
    print_status "Cleaning log files..."
    
    # Remove logs directory
    if [ -d "logs" ]; then
        rm -rf logs/
        print_success "Removed logs/ directory"
    else
        print_status "logs/ directory not found"
    fi
    
    # Remove individual log files
    find . -name "*.log" -type f -delete 2>/dev/null || true
    find . -name "*.log.*" -type f -delete 2>/dev/null || true
    
    print_success "Log files cleaned"
}

# Function to clean Docker resources
clean_docker() {
    print_status "Cleaning Docker resources..."
    
    # Stop and remove containers
    print_status "Stopping containers..."
    docker-compose down 2>/dev/null || true
    
    # Remove product service containers
    print_status "Removing product service containers..."
    docker ps -a --filter "name=product-service" --format "{{.ID}}" | xargs -r docker rm -f 2>/dev/null || true
    
    # Remove product service images
    print_status "Removing product service images..."
    docker images "product-service*" -q | xargs -r docker rmi -f 2>/dev/null || true
    
    # Remove unused images
    if [ "$FORCE" = true ]; then
        print_status "Removing unused images..."
        docker image prune -f
    fi
    
    # Remove unused volumes
    if [ "$FORCE" = true ]; then
        print_status "Removing unused volumes..."
        docker volume prune -f
    fi
    
    # Remove unused networks
    if [ "$FORCE" = true ]; then
        print_status "Removing unused networks..."
        docker network prune -f
    fi
    
    print_success "Docker resources cleaned"
}

# Function to clean temporary files
clean_temp() {
    print_status "Cleaning temporary files..."
    
    # Remove .DS_Store files (macOS)
    find . -name ".DS_Store" -type f -delete 2>/dev/null || true
    
    # Remove Thumbs.db files (Windows)
    find . -name "Thumbs.db" -type f -delete 2>/dev/null || true
    
    # Remove temporary Go files
    find . -name "*.tmp" -type f -delete 2>/dev/null || true
    
    # Remove backup files
    find . -name "*.bak" -type f -delete 2>/dev/null || true
    find . -name "*.backup" -type f -delete 2>/dev/null || true
    
    # Remove editor temporary files
    find . -name "*~" -type f -delete 2>/dev/null || true
    find . -name "*.swp" -type f -delete 2>/dev/null || true
    find . -name "*.swo" -type f -delete 2>/dev/null || true
    
    print_success "Temporary files cleaned"
}

# Function to clean backups
clean_backups() {
    print_status "Cleaning backup files..."
    
    if [ -d "backups" ]; then
        if [ "$FORCE" = true ]; then
            rm -rf backups/
            print_success "Removed backups/ directory"
        else
            print_warning "Backups directory found. Use --force to remove it."
        fi
    else
        print_status "backups/ directory not found"
    fi
}

# Function to show disk usage
show_disk_usage() {
    print_status "Current disk usage:"
    du -sh . 2>/dev/null || true
    echo ""
    
    if [ -d "bin" ]; then
        print_status "Build artifacts size:"
        du -sh bin/ 2>/dev/null || true
    fi
    
    if [ -d "logs" ]; then
        print_status "Log files size:"
        du -sh logs/ 2>/dev/null || true
    fi
    
    if [ -d "backups" ]; then
        print_status "Backup files size:"
        du -sh backups/ 2>/dev/null || true
    fi
}

# Main execution
if [ "$FORCE" != true ] && [ "$ACTION" = "all" ]; then
    print_warning "This will clean all build artifacts, logs, and temporary files."
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_status "Operation cancelled."
        exit 0
    fi
fi

case $ACTION in
    build)
        clean_build
        ;;
    logs)
        clean_logs
        ;;
    docker)
        clean_docker
        ;;
    all)
        clean_build
        clean_logs
        clean_temp
        clean_backups
        print_success "All cleanup completed!"
        ;;
esac

# Show disk usage after cleanup
show_disk_usage
