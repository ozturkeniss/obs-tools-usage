package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// HandleError handles errors and returns appropriate HTTP responses
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	errorMsg := err.Error()
	statusCode := http.StatusInternalServerError

	// Determine status code based on error message
	switch {
	case strings.Contains(errorMsg, "not found") || strings.Contains(errorMsg, "basket not found"):
		statusCode = http.StatusNotFound
	case strings.Contains(errorMsg, "validation") || strings.Contains(errorMsg, "invalid"):
		statusCode = http.StatusBadRequest
	case strings.Contains(errorMsg, "unauthorized"):
		statusCode = http.StatusUnauthorized
	case strings.Contains(errorMsg, "forbidden"):
		statusCode = http.StatusForbidden
	case strings.Contains(errorMsg, "conflict"):
		statusCode = http.StatusConflict
	case strings.Contains(errorMsg, "not available") || strings.Contains(errorMsg, "insufficient stock"):
		statusCode = http.StatusBadRequest
	case strings.Contains(errorMsg, "expired"):
		statusCode = http.StatusGone
	}

	c.JSON(statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: errorMsg,
	})
}
