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
	Products []Product `json:"products"`
	Count    int       `json:"count"`
}


// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ToResponse converts a Product to ProductResponse
func (p *Product) ToResponse() ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
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

// ToResponse converts a Product to ProductResponse
func (p *Product) ToResponse() interface{} {
	return struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Price       float64   `json:"price"`
		Stock       int       `json:"stock"`
		Category    string    `json:"category"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
