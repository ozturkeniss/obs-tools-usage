package command

import (
	"obs-tools-usage/internal/product/application/dto"
)

// CreateProductCommand represents a command to create a product
type CreateProductCommand struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	Category    string  `json:"category"`
}

// ToDTO converts command to DTO
func (c *CreateProductCommand) ToDTO() dto.CreateProductRequest {
	return dto.CreateProductRequest{
		Name:        c.Name,
		Description: c.Description,
		Price:       c.Price,
		Stock:       c.Stock,
		Category:    c.Category,
	}
}
