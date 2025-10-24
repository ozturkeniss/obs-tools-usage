#!/bin/bash

# OBS Tools Usage - Backup Script
# This script creates backups of the application and data

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
BACKUP_DIR="./backups"
NAMESPACE="obs-tools-usage"
ENVIRONMENT="dev"
BACKUP_TYPE="all"
RETENTION_DAYS=7
COMPRESS=true

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
    echo "  -d, --dir DIR           Backup directory (default: ./backups)"
    echo "  -n, --namespace NS      Kubernetes namespace (default: obs-tools-usage)"
    echo "  -e, --env ENV           Environment (dev, staging, prod)"
    echo "  -t, --type TYPE         Backup type (all, config, data, logs)"
    echo "  -r, --retention DAYS    Retention days (default: 7)"
    echo "  -c, --compress          Compress backups (default: true)"
    echo "  -h, --help              Show this help message"
    echo ""
    echo "Backup Types:"
    echo "  all                     Backup everything"
    echo "  config                  Backup configuration only"
    echo "  data                    Backup data only"
    echo "  logs                    Backup logs only"
    echo ""
    echo "Examples:"
    echo "  $0 -t all -r 30         # Full backup with 30-day retention"
    echo "  $0 -t config -n prod    # Config backup from production"
    echo "  $0 -t data -e staging   # Data backup from staging"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--dir)
                BACKUP_DIR="$2"
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
            -t|--type)
                BACKUP_TYPE="$2"
                shift 2
                ;;
            -r|--retention)
                RETENTION_DAYS="$2"
                shift 2
                ;;
            -c|--compress)
                COMPRESS=true
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

# Create backup directory
create_backup_dir() {
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local backup_path="$BACKUP_DIR/$ENVIRONMENT-$timestamp"
    
    mkdir -p "$backup_path"
    echo "$backup_path"
}

# Backup Kubernetes resources
backup_k8s_resources() {
    local backup_path="$1"
    
    log_info "Backing up Kubernetes resources..."
    
    # Create namespace backup
    kubectl get namespace $NAMESPACE -o yaml > "$backup_path/namespace.yaml"
    
    # Backup all resources in namespace
    kubectl get all -n $NAMESPACE -o yaml > "$backup_path/resources.yaml"
    
    # Backup configmaps
    kubectl get configmaps -n $NAMESPACE -o yaml > "$backup_path/configmaps.yaml"
    
    # Backup secrets
    kubectl get secrets -n $NAMESPACE -o yaml > "$backup_path/secrets.yaml"
    
    # Backup persistent volumes
    kubectl get pv,pvc -n $NAMESPACE -o yaml > "$backup_path/volumes.yaml"
    
    # Backup ingress
    kubectl get ingress -n $NAMESPACE -o yaml > "$backup_path/ingress.yaml"
    
    # Backup services
    kubectl get services -n $NAMESPACE -o yaml > "$backup_path/services.yaml"
    
    log_success "Kubernetes resources backed up"
}

# Backup Helm releases
backup_helm_releases() {
    local backup_path="$1"
    
    log_info "Backing up Helm releases..."
    
    # Get all releases in namespace
    helm list -n $NAMESPACE -o yaml > "$backup_path/helm-releases.yaml"
    
    # Backup each release
    for release in $(helm list -n $NAMESPACE -q); do
        log_info "Backing up Helm release: $release"
        helm get all $release -n $NAMESPACE > "$backup_path/helm-$release.yaml"
    done
    
    log_success "Helm releases backed up"
}

# Backup application data
backup_application_data() {
    local backup_path="$1"
    
    log_info "Backing up application data..."
    
    # Create data directory
    mkdir -p "$backup_path/data"
    
    # Backup PostgreSQL data
    if kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=postgresql >/dev/null 2>&1; then
        log_info "Backing up PostgreSQL data..."
        kubectl exec -n $NAMESPACE deployment/postgresql -- pg_dumpall -U postgres > "$backup_path/data/postgresql.sql"
    fi
    
    # Backup MariaDB data
    if kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=mariadb >/dev/null 2>&1; then
        log_info "Backing up MariaDB data..."
        kubectl exec -n $NAMESPACE deployment/mariadb -- mysqldump -u root -p --all-databases > "$backup_path/data/mariadb.sql"
    fi
    
    # Backup Redis data
    if kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=redis >/dev/null 2>&1; then
        log_info "Backing up Redis data..."
        kubectl exec -n $NAMESPACE deployment/redis -- redis-cli --rdb /tmp/redis-backup.rdb
        kubectl cp $NAMESPACE/$(kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=redis -o jsonpath='{.items[0].metadata.name}'):/tmp/redis-backup.rdb "$backup_path/data/redis.rdb"
    fi
    
    log_success "Application data backed up"
}

# Backup application logs
backup_application_logs() {
    local backup_path="$1"
    
    log_info "Backing up application logs..."
    
    # Create logs directory
    mkdir -p "$backup_path/logs"
    
    # Get all pods
    for pod in $(kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=obs-tools-usage -o jsonpath='{.items[*].metadata.name}'); do
        log_info "Backing up logs for pod: $pod"
        kubectl logs -n $NAMESPACE $pod > "$backup_path/logs/$pod.log"
    done
    
    # Backup system logs
    kubectl get events -n $NAMESPACE > "$backup_path/logs/events.log"
    
    log_success "Application logs backed up"
}

# Backup configuration files
backup_configuration() {
    local backup_path="$1"
    
    log_info "Backing up configuration files..."
    
    # Create config directory
    mkdir -p "$backup_path/config"
    
    # Backup Terraform state
    if [[ -f "terraform/terraform.tfstate" ]]; then
        cp terraform/terraform.tfstate "$backup_path/config/terraform.tfstate"
    fi
    
    # Backup Ansible inventory
    if [[ -f "ansible/inventory.yml" ]]; then
        cp ansible/inventory.yml "$backup_path/config/inventory.yml"
    fi
    
    # Backup Helm values
    if [[ -f "helm/values.yaml" ]]; then
        cp helm/values.yaml "$backup_path/config/helm-values.yaml"
    fi
    
    # Backup environment files
    if [[ -f ".env.dev" ]]; then
        cp .env.dev "$backup_path/config/env-dev"
    fi
    
    if [[ -f ".env.staging" ]]; then
        cp .env.staging "$backup_path/config/env-staging"
    fi
    
    if [[ -f ".env.prod" ]]; then
        cp .env.prod "$backup_path/config/env-prod"
    fi
    
    log_success "Configuration files backed up"
}

# Compress backup
compress_backup() {
    local backup_path="$1"
    
    if [[ "$COMPRESS" == "true" ]]; then
        log_info "Compressing backup..."
        tar -czf "$backup_path.tar.gz" -C "$(dirname "$backup_path")" "$(basename "$backup_path")"
        rm -rf "$backup_path"
        log_success "Backup compressed: $backup_path.tar.gz"
    fi
}

# Clean old backups
clean_old_backups() {
    log_info "Cleaning old backups (older than $RETENTION_DAYS days)..."
    
    find "$BACKUP_DIR" -name "*.tar.gz" -mtime +$RETENTION_DAYS -delete
    find "$BACKUP_DIR" -type d -mtime +$RETENTION_DAYS -exec rm -rf {} + 2>/dev/null || true
    
    log_success "Old backups cleaned"
}

# Create backup manifest
create_backup_manifest() {
    local backup_path="$1"
    
    log_info "Creating backup manifest..."
    
    cat > "$backup_path/backup-manifest.txt" << EOF
Backup Information
==================
Date: $(date)
Environment: $ENVIRONMENT
Namespace: $NAMESPACE
Backup Type: $BACKUP_TYPE
Retention Days: $RETENTION_DAYS
Compressed: $COMPRESS

Contents:
- Kubernetes resources
- Helm releases
- Application data
- Application logs
- Configuration files

Restore Instructions:
1. Extract backup: tar -xzf backup.tar.gz
2. Apply Kubernetes resources: kubectl apply -f resources.yaml
3. Restore Helm releases: helm install -f helm-*.yaml
4. Restore data: kubectl cp data/ <pod>:/tmp/
5. Restore logs: kubectl cp logs/ <pod>:/tmp/

Generated by OBS Tools Usage Backup Script
EOF

    log_success "Backup manifest created"
}

# Main backup function
main() {
    log_info "Starting backup process..."
    
    # Parse arguments
    parse_args "$@"
    
    # Check prerequisites
    check_prerequisites
    
    # Create backup directory
    local backup_path=$(create_backup_dir)
    log_info "Backup directory: $backup_path"
    
    # Perform backup based on type
    case $BACKUP_TYPE in
        all)
            backup_k8s_resources "$backup_path"
            backup_helm_releases "$backup_path"
            backup_application_data "$backup_path"
            backup_application_logs "$backup_path"
            backup_configuration "$backup_path"
            ;;
        config)
            backup_k8s_resources "$backup_path"
            backup_helm_releases "$backup_path"
            backup_configuration "$backup_path"
            ;;
        data)
            backup_application_data "$backup_path"
            ;;
        logs)
            backup_application_logs "$backup_path"
            ;;
        *)
            log_error "Invalid backup type: $BACKUP_TYPE"
            usage
            exit 1
            ;;
    esac
    
    # Create backup manifest
    create_backup_manifest "$backup_path"
    
    # Compress backup
    compress_backup "$backup_path"
    
    # Clean old backups
    clean_old_backups
    
    log_success "Backup completed successfully!"
    log_info "Backup location: $backup_path"
}

# Run main function
main "$@"
