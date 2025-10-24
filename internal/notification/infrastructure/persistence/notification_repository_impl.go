package persistence

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"obs-tools-usage/internal/notification/domain/repository"
	"obs-tools-usage/internal/notification/infrastructure/persistence"
)

// NewNotificationRepositoryImpl creates a new notification repository implementation
func NewNotificationRepositoryImpl(db *gorm.DB, logger *logrus.Logger) repository.NotificationRepository {
	return persistence.NewNotificationRepository(db, logger)
}
