package persistence

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"obs-tools-usage/internal/product/domain/entity"
	"obs-tools-usage/internal/product/infrastructure/config"
)

// gormLogWriter implements logger.Writer interface for GORM
type gormLogWriter struct {
	logger *logrus.Logger
}

func (w *gormLogWriter) Printf(format string, args ...interface{}) {
	w.logger.Printf(format, args...)
}

// Database represents the database connection
type Database struct {
	DB     *gorm.DB
	Config *config.DatabaseConfig
	Logger *logrus.Logger
}

// NewDatabase creates a new database connection
func NewDatabase(config *config.DatabaseConfig) (*Database, error) {
	// Create GORM logger
	gormLogger := logger.New(
		&gormLogWriter{logger: config.GetLogger()},
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Build DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger := config.GetLogger()
	logger.WithFields(logrus.Fields{
		"host":     config.Host,
		"port":     config.Port,
		"database": config.DBName,
		"user":     config.User,
	}).Info("Database connected successfully")

	return &Database{
		DB:     db,
		Config: config,
		Logger: logger,
	}, nil
}

// Migrate runs database migrations
func (d *Database) Migrate() error {
	d.Logger.Info("Running database migrations...")

	// Auto migrate Product model
	if err := d.DB.AutoMigrate(&entity.Product{}); err != nil {
		d.Logger.WithError(err).Error("Failed to migrate Product model")
		return fmt.Errorf("failed to migrate Product model: %w", err)
	}

	d.Logger.Info("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		d.Logger.WithError(err).Error("Failed to close database connection")
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	d.Logger.Info("Database connection closed")
	return nil
}

// Health checks database health
func (d *Database) Health() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// SeedData seeds the database with initial data
func (d *Database) SeedData() error {
	d.Logger.Info("Seeding database with initial data...")

	// Check if data already exists
	var count int64
	d.DB.Model(&entity.Product{}).Count(&count)
	if count > 0 {
		d.Logger.Info("Database already contains data, skipping seed")
		return nil
	}

	// Sample products
	products := []entity.Product{
		{
			Name:        "MacBook Pro 16-inch",
			Description: "Apple MacBook Pro with M2 Pro chip",
			Price:       2499.99,
			Stock:       10,
			Category:    "Electronics",
		},
		{
			Name:        "iPhone 15 Pro",
			Description: "Apple iPhone 15 Pro with titanium design",
			Price:       999.99,
			Stock:       25,
			Category:    "Electronics",
		},
		{
			Name:        "Sony WH-1000XM5",
			Description: "Sony noise-canceling headphones",
			Price:       399.99,
			Stock:       15,
			Category:    "Electronics",
		},
		{
			Name:        "Nike Air Max 270",
			Description: "Nike Air Max 270 running shoes",
			Price:       150.00,
			Stock:       50,
			Category:    "Clothing",
		},
		{
			Name:        "Adidas Ultraboost 22",
			Description: "Adidas Ultraboost 22 running shoes",
			Price:       180.00,
			Stock:       30,
			Category:    "Clothing",
		},
		{
			Name:        "The Great Gatsby",
			Description: "Classic novel by F. Scott Fitzgerald",
			Price:       12.99,
			Stock:       100,
			Category:    "Books",
		},
		{
			Name:        "1984",
			Description: "Dystopian novel by George Orwell",
			Price:       14.99,
			Stock:       75,
			Category:    "Books",
		},
		{
			Name:        "Coffee Maker Deluxe",
			Description: "Programmable coffee maker with timer",
			Price:       89.99,
			Stock:       5,
			Category:    "Home & Kitchen",
		},
		{
			Name:        "Bluetooth Speaker",
			Description: "Portable Bluetooth speaker with 360° sound",
			Price:       79.99,
			Stock:       20,
			Category:    "Electronics",
		},
		{
			Name:        "Yoga Mat Premium",
			Description: "Non-slip yoga mat with carrying strap",
			Price:       45.00,
			Stock:       40,
			Category:    "Sports",
		},
	}

	// Create products
	for _, product := range products {
		product.CreatedAt = time.Now()
		product.UpdatedAt = time.Now()
		
		if err := d.DB.Create(&product).Error; err != nil {
			d.Logger.WithError(err).WithField("product", product.Name).Error("Failed to seed product")
			return fmt.Errorf("failed to seed product %s: %w", product.Name, err)
		}
	}

	d.Logger.WithField("count", len(products)).Info("Database seeded successfully")
	return nil
}
