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

// BasketTotalResponse represents basket total response
type BasketTotalResponse struct {
	UserID    string  `json:"user_id"`
	Total     float64 `json:"total"`
	ItemCount int     `json:"item_count"`
	Currency  string  `json:"currency"`
}

// BasketItemCountResponse represents basket item count response
type BasketItemCountResponse struct {
	UserID    string `json:"user_id"`
	ItemCount int    `json:"item_count"`
	UniqueItems int  `json:"unique_items"`
}

// BasketStatsResponse represents basket statistics response
type BasketStatsResponse struct {
	UserID           string  `json:"user_id"`
	TotalItems       int     `json:"total_items"`
	UniqueItems      int     `json:"unique_items"`
	TotalValue       float64 `json:"total_value"`
	AverageItemPrice float64 `json:"average_item_price"`
	Categories       int     `json:"categories"`
	MostExpensiveItem float64 `json:"most_expensive_item"`
	LeastExpensiveItem float64 `json:"least_expensive_item"`
}

// BasketExpiryResponse represents basket expiry response
type BasketExpiryResponse struct {
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	IsExpired bool      `json:"is_expired"`
	TimeLeft  string    `json:"time_left"`
}

// BasketHistoryResponse represents basket history response
type BasketHistoryResponse struct {
	UserID    string              `json:"user_id"`
	History   []BasketItemResponse `json:"history"`
	TotalOperations int           `json:"total_operations"`
}

// BasketRecommendationsResponse represents basket recommendations response
type BasketRecommendationsResponse struct {
	UserID         string              `json:"user_id"`
	Recommendations []BasketItemResponse `json:"recommendations"`
	Reason         string              `json:"reason"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Service   string `json:"service"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}
