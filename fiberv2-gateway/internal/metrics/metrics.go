package metrics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all the metrics for the gateway
type Metrics struct {
	RequestDuration *prometheus.HistogramVec
	RequestTotal    *prometheus.CounterVec
	ActiveRequests  prometheus.Gauge
	BackendHealth   *prometheus.GaugeVec
	CircuitBreaker  *prometheus.GaugeVec
}

// GatewayMetrics holds the global metrics instance
var GatewayMetrics *Metrics

// SetupMetrics sets up Prometheus metrics
func SetupMetrics(app *fiber.App) {
	// Create metrics
	GatewayMetrics = &Metrics{
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gateway_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status", "service"},
		),
		RequestTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status", "service"},
		),
		ActiveRequests: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "gateway_active_requests",
				Help: "Number of active requests being processed",
			},
		),
		BackendHealth: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gateway_backend_health",
				Help: "Health status of backend services (1=healthy, 0=unhealthy)",
			},
			[]string{"service", "backend"},
		),
		CircuitBreaker: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gateway_circuit_breaker_state",
				Help: "Circuit breaker state (0=closed, 1=open, 2=half_open)",
			},
			[]string{"service"},
		),
	}

	// Custom metrics middleware
	app.Use(func(c *fiber.Ctx) error {
		// Increment active requests
		GatewayMetrics.ActiveRequests.Inc()
		
		// Decrement active requests when done
		defer GatewayMetrics.ActiveRequests.Dec()

		// Continue to next middleware
		return c.Next()
	})
}

// RecordRequestDuration records the duration of a request
func RecordRequestDuration(method, path, status, service string, duration float64) {
	GatewayMetrics.RequestDuration.WithLabelValues(method, path, status, service).Observe(duration)
}

// RecordRequestTotal records the total number of requests
func RecordRequestTotal(method, path, status, service string) {
	GatewayMetrics.RequestTotal.WithLabelValues(method, path, status, service).Inc()
}

// UpdateBackendHealth updates the health status of a backend
func UpdateBackendHealth(service, backend string, healthy bool) {
	value := 0.0
	if healthy {
		value = 1.0
	}
	GatewayMetrics.BackendHealth.WithLabelValues(service, backend).Set(value)
}

// UpdateCircuitBreakerState updates the circuit breaker state
func UpdateCircuitBreakerState(service string, state int) {
	GatewayMetrics.CircuitBreaker.WithLabelValues(service).Set(float64(state))
}
