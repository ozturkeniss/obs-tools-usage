package product

import (
	"fmt"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

// Prometheus metrics
var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "endpoint"},
	)

	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "endpoint"},
	)

	// Business metrics
	productsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "products_total",
			Help: "Total number of products",
		},
	)

	productsCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "products_created_total",
			Help: "Total number of products created",
		},
	)

	productsUpdatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "products_updated_total",
			Help: "Total number of products updated",
		},
	)

	productsDeletedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "products_deleted_total",
			Help: "Total number of products deleted",
		},
	)

	// System metrics
	memoryAllocBytes = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_alloc_bytes",
			Help: "Current memory allocation in bytes",
		},
	)

	goroutinesTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines_total",
			Help: "Current number of goroutines",
		},
	)

	// Database metrics
	databaseOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "status"},
	)

	databaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_operation_duration_seconds",
			Help:    "Database operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

// PerformanceMetrics holds performance-related metrics
type PerformanceMetrics struct {
	ResponseTime    int64   `json:"response_time_ms"`
	MemoryAlloc     uint64  `json:"memory_alloc_bytes"`
	MemorySys       uint64  `json:"memory_sys_bytes"`
	NumGoroutines   int     `json:"num_goroutines"`
	NumGC           uint32  `json:"num_gc"`
	GCForcedRuns    uint32  `json:"gc_forced_runs"`
	Endpoint        string  `json:"endpoint"`
	Method          string  `json:"method"`
	StatusCode      int     `json:"status_code"`
	RequestSize     int     `json:"request_size_bytes"`
	ResponseSize    int     `json:"response_size_bytes"`
}

// LogPerformanceMetrics logs performance metrics
func LogPerformanceMetrics(logger *logrus.Entry, metrics PerformanceMetrics) {
	logger.WithFields(logrus.Fields{
		"performance": true,
		"endpoint":    metrics.Endpoint,
		"method":      metrics.Method,
		"status_code": metrics.StatusCode,
		"response_time_ms": metrics.ResponseTime,
		"memory_alloc_bytes": metrics.MemoryAlloc,
		"memory_sys_bytes":   metrics.MemorySys,
		"num_goroutines":     metrics.NumGoroutines,
		"num_gc":             metrics.NumGC,
		"gc_forced_runs":     metrics.GCForcedRuns,
		"request_size_bytes": metrics.RequestSize,
		"response_size_bytes": metrics.ResponseSize,
	}).Info("Performance metrics")
}

// GetSystemMetrics returns current system metrics
func GetSystemMetrics() (uint64, uint64, int, uint32, uint32) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return m.Alloc, m.Sys, runtime.NumGoroutine(), m.NumGC, m.NumForcedGC
}

// PerformanceMiddleware logs performance metrics for each request
func PerformanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get initial metrics
		startMemAlloc, startMemSys, startGoroutines, startNumGC, startGCForced := GetSystemMetrics()
		startTime := time.Now()
		
		// Process request
		c.Next()
		
		// Calculate metrics
		duration := time.Since(startTime)
		endMemAlloc, endMemSys, endGoroutines, endNumGC, endGCForced := GetSystemMetrics()
		
		// Calculate request/response sizes
		requestSize := int(c.Request.ContentLength)
		if requestSize < 0 {
			requestSize = 0
		}
		responseSize := c.Writer.Size()
		
		// Create performance metrics
		metrics := PerformanceMetrics{
			ResponseTime:    duration.Milliseconds(),
			MemoryAlloc:     endMemAlloc,
			MemorySys:       endMemSys,
			NumGoroutines:   endGoroutines,
			NumGC:           endNumGC - startNumGC,
			GCForcedRuns:    endGCForced - startGCForced,
			Endpoint:        c.Request.URL.Path,
			Method:          c.Request.Method,
			StatusCode:      c.Writer.Status(),
			RequestSize:     requestSize,
			ResponseSize:    responseSize,
		}
		
		// Get logger with context
		logger := GetLoggerFromContext(c)
		
		// Log performance metrics
		LogPerformanceMetrics(logger, metrics)
	}
}

// LogSlowQueries logs queries that exceed a threshold
func LogSlowQueries(logger *logrus.Entry, operation string, duration time.Duration, threshold time.Duration) {
	if duration > threshold {
		logger.WithFields(logrus.Fields{
			"slow_query": true,
			"operation":  operation,
			"duration_ms": duration.Milliseconds(),
			"threshold_ms": threshold.Milliseconds(),
		}).Warn("Slow query detected")
	}
}

// Prometheus metrics functions

// RecordHTTPRequest records HTTP request metrics
func RecordHTTPRequest(method, endpoint string, statusCode int, duration time.Duration, requestSize, responseSize int) {
	statusCodeStr := fmt.Sprintf("%d", statusCode)
	
	httpRequestsTotal.WithLabelValues(method, endpoint, statusCodeStr).Inc()
	httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	httpRequestSize.WithLabelValues(method, endpoint).Observe(float64(requestSize))
	httpResponseSize.WithLabelValues(method, endpoint).Observe(float64(responseSize))
}

// RecordDatabaseOperation records database operation metrics
func RecordDatabaseOperation(operation, status string, duration time.Duration) {
	databaseOperationsTotal.WithLabelValues(operation, status).Inc()
	databaseOperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordProductCreated records product creation metric
func RecordProductCreated() {
	productsCreatedTotal.Inc()
	UpdateProductsTotal()
}

// RecordProductUpdated records product update metric
func RecordProductUpdated() {
	productsUpdatedTotal.Inc()
	UpdateProductsTotal()
}

// RecordProductDeleted records product deletion metric
func RecordProductDeleted() {
	productsDeletedTotal.Inc()
	UpdateProductsTotal()
}

// UpdateProductsTotal updates the total products count
func UpdateProductsTotal() {
	// This would typically query the database for actual count
	// For now, we'll use a simple counter approach
	// In a real implementation, you'd query the repository
}

// UpdateSystemMetrics updates system-level metrics
func UpdateSystemMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	memoryAllocBytes.Set(float64(memStats.Alloc))
	goroutinesTotal.Set(float64(runtime.NumGoroutine()))
}

// GetPrometheusMetrics returns the Prometheus registry
func GetPrometheusMetrics() *prometheus.Registry {
	return prometheus.DefaultRegisterer.(*prometheus.Registry)
}
