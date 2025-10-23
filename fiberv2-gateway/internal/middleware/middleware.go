package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.uber.org/ratelimit"
)

// RateLimiter holds rate limiting configuration
type RateLimiter struct {
	limiter ratelimit.Limiter
	logger  *logrus.Logger
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requests int, window time.Duration, logger *logrus.Logger) *RateLimiter {
	return &RateLimiter{
		limiter: ratelimit.New(requests, ratelimit.Per(window)),
		logger:  logger,
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func (rl *RateLimiter) RateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Take a token from the rate limiter
		rl.limiter.Take()

		// Continue to next middleware
		return c.Next()
	}
}

// RequestLoggerMiddleware creates a request logging middleware
func RequestLoggerMiddleware(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Continue to next middleware
		err := c.Next()

		// Log the request
		duration := time.Since(start)
		logger.WithFields(logrus.Fields{
			"method":   c.Method(),
			"path":     c.Path(),
			"status":   c.Response().StatusCode(),
			"duration": duration,
			"ip":       c.IP(),
			"user_agent": c.Get("User-Agent"),
		}).Info("Request processed")

		return err
	}
}

// CORSMiddleware creates a CORS middleware
func CORSMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set CORS headers
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")

		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}

		return c.Next()
	}
}

// SecurityMiddleware creates a security middleware
func SecurityMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set security headers
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		return c.Next()
	}
}

// TimeoutMiddleware creates a timeout middleware
func TimeoutMiddleware(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()

		// Set the context
		c.SetUserContext(ctx)

		// Continue to next middleware
		return c.Next()
	}
}
