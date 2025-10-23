package repository

import (
	"obs-tools-usage/internal/basket/domain/entity"
)

// BasketRepository defines the interface for basket data access
type BasketRepository interface {
	// Basic CRUD operations
	GetBasket(userID string) (*entity.Basket, error)
	SaveBasket(basket *entity.Basket) error
	DeleteBasket(userID string) error
	
	// Basket operations
	CreateBasket(userID string) (*entity.Basket, error)
	UpdateBasket(basket *entity.Basket) error
	
	// Utility operations
	BasketExists(userID string) (bool, error)
	GetAllBaskets() ([]*entity.Basket, error)
	ClearExpiredBaskets() error
	
	// Health check
	Ping() error
}
