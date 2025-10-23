package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"obs-tools-usage/internal/payment/application/handler"
	"obs-tools-usage/internal/payment/application/usecase"
	"obs-tools-usage/internal/payment/infrastructure/client"
	"obs-tools-usage/internal/payment/infrastructure/config"
	"obs-tools-usage/internal/payment/infrastructure/persistence"
	httpInterface "obs-tools-usage/internal/payment/interfaces/http"
	grpcInterface "obs-tools-usage/internal/payment/interfaces/grpc"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	logger := logrus.New()
	logger.SetLevel(getLogLevel(cfg.LogLevel))
	logger.SetFormatter(getLogFormatter(cfg.LogFormat))
	
	logger.Info("Payment service starting...")
	
	// Initialize database
	database, err := persistence.NewDatabase(cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer database.Close()
	
	// Run migrations
	if err := database.Migrate(); err != nil {
		logger.WithError(err).Fatal("Failed to run migrations")
	}
	
	// Seed data (only in development)
	if cfg.IsDevelopment() {
		if err := database.SeedData(); err != nil {
			logger.WithError(err).Warn("Failed to seed data")
		}
	}
	
	logger.Info("Connected to MariaDB database")
	
	// Initialize gRPC clients
	basketClient, err := client.NewBasketClientImpl(cfg.Basket.ServiceURL, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize basket client")
	}
	defer basketClient.Close()
	logger.Info("Connected to basket service")
	
	productClient, err := client.NewProductClientImpl(cfg.Product.ServiceURL, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize product client")
	}
	defer productClient.Close()
	logger.Info("Connected to product service")
	
	// Initialize repository
	paymentRepo := persistence.NewPaymentRepositoryImpl(database.DB, logger)
	
	// Initialize use case
	paymentUseCase := usecase.NewPaymentUseCase(paymentRepo, basketClient, productClient, logger)
	
	// Initialize handlers
	commandHandler := handler.NewCommandHandler(paymentUseCase)
	queryHandler := handler.NewQueryHandler(paymentUseCase)
	
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

	// Create gRPC server
	grpcPort := "50052" // Payment service gRPC port
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logger.WithError(err).Fatal("Failed to listen on gRPC port")
	}

	grpcServer := grpc.NewServer()
	grpcInterface.RegisterServer(grpcServer, commandHandler, queryHandler, logger)

	// Start gRPC server in a goroutine
	go func() {
		logger.WithField("port", grpcPort).Info("Starting gRPC server")
		if err := grpcServer.Serve(lis); err != nil {
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
	
	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("HTTP server forced to shutdown")
	}

	// Shutdown gRPC server
	logger.Info("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	
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

// getLogLevel converts string to logrus level
func getLogLevel(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

// getLogFormatter returns the appropriate log formatter
func getLogFormatter(format string) logrus.Formatter {
	switch format {
	case "json":
		return &logrus.JSONFormatter{}
	default:
		return &logrus.TextFormatter{
			FullTimestamp: true,
		}
	}
}
