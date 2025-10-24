#!/bin/bash

# OBS Tools Usage - Setup Script
# This script sets up the development environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install dependencies
install_dependencies() {
    log_info "Installing dependencies..."
    
    # Update package manager
    if command_exists apt-get; then
        sudo apt-get update
        sudo apt-get install -y curl wget git unzip jq tree htop vim
    elif command_exists yum; then
        sudo yum update -y
        sudo yum install -y curl wget git unzip jq tree htop vim
    elif command_exists brew; then
        brew update
        brew install curl wget git unzip jq tree htop vim
    else
        log_warning "Package manager not found. Please install dependencies manually."
    fi
}

# Install Go
install_go() {
    if command_exists go; then
        log_info "Go is already installed: $(go version)"
        return
    fi
    
    log_info "Installing Go..."
    
    # Download and install Go
    GO_VERSION="1.22.0"
    GO_ARCH="linux-amd64"
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        GO_ARCH="darwin-amd64"
    fi
    
    wget -q "https://golang.org/dl/go${GO_VERSION}.${GO_ARCH}.tar.gz"
    sudo tar -C /usr/local -xzf "go${GO_VERSION}.${GO_ARCH}.tar.gz"
    rm "go${GO_VERSION}.${GO_ARCH}.tar.gz"
    
    # Add Go to PATH
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
    
    log_success "Go installed successfully"
}

# Install Docker
install_docker() {
    if command_exists docker; then
        log_info "Docker is already installed: $(docker --version)"
        return
    fi
    
    log_info "Installing Docker..."
    
    # Install Docker
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    rm get-docker.sh
    
    # Add user to docker group
    sudo usermod -aG docker $USER
    
    # Install Docker Compose
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    
    log_success "Docker installed successfully"
}

# Install kubectl
install_kubectl() {
    if command_exists kubectl; then
        log_info "kubectl is already installed: $(kubectl version --client)"
        return
    fi
    
    log_info "Installing kubectl..."
    
    # Download kubectl
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    rm kubectl
    
    log_success "kubectl installed successfully"
}

# Install Helm
install_helm() {
    if command_exists helm; then
        log_info "Helm is already installed: $(helm version)"
        return
    fi
    
    log_info "Installing Helm..."
    
    # Download and install Helm
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
    
    log_success "Helm installed successfully"
}

# Install Terraform
install_terraform() {
    if command_exists terraform; then
        log_info "Terraform is already installed: $(terraform version)"
        return
    fi
    
    log_info "Installing Terraform..."
    
    # Download Terraform
    TERRAFORM_VERSION="1.6.0"
    wget -q "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip"
    unzip "terraform_${TERRAFORM_VERSION}_linux_amd64.zip"
    sudo mv terraform /usr/local/bin/
    rm "terraform_${TERRAFORM_VERSION}_linux_amd64.zip"
    
    log_success "Terraform installed successfully"
}

# Install Ansible
install_ansible() {
    if command_exists ansible; then
        log_info "Ansible is already installed: $(ansible --version)"
        return
    fi
    
    log_info "Installing Ansible..."
    
    # Install Ansible
    if command_exists pip3; then
        pip3 install ansible
    elif command_exists pip; then
        pip install ansible
    else
        log_error "pip not found. Please install pip first."
        exit 1
    fi
    
    log_success "Ansible installed successfully"
}

# Install AWS CLI
install_aws_cli() {
    if command_exists aws; then
        log_info "AWS CLI is already installed: $(aws --version)"
        return
    fi
    
    log_info "Installing AWS CLI..."
    
    # Download and install AWS CLI
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    unzip awscliv2.zip
    sudo ./aws/install
    rm -rf aws awscliv2.zip
    
    log_success "AWS CLI installed successfully"
}

# Install k6
install_k6() {
    if command_exists k6; then
        log_info "k6 is already installed: $(k6 version)"
        return
    fi
    
    log_info "Installing k6..."
    
    # Install k6
    sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
    echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
    sudo apt-get update
    sudo apt-get install k6
    
    log_success "k6 installed successfully"
}

# Setup Go modules
setup_go_modules() {
    log_info "Setting up Go modules..."
    
    # Initialize Go modules
    go mod init obs-tools-usage || true
    go mod tidy
    
    # Install Go tools
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install golang.org/x/tools/cmd/goimports@latest
    go install golang.org/x/vuln/cmd/govulncheck@latest
    
    log_success "Go modules setup completed"
}

# Setup Git hooks
setup_git_hooks() {
    log_info "Setting up Git hooks..."
    
    # Create pre-commit hook
    cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
# Pre-commit hook

echo "Running pre-commit checks..."

# Run gofmt
gofmt -w .
git add .

# Run goimports
goimports -w .
git add .

# Run tests
go test ./...

# Run linting
golangci-lint run

echo "Pre-commit checks completed"
EOF

    chmod +x .git/hooks/pre-commit
    
    log_success "Git hooks setup completed"
}

# Create development environment file
create_env_file() {
    log_info "Creating development environment file..."
    
    cat > .env.dev << 'EOF'
# Development Environment Variables
ENVIRONMENT=development
LOG_LEVEL=debug
LOG_FORMAT=text
LOG_OUTPUT=console

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=product_service
DB_SSL_MODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Service URLs
PRODUCT_SERVICE_URL=http://localhost:8080
BASKET_SERVICE_URL=http://localhost:8081
PAYMENT_SERVICE_URL=http://localhost:8082
GATEWAY_URL=http://localhost:8083

# Kafka Configuration
KAFKA_BROKERS=localhost:9092

# Monitoring
PROMETHEUS_ENABLED=true
GRAFANA_ENABLED=true
EOF

    log_success "Environment file created"
}

# Main setup function
main() {
    log_info "Starting OBS Tools Usage setup..."
    
    # Check if running as root
    if [[ $EUID -eq 0 ]]; then
        log_error "This script should not be run as root"
        exit 1
    fi
    
    # Install dependencies
    install_dependencies
    install_go
    install_docker
    install_kubectl
    install_helm
    install_terraform
    install_ansible
    install_aws_cli
    install_k6
    
    # Setup project
    setup_go_modules
    setup_git_hooks
    create_env_file
    
    log_success "Setup completed successfully!"
    log_info "Please restart your terminal or run 'source ~/.bashrc' to apply changes"
    log_info "You can now run 'make dev' to start the development environment"
}

# Run main function
main "$@"
