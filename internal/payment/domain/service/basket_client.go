package service

import (
	"context"
)

// BasketClient defines the interface for basket service communication
type BasketClient interface {
	// Get basket information
	GetBasket(ctx context.Context, userID string) (*BasketInfo, error)
	
	// Clear basket after successful payment
	ClearBasket(ctx context.Context, userID string) error
	
	// Health check
	Ping(ctx context.Context) error
}

// BasketInfo represents basket information from basket service
type BasketInfo struct {
	ID        string        `json:"id"`
	UserID    string        `json:"user_id"`
	Items     []BasketItem  `json:"items"`
	Total     float64       `json:"total"`
	ItemCount int           `json:"item_count"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
	ExpiresAt string        `json:"expires_at"`
}

// BasketItem represents a basket item
type BasketItem struct {
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
	Category  string  `json:"category"`
}
