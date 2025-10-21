package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"obs-tools-usage/internal/product"
)

func main() {
	// Load configuration
	config := product.LoadConfig()

	// Set Gin mode based on environment
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router with middleware
	r := gin.New()
	
	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())
	r.Use(product.CorrelationIDMiddleware())
	r.Use(product.RequestIDMiddleware())
	r.Use(product.HTTPLoggingMiddleware())
	r.Use(product.PerformanceMiddleware())

	// Initialize product service
	productService := product.NewService()

	// Register health checkers
	registerHealthCheckers()

	// Setup routes
	product.SetupRoutes(r, productService)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Product service starting on port %d", config.GetPort())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// registerHealthCheckers registers all health checkers
func registerHealthCheckers() {
	// Register memory health checker
	product.RegisterHealthChecker("memory", &MemoryHealthChecker{})
	
	// Register goroutine health checker
	product.RegisterHealthChecker("goroutines", &GoroutineHealthChecker{})
	
	// Register service health checker
	product.RegisterHealthChecker("service", &ServiceHealthChecker{})
}

// MemoryHealthChecker checks memory usage
type MemoryHealthChecker struct{}

func (m *MemoryHealthChecker) Check(ctx context.Context) product.CheckResult {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	// Check if memory usage is reasonable (less than 200MB)
	if memStats.Alloc > 200*1024*1024 {
		return product.CheckResult{
			Status:  "unhealthy",
			Message: "Memory usage too high",
		}
	}
	
	return product.CheckResult{
		Status:  "healthy",
		Message: "Memory usage normal",
	}
}

func (m *MemoryHealthChecker) Name() string {
	return "memory"
}

// GoroutineHealthChecker checks goroutine count
type GoroutineHealthChecker struct{}

func (g *GoroutineHealthChecker) Check(ctx context.Context) product.CheckResult {
	goroutineCount := runtime.NumGoroutine()
	
	// Check if goroutine count is reasonable (less than 500)
	if goroutineCount > 500 {
		return product.CheckResult{
			Status:  "unhealthy",
			Message: "Too many goroutines",
		}
	}
	
	return product.CheckResult{
		Status:  "healthy",
		Message: "Goroutine count normal",
	}
}

func (g *GoroutineHealthChecker) Name() string {
	return "goroutines"
}

// ServiceHealthChecker checks basic service functionality
type ServiceHealthChecker struct{}

func (s *ServiceHealthChecker) Check(ctx context.Context) product.CheckResult {
	// This is a simple check - in a real application you might check
	// database connectivity, external service availability, etc.
	
	// Simulate a quick check
	time.Sleep(10 * time.Millisecond)
	
	return product.CheckResult{
		Status:  "healthy",
		Message: "Service is functioning normally",
	}
}

func (s *ServiceHealthChecker) Name() string {
	return "service"
}
