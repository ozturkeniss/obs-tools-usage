package http

import (
	"net/http"
	"strconv"
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

// GetPaymentsByStatus handles GET /payments/status/:status
func (h *Handler) GetPaymentsByStatus(c *gin.Context) {
	status := c.Param("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid status",
			Message: "Status is required",
		})
		return
	}

	payments, err := h.queryHandler.HandleGetPaymentsByStatus(query.GetPaymentsByStatusQuery{Status: status})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetPaymentsByDateRange handles GET /payments/date/:start/:end
func (h *Handler) GetPaymentsByDateRange(c *gin.Context) {
	startDate := c.Param("start")
	endDate := c.Param("end")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid date range",
			Message: "Start and end dates are required",
		})
		return
	}

	payments, err := h.queryHandler.HandleGetPaymentsByDateRange(query.GetPaymentsByDateRangeQuery{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetPaymentsByAmountRange handles GET /payments/amount/:min/:max
func (h *Handler) GetPaymentsByAmountRange(c *gin.Context) {
	minAmountStr := c.Param("min")
	maxAmountStr := c.Param("max")

	minAmount, err := strconv.ParseFloat(minAmountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid min amount",
			Message: "Min amount must be a valid number",
		})
		return
	}

	maxAmount, err := strconv.ParseFloat(maxAmountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid max amount",
			Message: "Max amount must be a valid number",
		})
		return
	}

	payments, err := h.queryHandler.HandleGetPaymentsByAmountRange(query.GetPaymentsByAmountRangeQuery{
		MinAmount: minAmount,
		MaxAmount: maxAmount,
	})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetPaymentsByMethod handles GET /payments/method/:method
func (h *Handler) GetPaymentsByMethod(c *gin.Context) {
	method := c.Param("method")
	if method == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid method",
			Message: "Payment method is required",
		})
		return
	}

	payments, err := h.queryHandler.HandleGetPaymentsByMethod(query.GetPaymentsByMethodQuery{Method: method})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetPaymentsByProvider handles GET /payments/provider/:provider
func (h *Handler) GetPaymentsByProvider(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid provider",
			Message: "Payment provider is required",
		})
		return
	}

	payments, err := h.queryHandler.HandleGetPaymentsByProvider(query.GetPaymentsByProviderQuery{Provider: provider})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetPaymentItems handles GET /payments/:id/items
func (h *Handler) GetPaymentItems(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid payment ID",
			Message: "Payment ID is required",
		})
		return
	}

	items, err := h.queryHandler.HandleGetPaymentItems(query.GetPaymentItemsQuery{PaymentID: paymentID})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetPaymentAnalytics handles GET /payments/analytics
func (h *Handler) GetPaymentAnalytics(c *gin.Context) {
	analytics, err := h.queryHandler.HandleGetPaymentAnalytics(query.GetPaymentAnalyticsQuery{})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetPaymentMethods handles GET /payments/methods
func (h *Handler) GetPaymentMethods(c *gin.Context) {
	methods, err := h.queryHandler.HandleGetPaymentMethods(query.GetPaymentMethodsQuery{})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, methods)
}

// GetPaymentProviders handles GET /payments/providers
func (h *Handler) GetPaymentProviders(c *gin.Context) {
	providers, err := h.queryHandler.HandleGetPaymentProviders(query.GetPaymentProvidersQuery{})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, providers)
}

// GetPaymentSummary handles GET /payments/summary
func (h *Handler) GetPaymentSummary(c *gin.Context) {
	summary, err := h.queryHandler.HandleGetPaymentSummary(query.GetPaymentSummaryQuery{})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, summary)
}

// CancelPayment handles POST /payments/:id/cancel
func (h *Handler) CancelPayment(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid payment ID",
			Message: "Payment ID is required",
		})
		return
	}

	cmd := command.CancelPaymentCommand{PaymentID: paymentID}

	payment, err := h.commandHandler.HandleCancelPayment(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payment)
}

// RetryPayment handles POST /payments/:id/retry
func (h *Handler) RetryPayment(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid payment ID",
			Message: "Payment ID is required",
		})
		return
	}

	cmd := command.RetryPaymentCommand{PaymentID: paymentID}

	payment, err := h.commandHandler.HandleRetryPayment(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payment)
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
	r.POST("/payments/:id/cancel", handler.CancelPayment)
	r.POST("/payments/:id/retry", handler.RetryPayment)
	r.GET("/payments/user/:user_id", handler.GetPaymentsByUser)
	r.GET("/payments/stats/:user_id", handler.GetPaymentStats)

	// Query routes
	r.GET("/payments/status/:status", handler.GetPaymentsByStatus)
	r.GET("/payments/date/:start/:end", handler.GetPaymentsByDateRange)
	r.GET("/payments/amount/:min/:max", handler.GetPaymentsByAmountRange)
	r.GET("/payments/method/:method", handler.GetPaymentsByMethod)
	r.GET("/payments/provider/:provider", handler.GetPaymentsByProvider)
	r.GET("/payments/:id/items", handler.GetPaymentItems)
	r.GET("/payments/analytics", handler.GetPaymentAnalytics)
	r.GET("/payments/methods", handler.GetPaymentMethods)
	r.GET("/payments/providers", handler.GetPaymentProviders)
	r.GET("/payments/summary", handler.GetPaymentSummary)

	// Health check
	r.GET("/health", handler.HealthCheck)
}
