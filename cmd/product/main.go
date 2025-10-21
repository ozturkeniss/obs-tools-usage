package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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


	// Setup routes
	product.SetupRoutes(r, productService)
	
	// Add Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: r,
	}

	// Start metrics updater
	go updateMetricsPeriodically(productService)

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

// updateMetricsPeriodically updates business metrics every 30 seconds
func updateMetricsPeriodically(service *product.Service) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			// Update system metrics
			product.UpdateSystemMetrics()
			
			// Update business metrics by getting all products
			products, err := service.GetAllProducts()
			if err == nil {
				product.UpdateBusinessMetrics(products)
			}
		}
	}
}

