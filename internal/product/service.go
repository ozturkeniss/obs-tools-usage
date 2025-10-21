package product

import (
	"errors"
	"fmt"
)

type Service struct {
	repository *Repository
}

func NewService() *Service {
	return &Service{
		repository: NewRepository(),
	}
}

// GetAllProducts returns all products
func (s *Service) GetAllProducts() ([]Product, error) {
	return s.repository.GetAllProducts()
}

// GetProductByID returns a product by ID
func (s *Service) GetProductByID(id int) (*Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid product ID")
	}
	
	product, err := s.repository.GetProductByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	
	return product, nil
}

// CreateProduct creates a new product
func (s *Service) CreateProduct(product Product) (*Product, error) {
	// Validate product data
	if err := s.validateProduct(product); err != nil {
		return nil, err
	}
	
	// Generate new ID
	product.ID = s.repository.GetNextID()
	
	// Create product
	createdProduct, err := s.repository.CreateProduct(product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	
	return createdProduct, nil
}

// UpdateProduct updates an existing product
func (s *Service) UpdateProduct(product Product) (*Product, error) {
	// Validate product data
	if err := s.validateProduct(product); err != nil {
		return nil, err
	}
	
	// Check if product exists
	_, err := s.repository.GetProductByID(product.ID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	
	// Update product
	updatedProduct, err := s.repository.UpdateProduct(product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	
	return updatedProduct, nil
}

// DeleteProduct deletes a product
func (s *Service) DeleteProduct(id int) error {
	if id <= 0 {
		return errors.New("invalid product ID")
	}
	
	// Check if product exists
	_, err := s.repository.GetProductByID(id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}
	
	// Delete product
	err = s.repository.DeleteProduct(id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	
	return nil
}

// GetTopMostExpensive returns the most expensive products
func (s *Service) GetTopMostExpensive(limit int) ([]Product, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be greater than 0")
	}
	
	return s.repository.GetTopMostExpensive(limit)
}

// GetLowStockProducts returns products with low stock
func (s *Service) GetLowStockProducts(maxStock int) ([]Product, error) {
	if maxStock < 0 {
		return nil, errors.New("max stock must be non-negative")
	}
	
	return s.repository.GetLowStockProducts(maxStock)
}

// GetProductsByCategory returns products by category
func (s *Service) GetProductsByCategory(category string) ([]Product, error) {
	if category == "" {
		return nil, errors.New("category cannot be empty")
	}
	
	return s.repository.GetProductsByCategory(category)
}

// validateProduct validates product data
func (s *Service) validateProduct(product Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}
	
	if product.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	
	if product.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}
	
	return nil
}
