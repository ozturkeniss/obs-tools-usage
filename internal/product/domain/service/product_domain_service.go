package service

import (
	"errors"
	"obs-tools-usage/internal/product/domain/entity"
)

// ProductDomainService handles domain-specific business logic
type ProductDomainService struct{}

// NewProductDomainService creates a new domain service
func NewProductDomainService() *ProductDomainService {
	return &ProductDomainService{}
}

// ValidateProduct performs domain validation on product data
func (s *ProductDomainService) ValidateProduct(product entity.Product) error {
	if product.Name == "" {
		return errors.New("product name cannot be empty")
	}
	if product.Price < 0 {
		return errors.New("product price cannot be negative")
	}
	if product.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}
	return nil
}

// IsLowStock checks if a product has low stock
func (s *ProductDomainService) IsLowStock(product entity.Product, threshold int) bool {
	return product.Stock <= threshold
}

// IsHighValue checks if a product is high value
func (s *ProductDomainService) IsHighValue(product entity.Product, threshold float64) bool {
	return product.Price > threshold
}
