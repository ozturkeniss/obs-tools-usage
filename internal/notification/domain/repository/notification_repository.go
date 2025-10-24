package repository

import (
	"context"
	"obs-tools-usage/internal/notification/domain/entity"
)

// NotificationRepository defines the interface for notification data operations
type NotificationRepository interface {
	// Create operations
	Create(ctx context.Context, notification *entity.Notification) error
	
	// Read operations
	GetByID(ctx context.Context, id string) (*entity.Notification, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.Notification, error)
	GetByUserIDAndStatus(ctx context.Context, userID string, status entity.NotificationStatus, limit, offset int) ([]*entity.Notification, error)
	GetByUserIDAndType(ctx context.Context, userID string, notificationType entity.NotificationType, limit, offset int) ([]*entity.Notification, error)
	GetUnreadByUserID(ctx context.Context, userID string) ([]*entity.Notification, error)
	GetExpired(ctx context.Context) ([]*entity.Notification, error)
	
	// Update operations
	Update(ctx context.Context, notification *entity.Notification) error
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) (int64, error)
	MarkAsSent(ctx context.Context, id string) error
	MarkAsDelivered(ctx context.Context, id string) error
	MarkAsFailed(ctx context.Context, id string) error
	
	// Delete operations
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
	DeleteExpired(ctx context.Context) (int64, error)
	
	// Statistics
	GetStatsByUserID(ctx context.Context, userID string) (*entity.NotificationStats, error)
	GetCountByUserID(ctx context.Context, userID string) (int64, error)
	GetUnreadCountByUserID(ctx context.Context, userID string) (int64, error)
	GetCountByStatus(ctx context.Context, status entity.NotificationStatus) (int64, error)
	GetCountByType(ctx context.Context, notificationType entity.NotificationType) (int64, error)
	GetCountByChannel(ctx context.Context, channel entity.NotificationChannel) (int64, error)
	
	// Health check
	Ping(ctx context.Context) error
}

// NotificationStats represents notification statistics
type NotificationStats struct {
	TotalNotifications    int64                        `json:"total_notifications"`
	UnreadNotifications   int64                        `json:"unread_notifications"`
	SentNotifications     int64                        `json:"sent_notifications"`
	FailedNotifications   int64                        `json:"failed_notifications"`
	PendingNotifications  int64                        `json:"pending_notifications"`
	ByType                map[string]int64             `json:"by_type"`
	ByChannel             map[string]int64             `json:"by_channel"`
	ByStatus              map[string]int64             `json:"by_status"`
}
