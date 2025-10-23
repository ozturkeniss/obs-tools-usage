package entity

import (
	"time"
)

// Product represents a product in the system
type Product struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" binding:"required"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price" binding:"required,min=0"`
	Stock       int       `json:"stock" db:"stock" binding:"min=0"`
	Category    string    `json:"category" db:"category"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

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

// ToDTO converts a Product entity to a DTO-compatible struct
func (p *Product) ToDTO() map[string]interface{} {
	return map[string]interface{}{
		"id":          p.ID,
		"name":        p.Name,
		"description": p.Description,
		"price":       p.Price,
		"stock":       p.Stock,
		"category":    p.Category,
		"created_at":  p.CreatedAt,
		"updated_at":  p.UpdatedAt,
	}
}

// ToResponse converts a Product to response format
func (p *Product) ToResponse() interface{} {
	return p.ToDTO()
}

// FromCreateRequest converts CreateProductRequest to Product
func (p *Product) FromCreateRequest(req CreateProductRequest) {
	p.Name = req.Name
	p.Description = req.Description
	p.Price = req.Price
	p.Stock = req.Stock
	p.Category = req.Category
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// FromUpdateRequest converts UpdateProductRequest to Product
func (p *Product) FromUpdateRequest(req UpdateProductRequest) {
	p.Name = req.Name
	p.Description = req.Description
	p.Price = req.Price
	p.Stock = req.Stock
	p.Category = req.Category
	p.UpdatedAt = time.Now()
}

// ProductStats represents product statistics
type ProductStats struct {
	TotalProducts      int64   `json:"total_products"`
	TotalCategories    int64   `json:"total_categories"`
	AveragePrice       float64 `json:"average_price"`
	TotalValue         float64 `json:"total_value"`
	LowStockProducts   int64   `json:"low_stock_products"`
	OutOfStockProducts int64   `json:"out_of_stock_products"`
}

// Category represents a product category
type Category struct {
	Name         string  `json:"name"`
	ProductCount int64   `json:"product_count"`
	AveragePrice float64 `json:"average_price"`
}
