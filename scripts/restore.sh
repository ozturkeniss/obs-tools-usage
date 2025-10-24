#!/bin/bash

# OBS Tools Usage - Restore Script
# This script restores backups of the application and data

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
BACKUP_FILE=""
NAMESPACE="obs-tools-usage"
ENVIRONMENT="dev"
RESTORE_TYPE="all"
DRY_RUN=false
FORCE=false

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
    echo "Usage: $0 [OPTIONS] BACKUP_FILE"
    echo ""
    echo "Arguments:"
    echo "  BACKUP_FILE            Path to backup file (tar.gz or directory)"
    echo ""
    echo "Options:"
    echo "  -n, --namespace NS     Kubernetes namespace (default: obs-tools-usage)"
    echo "  -e, --env ENV          Environment (dev, staging, prod)"
    echo "  -t, --type TYPE        Restore type (all, config, data, logs)"
    echo "  -d, --dry-run          Show what would be restored without restoring"
    echo "  -f, --force            Force restore even if namespace exists"
    echo "  -h, --help             Show this help message"
    echo ""
    echo "Restore Types:"
    echo "  all                    Restore everything"
    echo "  config                 Restore configuration only"
    echo "  data                   Restore data only"
    echo "  logs                   Restore logs only"
    echo ""
    echo "Examples:"
    echo "  $0 backup.tar.gz                    # Restore full backup"
    echo "  $0 -t config backup.tar.gz          # Restore config only"
    echo "  $0 -d backup.tar.gz                 # Dry run restore"
    echo "  $0 -f backup.tar.gz                 # Force restore"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -n|--namespace)
                NAMESPACE="$2"
                shift 2
                ;;
            -e|--env)
                ENVIRONMENT="$2"
                shift 2
                ;;
            -t|--type)
                RESTORE_TYPE="$2"
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
            -h|--help)
                usage
                exit 0
                ;;
            *.tar.gz|*.tgz)
                BACKUP_FILE="$1"
                shift
                ;;
            *)
                if [[ -d "$1" ]]; then
                    BACKUP_FILE="$1"
                    shift
                else
                    log_error "Unknown option: $1"
                    usage
                    exit 1
                fi
                ;;
        esac
    done
    
    if [[ -z "$BACKUP_FILE" ]]; then
        log_error "Backup file is required"
        usage
        exit 1
    fi
    
    if [[ ! -f "$BACKUP_FILE" ]] && [[ ! -d "$BACKUP_FILE" ]]; then
        log_error "Backup file does not exist: $BACKUP_FILE"
        exit 1
    fi
}

# Check prerequisites
check_prerequisites() {
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
    
    # Check if tar is installed
    if ! command -v tar >/dev/null 2>&1; then
        log_error "tar is not installed"
        exit 1
    fi
}

# Extract backup
extract_backup() {
    local backup_file="$1"
    local extract_dir="./restore-temp"
    
    log_info "Extracting backup..."
    
    # Remove existing extract directory
    rm -rf "$extract_dir"
    
    if [[ -f "$backup_file" ]]; then
        # Extract compressed backup
        tar -xzf "$backup_file" -C "$(dirname "$extract_dir")" --strip-components=1
    elif [[ -d "$backup_file" ]]; then
        # Copy directory backup
        cp -r "$backup_file" "$extract_dir"
    fi
    
    echo "$extract_dir"
}

# Restore Kubernetes resources
restore_k8s_resources() {
    local restore_dir="$1"
    
    log_info "Restoring Kubernetes resources..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "Dry run: Would restore Kubernetes resources from $restore_dir"
        return
    fi
    
    # Create namespace if it doesn't exist
    if ! kubectl get namespace $NAMESPACE >/dev/null 2>&1; then
        kubectl create namespace $NAMESPACE
    fi
    
    # Restore resources
    if [[ -f "$restore_dir/resources.yaml" ]]; then
        kubectl apply -f "$restore_dir/resources.yaml"
    fi
    
    # Restore configmaps
    if [[ -f "$restore_dir/configmaps.yaml" ]]; then
        kubectl apply -f "$restore_dir/configmaps.yaml"
    fi
    
    # Restore secrets
    if [[ -f "$restore_dir/secrets.yaml" ]]; then
        kubectl apply -f "$restore_dir/secrets.yaml"
    fi
    
    # Restore persistent volumes
    if [[ -f "$restore_dir/volumes.yaml" ]]; then
        kubectl apply -f "$restore_dir/volumes.yaml"
    fi
    
    # Restore ingress
    if [[ -f "$restore_dir/ingress.yaml" ]]; then
        kubectl apply -f "$restore_dir/ingress.yaml"
    fi
    
    # Restore services
    if [[ -f "$restore_dir/services.yaml" ]]; then
        kubectl apply -f "$restore_dir/services.yaml"
    fi
    
    log_success "Kubernetes resources restored"
}

# Restore Helm releases
restore_helm_releases() {
    local restore_dir="$1"
    
    log_info "Restoring Helm releases..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "Dry run: Would restore Helm releases from $restore_dir"
        return
    fi
    
    # Restore each Helm release
    for helm_file in "$restore_dir"/helm-*.yaml; do
        if [[ -f "$helm_file" ]]; then
            local release_name=$(basename "$helm_file" .yaml | sed 's/helm-//')
            log_info "Restoring Helm release: $release_name"
            helm install "$release_name" "$helm_file" -n $NAMESPACE
        fi
    done
    
    log_success "Helm releases restored"
}

# Restore application data
restore_application_data() {
    local restore_dir="$1"
    
    log_info "Restoring application data..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "Dry run: Would restore application data from $restore_dir"
        return
    fi
    
    # Restore PostgreSQL data
    if [[ -f "$restore_dir/data/postgresql.sql" ]]; then
        log_info "Restoring PostgreSQL data..."
        kubectl exec -n $NAMESPACE deployment/postgresql -- psql -U postgres < "$restore_dir/data/postgresql.sql"
    fi
    
    # Restore MariaDB data
    if [[ -f "$restore_dir/data/mariadb.sql" ]]; then
        log_info "Restoring MariaDB data..."
        kubectl exec -n $NAMESPACE deployment/mariadb -- mysql -u root -p < "$restore_dir/data/mariadb.sql"
    fi
    
    # Restore Redis data
    if [[ -f "$restore_dir/data/redis.rdb" ]]; then
        log_info "Restoring Redis data..."
        kubectl cp "$restore_dir/data/redis.rdb" $NAMESPACE/$(kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=redis -o jsonpath='{.items[0].metadata.name}'):/tmp/redis-backup.rdb
        kubectl exec -n $NAMESPACE deployment/redis -- redis-cli --pipe < /tmp/redis-backup.rdb
    fi
    
    log_success "Application data restored"
}

# Restore application logs
restore_application_logs() {
    local restore_dir="$1"
    
    log_info "Restoring application logs..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "Dry run: Would restore application logs from $restore_dir"
        return
    fi
    
    # Restore logs for each pod
    for log_file in "$restore_dir/logs"/*.log; do
        if [[ -f "$log_file" ]]; then
            local pod_name=$(basename "$log_file" .log)
            log_info "Restoring logs for pod: $pod_name"
            kubectl cp "$log_file" $NAMESPACE/$pod_name:/tmp/restored-logs.log
        fi
    done
    
    log_success "Application logs restored"
}

# Restore configuration files
restore_configuration() {
    local restore_dir="$1"
    
    log_info "Restoring configuration files..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "Dry run: Would restore configuration files from $restore_dir"
        return
    fi
    
    # Restore Terraform state
    if [[ -f "$restore_dir/config/terraform.tfstate" ]]; then
        cp "$restore_dir/config/terraform.tfstate" terraform/terraform.tfstate
    fi
    
    # Restore Ansible inventory
    if [[ -f "$restore_dir/config/inventory.yml" ]]; then
        cp "$restore_dir/config/inventory.yml" ansible/inventory.yml
    fi
    
    # Restore Helm values
    if [[ -f "$restore_dir/config/helm-values.yaml" ]]; then
        cp "$restore_dir/config/helm-values.yaml" helm/values.yaml
    fi
    
    # Restore environment files
    if [[ -f "$restore_dir/config/env-dev" ]]; then
        cp "$restore_dir/config/env-dev" .env.dev
    fi
    
    if [[ -f "$restore_dir/config/env-staging" ]]; then
        cp "$restore_dir/config/env-staging" .env.staging
    fi
    
    if [[ -f "$restore_dir/config/env-prod" ]]; then
        cp "$restore_dir/config/env-prod" .env.prod
    fi
    
    log_success "Configuration files restored"
}

# Verify restore
verify_restore() {
    log_info "Verifying restore..."
    
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
    
    # Check if ingress is available
    local ingress=$(kubectl get ingress -n $NAMESPACE -o name)
    if [[ -n "$ingress" ]]; then
        log_success "Ingress is available"
    else
        log_warning "No ingress found"
    fi
}

# Cleanup
cleanup() {
    log_info "Cleaning up temporary files..."
    rm -rf "./restore-temp"
    log_success "Cleanup completed"
}

# Main restore function
main() {
    log_info "Starting restore process..."
    
    # Parse arguments
    parse_args "$@"
    
    # Check prerequisites
    check_prerequisites
    
    # Check if namespace exists and force is not set
    if kubectl get namespace $NAMESPACE >/dev/null 2>&1 && [[ "$FORCE" != "true" ]]; then
        log_error "Namespace $NAMESPACE already exists. Use -f to force restore."
        exit 1
    fi
    
    # Extract backup
    local restore_dir=$(extract_backup "$BACKUP_FILE")
    log_info "Restore directory: $restore_dir"
    
    # Perform restore based on type
    case $RESTORE_TYPE in
        all)
            restore_k8s_resources "$restore_dir"
            restore_helm_releases "$restore_dir"
            restore_application_data "$restore_dir"
            restore_application_logs "$restore_dir"
            restore_configuration "$restore_dir"
            ;;
        config)
            restore_k8s_resources "$restore_dir"
            restore_helm_releases "$restore_dir"
            restore_configuration "$restore_dir"
            ;;
        data)
            restore_application_data "$restore_dir"
            ;;
        logs)
            restore_application_logs "$restore_dir"
            ;;
        *)
            log_error "Invalid restore type: $RESTORE_TYPE"
            usage
            exit 1
            ;;
    esac
    
    # Verify restore
    verify_restore
    
    # Cleanup
    cleanup
    
    log_success "Restore completed successfully!"
    log_info "Application is available in namespace: $NAMESPACE"
}

# Handle errors
trap 'log_error "Restore failed. Cleaning up..."; cleanup; exit 1' ERR

# Run main function
main "$@"
