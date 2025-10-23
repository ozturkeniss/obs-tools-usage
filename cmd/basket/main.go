package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"obs-tools-usage/internal/basket/application/handler"
	"obs-tools-usage/internal/basket/application/usecase"
	"obs-tools-usage/internal/basket/infrastructure/client"
	"obs-tools-usage/internal/basket/infrastructure/config"
	"obs-tools-usage/internal/basket/infrastructure/metrics"
	"obs-tools-usage/internal/basket/infrastructure/persistence"
	httpInterface "obs-tools-usage/internal/basket/interfaces/http"
)

//go:generate wire

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	logger := logrus.New()
	logger.SetLevel(getLogLevel(cfg.LogLevel))
	logger.SetFormatter(getLogFormatter(cfg.LogFormat))
	
	logger.Info("Basket service starting...")
	
	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	defer redisClient.Close()
	
	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.WithError(err).Fatal("Failed to connect to Redis")
	}
	logger.Info("Connected to Redis")
	
	// Initialize product client
	productClient, err := client.NewProductClientImpl(cfg.Product.ServiceURL, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize product client")
	}
	defer productClient.Close()
	logger.Info("Connected to product service")
	
	// Initialize repository
	basketRepo := persistence.NewBasketRepositoryImpl(redisClient, logger)
	
	// Initialize use case
	basketUseCase := usecase.NewBasketUseCase(basketRepo, productClient, logger)
	
	// Initialize handlers
	commandHandler := handler.NewCommandHandler(basketUseCase)
	queryHandler := handler.NewQueryHandler(basketUseCase)
	
	// Initialize Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	
	// Add CORS middleware
	r.Use(corsMiddleware())
	
	// Add metrics middleware
	r.Use(metrics.HTTPLoggingMiddleware())
	
	// Add Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	
	// Setup HTTP routes
	httpInterface.SetupRoutes(r, commandHandler, queryHandler)
	
	// Start cleanup goroutine for expired baskets
	go startCleanupRoutine(basketRepo, logger)
	
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

// startCleanupRoutine starts a background routine to clean up expired baskets
func startCleanupRoutine(repo persistence.BasketRepository, logger *logrus.Logger) {
	ticker := time.NewTicker(1 * time.Hour) // Run every hour
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			logger.Info("Starting cleanup of expired baskets")
			if err := repo.ClearExpiredBaskets(); err != nil {
				logger.WithError(err).Error("Failed to clear expired baskets")
			} else {
				logger.Info("Successfully cleared expired baskets")
			}
		}
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
