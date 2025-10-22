#!/bin/bash

# Linting and formatting script
# This script runs various code quality checks

set -e

echo "🔍 Code Quality Check Script..."

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
FIX=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        format)
            ACTION="format"
            shift
            ;;
        lint)
            ACTION="lint"
            shift
            ;;
        vet)
            ACTION="vet"
            shift
            ;;
        all)
            ACTION="all"
            shift
            ;;
        --fix)
            FIX=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [ACTION] [OPTIONS]"
            echo ""
            echo "ACTIONS:"
            echo "  format    Format Go code with gofmt"
            echo "  lint      Run linter (golangci-lint)"
            echo "  vet       Run go vet"
            echo "  all       Run all checks (default)"
            echo ""
            echo "OPTIONS:"
            echo "  --fix     Fix issues automatically where possible"
            echo "  -h, --help Show this help message"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Function to format code
format_code() {
    print_status "Formatting Go code..."
    
    if [ "$FIX" = true ]; then
        go fmt ./...
        print_success "Code formatted!"
    else
        # Check if code is formatted
        UNFORMATTED=$(gofmt -l .)
        if [ -z "$UNFORMATTED" ]; then
            print_success "Code is properly formatted!"
        else
            print_error "Code is not formatted properly!"
            echo "Unformatted files:"
            echo "$UNFORMATTED"
            exit 1
        fi
    fi
}

# Function to run go vet
run_vet() {
    print_status "Running go vet..."
    
    go vet ./...
    
    if [ $? -eq 0 ]; then
        print_success "go vet passed!"
    else
        print_error "go vet found issues!"
        exit 1
    fi
}

# Function to run linter
run_linter() {
    print_status "Running linter..."
    
    if ! command -v golangci-lint &> /dev/null; then
        print_warning "golangci-lint not found, installing..."
        
        # Install golangci-lint
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
        
        if [ $? -eq 0 ]; then
            print_success "golangci-lint installed successfully!"
        else
            print_error "Failed to install golangci-lint!"
            exit 1
        fi
    fi
    
    # Create golangci-lint config if it doesn't exist
    if [ ! -f .golangci.yml ]; then
        print_status "Creating golangci-lint configuration..."
        cat > .golangci.yml << EOF
linters-settings:
  gofmt:
    simplify: true
  goimports:
    local-prefixes: obs-tools-usage
  govet:
    check-shadowing: true
  gocyclo:
    min-complexity: 15
  maligned:
    suggest-new: true
  dupl:
    threshold: 100
  goconst:
    min-len: 2
    min-occurrences: 2
  misspell:
    locale: US
  lll:
    line-length: 140
  unused:
    check-exported: false
  unparam:
    check-exported: false
  nakedret:
    max-func-lines: 30
  prealloc:
    simple: true
    range-loops: true
    for-loops: false
  gocritic:
    enabled-tags:
      - diagnostic
      - experimental
      - opinionated
      - performance
      - style
    disabled-checks:
      - dupImport
      - ifElseChain
      - octalLiteral
      - whyNoLint
      - wrapperFunc

linters:
  enable:
    - gofmt
    - goimports
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - typecheck
    - gocyclo
    - goconst
    - misspell
    - lll
    - unparam
    - nakedret
    - prealloc
    - gocritic

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - gomnd
        - goconst
        - gocritic
        - gocyclo
        - dupl
        - funlen
        - gocognit
        - lll
        - nakedret
        - unparam
        - wsl
EOF
    fi
    
    if [ "$FIX" = true ]; then
        golangci-lint run --fix
    else
        golangci-lint run
    fi
    
    if [ $? -eq 0 ]; then
        print_success "Linter passed!"
    else
        print_error "Linter found issues!"
        exit 1
    fi
}

# Function to check imports
check_imports() {
    print_status "Checking imports..."
    
    if ! command -v goimports &> /dev/null; then
        print_warning "goimports not found, installing..."
        go install golang.org/x/tools/cmd/goimports@latest
    fi
    
    UNFORMATTED=$(goimports -l .)
    if [ -z "$UNFORMATTED" ]; then
        print_success "Imports are properly formatted!"
    else
        if [ "$FIX" = true ]; then
            goimports -w .
            print_success "Imports fixed!"
        else
            print_error "Imports are not properly formatted!"
            echo "Files with import issues:"
            echo "$UNFORMATTED"
            exit 1
        fi
    fi
}

# Function to check for TODO/FIXME comments
check_todos() {
    print_status "Checking for TODO/FIXME comments..."
    
    TODOS=$(grep -r "TODO\|FIXME\|XXX\|HACK" --include="*.go" . || true)
    if [ -z "$TODOS" ]; then
        print_success "No TODO/FIXME comments found!"
    else
        print_warning "Found TODO/FIXME comments:"
        echo "$TODOS"
    fi
}

# Function to check for security issues
check_security() {
    print_status "Checking for security issues..."
    
    if ! command -v gosec &> /dev/null; then
        print_warning "gosec not found, installing..."
        go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
    fi
    
    gosec ./...
    
    if [ $? -eq 0 ]; then
        print_success "Security check passed!"
    else
        print_warning "Security check found potential issues!"
    fi
}

# Main execution
case $ACTION in
    format)
        format_code
        check_imports
        ;;
    lint)
        run_linter
        ;;
    vet)
        run_vet
        ;;
    all)
        format_code
        check_imports
        run_vet
        run_linter
        check_todos
        check_security
        print_success "All code quality checks completed!"
        ;;
esac
