package persistence

import (
	"errors"
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
// GetProductsByPriceRange returns products by price range
func (r *ProductRepositoryImpl) GetProductsByPriceRange(minPrice, maxPrice float64) ([]entity.Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductsByPriceRange",
		"min_price": minPrice,
		"max_price": maxPrice,
	}).Debug("Database operation started")

	var products []entity.Product
	result := r.db.Where("price BETWEEN ? AND ?", minPrice, maxPrice).Find(&products)
	duration := time.Since(start)

	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "GetProductsByPriceRange",
			"action":    "SELECT",
			"min_price": minPrice,
			"max_price": maxPrice,
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")

		external.RecordDatabaseOperation("GetProductsByPriceRange", "SELECT", duration)
		return nil, result.Error
	}

	external.RecordDatabaseOperation("GetProductsByPriceRange", "SELECT", duration)

	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductsByPriceRange",
		"action":    "SELECT",
		"min_price": minPrice,
		"max_price": maxPrice,
		"duration_ms": duration.Milliseconds(),
		"record_count": len(products),
	}).Info("Database operation completed")

	return products, nil
}

// GetProductsByName returns products by name
func (r *ProductRepositoryImpl) GetProductsByName(name string) ([]entity.Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductsByName",
		"name":      name,
	}).Debug("Database operation started")

	var products []entity.Product
	result := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&products)
	duration := time.Since(start)

	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "GetProductsByName",
			"action":    "SELECT",
			"name":      name,
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")

		external.RecordDatabaseOperation("GetProductsByName", "SELECT", duration)
		return nil, result.Error
	}

	external.RecordDatabaseOperation("GetProductsByName", "SELECT", duration)

	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductsByName",
		"action":    "SELECT",
		"name":      name,
		"duration_ms": duration.Milliseconds(),
		"record_count": len(products),
	}).Info("Database operation completed")

	return products, nil
}

// GetProductStats returns product statistics
func (r *ProductRepositoryImpl) GetProductStats() (*entity.ProductStats, error) {
	start := time.Now()
	r.logger.WithField("operation", "GetProductStats").Debug("Database operation started")

	var stats entity.ProductStats
	
	// Get total products count
	if err := r.db.Model(&entity.Product{}).Count(&stats.TotalProducts).Error; err != nil {
		return nil, err
	}

	// Get total categories count
	if err := r.db.Model(&entity.Product{}).Distinct("category").Count(&stats.TotalCategories).Error; err != nil {
		return nil, err
	}

	// Get average price
	if err := r.db.Model(&entity.Product{}).Select("AVG(price)").Scan(&stats.AveragePrice).Error; err != nil {
		return nil, err
	}

	// Get total value
	if err := r.db.Model(&entity.Product{}).Select("SUM(price * stock)").Scan(&stats.TotalValue).Error; err != nil {
		return nil, err
	}

	// Get low stock products count
	if err := r.db.Model(&entity.Product{}).Where("stock <= 10").Count(&stats.LowStockProducts).Error; err != nil {
		return nil, err
	}

	// Get out of stock products count
	if err := r.db.Model(&entity.Product{}).Where("stock = 0").Count(&stats.OutOfStockProducts).Error; err != nil {
		return nil, err
	}

	duration := time.Since(start)
	external.RecordDatabaseOperation("GetProductStats", "SELECT", duration)

	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductStats",
		"action":    "SELECT",
		"duration_ms": duration.Milliseconds(),
		"total_products": stats.TotalProducts,
		"total_categories": stats.TotalCategories,
	}).Info("Database operation completed")

	return &stats, nil
}

// GetCategories returns all categories
func (r *ProductRepositoryImpl) GetCategories() ([]entity.Category, error) {
	start := time.Now()
	r.logger.WithField("operation", "GetCategories").Debug("Database operation started")

	var categories []entity.Category
	result := r.db.Model(&entity.Product{}).
		Select("category as name, COUNT(*) as product_count, AVG(price) as average_price").
		Group("category").
		Find(&categories)
	duration := time.Since(start)

	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "GetCategories",
			"action":    "SELECT",
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")

		external.RecordDatabaseOperation("GetCategories", "SELECT", duration)
		return nil, result.Error
	}

	external.RecordDatabaseOperation("GetCategories", "SELECT", duration)

	r.logger.WithFields(logrus.Fields{
		"operation": "GetCategories",
		"action":    "SELECT",
		"duration_ms": duration.Milliseconds(),
		"record_count": len(categories),
	}).Info("Database operation completed")

	return categories, nil
}

// GetProductsByStock returns products by stock level
func (r *ProductRepositoryImpl) GetProductsByStock(stock int) ([]entity.Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductsByStock",
		"stock":     stock,
	}).Debug("Database operation started")

	var products []entity.Product
	result := r.db.Where("stock = ?", stock).Find(&products)
	duration := time.Since(start)

	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "GetProductsByStock",
			"action":    "SELECT",
			"stock":     stock,
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")

		external.RecordDatabaseOperation("GetProductsByStock", "SELECT", duration)
		return nil, result.Error
	}

	external.RecordDatabaseOperation("GetProductsByStock", "SELECT", duration)

	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductsByStock",
		"action":    "SELECT",
		"stock":     stock,
		"duration_ms": duration.Milliseconds(),
		"record_count": len(products),
	}).Info("Database operation completed")

	return products, nil
}

// GetRandomProducts returns random products
func (r *ProductRepositoryImpl) GetRandomProducts(count int) ([]entity.Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetRandomProducts",
		"count":     count,
	}).Debug("Database operation started")

	var products []entity.Product
	result := r.db.Order("RANDOM()").Limit(count).Find(&products)
	duration := time.Since(start)

	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "GetRandomProducts",
			"action":    "SELECT",
			"count":     count,
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")

		external.RecordDatabaseOperation("GetRandomProducts", "SELECT", duration)
		return nil, result.Error
	}

	external.RecordDatabaseOperation("GetRandomProducts", "SELECT", duration)

	r.logger.WithFields(logrus.Fields{
		"operation": "GetRandomProducts",
		"action":    "SELECT",
		"count":     count,
		"duration_ms": duration.Milliseconds(),
		"record_count": len(products),
	}).Info("Database operation completed")

	return products, nil
}

// GetProductsByDateRange returns products by date range
func (r *ProductRepositoryImpl) GetProductsByDateRange(startDate, endDate string) ([]entity.Product, error) {
	start := time.Now()
	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductsByDateRange",
		"start_date": startDate,
		"end_date":   endDate,
	}).Debug("Database operation started")

	var products []entity.Product
	result := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&products)
	duration := time.Since(start)

	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"operation": "GetProductsByDateRange",
			"action":    "SELECT",
			"start_date": startDate,
			"end_date":   endDate,
			"error":     result.Error.Error(),
			"duration_ms": duration.Milliseconds(),
		}).Error("Database operation failed")

		external.RecordDatabaseOperation("GetProductsByDateRange", "SELECT", duration)
		return nil, result.Error
	}

	external.RecordDatabaseOperation("GetProductsByDateRange", "SELECT", duration)

	r.logger.WithFields(logrus.Fields{
		"operation": "GetProductsByDateRange",
		"action":    "SELECT",
		"start_date": startDate,
		"end_date":   endDate,
		"duration_ms": duration.Milliseconds(),
		"record_count": len(products),
	}).Info("Database operation completed")

	return products, nil
}
