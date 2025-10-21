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
	"github.com/joho/godotenv"
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

	// Initialize database
	db, err := product.NewDatabase(&config.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Seed initial data
	if err := db.SeedData(); err != nil {
		log.Fatal("Failed to seed data:", err)
	}

	// Initialize repository with database
	productRepo := product.NewRepository(db.DB)

	// Initialize product service
	productService := product.NewService(productRepo)

	// Initialize gRPC server
	grpcServer := product.NewGRPCServer(productService, productRepo, product.GetLogger())

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

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Product HTTP service starting on port %d", config.GetPort())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start HTTP server:", err)
		}
	}()

	// Start gRPC server in a goroutine
	go func() {
		grpcPort := config.GetPort() + 1 // gRPC port is HTTP port + 1
		log.Printf("Product gRPC service starting on port %d", grpcPort)
		if err := grpcServer.Start(grpcPort); err != nil {
			log.Fatal("Failed to start gRPC server:", err)
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

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("HTTP server forced to shutdown:", err)
	}

	// Shutdown gRPC server
	grpcServer.Stop()

	log.Println("Servers exited")
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

// updateMetricsPeriodically updates metrics periodically
func updateMetricsPeriodically(service *product.Service) {
	// System metrics every 10 seconds
	systemTicker := time.NewTicker(10 * time.Second)
	defer systemTicker.Stop()
	
	// Business metrics every 30 seconds
	businessTicker := time.NewTicker(30 * time.Second)
	defer businessTicker.Stop()
	
	for {
		select {
		case <-systemTicker.C:
			// Update system metrics more frequently
			product.UpdateSystemMetrics()
			
		case <-businessTicker.C:
			// Update business metrics less frequently
			products, err := service.GetAllProducts()
			if err == nil {
				product.UpdateBusinessMetrics(products)
			}
		}
	}
}

