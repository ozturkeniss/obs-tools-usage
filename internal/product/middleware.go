package product

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const RequestIDKey = "request_id"

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
