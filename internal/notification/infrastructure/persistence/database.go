package persistence

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"obs-tools-usage/internal/notification/domain/entity"
	"obs-tools-usage/internal/notification/infrastructure/config"
)

// Database wraps GORM database connection
type Database struct {
	DB     *gorm.DB
	config *config.Config
	logger *logrus.Logger
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.Config, logger *logrus.Logger) (*Database, error) {
	// Build DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	// Configure GORM logger
	var gormLogger logger.Interface
	if cfg.LogLevel == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

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

	return &Database{
		DB:     db,
		config: cfg,
		logger: logger,
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate runs database migrations
func (d *Database) Migrate() error {
	d.logger.Info("Running database migrations...")
	
	// Auto-migrate notification table
	if err := d.DB.AutoMigrate(&entity.Notification{}); err != nil {
		return fmt.Errorf("failed to migrate notification table: %w", err)
	}

	d.logger.Info("Database migrations completed successfully")
	return nil
}

// SeedData seeds the database with initial data
func (d *Database) SeedData() error {
	d.logger.Info("Seeding database with initial data...")

	// Check if data already exists
	var count int64
	if err := d.DB.Model(&entity.Notification{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count existing notifications: %w", err)
	}

	if count > 0 {
		d.logger.Info("Database already has data, skipping seed")
		return nil
	}

	// Create sample notifications
	sampleNotifications := []*entity.Notification{
		{
			ID:        "sample-1",
			UserID:    "user-1",
			Title:     "Welcome!",
			Message:   "Welcome to our platform!",
			Type:      entity.NotificationTypeSuccess,
			Status:    entity.NotificationStatusSent,
			Priority:  entity.NotificationPriorityNormal,
			Channel:   entity.NotificationChannelInApp,
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
			SentAt:    &[]time.Time{time.Now().Add(-24 * time.Hour)}[0],
		},
		{
			ID:        "sample-2",
			UserID:    "user-1",
			Title:     "New Product Available",
			Message:   "Check out our latest product!",
			Type:      entity.NotificationTypeInfo,
			Status:    entity.NotificationStatusPending,
			Priority:  entity.NotificationPriorityLow,
			Channel:   entity.NotificationChannelEmail,
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UpdatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        "sample-3",
			UserID:    "user-2",
			Title:     "Payment Successful",
			Message:   "Your payment has been processed successfully.",
			Type:      entity.NotificationTypePayment,
			Status:    entity.NotificationStatusDelivered,
			Priority:  entity.NotificationPriorityHigh,
			Channel:   entity.NotificationChannelEmail,
			CreatedAt: time.Now().Add(-2 * time.Hour),
			UpdatedAt: time.Now().Add(-2 * time.Hour),
			SentAt:    &[]time.Time{time.Now().Add(-2 * time.Hour)}[0],
		},
	}

	// Insert sample data
	for _, notification := range sampleNotifications {
		if err := d.DB.Create(notification).Error; err != nil {
			return fmt.Errorf("failed to create sample notification %s: %w", notification.ID, err)
		}
	}

	d.logger.Info("Database seeded successfully")
	return nil
}

// IsDevelopment checks if the environment is development
func (c *config.Config) IsDevelopment() bool {
	return c.Environment == "development"
}
