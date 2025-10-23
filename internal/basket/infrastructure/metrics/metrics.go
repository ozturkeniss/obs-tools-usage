package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

// Prometheus metrics for basket service
var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "basket_http_requests_total",
			Help: "Total number of HTTP requests to basket service",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "basket_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Business metrics
	basketsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "baskets_total",
			Help: "Total number of active baskets",
		},
	)

	basketItemsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "basket_items_total",
			Help: "Total number of items across all baskets",
		},
	)

	basketOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "basket_operations_total",
			Help: "Total number of basket operations",
		},
		[]string{"operation"},
	)

	// Redis metrics
	redisOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "basket_redis_operations_total",
			Help: "Total number of Redis operations",
		},
		[]string{"operation", "status"},
	)

	redisOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "basket_redis_operation_duration_seconds",
			Help:    "Redis operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// Product service metrics
	productServiceRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "basket_product_service_requests_total",
			Help: "Total number of requests to product service",
		},
		[]string{"operation", "status"},
	)

	productServiceRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "basket_product_service_request_duration_seconds",
			Help:    "Product service request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

// RecordHTTPRequest records HTTP request metrics
func RecordHTTPRequest(method, endpoint string, statusCode int, duration time.Duration) {
	statusCodeStr := string(rune(statusCode))
	httpRequestsTotal.WithLabelValues(method, endpoint, statusCodeStr).Inc()
	httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// RecordBasketOperation records basket operation metrics
func RecordBasketOperation(operation string) {
	basketOperationsTotal.WithLabelValues(operation).Inc()
}

// RecordRedisOperation records Redis operation metrics
func RecordRedisOperation(operation, status string, duration time.Duration) {
	redisOperationsTotal.WithLabelValues(operation, status).Inc()
	redisOperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordProductServiceRequest records product service request metrics
func RecordProductServiceRequest(operation, status string, duration time.Duration) {
	productServiceRequestsTotal.WithLabelValues(operation, status).Inc()
	productServiceRequestDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// UpdateBasketMetrics updates basket-related metrics
func UpdateBasketMetrics(basketCount, itemCount int) {
	basketsTotal.Set(float64(basketCount))
	basketItemsTotal.Set(float64(itemCount))
}

// HTTPLoggingMiddleware logs HTTP requests and responses
func HTTPLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Process request
		c.Next()
		
		// Calculate duration
		duration := time.Since(start)
		
		// Record metrics
		RecordHTTPRequest(c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
		
		// Log request
		logger := logrus.WithFields(logrus.Fields{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": c.Writer.Status(),
			"duration_ms": duration.Milliseconds(),
			"ip":          c.ClientIP(),
		})
		
		if c.Writer.Status() >= 400 {
			logger.Error("HTTP request completed with error")
		} else {
			logger.Info("HTTP request completed")
		}
	}
}
