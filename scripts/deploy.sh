#!/bin/bash

# OBS Tools Usage - Deployment Script
# This script deploys the application to different environments

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
ENVIRONMENT=""
NAMESPACE=""
DRY_RUN=false
FORCE=false
SKIP_TESTS=false

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
    echo "  -e, --environment ENV    Target environment (dev, staging, prod)"
    echo "  -n, --namespace NS      Kubernetes namespace"
    echo "  -d, --dry-run           Show what would be deployed without deploying"
    echo "  -f, --force             Force deployment even if tests fail"
    echo "  -s, --skip-tests        Skip running tests before deployment"
    echo "  -h, --help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 -e dev                    # Deploy to development"
    echo "  $0 -e staging -d             # Dry run for staging"
    echo "  $0 -e prod -f                # Force deploy to production"
    echo "  $0 -e dev -s                 # Deploy to dev without tests"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -e|--environment)
                ENVIRONMENT="$2"
                shift 2
                ;;
            -n|--namespace)
                NAMESPACE="$2"
                shift 2
                ;;
            -d|--dry-run)
                DRY_RUN=true
                shift
                ;;
            -f|--force)
                FORCE=true
                shift
                ;;
            -s|--skip-tests)
                SKIP_TESTS=true
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

# Validate environment
validate_environment() {
    if [[ -z "$ENVIRONMENT" ]]; then
        log_error "Environment is required"
        usage
        exit 1
    fi
    
    case $ENVIRONMENT in
        dev|development)
            ENVIRONMENT="dev"
            NAMESPACE=${NAMESPACE:-"obs-tools-usage-dev"}
            ;;
        staging)
            ENVIRONMENT="staging"
            NAMESPACE=${NAMESPACE:-"obs-tools-usage-staging"}
            ;;
        prod|production)
            ENVIRONMENT="prod"
            NAMESPACE=${NAMESPACE:-"obs-tools-usage"}
            ;;
        *)
            log_error "Invalid environment: $ENVIRONMENT"
            log_error "Valid environments: dev, staging, prod"
            exit 1
            ;;
    esac
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if kubectl is installed
    if ! command -v kubectl >/dev/null 2>&1; then
        log_error "kubectl is not installed"
        exit 1
    fi
    
    # Check if helm is installed
    if ! command -v helm >/dev/null 2>&1; then
        log_error "helm is not installed"
        exit 1
    fi
    
    # Check if terraform is installed (for infrastructure)
    if ! command -v terraform >/dev/null 2>&1; then
        log_error "terraform is not installed"
        exit 1
    fi
    
    # Check if ansible is installed (for configuration)
    if ! command -v ansible >/dev/null 2>&1; then
        log_error "ansible is not installed"
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

# Run tests
run_tests() {
    if [[ "$SKIP_TESTS" == "true" ]]; then
        log_warning "Skipping tests"
        return
    fi
    
    log_info "Running tests..."
    
    # Run Go tests
    go test -v -race -coverprofile=coverage.out ./...
    
    # Run linting
    golangci-lint run
    
    # Run security scan
    if command -v govulncheck >/dev/null 2>&1; then
        govulncheck ./...
    fi
    
    log_success "Tests passed"
}

# Build Docker images
build_images() {
    log_info "Building Docker images..."
    
    # Build product service
    docker build -f dockerfiles/product.dockerfile -t obs-tools-usage/product-service:latest .
    
    # Build basket service
    docker build -f dockerfiles/basket.dockerfile -t obs-tools-usage/basket-service:latest .
    
    # Build payment service
    docker build -f dockerfiles/payment.dockerfile -t obs-tools-usage/payment-service:latest .
    
    # Build gateway
    docker build -f dockerfiles/gateway.dockerfile -t obs-tools-usage/gateway:latest .
    
    log_success "Docker images built successfully"
}

# Deploy infrastructure
deploy_infrastructure() {
    log_info "Deploying infrastructure with Terraform..."
    
    cd terraform
    
    # Initialize Terraform
    terraform init
    
    # Plan deployment
    terraform plan -var-file="terraform.tfvars.${ENVIRONMENT}" -out=tfplan
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "Dry run: Infrastructure changes planned"
        return
    fi
    
    # Apply changes
    terraform apply tfplan
    
    cd ..
    
    log_success "Infrastructure deployed successfully"
}

# Configure with Ansible
configure_with_ansible() {
    log_info "Configuring with Ansible..."
    
    cd ansible
    
    # Run Ansible playbooks
    ansible-playbook -i inventory.yml main.yml --limit $ENVIRONMENT
    ansible-playbook -i inventory.yml k8s-setup.yml --limit $ENVIRONMENT
    ansible-playbook -i inventory.yml app-deploy.yml --limit $ENVIRONMENT
    
    cd ..
    
    log_success "Configuration completed successfully"
}

# Deploy application
deploy_application() {
    log_info "Deploying application with Helm..."
    
    # Update Helm dependencies
    helm dependency update helm/
    
    # Deploy with Helm
    if [[ "$DRY_RUN" == "true" ]]; then
        helm install obs-tools-usage helm/ --namespace $NAMESPACE --create-namespace --dry-run --debug
        log_info "Dry run: Application deployment planned"
        return
    fi
    
    # Install or upgrade
    helm upgrade --install obs-tools-usage helm/ \
        --namespace $NAMESPACE \
        --create-namespace \
        --values helm/values-${ENVIRONMENT}.yaml \
        --wait \
        --timeout=10m
    
    log_success "Application deployed successfully"
}

# Verify deployment
verify_deployment() {
    log_info "Verifying deployment..."
    
    # Check pods
    kubectl get pods -n $NAMESPACE
    
    # Check services
    kubectl get services -n $NAMESPACE
    
    # Check ingress
    kubectl get ingress -n $NAMESPACE
    
    # Wait for pods to be ready
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=obs-tools-usage -n $NAMESPACE --timeout=300s
    
    # Check application health
    if kubectl get service obs-tools-usage-gateway -n $NAMESPACE >/dev/null 2>&1; then
        GATEWAY_URL=$(kubectl get service obs-tools-usage-gateway -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
        if [[ -n "$GATEWAY_URL" ]]; then
            log_info "Testing application health at http://$GATEWAY_URL/health"
            curl -f http://$GATEWAY_URL/health || log_warning "Health check failed"
        fi
    fi
    
    log_success "Deployment verification completed"
}

# Rollback deployment
rollback_deployment() {
    log_info "Rolling back deployment..."
    
    # Rollback Helm release
    helm rollback obs-tools-usage -n $NAMESPACE
    
    log_success "Rollback completed"
}

# Cleanup
cleanup() {
    log_info "Cleaning up..."
    
    # Remove temporary files
    rm -f terraform/tfplan
    rm -f coverage.out
    
    log_success "Cleanup completed"
}

# Main deployment function
main() {
    log_info "Starting deployment to $ENVIRONMENT environment..."
    
    # Parse arguments
    parse_args "$@"
    
    # Validate environment
    validate_environment
    
    # Check prerequisites
    check_prerequisites
    
    # Run tests (unless skipped)
    if [[ "$SKIP_TESTS" == "false" ]]; then
        run_tests
    fi
    
    # Build images
    build_images
    
    # Deploy infrastructure
    deploy_infrastructure
    
    # Configure with Ansible
    configure_with_ansible
    
    # Deploy application
    deploy_application
    
    # Verify deployment
    verify_deployment
    
    # Cleanup
    cleanup
    
    log_success "Deployment to $ENVIRONMENT completed successfully!"
    log_info "Application is available in namespace: $NAMESPACE"
}

# Handle errors
trap 'log_error "Deployment failed. Rolling back..."; rollback_deployment; exit 1' ERR

# Run main function
main "$@"
