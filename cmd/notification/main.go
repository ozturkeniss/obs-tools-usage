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
	"github.com/sirupsen/logrus"

	"obs-tools-usage/internal/notification/application/handler"
	"obs-tools-usage/internal/notification/application/usecase"
	"obs-tools-usage/internal/notification/infrastructure/config"
	"obs-tools-usage/internal/notification/infrastructure/metrics"
	"obs-tools-usage/internal/notification/infrastructure/persistence"
	httpInterface "obs-tools-usage/internal/notification/interfaces/http"
	"obs-tools-usage/kafka/consumer"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	logger := logrus.New()
	logger.SetLevel(getLogLevel(cfg.LogLevel))
	logger.SetFormatter(getLogFormatter(cfg.LogFormat))
	
	logger.Info("Notification service starting...")
	
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
	
	logger.Info("Connected to database")
	
	// Initialize repository
	notificationRepo := persistence.NewNotificationRepositoryImpl(database.DB, logger)
	
	// Initialize Kafka consumer for events
	kafkaBrokers := []string{"localhost:9092"} // In production, this should come from config
	eventHandler := consumer.NewNotificationEventHandler(logger)
	
	// Start Kafka consumer in background
	go func() {
		consumer, err := consumer.NewNotificationConsumer(kafkaBrokers, "notification-service", eventHandler, logger)
		if err != nil {
			logger.WithError(err).Fatal("Failed to initialize Kafka consumer")
		}
		
		ctx := context.Background()
		if err := consumer.Start(ctx); err != nil {
			logger.WithError(err).Error("Kafka consumer error")
		}
	}()
	logger.Info("Connected to Kafka")
	
	// Initialize use case
	notificationUseCase := usecase.NewNotificationUseCase(notificationRepo, logger)
	
	// Initialize handlers
	commandHandler := handler.NewCommandHandler(notificationUseCase)
	queryHandler := handler.NewQueryHandler(notificationUseCase)
	
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
