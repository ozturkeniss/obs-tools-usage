#!/bin/bash

# OBS Tools Usage - Security Script
# This script performs security scans and audits

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
SCAN_TYPE="all"
OUTPUT_DIR="./security-reports"
NAMESPACE="obs-tools-usage"
ENVIRONMENT="dev"
SEVERITY_THRESHOLD="HIGH"

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
    echo "  -t, --type TYPE         Scan type (all, code, container, infra, secrets)"
    echo "  -o, --output DIR         Output directory (default: ./security-reports)"
    echo "  -n, --namespace NS       Kubernetes namespace (default: obs-tools-usage)"
    echo "  -e, --env ENV            Environment (dev, staging, prod)"
    echo "  -s, --severity LEVEL     Severity threshold (LOW, MEDIUM, HIGH, CRITICAL)"
    echo "  -h, --help              Show this help message"
    echo ""
    echo "Scan Types:"
    echo "  all                     Run all security scans"
    echo "  code                    Scan source code for vulnerabilities"
    echo "  container               Scan container images for vulnerabilities"
    echo "  infra                   Scan infrastructure for security issues"
    echo "  secrets                 Scan for secrets and sensitive data"
    echo ""
    echo "Examples:"
    echo "  $0 -t all -s HIGH       # Run all scans with HIGH severity threshold"
    echo "  $0 -t code -o reports   # Scan code and save to reports directory"
    echo "  $0 -t container -n prod # Scan containers in production namespace"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -t|--type)
                SCAN_TYPE="$2"
                shift 2
                ;;
            -o|--output)
                OUTPUT_DIR="$2"
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
            -s|--severity)
                SEVERITY_THRESHOLD="$2"
                shift 2
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
    
    if ! command -v trivy >/dev/null 2>&1; then
        missing_tools+=("trivy")
    fi
    
    if ! command -v gosec >/dev/null 2>&1; then
        missing_tools+=("gosec")
    fi
    
    if ! command -v checkov >/dev/null 2>&1; then
        missing_tools+=("checkov")
    fi
    
    if ! command -v trufflehog >/dev/null 2>&1; then
        missing_tools+=("trufflehog")
    fi
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        log_warning "Missing tools: ${missing_tools[*]}"
        log_info "Installing missing tools..."
        
        # Install Trivy
        if [[ " ${missing_tools[*]} " =~ " trivy " ]]; then
            curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin
        fi
        
        # Install gosec
        if [[ " ${missing_tools[*]} " =~ " gosec " ]]; then
            go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
        fi
        
        # Install Checkov
        if [[ " ${missing_tools[*]} " =~ " checkov " ]]; then
            pip install checkov
        fi
        
        # Install TruffleHog
        if [[ " ${missing_tools[*]} " =~ " trufflehog " ]]; then
            go install github.com/trufflesecurity/trufflehog/v3@latest
        fi
    fi
    
    log_success "Prerequisites check completed"
}

# Create output directory
create_output_dir() {
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local scan_dir="$OUTPUT_DIR/$ENVIRONMENT-$timestamp"
    
    mkdir -p "$scan_dir"
    echo "$scan_dir"
}

# Scan source code
scan_source_code() {
    local output_dir="$1"
    
    log_info "Scanning source code for vulnerabilities..."
    
    # Create code scan directory
    mkdir -p "$output_dir/code"
    
    # Run gosec
    log_info "Running gosec security scan..."
    gosec -fmt json -out "$output_dir/code/gosec-report.json" ./...
    gosec -fmt sarif -out "$output_dir/code/gosec-report.sarif" ./...
    
    # Run govulncheck
    log_info "Running govulncheck..."
    govulncheck ./... > "$output_dir/code/govulncheck-report.txt" 2>&1 || true
    
    # Run go audit
    log_info "Running go audit..."
    go list -json -deps ./... | nancy sleuth > "$output_dir/code/nancy-report.txt" 2>&1 || true
    
    log_success "Source code scan completed"
}

# Scan container images
scan_container_images() {
    local output_dir="$1"
    
    log_info "Scanning container images for vulnerabilities..."
    
    # Create container scan directory
    mkdir -p "$output_dir/container"
    
    # List of images to scan
    local images=(
        "obs-tools-usage/product-service:latest"
        "obs-tools-usage/basket-service:latest"
        "obs-tools-usage/payment-service:latest"
        "obs-tools-usage/gateway:latest"
    )
    
    for image in "${images[@]}"; do
        log_info "Scanning image: $image"
        local image_name=$(echo "$image" | sed 's/[:/]/_/g')
        
        # Run Trivy scan
        trivy image --format json --output "$output_dir/container/trivy-$image_name.json" "$image"
        trivy image --format table --output "$output_dir/container/trivy-$image_name.txt" "$image"
        trivy image --format sarif --output "$output_dir/container/trivy-$image_name.sarif" "$image"
    done
    
    log_success "Container image scan completed"
}

# Scan infrastructure
scan_infrastructure() {
    local output_dir="$1"
    
    log_info "Scanning infrastructure for security issues..."
    
    # Create infrastructure scan directory
    mkdir -p "$output_dir/infrastructure"
    
    # Scan Terraform files
    if [[ -d "terraform" ]]; then
        log_info "Scanning Terraform files..."
        checkov -d terraform --output json --output-file-path "$output_dir/infrastructure/checkov-terraform.json"
        checkov -d terraform --output sarif --output-file-path "$output_dir/infrastructure/checkov-terraform.sarif"
    fi
    
    # Scan Kubernetes manifests
    if [[ -d "helm" ]]; then
        log_info "Scanning Kubernetes manifests..."
        checkov -d helm --output json --output-file-path "$output_dir/infrastructure/checkov-k8s.json"
        checkov -d helm --output sarif --output-file-path "$output_dir/infrastructure/checkov-k8s.sarif"
    fi
    
    # Scan Docker files
    if [[ -d "dockerfiles" ]]; then
        log_info "Scanning Docker files..."
        checkov -d dockerfiles --output json --output-file-path "$output_dir/infrastructure/checkov-docker.json"
        checkov -d dockerfiles --output sarif --output-file-path "$output_dir/infrastructure/checkov-docker.sarif"
    fi
    
    log_success "Infrastructure scan completed"
}

# Scan for secrets
scan_secrets() {
    local output_dir="$1"
    
    log_info "Scanning for secrets and sensitive data..."
    
    # Create secrets scan directory
    mkdir -p "$output_dir/secrets"
    
    # Run TruffleHog
    log_info "Running TruffleHog scan..."
    trufflehog filesystem . --output json --output-file "$output_dir/secrets/trufflehog-report.json"
    trufflehog filesystem . --output table --output-file "$output_dir/secrets/trufflehog-report.txt"
    
    # Run git secrets
    log_info "Running git secrets scan..."
    if command -v git-secrets >/dev/null 2>&1; then
        git secrets --scan > "$output_dir/secrets/git-secrets-report.txt" 2>&1 || true
    fi
    
    # Scan for hardcoded secrets in code
    log_info "Scanning for hardcoded secrets..."
    grep -r -i "password\|secret\|key\|token" --include="*.go" --include="*.yaml" --include="*.yml" . > "$output_dir/secrets/hardcoded-secrets.txt" 2>/dev/null || true
    
    log_success "Secrets scan completed"
}

# Scan Kubernetes cluster
scan_k8s_cluster() {
    local output_dir="$1"
    
    log_info "Scanning Kubernetes cluster for security issues..."
    
    # Create k8s scan directory
    mkdir -p "$output_dir/k8s"
    
    # Run Trivy on Kubernetes
    log_info "Running Trivy on Kubernetes cluster..."
    trivy k8s cluster --format json --output "$output_dir/k8s/trivy-k8s.json"
    trivy k8s cluster --format table --output "$output_dir/k8s/trivy-k8s.txt"
    
    # Check RBAC permissions
    log_info "Checking RBAC permissions..."
    kubectl auth can-i --list --all-namespaces > "$output_dir/k8s/rbac-permissions.txt"
    
    # Check network policies
    log_info "Checking network policies..."
    kubectl get networkpolicies -A -o yaml > "$output_dir/k8s/network-policies.yaml"
    
    # Check pod security
    log_info "Checking pod security..."
    kubectl get pods -A -o json | jq '.items[] | select(.spec.securityContext == null) | .metadata.name' > "$output_dir/k8s/pods-without-security-context.txt"
    
    log_success "Kubernetes cluster scan completed"
}

# Generate security report
generate_security_report() {
    local output_dir="$1"
    
    log_info "Generating security report..."
    
    # Create summary report
    cat > "$output_dir/security-summary.md" << EOF
# Security Scan Report

## Scan Information
- Date: $(date)
- Environment: $ENVIRONMENT
- Namespace: $NAMESPACE
- Scan Type: $SCAN_TYPE
- Severity Threshold: $SEVERITY_THRESHOLD

## Scan Results

### Source Code Security
- Gosec Report: [gosec-report.json](code/gosec-report.json)
- Vulnerability Check: [govulncheck-report.txt](code/govulncheck-report.txt)
- Dependency Audit: [nancy-report.txt](code/nancy-report.txt)

### Container Security
- Trivy Reports: [container/](container/)
- Image Vulnerabilities: See individual reports

### Infrastructure Security
- Terraform Security: [checkov-terraform.json](infrastructure/checkov-terraform.json)
- Kubernetes Security: [checkov-k8s.json](infrastructure/checkov-k8s.json)
- Docker Security: [checkov-docker.json](infrastructure/checkov-docker.json)

### Secrets Detection
- TruffleHog Report: [trufflehog-report.json](secrets/trufflehog-report.json)
- Git Secrets: [git-secrets-report.txt](secrets/git-secrets-report.txt)
- Hardcoded Secrets: [hardcoded-secrets.txt](secrets/hardcoded-secrets.txt)

### Kubernetes Cluster Security
- Cluster Vulnerabilities: [trivy-k8s.json](k8s/trivy-k8s.json)
- RBAC Permissions: [rbac-permissions.txt](k8s/rbac-permissions.txt)
- Network Policies: [network-policies.yaml](k8s/network-policies.yaml)
- Pod Security: [pods-without-security-context.txt](k8s/pods-without-security-context.txt)

## Recommendations

1. **High Priority**: Address all CRITICAL and HIGH severity findings
2. **Medium Priority**: Review and fix MEDIUM severity issues
3. **Low Priority**: Consider LOW severity findings for future improvements
4. **Secrets**: Remove or rotate any exposed secrets
5. **RBAC**: Review and minimize RBAC permissions
6. **Network**: Implement network policies for traffic isolation
7. **Pod Security**: Apply security contexts to all pods

## Next Steps

1. Review all findings and prioritize fixes
2. Update dependencies to latest secure versions
3. Implement security best practices
4. Schedule regular security scans
5. Monitor for new vulnerabilities

---
*Generated by OBS Tools Usage Security Script*
EOF

    log_success "Security report generated"
}

# Main security scan function
main() {
    log_info "Starting security scan..."
    
    # Parse arguments
    parse_args "$@"
    
    # Check prerequisites
    check_prerequisites
    
    # Create output directory
    local output_dir=$(create_output_dir)
    log_info "Output directory: $output_dir"
    
    # Perform scan based on type
    case $SCAN_TYPE in
        all)
            scan_source_code "$output_dir"
            scan_container_images "$output_dir"
            scan_infrastructure "$output_dir"
            scan_secrets "$output_dir"
            scan_k8s_cluster "$output_dir"
            ;;
        code)
            scan_source_code "$output_dir"
            ;;
        container)
            scan_container_images "$output_dir"
            ;;
        infra)
            scan_infrastructure "$output_dir"
            ;;
        secrets)
            scan_secrets "$output_dir"
            ;;
        *)
            log_error "Invalid scan type: $SCAN_TYPE"
            usage
            exit 1
            ;;
    esac
    
    # Generate security report
    generate_security_report "$output_dir"
    
    log_success "Security scan completed successfully!"
    log_info "Reports saved to: $output_dir"
    log_info "Summary report: $output_dir/security-summary.md"
}

# Run main function
main "$@"
