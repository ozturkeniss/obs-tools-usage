package repository

import (
	"obs-tools-usage/internal/product/domain/entity"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	GetAllProducts() ([]entity.Product, error)
	GetProductByID(id int) (*entity.Product, error)
	CreateProduct(product entity.Product) (*entity.Product, error)
	UpdateProduct(product entity.Product) (*entity.Product, error)
	DeleteProduct(id int) error
	GetTopMostExpensive(limit int) ([]entity.Product, error)
	GetLowStockProducts(maxStock int) ([]entity.Product, error)
	GetProductsByCategory(category string) ([]entity.Product, error)
	GetProductsByPriceRange(minPrice, maxPrice float64) ([]entity.Product, error)
	GetProductsByName(name string) ([]entity.Product, error)
	GetProductStats() (*entity.ProductStats, error)
	GetCategories() ([]entity.Category, error)
	GetProductsByStock(stock int) ([]entity.Product, error)
	GetRandomProducts(count int) ([]entity.Product, error)
	GetProductsByDateRange(startDate, endDate string) ([]entity.Product, error)
}
