package product

import (
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Repository struct {
	products map[int]*Product
	nextID   int
	mu       sync.RWMutex
	logger   *logrus.Logger
}

func NewRepository() *Repository {
	return &Repository{
		products: make(map[int]*Product),
		nextID:   1,
		logger:   Logger,
	}
}

// GetAllProducts returns all products
func (r *Repository) GetAllProducts() ([]Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetAllProducts",
		"action":    "SELECT",
	}).Info("Database operation started")
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	products := make([]Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, *product)
	}
	
	duration := time.Since(start)
	
	// Record Prometheus metrics
	RecordDatabaseOperation("GetAllProducts", "success", duration)
	
	// Update business metrics
	UpdateBusinessMetrics(products)
	
	// Log slow queries
	LogSlowQueries(r.logger, "GetAllProducts", duration, 100*time.Millisecond)
	
	r.logger.WithFields(logrus.Fields{
		"operation": "GetAllProducts",
		"action":    "SELECT",
		"duration_ms": duration.Milliseconds(),
		"record_count": len(products),
	}).Info("Database operation completed")
	
	return products, nil
}

// GetProductByID returns a product by ID
func (r *Repository) GetProductByID(id int) (*Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductByID",
		"action":    "SELECT",
		"product_id": id,
	}).Info("Database operation started")
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	product, exists := r.products[id]
	if !exists {
		duration := time.Since(start)
		r.logger.WithFields(logrus.Fields{
			"operation": "GetProductByID",
			"action":    "SELECT",
			"product_id": id,
			"duration_ms": duration.Milliseconds(),
			"error": "product not found",
		}).Warn("Database operation failed")
		return nil, errors.New("product not found")
	}
	
	// Return a copy to avoid external modifications
	productCopy := *product
	
	duration := time.Since(start)
	
	// Log slow queries
	LogSlowQueries(r.logger, "GetProductByID", duration, 50*time.Millisecond)
	
	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductByID",
		"action":    "SELECT",
		"product_id": id,
		"duration_ms": duration.Milliseconds(),
		"found": true,
	}).Info("Database operation completed")
	
	return &productCopy, nil
}

// CreateProduct creates a new product
func (r *Repository) CreateProduct(product Product) (*Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "CreateProduct",
		"action":    "INSERT",
		"product_name": product.Name,
		"product_id": product.ID,
	}).Info("Database operation started")
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Check if product with same name already exists
	for _, existingProduct := range r.products {
		if existingProduct.Name == product.Name {
			duration := time.Since(start)
			r.logger.WithFields(logrus.Fields{
				"operation": "CreateProduct",
				"action":    "INSERT",
				"product_name": product.Name,
				"duration_ms": duration.Milliseconds(),
				"error": "product with this name already exists",
			}).Warn("Database operation failed")
			return nil, errors.New("product with this name already exists")
		}
	}
	
	// Create new product
	newProduct := &Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
	
	r.products[product.ID] = newProduct
	
	// Return a copy
	productCopy := *newProduct
	
	duration := time.Since(start)
	r.logger.WithFields(logrus.Fields{
		"operation": "CreateProduct",
		"action":    "INSERT",
		"product_name": product.Name,
		"product_id": product.ID,
		"duration_ms": duration.Milliseconds(),
		"created": true,
	}).Info("Database operation completed")
	
	return &productCopy, nil
}

// UpdateProduct updates an existing product
func (r *Repository) UpdateProduct(product Product) (*Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	existingProduct, exists := r.products[product.ID]
	if !exists {
		return nil, errors.New("product not found")
	}
	
	// Check if another product with same name exists (excluding current product)
	for id, existingProduct := range r.products {
		if id != product.ID && existingProduct.Name == product.Name {
			return nil, errors.New("product with this name already exists")
		}
	}
	
	// Update product
	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Price = product.Price
	existingProduct.Stock = product.Stock
	existingProduct.Category = product.Category
	existingProduct.UpdatedAt = product.UpdatedAt
	
	// Return a copy
	productCopy := *existingProduct
	return &productCopy, nil
}

// DeleteProduct deletes a product
func (r *Repository) DeleteProduct(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	_, exists := r.products[id]
	if !exists {
		return errors.New("product not found")
	}
	
	delete(r.products, id)
	return nil
}

// GetNextID returns the next available ID
func (r *Repository) GetNextID() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	id := r.nextID
	r.nextID++
	return id
}

// GetTopMostExpensive returns the most expensive products
func (r *Repository) GetTopMostExpensive(limit int) ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Convert map to slice
	products := make([]Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, *product)
	}
	
	// Sort by price in descending order
	for i := 0; i < len(products); i++ {
		for j := i + 1; j < len(products); j++ {
			if products[i].Price < products[j].Price {
				products[i], products[j] = products[j], products[i]
			}
		}
	}
	
	// Return only the requested number of products
	if limit > len(products) {
		limit = len(products)
	}
	
	return products[:limit], nil
}

// GetLowStockProducts returns products with low stock
func (r *Repository) GetLowStockProducts(maxStock int) ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var lowStockProducts []Product
	
	for _, product := range r.products {
		if maxStock == 1 {
			// For stock = 1
			if product.Stock == 1 {
				lowStockProducts = append(lowStockProducts, *product)
			}
		} else {
			// For stock < maxStock
			if product.Stock < maxStock {
				lowStockProducts = append(lowStockProducts, *product)
			}
		}
	}
	
	return lowStockProducts, nil
}

// GetProductsByCategory returns products by category
func (r *Repository) GetProductsByCategory(category string) ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var categoryProducts []Product
	
	for _, product := range r.products {
		if product.Category == category {
			categoryProducts = append(categoryProducts, *product)
		}
	}
	
	return categoryProducts, nil
}

// GetProductCount returns the total number of products
func (r *Repository) GetProductCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return len(r.products)
}
