package product

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HealthStatus represents the health status of the service
type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	Uptime    string                 `json:"uptime"`
	Checks    map[string]CheckResult `json:"checks"`
}

// CheckResult represents the result of a health check
type CheckResult struct {
	Status      string        `json:"status"`
	Message     string        `json:"message,omitempty"`
	Duration    time.Duration `json:"duration_ms"`
	LastChecked time.Time     `json:"last_checked"`
}

// HealthChecker interface for health checks
type HealthChecker interface {
	Check(ctx context.Context) CheckResult
	Name() string
}

// ServiceStartTime tracks when the service started
var ServiceStartTime = time.Now()

// HealthCheckers holds all registered health checkers
var HealthCheckers = make(map[string]HealthChecker)

// RegisterHealthChecker registers a health checker
func RegisterHealthChecker(name string, checker HealthChecker) {
	HealthCheckers[name] = checker
}

// LivenessProbe checks if the service is alive
func LivenessProbe(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Basic liveness checks
	checks := make(map[string]CheckResult)
	
	// Check if service is responsive
	checks["service_responsive"] = CheckResult{
		Status:      "healthy",
		Message:     "Service is responding",
		Duration:    0,
		LastChecked: time.Now(),
	}
	
	// Check memory usage
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)
	
	// If memory usage is too high, mark as unhealthy
	if memStats.Alloc > 100*1024*1024 { // 100MB threshold
		checks["memory_usage"] = CheckResult{
			Status:      "unhealthy",
			Message:     "Memory usage too high",
			Duration:    0,
			LastChecked: time.Now(),
		}
	} else {
		checks["memory_usage"] = CheckResult{
			Status:      "healthy",
			Message:     "Memory usage normal",
			Duration:    0,
			LastChecked: time.Now(),
		}
	}
	
	// Check goroutine count
	goroutineCount := runtime.NumGoroutine()
	if goroutineCount > 1000 { // 1000 goroutine threshold
		checks["goroutine_count"] = CheckResult{
			Status:      "unhealthy",
			Message:     "Too many goroutines",
			Duration:    0,
			LastChecked: time.Now(),
		}
	} else {
		checks["goroutine_count"] = CheckResult{
			Status:      "healthy",
			Message:     "Goroutine count normal",
			Duration:    0,
			LastChecked: time.Now(),
		}
	}
	
	// Run registered health checkers
	for name, checker := range HealthCheckers {
		start := time.Now()
		result := checker.Check(ctx)
		result.Duration = time.Since(start)
		result.LastChecked = time.Now()
		checks[name] = result
	}
	
	// Determine overall status
	overallStatus := "healthy"
	for _, check := range checks {
		if check.Status == "unhealthy" {
			overallStatus = "unhealthy"
			break
		}
	}
	
	// Calculate uptime
	uptime := time.Since(ServiceStartTime)
	
	healthStatus := HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   "1.0.0", // This could be read from build info
		Uptime:    uptime.String(),
		Checks:    checks,
	}
	
	// Set appropriate HTTP status code
	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}
	
	// Log health check
	logger := GetLoggerFromContext(c)
	logger.WithFields(logrus.Fields{
		"status": overallStatus,
		"checks": len(checks),
		"uptime": uptime.String(),
	}).Info("Liveness probe executed")
	
	c.JSON(statusCode, healthStatus)
}

// ReadinessProbe checks if the service is ready to accept traffic
func ReadinessProbe(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checks := make(map[string]CheckResult)
	
	// Check if service is ready to accept requests
	checks["service_ready"] = CheckResult{
		Status:      "healthy",
		Message:     "Service is ready",
		Duration:    0,
		LastChecked: time.Now(),
	}
	
	// Check if all dependencies are available
	// This is where you would check database, external services, etc.
	checks["dependencies"] = CheckResult{
		Status:      "healthy",
		Message:     "All dependencies available",
		Duration:    0,
		LastChecked: time.Now(),
	}
	
	// Run registered health checkers
	for name, checker := range HealthCheckers {
		start := time.Now()
		result := checker.Check(ctx)
		result.Duration = time.Since(start)
		result.LastChecked = time.Now()
		checks[name] = result
	}
	
	// Determine overall status
	overallStatus := "healthy"
	for _, check := range checks {
		if check.Status == "unhealthy" {
			overallStatus = "unhealthy"
			break
		}
	}
	
	uptime := time.Since(ServiceStartTime)
	
	healthStatus := HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    uptime.String(),
		Checks:    checks,
	}
	
	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}
	
	// Log readiness check
	logger := GetLoggerFromContext(c)
	logger.WithFields(logrus.Fields{
		"status": overallStatus,
		"checks": len(checks),
		"uptime": uptime.String(),
	}).Info("Readiness probe executed")
	
	c.JSON(statusCode, healthStatus)
}

// SimpleHealthCheck provides a simple health check endpoint
func SimpleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "product-service",
	})
}
