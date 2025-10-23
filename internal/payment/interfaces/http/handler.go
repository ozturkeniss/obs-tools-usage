package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"obs-tools-usage/internal/payment/application/command"
	"obs-tools-usage/internal/payment/application/dto"
	"obs-tools-usage/internal/payment/application/handler"
	"obs-tools-usage/internal/payment/application/query"
)

// Handler handles HTTP requests using CQRS pattern
type Handler struct {
	commandHandler *handler.CommandHandler
	queryHandler   *handler.QueryHandler
}

// NewHandler creates a new HTTP handler
func NewHandler(commandHandler *handler.CommandHandler, queryHandler *handler.QueryHandler) *Handler {
	return &Handler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
	}
}

// CreatePayment handles POST /payments
func (h *Handler) CreatePayment(c *gin.Context) {
	var cmd command.CreatePaymentCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	payment, err := h.commandHandler.HandleCreatePayment(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// GetPayment handles GET /payments/:id
func (h *Handler) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid payment ID",
			Message: "Payment ID is required",
		})
		return
	}

	payment, err := h.queryHandler.HandleGetPayment(query.GetPaymentQuery{PaymentID: paymentID})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payment)
}

// UpdatePayment handles PUT /payments/:id
func (h *Handler) UpdatePayment(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid payment ID",
			Message: "Payment ID is required",
		})
		return
	}

	var cmd command.UpdatePaymentCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	cmd.PaymentID = paymentID

	payment, err := h.commandHandler.HandleUpdatePayment(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payment)
}

// ProcessPayment handles POST /payments/:id/process
func (h *Handler) ProcessPayment(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid payment ID",
			Message: "Payment ID is required",
		})
		return
	}

	var cmd command.ProcessPaymentCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	cmd.PaymentID = paymentID

	payment, err := h.commandHandler.HandleProcessPayment(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payment)
}

// RefundPayment handles POST /payments/:id/refund
func (h *Handler) RefundPayment(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid payment ID",
			Message: "Payment ID is required",
		})
		return
	}

	var cmd command.RefundPaymentCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	cmd.PaymentID = paymentID

	payment, err := h.commandHandler.HandleRefundPayment(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payment)
}

// GetPaymentsByUser handles GET /payments/user/:user_id
func (h *Handler) GetPaymentsByUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID is required",
		})
		return
	}

	payments, err := h.queryHandler.HandleGetPaymentsByUser(query.GetPaymentsByUserQuery{UserID: userID})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetPaymentStats handles GET /payments/stats/:user_id
func (h *Handler) GetPaymentStats(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID is required",
		})
		return
	}

	stats, err := h.queryHandler.HandleGetPaymentStats(query.GetPaymentStatsQuery{UserID: userID})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, dto.HealthResponse{
		Service:   "payment-service",
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	})
}

// SetupRoutes sets up all routes
func SetupRoutes(r *gin.Engine, commandHandler *handler.CommandHandler, queryHandler *handler.QueryHandler) {
	handler := NewHandler(commandHandler, queryHandler)

	// Payment routes
	r.POST("/payments", handler.CreatePayment)
	r.GET("/payments/:id", handler.GetPayment)
	r.PUT("/payments/:id", handler.UpdatePayment)
	r.POST("/payments/:id/process", handler.ProcessPayment)
	r.POST("/payments/:id/refund", handler.RefundPayment)
	r.GET("/payments/user/:user_id", handler.GetPaymentsByUser)
	r.GET("/payments/stats/:user_id", handler.GetPaymentStats)

	// Health check
	r.GET("/health", handler.HealthCheck)
}
