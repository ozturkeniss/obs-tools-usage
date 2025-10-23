package dto

import "time"

// CreateBasketRequest represents the request payload for creating a basket
type CreateBasketRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// AddItemRequest represents the request payload for adding an item to basket
type AddItemRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

// UpdateItemRequest represents the request payload for updating basket item quantity
type UpdateItemRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=0"`
}

// BasketItemResponse represents a basket item in response
type BasketItemResponse struct {
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
	Category  string  `json:"category"`
}

// BasketResponse represents the response payload for basket operations
type BasketResponse struct {
	ID        string              `json:"id"`
	UserID    string              `json:"user_id"`
	Items     []BasketItemResponse `json:"items"`
	Total     float64             `json:"total"`
	ItemCount int                 `json:"item_count"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	ExpiresAt time.Time           `json:"expires_at"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Service   string `json:"service"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}
