package product

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	RequestIDKey     = "request_id"
	CorrelationIDKey = "correlation_id"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a unique request ID
		requestID := generateRequestID()
		
		// Add request ID to context
		c.Set(RequestIDKey, requestID)
		
		// Add request ID to response headers
		c.Header("X-Request-ID", requestID)
		
		// Add request ID to logger context
		logger := Logger.WithField("request_id", requestID)
		c.Set("logger", logger)
		
		// Continue to next handler
		c.Next()
	}
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetLoggerFromContext returns logger with request ID from context
func GetLoggerFromContext(c *gin.Context) *logrus.Entry {
	if logger, exists := c.Get("logger"); exists {
		if logEntry, ok := logger.(*logrus.Entry); ok {
			return logEntry
		}
	}
	return Logger
}

// GetRequestIDFromContext returns request ID from context
func GetRequestIDFromContext(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// GetCorrelationIDFromContext returns correlation ID from context
func GetCorrelationIDFromContext(c *gin.Context) string {
	if correlationID, exists := c.Get(CorrelationIDKey); exists {
		if id, ok := correlationID.(string); ok {
			return id
		}
	}
	return ""
}

// HTTPLoggingMiddleware logs HTTP requests and responses
func HTTPLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get logger with request ID
		logger := GetLoggerFromContext(c)
		
		// Log incoming request
		logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Request.URL.RawQuery,
			"user_agent": c.Request.UserAgent(),
			"ip":         c.ClientIP(),
		}).Info("Incoming HTTP request")
		
		// Start timer
		start := time.Now()
		
		// Process request
		c.Next()
		
		// Calculate duration
		duration := time.Since(start)
		
		// Log response
		logger.WithFields(logrus.Fields{
			"status_code": c.Writer.Status(),
			"duration_ms": duration.Milliseconds(),
			"response_size": c.Writer.Size(),
		}).Info("HTTP request completed")
	}
}

// CorrelationIDMiddleware extracts and sets correlation ID from headers
func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get correlation ID from various headers
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = c.GetHeader("X-Request-ID")
		}
		if correlationID == "" {
			correlationID = c.GetHeader("X-Trace-ID")
		}
		
		// If no correlation ID found, generate one
		if correlationID == "" {
			correlationID = generateRequestID()
		}
		
		// Set correlation ID in context
		c.Set(CorrelationIDKey, correlationID)
		
		// Add correlation ID to response headers
		c.Header("X-Correlation-ID", correlationID)
		
		// Update logger with correlation ID
		logger := GetLoggerFromContext(c)
		logger = logger.WithField("correlation_id", correlationID)
		c.Set("logger", logger)
		
		// Continue to next handler
		c.Next()
	}
}
