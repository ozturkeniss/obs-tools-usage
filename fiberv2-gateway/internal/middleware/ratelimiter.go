package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"fiberv2-gateway/internal/ratelimiter"
)

// RateLimitMiddleware creates a rate limiting middleware using Redis sliding window
func RateLimitMiddleware(rateLimiter *ratelimiter.SlidingWindowRateLimiter, config ratelimiter.RateLimitConfig, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get client identifier (IP address or user ID)
		identifier := getClientIdentifier(c)
		
		// Check rate limit
		result, err := rateLimiter.CheckRateLimitWithSlidingWindow(c.Context(), config, identifier)
		if err != nil {
			logger.WithError(err).Error("Failed to check rate limit")
			// On error, allow the request but log the error
			return c.Next()
		}
		
		// Set rate limit headers
		setRateLimitHeaders(c, result, config)
		
		// If rate limit exceeded, return 429 Too Many Requests
		if !result.Allowed {
			logger.WithFields(logrus.Fields{
				"identifier":   identifier,
				"remaining":    result.Remaining,
				"retry_after":  result.RetryAfter,
				"reset_time":   result.ResetTime,
			}).Warn("Rate limit exceeded")
			
			c.Status(429).JSON(fiber.Map{
				"error":       "Rate limit exceeded",
				"retry_after": result.RetryAfter.Seconds(),
				"reset_time":  result.ResetTime,
			})
			return nil
		}
		
		logger.WithFields(logrus.Fields{
			"identifier": identifier,
			"remaining":  result.Remaining,
			"reset_time": result.ResetTime,
		}).Debug("Rate limit check passed")
		
		return c.Next()
	}
}

// AdaptiveRateLimitMiddleware creates an adaptive rate limiting middleware
func AdaptiveRateLimitMiddleware(rateLimiter *ratelimiter.SlidingWindowRateLimiter, configs map[string]ratelimiter.RateLimitConfig, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get client identifier
		identifier := getClientIdentifier(c)
		
		// Determine which rate limit config to use based on request characteristics
		config := selectRateLimitConfig(c, configs)
		
		// Check rate limit
		result, err := rateLimiter.CheckRateLimitWithSlidingWindow(c.Context(), config, identifier)
		if err != nil {
			logger.WithError(err).Error("Failed to check adaptive rate limit")
			return c.Next()
		}
		
		// Set rate limit headers
		setRateLimitHeaders(c, result, config)
		
		// If rate limit exceeded
		if !result.Allowed {
			logger.WithFields(logrus.Fields{
				"identifier":   identifier,
				"path":         c.Path(),
				"method":       c.Method(),
				"remaining":    result.Remaining,
				"retry_after":  result.RetryAfter,
				"reset_time":   result.ResetTime,
			}).Warn("Adaptive rate limit exceeded")
			
			c.Status(429).JSON(fiber.Map{
				"error":       "Rate limit exceeded",
				"retry_after": result.RetryAfter.Seconds(),
				"reset_time":  result.ResetTime,
			})
			return nil
		}
		
		return c.Next()
	}
}

// PerServiceRateLimitMiddleware creates rate limiting middleware per service
func PerServiceRateLimitMiddleware(rateLimiter *ratelimiter.SlidingWindowRateLimiter, serviceConfigs map[string]ratelimiter.RateLimitConfig, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get client identifier
		identifier := getClientIdentifier(c)
		
		// Determine service from path
		service := getServiceFromPath(c.Path())
		
		// Get rate limit config for this service
		config, exists := serviceConfigs[service]
		if !exists {
			// Use default config if service not found
			config = ratelimiter.RateLimitConfig{
				WindowSize:  time.Minute,
				MaxRequests: 100,
				KeyPrefix:   "gateway:rate_limit",
			}
		}
		
		// Create service-specific identifier
		serviceIdentifier := identifier + ":" + service
		
		// Check rate limit
		result, err := rateLimiter.CheckRateLimitWithSlidingWindow(c.Context(), config, serviceIdentifier)
		if err != nil {
			logger.WithError(err).Error("Failed to check per-service rate limit")
			return c.Next()
		}
		
		// Set rate limit headers
		setRateLimitHeaders(c, result, config)
		
		// If rate limit exceeded
		if !result.Allowed {
			logger.WithFields(logrus.Fields{
				"identifier":   identifier,
				"service":      service,
				"remaining":    result.Remaining,
				"retry_after":  result.RetryAfter,
				"reset_time":   result.ResetTime,
			}).Warn("Per-service rate limit exceeded")
			
			c.Status(429).JSON(fiber.Map{
				"error":       "Rate limit exceeded for service: " + service,
				"retry_after": result.RetryAfter.Seconds(),
				"reset_time":  result.ResetTime,
			})
			return nil
		}
		
		return c.Next()
	}
}

// getClientIdentifier extracts client identifier from request
func getClientIdentifier(c *fiber.Ctx) string {
	// Try to get user ID from header or JWT token
	userID := c.Get("X-User-ID")
	if userID != "" {
		return "user:" + userID
	}
	
	// Fallback to IP address
	ip := c.IP()
	return "ip:" + ip
}

// selectRateLimitConfig selects appropriate rate limit config based on request characteristics
func selectRateLimitConfig(c *fiber.Ctx, configs map[string]ratelimiter.RateLimitConfig) ratelimiter.RateLimitConfig {
	// Check for admin endpoints
	if c.Path()[:6] == "/admin" {
		if config, exists := configs["admin"]; exists {
			return config
		}
	}
	
	// Check for API endpoints
	if c.Path()[:4] == "/api" {
		if config, exists := configs["api"]; exists {
			return config
		}
	}
	
	// Check for health endpoints
	if c.Path()[:7] == "/health" {
		if config, exists := configs["health"]; exists {
			return config
		}
	}
	
	// Default config
	if config, exists := configs["default"]; exists {
		return config
	}
	
	// Fallback config
	return ratelimiter.RateLimitConfig{
		WindowSize:  time.Minute,
		MaxRequests: 100,
		KeyPrefix:   "gateway:rate_limit",
	}
}

// getServiceFromPath extracts service name from path
func getServiceFromPath(path string) string {
	// Extract service from /api/{service}/* pattern
	if len(path) > 5 && path[:5] == "/api/" {
		// Find next slash
		for i := 5; i < len(path); i++ {
			if path[i] == '/' {
				return path[5:i]
			}
		}
		// If no slash found, return the whole service part
		if len(path) > 5 {
			return path[5:]
		}
	}
	
	return "unknown"
}

// setRateLimitHeaders sets rate limit headers in response
func setRateLimitHeaders(c *fiber.Ctx, result *ratelimiter.RateLimitResult, config ratelimiter.RateLimitConfig) {
	// Standard rate limit headers
	c.Set("X-RateLimit-Limit", strconv.Itoa(config.MaxRequests))
	c.Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
	c.Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetTime.Unix(), 10))
	
	// Additional headers
	c.Set("X-RateLimit-Window", config.WindowSize.String())
	
	// If rate limit exceeded, set Retry-After header
	if !result.Allowed {
		c.Set("Retry-After", strconv.FormatInt(int64(result.RetryAfter.Seconds()), 10))
	}
}

// RateLimitStatusMiddleware provides rate limit status endpoint
func RateLimitStatusMiddleware(rateLimiter *ratelimiter.SlidingWindowRateLimiter, config ratelimiter.RateLimitConfig, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Query("identifier")
		if identifier == "" {
			identifier = getClientIdentifier(c)
		}
		
		// Get rate limit status
		result, err := rateLimiter.GetRateLimitStatus(c.Context(), config, identifier)
		if err != nil {
			logger.WithError(err).Error("Failed to get rate limit status")
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to get rate limit status",
			})
		}
		
		// Get detailed stats
		stats, err := rateLimiter.GetRateLimitStats(c.Context(), config, identifier)
		if err != nil {
			logger.WithError(err).Error("Failed to get rate limit stats")
		}
		
		return c.JSON(fiber.Map{
			"identifier":    identifier,
			"allowed":       result.Allowed,
			"remaining":     result.Remaining,
			"reset_time":    result.ResetTime,
			"retry_after":   result.RetryAfter.Seconds(),
			"stats":         stats,
		})
	}
}
