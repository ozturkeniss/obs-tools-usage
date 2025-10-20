package product

import (
	"errors"
	"sync"
)

type Repository struct {
	products map[int]*Product
	nextID   int
	mu       sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		products: make(map[int]*Product),
		nextID:   1,
	}
}

// GetAllProducts returns all products
func (r *Repository) GetAllProducts() ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	products := make([]Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, *product)
	}
	
	return products, nil
}

// GetProductByID returns a product by ID
func (r *Repository) GetProductByID(id int) (*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	product, exists := r.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	
	// Return a copy to avoid external modifications
	productCopy := *product
	return &productCopy, nil
}

// CreateProduct creates a new product
func (r *Repository) CreateProduct(product Product) (*Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Check if product with same name already exists
	for _, existingProduct := range r.products {
		if existingProduct.Name == product.Name {
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

// GetProductCount returns the total number of products
func (r *Repository) GetProductCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return len(r.products)
}
