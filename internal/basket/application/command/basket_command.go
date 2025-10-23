package command

import "obs-tools-usage/internal/basket/application/dto"

// CreateBasketCommand represents a command to create a basket
type CreateBasketCommand struct {
	UserID string `json:"user_id" binding:"required"`
}

// ToDTO converts command to DTO
func (c *CreateBasketCommand) ToDTO() dto.CreateBasketRequest {
	return dto.CreateBasketRequest{
		UserID: c.UserID,
	}
}

// AddItemCommand represents a command to add an item to basket
type AddItemCommand struct {
	UserID    string `json:"user_id" binding:"required"`
	ProductID int    `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// ToDTO converts command to DTO
func (c *AddItemCommand) ToDTO() dto.AddItemRequest {
	return dto.AddItemRequest{
		ProductID: c.ProductID,
		Quantity:  c.Quantity,
	}
}

// UpdateItemCommand represents a command to update basket item quantity
type UpdateItemCommand struct {
	UserID    string `json:"user_id" binding:"required"`
	ProductID int    `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=0"`
}

// ToDTO converts command to DTO
func (c *UpdateItemCommand) ToDTO() dto.UpdateItemRequest {
	return dto.UpdateItemRequest{
		ProductID: c.ProductID,
		Quantity:  c.Quantity,
	}
}

// RemoveItemCommand represents a command to remove an item from basket
type RemoveItemCommand struct {
	UserID    string `json:"user_id" binding:"required"`
	ProductID int    `json:"product_id" binding:"required"`
}

// ClearBasketCommand represents a command to clear the basket
type ClearBasketCommand struct {
	UserID string `json:"user_id" binding:"required"`
}
