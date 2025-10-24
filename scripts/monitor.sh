#!/bin/bash

# OBS Tools Usage - Monitoring Script
# This script provides monitoring and health check capabilities

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
NAMESPACE="obs-tools-usage"
ENVIRONMENT="dev"
WATCH_MODE=false
LOG_LEVEL="info"
EXPORT_METRICS=false

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
    echo "Usage: $0 [OPTIONS] [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  status              Show application status"
    echo "  health              Check application health"
    echo "  logs                Show application logs"
    echo "  metrics             Show application metrics"
    echo "  events              Show Kubernetes events"
    echo "  resources           Show resource usage"
    echo "  alerts              Check for alerts"
    echo "  dashboard           Open monitoring dashboard"
    echo ""
    echo "Options:"
    echo "  -n, --namespace NS  Kubernetes namespace (default: obs-tools-usage)"
    echo "  -e, --env ENV       Environment (dev, staging, prod)"
    echo "  -w, --watch         Watch mode (continuous monitoring)"
    echo "  -l, --log-level LVL Log level (debug, info, warn, error)"
    echo "  -x, --export        Export metrics to file"
    echo "  -h, --help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 status                    # Show status"
    echo "  $0 health -w                 # Watch health continuously"
    echo "  $0 logs -n obs-tools-usage  # Show logs from specific namespace"
    echo "  $0 metrics -x                # Export metrics"
}

# Parse command line arguments
parse_args() {
    COMMAND=""
    
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
            -w|--watch)
                WATCH_MODE=true
                shift
                ;;
            -l|--log-level)
                LOG_LEVEL="$2"
                shift 2
                ;;
            -x|--export)
                EXPORT_METRICS=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            status|health|logs|metrics|events|resources|alerts|dashboard)
                COMMAND="$1"
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
    
    if [[ -z "$COMMAND" ]]; then
        log_error "Command is required"
        usage
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
    
    # Check if curl is installed
    if ! command -v curl >/dev/null 2>&1; then
        log_error "curl is not installed"
        exit 1
    fi
}

# Show application status
show_status() {
    log_info "Application Status in namespace: $NAMESPACE"
    echo ""
    
    # Show pods
    echo "=== PODS ==="
    kubectl get pods -n $NAMESPACE -o wide
    echo ""
    
    # Show services
    echo "=== SERVICES ==="
    kubectl get services -n $NAMESPACE
    echo ""
    
    # Show ingress
    echo "=== INGRESS ==="
    kubectl get ingress -n $NAMESPACE
    echo ""
    
    # Show deployments
    echo "=== DEPLOYMENTS ==="
    kubectl get deployments -n $NAMESPACE
    echo ""
    
    # Show replicasets
    echo "=== REPLICASETS ==="
    kubectl get replicasets -n $NAMESPACE
    echo ""
}

# Check application health
check_health() {
    log_info "Checking application health..."
    
    # Check if namespace exists
    if ! kubectl get namespace $NAMESPACE >/dev/null 2>&1; then
        log_error "Namespace $NAMESPACE does not exist"
        return 1
    fi
    
    # Check pod health
    log_info "Checking pod health..."
    kubectl get pods -n $NAMESPACE --field-selector=status.phase!=Running
    
    # Check service endpoints
    log_info "Checking service endpoints..."
    kubectl get endpoints -n $NAMESPACE
    
    # Check application health endpoints
    log_info "Checking application health endpoints..."
    
    # Get gateway service
    GATEWAY_SERVICE=$(kubectl get service -n $NAMESPACE -l app.kubernetes.io/component=gateway -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
    
    if [[ -n "$GATEWAY_SERVICE" ]]; then
        # Port forward to gateway service
        kubectl port-forward service/$GATEWAY_SERVICE 8080:8080 -n $NAMESPACE &
        PORT_FORWARD_PID=$!
        sleep 5
        
        # Test health endpoint
        if curl -f http://localhost:8080/health >/dev/null 2>&1; then
            log_success "Gateway health check passed"
        else
            log_warning "Gateway health check failed"
        fi
        
        # Test metrics endpoint
        if curl -f http://localhost:8080/metrics >/dev/null 2>&1; then
            log_success "Metrics endpoint accessible"
        else
            log_warning "Metrics endpoint not accessible"
        fi
        
        # Kill port forward
        kill $PORT_FORWARD_PID 2>/dev/null || true
    else
        log_warning "Gateway service not found"
    fi
}

# Show application logs
show_logs() {
    log_info "Showing application logs..."
    
    # Show logs for all pods
    kubectl logs -l app.kubernetes.io/name=obs-tools-usage -n $NAMESPACE --tail=100
    
    if [[ "$WATCH_MODE" == "true" ]]; then
        log_info "Watching logs (Ctrl+C to stop)..."
        kubectl logs -l app.kubernetes.io/name=obs-tools-usage -n $NAMESPACE -f
    fi
}

# Show application metrics
show_metrics() {
    log_info "Showing application metrics..."
    
    # Get gateway service
    GATEWAY_SERVICE=$(kubectl get service -n $NAMESPACE -l app.kubernetes.io/component=gateway -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
    
    if [[ -n "$GATEWAY_SERVICE" ]]; then
        # Port forward to gateway service
        kubectl port-forward service/$GATEWAY_SERVICE 8080:8080 -n $NAMESPACE &
        PORT_FORWARD_PID=$!
        sleep 5
        
        # Get metrics
        METRICS=$(curl -s http://localhost:8080/metrics)
        
        if [[ "$EXPORT_METRICS" == "true" ]]; then
            echo "$METRICS" > metrics-$(date +%Y%m%d-%H%M%S).txt
            log_success "Metrics exported to file"
        else
            echo "$METRICS"
        fi
        
        # Kill port forward
        kill $PORT_FORWARD_PID 2>/dev/null || true
    else
        log_warning "Gateway service not found"
    fi
}

# Show Kubernetes events
show_events() {
    log_info "Showing Kubernetes events..."
    
    kubectl get events -n $NAMESPACE --sort-by='.lastTimestamp'
}

# Show resource usage
show_resources() {
    log_info "Showing resource usage..."
    
    # Show node resource usage
    echo "=== NODE RESOURCES ==="
    kubectl top nodes
    echo ""
    
    # Show pod resource usage
    echo "=== POD RESOURCES ==="
    kubectl top pods -n $NAMESPACE
    echo ""
    
    # Show resource quotas
    echo "=== RESOURCE QUOTAS ==="
    kubectl get resourcequotas -n $NAMESPACE
    echo ""
    
    # Show resource requests and limits
    echo "=== RESOURCE REQUESTS AND LIMITS ==="
    kubectl describe pods -n $NAMESPACE | grep -A 5 "Requests\|Limits"
}

# Check for alerts
check_alerts() {
    log_info "Checking for alerts..."
    
    # Check for failed pods
    FAILED_PODS=$(kubectl get pods -n $NAMESPACE --field-selector=status.phase=Failed -o name)
    if [[ -n "$FAILED_PODS" ]]; then
        log_warning "Failed pods found:"
        echo "$FAILED_PODS"
    else
        log_success "No failed pods found"
    fi
    
    # Check for pending pods
    PENDING_PODS=$(kubectl get pods -n $NAMESPACE --field-selector=status.phase=Pending -o name)
    if [[ -n "$PENDING_PODS" ]]; then
        log_warning "Pending pods found:"
        echo "$PENDING_PODS"
    else
        log_success "No pending pods found"
    fi
    
    # Check for events with warning or error
    WARNING_EVENTS=$(kubectl get events -n $NAMESPACE --field-selector=type=Warning --no-headers | wc -l)
    if [[ $WARNING_EVENTS -gt 0 ]]; then
        log_warning "Warning events found: $WARNING_EVENTS"
        kubectl get events -n $NAMESPACE --field-selector=type=Warning
    else
        log_success "No warning events found"
    fi
}

# Open monitoring dashboard
open_dashboard() {
    log_info "Opening monitoring dashboard..."
    
    # Check if Grafana is available
    GRAFANA_SERVICE=$(kubectl get service -n $NAMESPACE -l app.kubernetes.io/name=grafana -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
    
    if [[ -n "$GRAFANA_SERVICE" ]]; then
        log_info "Opening Grafana dashboard..."
        kubectl port-forward service/$GRAFANA_SERVICE 3000:3000 -n $NAMESPACE &
        PORT_FORWARD_PID=$!
        sleep 5
        
        # Open browser
        if command -v xdg-open >/dev/null 2>&1; then
            xdg-open http://localhost:3000
        elif command -v open >/dev/null 2>&1; then
            open http://localhost:3000
        else
            log_info "Please open http://localhost:3000 in your browser"
        fi
        
        log_info "Press Ctrl+C to stop port forwarding"
        wait $PORT_FORWARD_PID
    else
        log_warning "Grafana service not found"
        
        # Try Prometheus
        PROMETHEUS_SERVICE=$(kubectl get service -n $NAMESPACE -l app.kubernetes.io/name=prometheus -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
        
        if [[ -n "$PROMETHEUS_SERVICE" ]]; then
            log_info "Opening Prometheus dashboard..."
            kubectl port-forward service/$PROMETHEUS_SERVICE 9090:9090 -n $NAMESPACE &
            PORT_FORWARD_PID=$!
            sleep 5
            
            # Open browser
            if command -v xdg-open >/dev/null 2>&1; then
                xdg-open http://localhost:9090
            elif command -v open >/dev/null 2>&1; then
                open http://localhost:9090
            else
                log_info "Please open http://localhost:9090 in your browser"
            fi
            
            log_info "Press Ctrl+C to stop port forwarding"
            wait $PORT_FORWARD_PID
        else
            log_error "No monitoring dashboard found"
        fi
    fi
}

# Main monitoring function
main() {
    # Parse arguments
    parse_args "$@"
    
    # Check prerequisites
    check_prerequisites
    
    # Execute command
    case $COMMAND in
        status)
            show_status
            ;;
        health)
            check_health
            ;;
        logs)
            show_logs
            ;;
        metrics)
            show_metrics
            ;;
        events)
            show_events
            ;;
        resources)
            show_resources
            ;;
        alerts)
            check_alerts
            ;;
        dashboard)
            open_dashboard
            ;;
        *)
            log_error "Unknown command: $COMMAND"
            usage
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
