package service

import (
	"context"
)

// ProductClient defines the interface for product service communication
type ProductClient interface {
	// Get product information
	GetProduct(ctx context.Context, productID int) (*ProductInfo, error)
	GetProducts(ctx context.Context, productIDs []int) ([]*ProductInfo, error)
	
	// Health check
	Ping(ctx context.Context) error
}

// ProductInfo represents product information from product service
type ProductInfo struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
	Available   bool    `json:"available"`
}
