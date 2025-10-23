package usecase

import (
	"fmt"
	"obs-tools-usage/internal/product/application/dto"
	"obs-tools-usage/internal/product/domain/entity"
	"obs-tools-usage/internal/product/domain/repository"
	"obs-tools-usage/internal/product/domain/service"
)

// ProductUseCase handles product business logic
type ProductUseCase struct {
	productRepo       repository.ProductRepository
	domainService     *service.ProductDomainService
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(productRepo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo:   productRepo,
		domainService: service.NewProductDomainService(),
	}
}

// GetAllProducts returns all products
func (uc *ProductUseCase) GetAllProducts() ([]entity.Product, error) {
	return uc.productRepo.GetAllProducts()
}

// GetProductByID returns a product by its ID
func (uc *ProductUseCase) GetProductByID(id int) (*entity.Product, error) {
	product, err := uc.productRepo.GetProductByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return product, nil
}

// CreateProduct creates a new product
func (uc *ProductUseCase) CreateProduct(req dto.CreateProductRequest) (*entity.Product, error) {
	// Convert DTO to entity
	product := entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
	}

	// Validate using domain service
	if err := uc.domainService.ValidateProduct(product); err != nil {
		return nil, err
	}

	// Create product
	createdProduct, err := uc.productRepo.CreateProduct(product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return createdProduct, nil
}

// UpdateProduct updates an existing product
func (uc *ProductUseCase) UpdateProduct(id int, req dto.UpdateProductRequest) (*entity.Product, error) {
	// Check if product exists
	existingProduct, err := uc.productRepo.GetProductByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Update fields
	existingProduct.Name = req.Name
	existingProduct.Description = req.Description
	existingProduct.Price = req.Price
	existingProduct.Stock = req.Stock
	existingProduct.Category = req.Category

	// Validate using domain service
	if err := uc.domainService.ValidateProduct(*existingProduct); err != nil {
		return nil, err
	}

	// Update product
	updatedProduct, err := uc.productRepo.UpdateProduct(*existingProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return updatedProduct, nil
}

// DeleteProduct deletes a product by its ID
func (uc *ProductUseCase) DeleteProduct(id int) error {
	err := uc.productRepo.DeleteProduct(id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

// GetTopMostExpensive returns the top N most expensive products
func (uc *ProductUseCase) GetTopMostExpensive(limit int) ([]entity.Product, error) {
	return uc.productRepo.GetTopMostExpensive(limit)
}

// GetLowStockProducts returns products with stock less than or equal to maxStock
func (uc *ProductUseCase) GetLowStockProducts(maxStock int) ([]entity.Product, error) {
	return uc.productRepo.GetLowStockProducts(maxStock)
}

// GetProductsByCategory returns products belonging to a specific category
func (uc *ProductUseCase) GetProductsByCategory(category string) ([]entity.Product, error) {
	return uc.productRepo.GetProductsByCategory(category)
}

// GetProductsByPriceRange returns products by price range
func (uc *ProductUseCase) GetProductsByPriceRange(minPrice, maxPrice float64) ([]entity.Product, error) {
	return uc.productRepo.GetProductsByPriceRange(minPrice, maxPrice)
}

// GetProductsByName returns products by name
func (uc *ProductUseCase) GetProductsByName(name string) ([]entity.Product, error) {
	return uc.productRepo.GetProductsByName(name)
}

// GetProductStats returns product statistics
func (uc *ProductUseCase) GetProductStats() (*entity.ProductStats, error) {
	return uc.productRepo.GetProductStats()
}

// GetCategories returns all categories
func (uc *ProductUseCase) GetCategories() ([]entity.Category, error) {
	return uc.productRepo.GetCategories()
}

// GetProductsByStock returns products by stock level
func (uc *ProductUseCase) GetProductsByStock(stock int) ([]entity.Product, error) {
	return uc.productRepo.GetProductsByStock(stock)
}

// GetRandomProducts returns random products
func (uc *ProductUseCase) GetRandomProducts(count int) ([]entity.Product, error) {
	return uc.productRepo.GetRandomProducts(count)
}

// GetProductsByDateRange returns products by date range
func (uc *ProductUseCase) GetProductsByDateRange(startDate, endDate string) ([]entity.Product, error) {
	return uc.productRepo.GetProductsByDateRange(startDate, endDate)
}
