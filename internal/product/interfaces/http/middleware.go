package http

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"obs-tools-usage/internal/product/infrastructure/config"
	"obs-tools-usage/internal/product/infrastructure/external"
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
		logger := config.GetLogger().WithField("request_id", requestID)
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
	return config.GetLogger().WithField("source", "middleware")
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
		
		// Track active connections
		// httpConnections.Inc()
		// defer httpConnections.Dec()
		
		// Prepare request fields
		requestFields := map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Request.URL.RawQuery,
			"user_agent": c.Request.UserAgent(),
			"ip":         c.ClientIP(),
		}
		
		// Mask sensitive data in request fields
		// maskedRequestFields := MaskFields(requestFields)
		maskedRequestFields := requestFields
		
		// Log incoming request
		logger.WithFields(maskedRequestFields).Info("Incoming HTTP request")
		
		// Start timer
		start := time.Now()
		
		// Process request
		c.Next()
		
		// Calculate duration
		duration := time.Since(start)
		
		// Calculate request/response sizes
		requestSize := int(c.Request.ContentLength)
		if requestSize < 0 {
			requestSize = 0
		}
		responseSize := c.Writer.Size()
		
		// Prepare response fields
		responseFields := map[string]interface{}{
			"status_code":    c.Writer.Status(),
			"duration_ms":    duration.Milliseconds(),
			"response_size":  responseSize,
		}
		
		// Record Prometheus metrics
		external.RecordHTTPRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
			requestSize,
			responseSize,
		)
		
		// Log response
		logger.WithFields(responseFields).Info("HTTP request completed")
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
