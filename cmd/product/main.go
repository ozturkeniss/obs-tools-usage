package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"obs-tools-usage/internal/product/application/handler"
	"obs-tools-usage/internal/product/application/usecase"
	"obs-tools-usage/internal/product/infrastructure/config"
	"obs-tools-usage/internal/product/infrastructure/persistence"
	"obs-tools-usage/internal/product/interfaces/grpc"
	httpInterface "obs-tools-usage/internal/product/interfaces/http"
)

//go:generate wire

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	logger := config.GetLogger()
	
	logger.Info("Product service starting...")
	
	// Initialize database
	db, err := persistence.NewDatabase(&cfg.Database)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize database")
	}
	defer db.Close()
	
	// Run database migrations
	if err := db.Migrate(); err != nil {
		logger.WithError(err).Fatal("Failed to run database migrations")
	}
	
	// Seed database with initial data
	if err := db.SeedData(); err != nil {
		logger.WithError(err).Warn("Failed to seed database")
	}
	
	// Initialize repository
	productRepo := persistence.NewProductRepositoryImpl(db.DB)
	
	// Initialize use case
	productUseCase := usecase.NewProductUseCase(productRepo)
	
	// Initialize handlers
	commandHandler := handler.NewCommandHandler(productUseCase)
	queryHandler := handler.NewQueryHandler(productUseCase)
	
	// Initialize HTTP handler
	httpHandler := httpInterface.NewHandler(commandHandler, queryHandler)
	
	// Initialize gRPC server
	grpcServer := grpc.NewGRPCServer(commandHandler, queryHandler, productRepo)
	
	// Initialize Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	
	// Add CORS middleware
	r.Use(corsMiddleware())
	
	// Add Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	
	// Setup HTTP routes
	httpInterface.SetupRoutes(r, commandHandler, queryHandler)
	
	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}
	
	// Start HTTP server in a goroutine
	go func() {
		logger.WithField("port", cfg.Port).Info("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start HTTP server")
		}
	}()
	
	// Start gRPC server in a goroutine
	go func() {
		logger.WithField("port", 50050).Info("Starting gRPC server")
		if err := grpcServer.Start(50050); err != nil {
			logger.WithError(err).Fatal("Failed to start gRPC server")
		}
	}()
	
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")
	
	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Shutdown gRPC server
	grpcServer.Stop()
	
	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("HTTP server forced to shutdown")
	}
	
	logger.Info("Server exited")
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