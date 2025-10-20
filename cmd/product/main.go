package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"obs-tools-usage/internal/product"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Gin router
	r := gin.Default()

	// Initialize product service
	productService := product.NewService()

	// Setup routes
	product.SetupRoutes(r, productService)

	// Start server
	log.Printf("Product service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
