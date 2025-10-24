#!/bin/bash

# OBS Tools Usage - Update Script
# This script updates the application and dependencies

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
UPDATE_TYPE="all"
NAMESPACE="obs-tools-usage"
ENVIRONMENT="dev"
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
    echo "  -t, --type TYPE         Update type (all, deps, images, helm, terraform, ansible)"
    echo "  -n, --namespace NS      Kubernetes namespace (default: obs-tools-usage)"
    echo "  -e, --env ENV            Environment (dev, staging, prod)"
    echo "  -d, --dry-run            Show what would be updated without updating"
    echo "  -f, --force              Force update even if tests fail"
    echo "  -s, --skip-tests         Skip running tests before update"
    echo "  -h, --help               Show this help message"
    echo ""
    echo "Update Types:"
    echo "  all                     Update everything"
    echo "  deps                    Update dependencies only"
    echo "  images                  Update container images only"
    echo "  helm                    Update Helm charts only"
    echo "  terraform               Update Terraform configuration only"
    echo "  ansible                 Update Ansible playbooks only"
    echo ""
    echo "Examples:"
    echo "  $0 -t all -e prod       # Update everything in production"
    echo "  $0 -t deps -d            # Dry run dependency update"
    echo "  $0 -t images -f          # Force update container images"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -t|--type)
                UPDATE_TYPE="$2"
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

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if required tools are installed
    local missing_tools=()
    
    if ! command -v go >/dev/null 2>&1; then
        missing_tools+=("go")
    fi
    
    if ! command -v docker >/dev/null 2>&1; then
        missing_tools+=("docker")
    fi
    
    if ! command -v helm >/dev/null 2>&1; then
        missing_tools+=("helm")
    fi
    
    if ! command -v terraform >/dev/null 2>&1; then
        missing_tools+=("terraform")
    fi
    
    if ! command -v ansible >/dev/null 2>&1; then
        missing_tools+=("ansible")
    fi
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        log_error "Missing tools: ${missing_tools[*]}"
        log_info "Please install missing tools first"
        exit 1
    fi
    
    log_success "Prerequisites check completed"
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

# Update dependencies
update_dependencies() {
    log_info "Updating dependencies..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "Dry run: Would update dependencies"
        return
    fi
    
    # Update Go modules
    log_info "Updating Go modules..."
    go get -u ./...
    go mod tidy
    
    # Update Helm dependencies
    log_info "Updating Helm dependencies..."
    helm dependency update helm/
    
    # Update Terraform providers
    log_info "Updating Terraform providers..."
    cd terraform
    terraform init -upgrade
    cd ..
    
    # Update Ansible dependencies
    log_info "Updating Ansible dependencies..."
    cd ansible
    ansible-galaxy install -r requirements.yml --force
    cd ..
    
    log_success "Dependencies updated"
}

# Update container images
update_container_images() {
    log_info "Updating container images..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "Dry run: Would update container images"
        return
    fi
    
    # Build new images
    log_info "Building new container images..."
    docker build -f dockerfiles/product.dockerfile -t obs-tools-usage/product-service:latest .
    docker build -f dockerfiles/basket.dockerfile -t obs-tools-usage/basket-service:latest .
    docker build -f dockerfiles/payment.dockerfile -t obs-tools-usage/payment-service:latest .
    docker build -f dockerfiles/gateway.dockerfile -t obs-tools-usage/gateway:latest .
    
    # Tag images with environment
    docker tag obs-tools-usage/product-service:latest obs-tools-usage/product-service:$ENVIRONMENT
    docker tag obs-tools-usage/basket-service:latest obs-tools-usage/basket-service:$ENVIRONMENT
    docker tag obs-tools-usage/payment-service:latest obs-tools-usage/payment-service:$ENVIRONMENT
    docker tag obs-tools-usage/gateway:latest obs-tools-usage/gateway:$ENVIRONMENT
    
    log_success "Container images updated"
}

# Update Helm charts
update_helm_charts() {
    log_info "Updating Helm charts..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "Dry run: Would update Helm charts"
        return
    fi
    
    # Update Helm dependencies
    helm dependency update helm/
    
    # Upgrade Helm release
    helm upgrade --install obs-tools-usage helm/ \
        --namespace $NAMESPACE \
        --create-namespace \
        --values helm/values-${ENVIRONMENT}.yaml \
        --wait \
        --timeout=10m
    
    log_success "Helm charts updated"
}

# Update Terraform configuration
update_terraform() {
    log_info "Updating Terraform configuration..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "Dry run: Would update Terraform configuration"
        return
    fi
    
    cd terraform
    
    # Initialize Terraform
    terraform init
    
    # Plan changes
    terraform plan -var-file="terraform.tfvars.${ENVIRONMENT}" -out=tfplan
    
    # Apply changes
    terraform apply tfplan
    
    cd ..
    
    log_success "Terraform configuration updated"
}

# Update Ansible playbooks
update_ansible() {
    log_info "Updating Ansible playbooks..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "Dry run: Would update Ansible playbooks"
        return
    fi
    
    cd ansible
    
    # Run Ansible playbooks
    ansible-playbook -i inventory.yml main.yml --limit $ENVIRONMENT
    ansible-playbook -i inventory.yml k8s-setup.yml --limit $ENVIRONMENT
    ansible-playbook -i inventory.yml app-deploy.yml --limit $ENVIRONMENT
    
    cd ..
    
    log_success "Ansible playbooks updated"
}

# Verify update
verify_update() {
    log_info "Verifying update..."
    
    # Check if namespace exists
    if ! kubectl get namespace $NAMESPACE >/dev/null 2>&1; then
        log_error "Namespace $NAMESPACE does not exist"
        return 1
    fi
    
    # Check if pods are running
    local failed_pods=$(kubectl get pods -n $NAMESPACE --field-selector=status.phase!=Running -o name)
    if [[ -n "$failed_pods" ]]; then
        log_warning "Some pods are not running:"
        echo "$failed_pods"
    else
        log_success "All pods are running"
    fi
    
    # Check if services are available
    local services=$(kubectl get services -n $NAMESPACE -o name)
    if [[ -n "$services" ]]; then
        log_success "Services are available"
    else
        log_warning "No services found"
    fi
    
    # Check application health
    if kubectl get service obs-tools-usage-gateway -n $NAMESPACE >/dev/null 2>&1; then
        GATEWAY_URL=$(kubectl get service obs-tools-usage-gateway -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
        if [[ -n "$GATEWAY_URL" ]]; then
            log_info "Testing application health at http://$GATEWAY_URL/health"
            curl -f http://$GATEWAY_URL/health || log_warning "Health check failed"
        fi
    fi
}

# Rollback update
rollback_update() {
    log_info "Rolling back update..."
    
    # Rollback Helm release
    helm rollback obs-tools-usage -n $NAMESPACE
    
    log_success "Update rolled back"
}

# Cleanup
cleanup() {
    log_info "Cleaning up..."
    
    # Remove temporary files
    rm -f terraform/tfplan
    rm -f coverage.out
    
    log_success "Cleanup completed"
}

# Main update function
main() {
    log_info "Starting update process..."
    
    # Parse arguments
    parse_args "$@"
    
    # Check prerequisites
    check_prerequisites
    
    # Run tests (unless skipped)
    if [[ "$SKIP_TESTS" == "false" ]]; then
        run_tests
    fi
    
    # Perform update based on type
    case $UPDATE_TYPE in
        all)
            update_dependencies
            update_container_images
            update_helm_charts
            update_terraform
            update_ansible
            ;;
        deps)
            update_dependencies
            ;;
        images)
            update_container_images
            ;;
        helm)
            update_helm_charts
            ;;
        terraform)
            update_terraform
            ;;
        ansible)
            update_ansible
            ;;
        *)
            log_error "Invalid update type: $UPDATE_TYPE"
            usage
            exit 1
            ;;
    esac
    
    # Verify update
    verify_update
    
    # Cleanup
    cleanup
    
    log_success "Update completed successfully!"
    log_info "Application is available in namespace: $NAMESPACE"
}

# Handle errors
trap 'log_error "Update failed. Rolling back..."; rollback_update; exit 1' ERR

# Run main function
main "$@"
