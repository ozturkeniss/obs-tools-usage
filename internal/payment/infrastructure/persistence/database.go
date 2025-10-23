package persistence

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"obs-tools-usage/internal/payment/domain/entity"
	"obs-tools-usage/internal/payment/infrastructure/config"
)

// Database represents the database connection
type Database struct {
	DB     *gorm.DB
	Logger *logrus.Logger
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.Config, logger *logrus.Logger) (*Database, error) {
	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	// Configure GORM logger
	var gormLogger logger.Interface
	if cfg.IsDevelopment() {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Connect to database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxConn)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("Connected to MariaDB database")

	return &Database{
		DB:     db,
		Logger: logger,
	}, nil
}

// Migrate runs database migrations
func (d *Database) Migrate() error {
	d.Logger.Info("Running database migrations...")

	err := d.DB.AutoMigrate(
		&entity.Payment{},
		&entity.PaymentItem{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	d.Logger.Info("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Health checks database health
func (d *Database) Health() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// SeedData seeds the database with initial data
func (d *Database) SeedData() error {
	d.Logger.Info("Seeding database with initial data...")

	// Check if data already exists
	var count int64
	d.DB.Model(&entity.Payment{}).Count(&count)
	if count > 0 {
		d.Logger.Info("Database already has data, skipping seed")
		return nil
	}

	// Create sample payments for testing
	samplePayments := []entity.Payment{
		{
			ID:          "pay_1",
			UserID:      "user_1",
			BasketID:    "basket_1",
			Amount:      99.99,
			Currency:    "USD",
			Status:      entity.PaymentStatusCompleted,
			Method:      entity.PaymentMethodCreditCard,
			Provider:    "stripe",
			ProviderID:  "pi_1234567890",
			Description: "Sample payment 1",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "pay_2",
			UserID:      "user_2",
			BasketID:    "basket_2",
			Amount:      149.50,
			Currency:    "USD",
			Status:      entity.PaymentStatusPending,
			Method:      entity.PaymentMethodPayPal,
			Provider:    "paypal",
			Description: "Sample payment 2",
			CreatedAt:   time.Now().Add(-12 * time.Hour),
			UpdatedAt:   time.Now().Add(-12 * time.Hour),
		},
	}

	for _, payment := range samplePayments {
		if err := d.DB.Create(&payment).Error; err != nil {
			d.Logger.WithError(err).Error("Failed to create sample payment")
		}
	}

	d.Logger.Info("Database seeding completed successfully")
	return nil
}
