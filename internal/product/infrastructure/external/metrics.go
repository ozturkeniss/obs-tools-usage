package external

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"obs-tools-usage/internal/product/domain/entity"
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

	productsByCategory = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "products_by_category_total",
			Help: "Total number of products by category",
		},
		[]string{"category"},
	)

	productsLowStock = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "products_low_stock_total",
			Help: "Total number of products with low stock",
		},
	)

	productsOutOfStock = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "products_out_of_stock_total",
			Help: "Total number of products out of stock",
		},
	)

	productsHighValue = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "products_high_value_total",
			Help: "Total number of high-value products (>1000)",
		},
	)

	averageProductPrice = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "average_product_price",
			Help: "Average product price",
		},
	)

	totalInventoryValue = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "total_inventory_value",
			Help: "Total inventory value (price * stock)",
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

	// Stock level metrics
	stockLevels = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "product_stock_levels",
			Help:    "Distribution of product stock levels",
			Buckets: []float64{0, 1, 5, 10, 25, 50, 100, 250, 500, 1000},
		},
		[]string{"category"},
	)

	// Price range metrics
	priceRanges = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "product_price_ranges",
			Help:    "Distribution of product prices",
			Buckets: []float64{0, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
		},
		[]string{"category"},
	)

	// System metrics
	memoryAllocBytes = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_alloc_bytes",
			Help: "Current memory allocation in bytes",
		},
	)

	memorySysBytes = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_sys_bytes",
			Help: "Total memory obtained from OS in bytes",
		},
	)

	memoryHeapBytes = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_heap_bytes",
			Help: "Heap memory size in bytes",
		},
	)

	memoryStackBytes = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_stack_bytes",
			Help: "Stack memory size in bytes",
		},
	)

	gcDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "gc_duration_seconds",
			Help:    "GC duration in seconds",
			Buckets: prometheus.ExponentialBuckets(0.0001, 2, 15),
		},
	)

	gcCount = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "gc_count_total",
			Help: "Total number of GC cycles",
		},
	)

	goroutinesTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines_total",
			Help: "Current number of goroutines",
		},
	)

	cgoCalls = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cgo_calls_total",
			Help: "Total number of CGO calls",
		},
	)

	threadsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "threads_total",
			Help: "Current number of OS threads",
		},
	)

	// CPU metrics (approximated)
	cpuUsagePercent = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "CPU usage percentage (approximated)",
		},
	)

	// Application metrics
	httpConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_connections_active",
			Help: "Number of active HTTP connections",
		},
	)

	requestQueueSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "request_queue_size",
			Help: "Current request queue size",
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
		_, _, _, startNumGC, startGCForced := GetSystemMetrics()
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

// UpdateBusinessMetrics updates all business metrics
func UpdateBusinessMetrics(products []entity.Product) {
	// Reset category counters
	productsByCategory.Reset()
	
	// Counters
	totalCount := len(products)
	lowStockCount := 0
	outOfStockCount := 0
	highValueCount := 0
	totalPrice := 0.0
	totalInventoryValueCalc := 0.0
	
	// Category counters
	categoryCounts := make(map[string]int)
	
	for _, product := range products {
		// Count by category
		categoryCounts[product.Category]++
		
		// Stock level checks
		if product.Stock == 0 {
			outOfStockCount++
		} else if product.Stock < 10 {
			lowStockCount++
		}
		
		// High value products
		if product.Price > 1000 {
			highValueCount++
		}
		
		// Price calculations
		totalPrice += product.Price
		totalInventoryValueCalc += product.Price * float64(product.Stock)
		
		// Record stock level distribution
		stockLevels.WithLabelValues(product.Category).Observe(float64(product.Stock))
		
		// Record price distribution
		priceRanges.WithLabelValues(product.Category).Observe(product.Price)
	}
	
	// Update gauges
	productsTotal.Set(float64(totalCount))
	productsLowStock.Set(float64(lowStockCount))
	productsOutOfStock.Set(float64(outOfStockCount))
	productsHighValue.Set(float64(highValueCount))
	
	// Calculate and set average price
	if totalCount > 0 {
		averageProductPrice.Set(totalPrice / float64(totalCount))
	} else {
		averageProductPrice.Set(0)
	}
	
	// Set total inventory value
	totalInventoryValue.Set(totalInventoryValueCalc)
	
	// Update category counters
	for category, count := range categoryCounts {
		productsByCategory.WithLabelValues(category).Set(float64(count))
	}
}

// RecordProductStockLevel records individual product stock level
func RecordProductStockLevel(product entity.Product) {
	stockLevels.WithLabelValues(product.Category).Observe(float64(product.Stock))
	priceRanges.WithLabelValues(product.Category).Observe(product.Price)
}

// RecordLowStockAlert records low stock alert
func RecordLowStockAlert(product entity.Product) {
	// This could be used for alerting when stock is low
	// For now, we just record the metric
	productsLowStock.Inc()
}

// RecordOutOfStockAlert records out of stock alert
func RecordOutOfStockAlert(product entity.Product) {
	// This could be used for alerting when product is out of stock
	productsOutOfStock.Inc()
}

// UpdateSystemMetrics updates system-level metrics
func UpdateSystemMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	// Memory metrics
	memoryAllocBytes.Set(float64(memStats.Alloc))
	memorySysBytes.Set(float64(memStats.Sys))
	memoryHeapBytes.Set(float64(memStats.HeapAlloc))
	memoryStackBytes.Set(float64(memStats.StackInuse))
	
	// GC metrics
	gcCount.Add(float64(memStats.NumGC - lastGCCount))
	lastGCCount = memStats.NumGC
	
	// Record GC duration if available
	if memStats.PauseTotalNs > 0 {
		avgGCPause := float64(memStats.PauseTotalNs) / float64(memStats.NumGC) / 1e9
		gcDuration.Observe(avgGCPause)
	}
	
	// Goroutine and thread metrics
	goroutinesTotal.Set(float64(runtime.NumGoroutine()))
	
	// CGO calls (not available in runtime.MemStats)
	// cgoCalls.Add(float64(memStats.CGOCall - lastCGOCalls))
	// lastCGOCalls = memStats.CGOCall
	
	// Approximate CPU usage (this is a simple approximation)
	updateCPUUsage()
	
	// Application metrics
	updateApplicationMetrics()
}

// Global variables to track changes
var (
	lastGCCount  uint32
	lastCGOCalls uint64
	lastCPUTime  time.Time
	lastGCPause  time.Duration
)

// updateCPUUsage approximates CPU usage
func updateCPUUsage() {
	now := time.Now()
	if !lastCPUTime.IsZero() {
		// Simple approximation based on GC activity and goroutines
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		
		// Approximate CPU usage based on GC pressure and goroutine count
		gcPressure := float64(memStats.NumGC) / float64(now.Sub(lastCPUTime).Seconds())
		goroutinePressure := float64(runtime.NumGoroutine()) / 100.0
		
		// Combine factors for approximation (this is not precise)
		cpuUsage := (gcPressure * 10) + (goroutinePressure * 5)
		if cpuUsage > 100 {
			cpuUsage = 100
		}
		
		cpuUsagePercent.Set(cpuUsage)
	}
	lastCPUTime = now
}

// updateApplicationMetrics updates application-specific metrics
func updateApplicationMetrics() {
	// These would be updated based on actual application state
	// For now, we'll use simple approximations
	
	// Approximate active connections based on goroutines
	goroutineCount := runtime.NumGoroutine()
	estimatedConnections := float64(goroutineCount) * 0.1 // Rough estimate
	httpConnections.Set(estimatedConnections)
	
	// Request queue size (simplified)
	requestQueueSize.Set(0) // In a real app, this would track actual queue size
}

// GetPrometheusMetrics returns the Prometheus registry
func GetPrometheusMetrics() *prometheus.Registry {
	return prometheus.DefaultRegisterer.(*prometheus.Registry)
}

// GetLoggerFromContext returns logger from gin context
func GetLoggerFromContext(c *gin.Context) *logrus.Entry {
	// For now, return a basic logger entry
	// In a real implementation, you'd get the logger from context
	logger := logrus.New()
	return logger.WithFields(logrus.Fields{
		"request_id": c.GetString("request_id"),
		"user_id":    c.GetString("user_id"),
		"ip":         c.ClientIP(),
	})
}
