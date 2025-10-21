package product

import (
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
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
