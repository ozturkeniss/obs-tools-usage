package product

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Code      int    `json:"code"`
}

// LogError logs an error with context information
func LogError(logger *logrus.Entry, err error, context map[string]interface{}) {
	// Get caller information
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		context["caller"] = funcName
		context["file"] = file
		context["line"] = line
	}

	// Add error details
	context["error"] = MaskSensitiveData(err.Error())
	context["error_type"] = getErrorType(err)

	// Mask sensitive data in context
	maskedContext := MaskFields(context)

	logger.WithFields(maskedContext).Error("Application error occurred")
}

// getErrorType determines the type of error
func getErrorType(err error) string {
	errorStr := err.Error()
	
	switch {
	case strings.Contains(errorStr, "not found"):
		return "NotFoundError"
	case strings.Contains(errorStr, "already exists"):
		return "ConflictError"
	case strings.Contains(errorStr, "invalid"):
		return "ValidationError"
	case strings.Contains(errorStr, "unauthorized"):
		return "UnauthorizedError"
	case strings.Contains(errorStr, "forbidden"):
		return "ForbiddenError"
	case strings.Contains(errorStr, "timeout"):
		return "TimeoutError"
	default:
		return "InternalError"
	}
}

// HandleError handles and logs errors with appropriate HTTP response
func HandleError(c *gin.Context, err error, statusCode int, message string) {
	logger := GetLoggerFromContext(c)
	requestID := GetRequestIDFromContext(c)
	correlationID := GetCorrelationIDFromContext(c)

	// Prepare error context
	errorContext := map[string]interface{}{
		"request_id":     requestID,
		"correlation_id": correlationID,
		"status_code":    statusCode,
		"endpoint":       c.Request.URL.Path,
		"method":         c.Request.Method,
		"user_agent":     c.Request.UserAgent(),
		"ip":             c.ClientIP(),
	}

	// Log the error
	LogError(logger, err, errorContext)

	// Prepare error response
	errorResponse := ErrorResponse{
		Error:     err.Error(),
		Message:   message,
		RequestID: requestID,
		Code:      statusCode,
	}

	// Send error response
	c.JSON(statusCode, errorResponse)
}

// HandleValidationError handles validation errors
func HandleValidationError(c *gin.Context, err error) {
	HandleError(c, err, http.StatusBadRequest, "Validation failed")
}

// HandleNotFoundError handles not found errors
func HandleNotFoundError(c *gin.Context, err error) {
	HandleError(c, err, http.StatusNotFound, "Resource not found")
}

// HandleInternalError handles internal server errors
func HandleInternalError(c *gin.Context, err error) {
	HandleError(c, err, http.StatusInternalServerError, "Internal server error")
}

// HandleConflictError handles conflict errors
func HandleConflictError(c *gin.Context, err error) {
	HandleError(c, err, http.StatusConflict, "Resource conflict")
}
