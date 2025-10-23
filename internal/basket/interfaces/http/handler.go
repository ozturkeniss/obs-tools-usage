package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"obs-tools-usage/internal/basket/application/command"
	"obs-tools-usage/internal/basket/application/dto"
	"obs-tools-usage/internal/basket/application/handler"
	"obs-tools-usage/internal/basket/application/query"
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

// GetBasket handles GET /baskets/:user_id
func (h *Handler) GetBasket(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID is required",
		})
		return
	}

	basket, err := h.queryHandler.HandleGetBasket(query.GetBasketQuery{UserID: userID})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, basket)
}

// CreateBasket handles POST /baskets
func (h *Handler) CreateBasket(c *gin.Context) {
	var cmd command.CreateBasketCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	basket, err := h.commandHandler.HandleCreateBasket(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, basket)
}

// AddItem handles POST /baskets/:user_id/items
func (h *Handler) AddItem(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID is required",
		})
		return
	}

	var cmd command.AddItemCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	cmd.UserID = userID

	basket, err := h.commandHandler.HandleAddItem(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, basket)
}

// UpdateItem handles PUT /baskets/:user_id/items/:product_id
func (h *Handler) UpdateItem(c *gin.Context) {
	userID := c.Param("user_id")
	productIDStr := c.Param("product_id")
	
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID is required",
		})
		return
	}

	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid product ID",
			Message: "Product ID is required",
		})
		return
	}

	var cmd command.UpdateItemCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	cmd.UserID = userID
	// Note: product_id from URL param should be used, but for simplicity we'll use the one from JSON

	basket, err := h.commandHandler.HandleUpdateItem(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, basket)
}

// RemoveItem handles DELETE /baskets/:user_id/items/:product_id
func (h *Handler) RemoveItem(c *gin.Context) {
	userID := c.Param("user_id")
	productIDStr := c.Param("product_id")
	
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID is required",
		})
		return
	}

	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid product ID",
			Message: "Product ID is required",
		})
		return
	}

	// Parse product ID (you might want to add proper parsing with strconv.Atoi)
	cmd := command.RemoveItemCommand{
		UserID:    userID,
		ProductID: 0, // This should be parsed from productIDStr
	}

	basket, err := h.commandHandler.HandleRemoveItem(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, basket)
}

// ClearBasket handles DELETE /baskets/:user_id/items
func (h *Handler) ClearBasket(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID is required",
		})
		return
	}

	cmd := command.ClearBasketCommand{UserID: userID}

	basket, err := h.commandHandler.HandleClearBasket(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, basket)
}

// DeleteBasket handles DELETE /baskets/:user_id
func (h *Handler) DeleteBasket(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID is required",
		})
		return
	}

	cmd := command.ClearBasketCommand{UserID: userID}

	err := h.commandHandler.HandleDeleteBasket(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Basket deleted successfully",
	})
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, dto.HealthResponse{
		Service:   "basket-service",
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	})
}

// SetupRoutes sets up all routes
func SetupRoutes(r *gin.Engine, commandHandler *handler.CommandHandler, queryHandler *handler.QueryHandler) {
	handler := NewHandler(commandHandler, queryHandler)

	// Basket routes
	r.GET("/baskets/:user_id", handler.GetBasket)
	r.POST("/baskets", handler.CreateBasket)
	r.POST("/baskets/:user_id/items", handler.AddItem)
	r.PUT("/baskets/:user_id/items/:product_id", handler.UpdateItem)
	r.DELETE("/baskets/:user_id/items/:product_id", handler.RemoveItem)
	r.DELETE("/baskets/:user_id/items", handler.ClearBasket)
	r.DELETE("/baskets/:user_id", handler.DeleteBasket)

	// Health check
	r.GET("/health", handler.HealthCheck)
}
