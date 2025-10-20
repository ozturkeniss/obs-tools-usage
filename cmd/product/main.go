package main

import (
	"log"
	"os"

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

	// Initialize product service
	productService := product.NewService()

	// Setup routes
	product.SetupRoutes(r, productService)

	// Start server
	log.Printf("Product service starting on port %d", config.GetPort())
	if err := r.Run(":" + config.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
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
