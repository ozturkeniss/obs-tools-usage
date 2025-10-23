package dto

import "time"

// CreateProductRequest represents the request payload for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	Category    string  `json:"category"`
}

// UpdateProductRequest represents the request payload for updating a product
type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	Category    string  `json:"category"`
}

// ProductResponse represents the response payload for product operations
type ProductResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductsResponse represents the response payload for multiple products
type ProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Count    int               `json:"count"`
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

// ProductStatsResponse represents product statistics response
type ProductStatsResponse struct {
	TotalProducts      int64   `json:"total_products"`
	TotalCategories    int64   `json:"total_categories"`
	AveragePrice       float64 `json:"average_price"`
	TotalValue         float64 `json:"total_value"`
	LowStockProducts   int64   `json:"low_stock_products"`
	OutOfStockProducts int64   `json:"out_of_stock_products"`
}

// CategoryResponse represents a category response
type CategoryResponse struct {
	Name         string  `json:"name"`
	ProductCount int64   `json:"product_count"`
	AveragePrice float64 `json:"average_price"`
}

// CategoriesResponse represents categories response
type CategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Count      int                `json:"count"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Service   string `json:"service"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}
