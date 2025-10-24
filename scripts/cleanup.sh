#!/bin/bash

# OBS Tools Usage - Cleanup Script
# This script cleans up resources and temporary files

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
CLEANUP_TYPE="all"
NAMESPACE="obs-tools-usage"
ENVIRONMENT="dev"
FORCE=false
KEEP_LOGS=false

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Show usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -t, --type TYPE         Cleanup type (all, k8s, docker, files, logs, temp)"
    echo "  -n, --namespace NS       Kubernetes namespace (default: obs-tools-usage)"
    echo "  -e, --env ENV            Environment (dev, staging, prod)"
    echo "  -f, --force              Force cleanup without confirmation"
    echo "  -k, --keep-logs          Keep log files"
    echo "  -h, --help               Show this help message"
    echo ""
    echo "Cleanup Types:"
    echo "  all                     Clean up everything"
    echo "  k8s                     Clean up Kubernetes resources"
    echo "  docker                  Clean up Docker resources"
    echo "  files                   Clean up temporary files"
    echo "  logs                    Clean up log files"
    echo "  temp                    Clean up temporary directories"
    echo ""
    echo "Examples:"
    echo "  $0 -t all -f            # Force cleanup everything"
    echo "  $0 -t k8s -n prod      # Clean up production Kubernetes resources"
    echo "  $0 -t docker            # Clean up Docker resources"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -t|--type)
                CLEANUP_TYPE="$2"
                shift 2
                ;;
            -n|--namespace)
                NAMESPACE="$2"
                shift 2
                ;;
            -e|--env)
                ENVIRONMENT="$2"
                shift 2
                ;;
            -f|--force)
                FORCE=true
                shift
                ;;
            -k|--keep-logs)
                KEEP_LOGS=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if kubectl is installed
    if ! command -v kubectl >/dev/null 2>&1; then
        log_error "kubectl is not installed"
        exit 1
    fi
    
    # Check if docker is installed
    if ! command -v docker >/dev/null 2>&1; then
        log_error "docker is not installed"
        exit 1
    fi
    
    log_success "Prerequisites check completed"
}

# Confirm cleanup
confirm_cleanup() {
    if [[ "$FORCE" == "true" ]]; then
        return 0
    fi
    
    log_warning "This will clean up $CLEANUP_TYPE resources in $ENVIRONMENT environment"
    read -p "Are you sure you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Cleanup cancelled"
        exit 0
    fi
}

# Clean up Kubernetes resources
cleanup_k8s() {
    log_info "Cleaning up Kubernetes resources..."
    
    # Check if namespace exists
    if ! kubectl get namespace $NAMESPACE >/dev/null 2>&1; then
        log_warning "Namespace $NAMESPACE does not exist"
        return
    fi
    
    # Delete all resources in namespace
    log_info "Deleting all resources in namespace: $NAMESPACE"
    kubectl delete all --all -n $NAMESPACE
    
    # Delete configmaps
    kubectl delete configmaps --all -n $NAMESPACE
    
    # Delete secrets
    kubectl delete secrets --all -n $NAMESPACE
    
    # Delete persistent volume claims
    kubectl delete pvc --all -n $NAMESPACE
    
    # Delete ingress
    kubectl delete ingress --all -n $NAMESPACE
    
    # Delete network policies
    kubectl delete networkpolicies --all -n $NAMESPACE
    
    # Delete service accounts
    kubectl delete serviceaccounts --all -n $NAMESPACE
    
    # Delete namespace
    log_info "Deleting namespace: $NAMESPACE"
    kubectl delete namespace $NAMESPACE
    
    log_success "Kubernetes resources cleaned up"
}

# Clean up Docker resources
cleanup_docker() {
    log_info "Cleaning up Docker resources..."
    
    # Stop and remove containers
    log_info "Stopping and removing containers..."
    docker stop $(docker ps -aq) 2>/dev/null || true
    docker rm $(docker ps -aq) 2>/dev/null || true
    
    # Remove images
    log_info "Removing images..."
    docker rmi $(docker images -q) 2>/dev/null || true
    
    # Remove volumes
    log_info "Removing volumes..."
    docker volume rm $(docker volume ls -q) 2>/dev/null || true
    
    # Remove networks
    log_info "Removing networks..."
    docker network rm $(docker network ls -q) 2>/dev/null || true
    
    # Prune system
    log_info "Pruning Docker system..."
    docker system prune -af
    
    log_success "Docker resources cleaned up"
}

# Clean up temporary files
cleanup_files() {
    log_info "Cleaning up temporary files..."
    
    # Remove build artifacts
    log_info "Removing build artifacts..."
    rm -rf bin/
    rm -rf dist/
    rm -rf build/
    
    # Remove test artifacts
    log_info "Removing test artifacts..."
    rm -f coverage.out
    rm -f coverage.html
    rm -f *.log
    
    # Remove temporary files
    log_info "Removing temporary files..."
    rm -f *.tmp
    rm -f *.temp
    rm -f .DS_Store
    rm -f Thumbs.db
    
    # Remove Terraform files
    log_info "Removing Terraform files..."
    rm -f terraform/*.tfplan
    rm -f terraform/*.tfstate.backup
    rm -rf terraform/.terraform/
    
    # Remove Ansible files
    log_info "Removing Ansible files..."
    rm -rf ansible/.ansible/
    rm -f ansible/*.retry
    
    # Remove Helm files
    log_info "Removing Helm files..."
    rm -rf helm/charts/
    rm -f helm/Chart.lock
    
    log_success "Temporary files cleaned up"
}

# Clean up log files
cleanup_logs() {
    if [[ "$KEEP_LOGS" == "true" ]]; then
        log_info "Keeping log files as requested"
        return
    fi
    
    log_info "Cleaning up log files..."
    
    # Remove application logs
    log_info "Removing application logs..."
    rm -f logs/*.log
    rm -f logs/*.txt
    
    # Remove system logs
    log_info "Removing system logs..."
    rm -f /var/log/obs-tools-usage/*.log 2>/dev/null || true
    
    # Remove Docker logs
    log_info "Removing Docker logs..."
    docker logs $(docker ps -aq) 2>/dev/null | head -n 0 > /dev/null || true
    
    log_success "Log files cleaned up"
}

# Clean up temporary directories
cleanup_temp() {
    log_info "Cleaning up temporary directories..."
    
    # Remove temporary directories
    log_info "Removing temporary directories..."
    rm -rf /tmp/obs-tools-usage-*
    rm -rf ./temp/
    rm -rf ./tmp/
    rm -rf ./restore-temp/
    
    # Remove backup directories
    log_info "Removing backup directories..."
    rm -rf ./backups/
    rm -rf ./security-reports/
    
    # Remove cache directories
    log_info "Removing cache directories..."
    rm -rf .cache/
    rm -rf .tmp/
    
    log_success "Temporary directories cleaned up"
}

# Clean up all resources
cleanup_all() {
    log_info "Cleaning up all resources..."
    
    cleanup_k8s
    cleanup_docker
    cleanup_files
    cleanup_logs
    cleanup_temp
    
    log_success "All resources cleaned up"
}

# Main cleanup function
main() {
    log_info "Starting cleanup process..."
    
    # Parse arguments
    parse_args "$@"
    
    # Check prerequisites
    check_prerequisites
    
    # Confirm cleanup
    confirm_cleanup
    
    # Perform cleanup based on type
    case $CLEANUP_TYPE in
        all)
            cleanup_all
            ;;
        k8s)
            cleanup_k8s
            ;;
        docker)
            cleanup_docker
            ;;
        files)
            cleanup_files
            ;;
        logs)
            cleanup_logs
            ;;
        temp)
            cleanup_temp
            ;;
        *)
            log_error "Invalid cleanup type: $CLEANUP_TYPE"
            usage
            exit 1
            ;;
    esac
    
    log_success "Cleanup completed successfully!"
}

# Run main function
main "$@"
