package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"

	"fiberv2-gateway/internal/config"
	"fiberv2-gateway/internal/gateway"
	"fiberv2-gateway/internal/health"
	"fiberv2-gateway/internal/logging"
	"fiberv2-gateway/internal/metrics"
	"fiberv2-gateway/internal/ratelimiter"
	"fiberv2-gateway/internal/redis"
	"fiberv2-gateway/internal/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	
	// Setup logger
	logger := logging.SetupLogger(cfg.LogLevel, cfg.LogFormat)
	
	// Setup Redis client
	redisClient := redis.NewClient(redis.Config{
		Host:         cfg.Redis.Host,
		Port:         cfg.Redis.Port,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		MaxRetries:   cfg.Redis.MaxRetries,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
	}, logger)
	
	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx); err != nil {
		logger.WithError(err).Fatal("Failed to connect to Redis")
	}
	defer redisClient.Close()
	
	// Setup rate limiter
	rateLimiter := ratelimiter.NewSlidingWindowRateLimiter(redisClient.GetClient(), logger)
	
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "FiberV2 Gateway",
		ServerHeader: "FiberV2-Gateway",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.WithError(err).Error("Request error")
			return c.Status(500).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		},
	})

	// Setup middleware
	setupMiddleware(app, logger, rateLimiter, cfg)

	// Setup metrics
	metrics.SetupMetrics(app)

	// Setup health checks
	health.SetupHealthRoutes(app)

	// Setup gateway routes
	gateway.SetupRoutes(app, cfg, logger)

	// Start server
	startServer(app, cfg, logger)
}

func setupMiddleware(app *fiber.App, logger *logrus.Logger) {
	// Recovery middleware
	app.Use(recover.New())

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))

	// Custom request ID middleware
	app.Use(func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("req_%d", time.Now().UnixNano())
		}
		c.Set("X-Request-ID", requestID)
		c.Locals("requestID", requestID)
		return c.Next()
	})
}

func startServer(app *fiber.App, cfg *config.Config, logger *logrus.Logger) {
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server in goroutine
	go func() {
		logger.WithFields(logrus.Fields{
			"port": cfg.Port,
			"env":  cfg.Environment,
		}).Info("Starting FiberV2 Gateway server")

		if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	} else {
		logger.Info("Server shutdown completed")
	}
}
