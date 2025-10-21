package product

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db:     db,
		logger: Logger,
	}
}

// GetAllProducts returns all products
func (r *Repository) GetAllProducts() ([]Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetAllProducts",
		"action":    "SELECT",
	}).Debug("Database operation started")

	var products []Product
	result := r.db.Find(&products)
	
	duration := time.Since(start)
	
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "GetAllProducts",
			"action":    "SELECT",
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")
		
		// Record failed database operation
		RecordDatabaseOperation("GetAllProducts", "SELECT", duration)
		return nil, result.Error
	}

	// Record successful database operation
	RecordDatabaseOperation("GetAllProducts", "SELECT", duration)
	
	// Update business metrics
	UpdateBusinessMetrics(products)
	
	// Log slow queries
	LogSlowQueries(r.logger.WithField("source", "repository"), "GetAllProducts", duration, 100*time.Millisecond)
	
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
	}).Debug("Database operation started")

	var product Product
	result := r.db.First(&product, id)
	
	duration := time.Since(start)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.WithFields(logrus.Fields{
				"operation": "GetProductByID",
				"action":    "SELECT",
				"product_id": id,
				"duration_ms": duration.Milliseconds(),
			}).Warn("Product not found")
			
		// Record failed database operation
		RecordDatabaseOperation("GetProductByID", "SELECT", duration)
			return nil, errors.New("product not found")
		}
		
		r.logger.WithFields(logrus.Fields{
			"operation": "GetProductByID",
			"action":    "SELECT",
			"product_id": id,
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")
		
		// Record failed database operation
		RecordDatabaseOperation("GetProductByID", "SELECT", duration)
		return nil, result.Error
	}

	// Record successful database operation
	RecordDatabaseOperation("GetProductByID", "SELECT", duration)
	
	// Log slow queries
	LogSlowQueries(r.logger.WithField("source", "repository"), "GetProductByID", duration, 50*time.Millisecond)
	
	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductByID",
		"action":    "SELECT",
		"product_id": id,
		"duration_ms": duration.Milliseconds(),
	}).Info("Database operation completed")

	return &product, nil
}

// CreateProduct creates a new product
func (r *Repository) CreateProduct(product Product) (*Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "CreateProduct",
		"action":    "INSERT",
		"name":      product.Name,
		"category":  product.Category,
	}).Debug("Database operation started")

	// Set timestamps
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	result := r.db.Create(&product)
	
	duration := time.Since(start)
	
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "CreateProduct",
			"action":    "INSERT",
			"name":      product.Name,
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")
		
		// Record failed database operation
		RecordDatabaseOperation("CreateProduct", "INSERT", duration)
		return nil, result.Error
	}

	// Record successful database operation
	RecordDatabaseOperation("CreateProduct", "INSERT", duration)
	
	r.logger.WithFields(logrus.Fields{
		"operation": "CreateProduct",
		"action":    "INSERT",
		"product_id": product.ID,
		"name":      product.Name,
		"duration_ms": duration.Milliseconds(),
	}).Info("Database operation completed")

	return &product, nil
}

// UpdateProduct updates an existing product
func (r *Repository) UpdateProduct(product Product) (*Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "UpdateProduct",
		"action":    "UPDATE",
		"product_id": product.ID,
		"name":      product.Name,
	}).Debug("Database operation started")

	// Set updated timestamp
	product.UpdatedAt = time.Now()

	result := r.db.Save(&product)
	
	duration := time.Since(start)
	
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "UpdateProduct",
			"action":    "UPDATE",
			"product_id": product.ID,
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")
		
		// Record failed database operation
		RecordDatabaseOperation("UpdateProduct", "UPDATE", duration)
		return nil, result.Error
	}

	// Record successful database operation
	RecordDatabaseOperation("UpdateProduct", "UPDATE", duration)
	
	r.logger.WithFields(logrus.Fields{
		"operation": "UpdateProduct",
		"action":    "UPDATE",
		"product_id": product.ID,
		"duration_ms": duration.Milliseconds(),
	}).Info("Database operation completed")

	return &product, nil
}

// DeleteProduct deletes a product by ID
func (r *Repository) DeleteProduct(id int) error {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "DeleteProduct",
		"action":    "DELETE",
		"product_id": id,
	}).Debug("Database operation started")

	result := r.db.Delete(&Product{}, id)
	
	duration := time.Since(start)
	
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "DeleteProduct",
			"action":    "DELETE",
			"product_id": id,
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")
		
		// Record failed database operation
		RecordDatabaseOperation("DeleteProduct", "DELETE", duration)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.WithFields(logrus.Fields{
			"operation": "DeleteProduct",
			"action":    "DELETE",
			"product_id": id,
			"duration_ms": duration.Milliseconds(),
		}).Warn("Product not found for deletion")
		
		// Record failed database operation
		RecordDatabaseOperation("DeleteProduct", "DELETE", duration)
		return errors.New("product not found")
	}

	// Record successful database operation
	RecordDatabaseOperation("DeleteProduct", "DELETE", duration)
	
	r.logger.WithFields(logrus.Fields{
		"operation": "DeleteProduct",
		"action":    "DELETE",
		"product_id": id,
		"duration_ms": duration.Milliseconds(),
		"rows_affected": result.RowsAffected,
	}).Info("Database operation completed")

	return nil
}

// GetTopMostExpensive returns the most expensive products
func (r *Repository) GetTopMostExpensive(limit int) ([]Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetTopMostExpensive",
		"action":    "SELECT",
		"limit":     limit,
	}).Debug("Database operation started")

	var products []Product
	result := r.db.Order("price DESC").Limit(limit).Find(&products)
	
	duration := time.Since(start)
	
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "GetTopMostExpensive",
			"action":    "SELECT",
			"limit":     limit,
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")
		
		// Record failed database operation
		RecordDatabaseOperation("GetTopMostExpensive", "SELECT", duration)
		return nil, result.Error
	}

	// Record successful database operation
	RecordDatabaseOperation("GetTopMostExpensive", "SELECT", duration)
	
	r.logger.WithFields(logrus.Fields{
		"operation": "GetTopMostExpensive",
		"action":    "SELECT",
		"limit":     limit,
		"duration_ms": duration.Milliseconds(),
		"record_count": len(products),
	}).Info("Database operation completed")

	return products, nil
}

// GetLowStockProducts returns products with low stock
func (r *Repository) GetLowStockProducts(maxStock int) ([]Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetLowStockProducts",
		"action":    "SELECT",
		"max_stock": maxStock,
	}).Debug("Database operation started")

	var products []Product
	result := r.db.Where("stock <= ?", maxStock).Find(&products)
	
	duration := time.Since(start)
	
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "GetLowStockProducts",
			"action":    "SELECT",
			"max_stock": maxStock,
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")
		
		// Record failed database operation
		RecordDatabaseOperation("GetLowStockProducts", "SELECT", duration)
		return nil, result.Error
	}

	// Record successful database operation
	RecordDatabaseOperation("GetLowStockProducts", "SELECT", duration)
	
	r.logger.WithFields(logrus.Fields{
		"operation": "GetLowStockProducts",
		"action":    "SELECT",
		"max_stock": maxStock,
		"duration_ms": duration.Milliseconds(),
		"record_count": len(products),
	}).Info("Database operation completed")

	return products, nil
}

// GetProductsByCategory returns products by category
func (r *Repository) GetProductsByCategory(category string) ([]Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductsByCategory",
		"action":    "SELECT",
		"category":  category,
	}).Debug("Database operation started")

	var products []Product
	result := r.db.Where("category = ?", category).Find(&products)
	
	duration := time.Since(start)
	
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "GetProductsByCategory",
			"action":    "SELECT",
			"category":  category,
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")
		
		// Record failed database operation
		RecordDatabaseOperation("GetProductsByCategory", "SELECT", duration)
		return nil, result.Error
	}

	// Record successful database operation
	RecordDatabaseOperation("GetProductsByCategory", "SELECT", duration)
	
	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductsByCategory",
		"action":    "SELECT",
		"category":  category,
		"duration_ms": duration.Milliseconds(),
		"record_count": len(products),
	}).Info("Database operation completed")

	return products, nil
}