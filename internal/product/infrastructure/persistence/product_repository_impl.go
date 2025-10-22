package persistence

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"obs-tools-usage/internal/product/domain/entity"
	"obs-tools-usage/internal/product/infrastructure/config"
	"obs-tools-usage/internal/product/infrastructure/external"
)

// ProductRepositoryImpl implements the ProductRepository interface using GORM
type ProductRepositoryImpl struct {
	db     *gorm.DB
	logger *logrus.Entry
}

// NewProductRepositoryImpl creates a new product repository implementation
func NewProductRepositoryImpl(db *gorm.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{
		db:     db,
		logger: config.GetLogger().WithField("component", "repository"),
	}
}

// GetAllProducts returns all products
func (r *ProductRepositoryImpl) GetAllProducts() ([]entity.Product, error) {
	start := time.Now()
	r.logger.WithField("operation", "GetAllProducts").Debug("Database operation started")

	var products []entity.Product
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
		external.RecordDatabaseOperation("GetAllProducts", "SELECT", duration)
		return nil, result.Error
	}

	// Record successful database operation
	external.RecordDatabaseOperation("GetAllProducts", "SELECT", duration)

	// Update business metrics
	external.UpdateBusinessMetrics(products)

	// Log slow queries
	external.LogSlowQueries(r.logger.WithField("source", "repository"), "GetAllProducts", duration, 100*time.Millisecond)

	r.logger.WithFields(logrus.Fields{
		"operation": "GetAllProducts",
		"action":    "SELECT",
		"duration_ms": duration.Milliseconds(),
		"record_count": len(products),
	}).Info("Database operation completed")

	return products, nil
}

// GetProductByID returns a product by its ID
func (r *ProductRepositoryImpl) GetProductByID(id int) (*entity.Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductByID",
		"product_id": id,
	}).Debug("Database operation started")

	var product entity.Product
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
			external.RecordDatabaseOperation("GetProductByID", "SELECT", duration)
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
		external.RecordDatabaseOperation("GetProductByID", "SELECT", duration)
		return nil, result.Error
	}

	// Record successful database operation
	external.RecordDatabaseOperation("GetProductByID", "SELECT", duration)

	// Log slow queries
	external.LogSlowQueries(r.logger.WithField("source", "repository"), "GetProductByID", duration, 50*time.Millisecond)

	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductByID",
		"action":    "SELECT",
		"product_id": id,
		"duration_ms": duration.Milliseconds(),
	}).Info("Database operation completed")

	return &product, nil
}

// CreateProduct creates a new product
func (r *ProductRepositoryImpl) CreateProduct(product entity.Product) (*entity.Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "CreateProduct",
		"name":      product.Name,
		"category":  product.Category,
	}).Debug("Database operation started")

	result := r.db.Create(&product)
	duration := time.Since(start)

	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "CreateProduct",
			"action":    "INSERT",
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")

		// Record failed database operation
		external.RecordDatabaseOperation("CreateProduct", "INSERT", duration)
		return nil, result.Error
	}

	// Record successful database operation
	external.RecordDatabaseOperation("CreateProduct", "INSERT", duration)

	r.logger.WithFields(logrus.Fields{
		"operation": "CreateProduct",
		"action":    "INSERT",
		"product_id": product.ID,
		"name":      product.Name,
		"duration_ms": duration.Milliseconds(),
	}).Info("Database operation completed")

	external.RecordProductCreated()
	return &product, nil
}

// UpdateProduct updates an existing product
func (r *ProductRepositoryImpl) UpdateProduct(product entity.Product) (*entity.Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "UpdateProduct",
		"product_id": product.ID,
		"name":      product.Name,
	}).Debug("Database operation started")

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
		external.RecordDatabaseOperation("UpdateProduct", "UPDATE", duration)
		return nil, result.Error
	}

	// Record successful database operation
	external.RecordDatabaseOperation("UpdateProduct", "UPDATE", duration)

	r.logger.WithFields(logrus.Fields{
		"operation": "UpdateProduct",
		"action":    "UPDATE",
		"product_id": product.ID,
		"name":      product.Name,
		"duration_ms": duration.Milliseconds(),
	}).Info("Database operation completed")

	external.RecordProductUpdated()
	return &product, nil
}

// DeleteProduct deletes a product by its ID
func (r *ProductRepositoryImpl) DeleteProduct(id int) error {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "DeleteProduct",
		"product_id": id,
	}).Debug("Database operation started")

	result := r.db.Delete(&entity.Product{}, id)
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
		external.RecordDatabaseOperation("DeleteProduct", "DELETE", duration)
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
		external.RecordDatabaseOperation("DeleteProduct", "DELETE", duration)
		return errors.New("product not found")
	}

	// Record successful database operation
	external.RecordDatabaseOperation("DeleteProduct", "DELETE", duration)

	r.logger.WithFields(logrus.Fields{
		"operation": "DeleteProduct",
		"action":    "DELETE",
		"product_id": id,
		"duration_ms": duration.Milliseconds(),
	}).Info("Database operation completed")

	external.RecordProductDeleted()
	return nil
}

// GetTopMostExpensive returns the top N most expensive products
func (r *ProductRepositoryImpl) GetTopMostExpensive(limit int) ([]entity.Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetTopMostExpensive",
		"limit":     limit,
	}).Debug("Database operation started")

	var products []entity.Product
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
		external.RecordDatabaseOperation("GetTopMostExpensive", "SELECT", duration)
		return nil, result.Error
	}

	// Record successful database operation
	external.RecordDatabaseOperation("GetTopMostExpensive", "SELECT", duration)

	r.logger.WithFields(logrus.Fields{
		"operation": "GetTopMostExpensive",
		"action":    "SELECT",
		"limit":     limit,
		"duration_ms": duration.Milliseconds(),
		"record_count": len(products),
	}).Info("Database operation completed")

	return products, nil
}

// GetLowStockProducts returns products with stock less than or equal to maxStock
func (r *ProductRepositoryImpl) GetLowStockProducts(maxStock int) ([]entity.Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetLowStockProducts",
		"max_stock": maxStock,
	}).Debug("Database operation started")

	var products []entity.Product
	result := r.db.Where("stock <= ?", maxStock).Order("stock ASC").Find(&products)
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
		external.RecordDatabaseOperation("GetLowStockProducts", "SELECT", duration)
		return nil, result.Error
	}

	// Record successful database operation
	external.RecordDatabaseOperation("GetLowStockProducts", "SELECT", duration)

	r.logger.WithFields(logrus.Fields{
		"operation": "GetLowStockProducts",
		"action":    "SELECT",
		"max_stock": maxStock,
		"duration_ms": duration.Milliseconds(),
		"record_count": len(products),
	}).Info("Database operation completed")

	return products, nil
}

// GetProductsByCategory returns products belonging to a specific category
func (r *ProductRepositoryImpl) GetProductsByCategory(category string) ([]entity.Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductsByCategory",
		"category":  category,
	}).Debug("Database operation started")

	var products []entity.Product
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
		external.RecordDatabaseOperation("GetProductsByCategory", "SELECT", duration)
		return nil, result.Error
	}

	// Record successful database operation
	external.RecordDatabaseOperation("GetProductsByCategory", "SELECT", duration)

	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductsByCategory",
		"action":    "SELECT",
		"category":  category,
		"duration_ms": duration.Milliseconds(),
		"record_count": len(products),
	}).Info("Database operation completed")

	return products, nil
}