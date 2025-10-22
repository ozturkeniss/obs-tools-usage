package command

import (
	"obs-tools-usage/internal/product/application/dto"
)

// UpdateProductCommand represents a command to update a product
type UpdateProductCommand struct {
	ID          int     `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	Category    string  `json:"category"`
}

// ToDTO converts command to DTO
func (c *UpdateProductCommand) ToDTO() dto.UpdateProductRequest {
	return dto.UpdateProductRequest{
		Name:        c.Name,
		Description: c.Description,
		Price:       c.Price,
		Stock:       c.Stock,
		Category:    c.Category,
	}
}
