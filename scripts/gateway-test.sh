#!/bin/bash

# Gateway Testing Script
# This script tests the FiberV2 Gateway functionality including rate limiting, circuit breaker, and load balancing

set -e

echo "🧪 Gateway Testing Script..."

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
GATEWAY_URL="http://localhost:8083"
TEST_TYPE="all"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --url)
            GATEWAY_URL="$2"
            shift 2
            ;;
        --type)
            TEST_TYPE="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "OPTIONS:"
            echo "  --url URL    Gateway URL (default: http://localhost:8083)"
            echo "  --type TYPE  Test type (health|rate-limit|circuit-breaker|load-balancer|all)"
            echo "  -h, --help   Show this help message"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Function to make HTTP request
make_request() {
    local url="$1"
    local method="${2:-GET}"
    local headers="$3"
    local data="$4"
    
    if [ -n "$data" ]; then
        curl -s -X "$method" -H "$headers" -d "$data" "$url"
    else
        curl -s -X "$method" -H "$headers" "$url"
    fi
}

# Function to test gateway health
test_health() {
    print_status "Testing Gateway Health..."
    
    local response=$(make_request "$GATEWAY_URL/health")
    local status_code=$(curl -s -o /dev/null -w "%{http_code}" "$GATEWAY_URL/health")
    
    if [ "$status_code" = "200" ]; then
        print_success "Gateway health check passed!"
        echo "Response: $response"
    else
        print_error "Gateway health check failed! Status: $status_code"
        return 1
    fi
}

# Function to test rate limiting
test_rate_limiting() {
    print_status "Testing Rate Limiting..."
    
    # Test rate limit with multiple requests
    local success_count=0
    local rate_limit_count=0
    
    for i in {1..110}; do
        local status_code=$(curl -s -o /dev/null -w "%{http_code}" "$GATEWAY_URL/api/product")
        
        if [ "$status_code" = "200" ]; then
            success_count=$((success_count + 1))
        elif [ "$status_code" = "429" ]; then
            rate_limit_count=$((rate_limit_count + 1))
        fi
        
        if [ $((i % 10)) -eq 0 ]; then
            print_status "Request $i: Status $status_code"
        fi
    done
    
    print_status "Rate Limiting Test Results:"
    echo "  Successful requests: $success_count"
    echo "  Rate limited requests: $rate_limit_count"
    
    if [ $rate_limit_count -gt 0 ]; then
        print_success "Rate limiting is working!"
    else
        print_warning "Rate limiting might not be working properly"
    fi
}

# Function to test circuit breaker
test_circuit_breaker() {
    print_status "Testing Circuit Breaker..."
    
    # Test with a non-existent service to trigger circuit breaker
    local response=$(make_request "$GATEWAY_URL/api/nonexistent")
    local status_code=$(curl -s -o /dev/null -w "%{http_code}" "$GATEWAY_URL/api/nonexistent")
    
    print_status "Circuit Breaker Test:"
    echo "  URL: $GATEWAY_URL/api/nonexistent"
    echo "  Status Code: $status_code"
    echo "  Response: $response"
    
    if [ "$status_code" = "503" ] || [ "$status_code" = "502" ]; then
        print_success "Circuit breaker is working!"
    else
        print_warning "Circuit breaker behavior unclear"
    fi
}

# Function to test load balancing
test_load_balancing() {
    print_status "Testing Load Balancing..."
    
    # Test multiple requests to see if they're distributed
    print_status "Making 10 requests to test load balancing..."
    
    for i in {1..10}; do
        local response=$(make_request "$GATEWAY_URL/api/product")
        local status_code=$(curl -s -o /dev/null -w "%{http_code}" "$GATEWAY_URL/api/product")
        
        echo "Request $i: Status $status_code"
        
        if [ "$status_code" != "200" ] && [ "$status_code" != "429" ]; then
            print_warning "Request $i failed with status $status_code"
        fi
    done
    
    print_success "Load balancing test completed!"
}

# Function to test gateway metrics
test_metrics() {
    print_status "Testing Gateway Metrics..."
    
    local response=$(make_request "$GATEWAY_URL/metrics")
    local status_code=$(curl -s -o /dev/null -w "%{http_code}" "$GATEWAY_URL/metrics")
    
    if [ "$status_code" = "200" ]; then
        print_success "Metrics endpoint is accessible!"
        echo "Metrics response length: $(echo "$response" | wc -c) characters"
        
        # Check for specific metrics
        if echo "$response" | grep -q "gateway_request_duration_seconds"; then
            print_success "Request duration metrics found!"
        fi
        
        if echo "$response" | grep -q "gateway_requests_total"; then
            print_success "Request total metrics found!"
        fi
        
        if echo "$response" | grep -q "gateway_circuit_breaker_state"; then
            print_success "Circuit breaker metrics found!"
        fi
    else
        print_error "Metrics endpoint failed! Status: $status_code"
        return 1
    fi
}

# Function to test gateway admin endpoints
test_admin() {
    print_status "Testing Gateway Admin Endpoints..."
    
    # Test admin health
    local admin_health=$(make_request "$GATEWAY_URL/admin/health")
    local admin_status=$(curl -s -o /dev/null -w "%{http_code}" "$GATEWAY_URL/admin/health")
    
    if [ "$admin_status" = "200" ]; then
        print_success "Admin health endpoint accessible!"
        echo "Admin health response: $admin_health"
    else
        print_warning "Admin health endpoint not accessible (Status: $admin_status)"
    fi
    
    # Test rate limit status
    local rate_limit_status=$(make_request "$GATEWAY_URL/admin/rate-limit/status")
    local rate_limit_code=$(curl -s -o /dev/null -w "%{http_code}" "$GATEWAY_URL/admin/rate-limit/status")
    
    if [ "$rate_limit_code" = "200" ]; then
        print_success "Rate limit status endpoint accessible!"
        echo "Rate limit status: $rate_limit_status"
    else
        print_warning "Rate limit status endpoint not accessible (Status: $rate_limit_code)"
    fi
}

# Function to test service routing
test_service_routing() {
    print_status "Testing Service Routing..."
    
    # Test product service routing
    local product_response=$(make_request "$GATEWAY_URL/api/product")
    local product_status=$(curl -s -o /dev/null -w "%{http_code}" "$GATEWAY_URL/api/product")
    
    print_status "Product Service Routing:"
    echo "  Status Code: $product_status"
    
    # Test basket service routing
    local basket_response=$(make_request "$GATEWAY_URL/api/basket")
    local basket_status=$(curl -s -o /dev/null -w "%{http_code}" "$GATEWAY_URL/api/basket")
    
    print_status "Basket Service Routing:"
    echo "  Status Code: $basket_status"
    
    # Test payment service routing
    local payment_response=$(make_request "$GATEWAY_URL/api/payment")
    local payment_status=$(curl -s -o /dev/null -w "%{http_code}" "$GATEWAY_URL/api/payment")
    
    print_status "Payment Service Routing:"
    echo "  Status Code: $payment_status"
    
    print_success "Service routing test completed!"
}

# Main test execution
print_status "Starting Gateway Tests..."
print_status "Gateway URL: $GATEWAY_URL"
print_status "Test Type: $TEST_TYPE"
echo ""

case $TEST_TYPE in
    health)
        test_health
        ;;
    rate-limit)
        test_rate_limiting
        ;;
    circuit-breaker)
        test_circuit_breaker
        ;;
    load-balancer)
        test_load_balancing
        ;;
    metrics)
        test_metrics
        ;;
    admin)
        test_admin
        ;;
    routing)
        test_service_routing
        ;;
    all)
        test_health
        echo ""
        test_rate_limiting
        echo ""
        test_circuit_breaker
        echo ""
        test_load_balancing
        echo ""
        test_metrics
        echo ""
        test_admin
        echo ""
        test_service_routing
        ;;
    *)
        print_error "Unknown test type: $TEST_TYPE"
        exit 1
        ;;
esac

echo ""
print_success "Gateway testing completed!"
print_status "For more detailed testing, use specific test types:"
echo $0 --help
