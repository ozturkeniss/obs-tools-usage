package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"obs-tools-usage/internal/notification/domain/entity"
	"obs-tools-usage/internal/notification/domain/repository"
)

// NotificationRepository implements the notification repository interface
type NotificationRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *gorm.DB, logger *logrus.Logger) repository.NotificationRepository {
	return &NotificationRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new notification
func (r *NotificationRepository) Create(ctx context.Context, notification *entity.Notification) error {
	if err := r.db.WithContext(ctx).Create(notification).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create notification")
		return err
	}
	return nil
}

// GetByID gets a notification by ID
func (r *NotificationRepository) GetByID(ctx context.Context, id string) (*entity.Notification, error) {
	var notification entity.Notification
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&notification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("notification not found")
		}
		r.logger.WithError(err).Error("Failed to get notification by ID")
		return nil, err
	}
	return &notification, nil
}

// GetByUserID gets notifications by user ID
func (r *NotificationRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.Notification, error) {
	var notifications []*entity.Notification
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Find(&notifications).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get notifications by user ID")
		return nil, err
	}
	return notifications, nil
}

// GetByUserIDAndStatus gets notifications by user ID and status
func (r *NotificationRepository) GetByUserIDAndStatus(ctx context.Context, userID string, status entity.NotificationStatus, limit, offset int) ([]*entity.Notification, error) {
	var notifications []*entity.Notification
	query := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, status).Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Find(&notifications).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get notifications by user ID and status")
		return nil, err
	}
	return notifications, nil
}

// GetByUserIDAndType gets notifications by user ID and type
func (r *NotificationRepository) GetByUserIDAndType(ctx context.Context, userID string, notificationType entity.NotificationType, limit, offset int) ([]*entity.Notification, error) {
	var notifications []*entity.Notification
	query := r.db.WithContext(ctx).Where("user_id = ? AND type = ?", userID, notificationType).Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Find(&notifications).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get notifications by user ID and type")
		return nil, err
	}
	return notifications, nil
}

// GetUnreadByUserID gets unread notifications by user ID
func (r *NotificationRepository) GetUnreadByUserID(ctx context.Context, userID string) ([]*entity.Notification, error) {
	var notifications []*entity.Notification
	if err := r.db.WithContext(ctx).Where("user_id = ? AND read_at IS NULL", userID).Order("created_at DESC").Find(&notifications).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get unread notifications by user ID")
		return nil, err
	}
	return notifications, nil
}

// GetExpired gets expired notifications
func (r *NotificationRepository) GetExpired(ctx context.Context) ([]*entity.Notification, error) {
	var notifications []*entity.Notification
	if err := r.db.WithContext(ctx).Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).Find(&notifications).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get expired notifications")
		return nil, err
	}
	return notifications, nil
}

// Update updates a notification
func (r *NotificationRepository) Update(ctx context.Context, notification *entity.Notification) error {
	if err := r.db.WithContext(ctx).Save(notification).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update notification")
		return err
	}
	return nil
}

// MarkAsRead marks a notification as read
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"read_at":   &now,
		"status":    entity.NotificationStatusRead,
		"updated_at": now,
	}).Error; err != nil {
		r.logger.WithError(err).Error("Failed to mark notification as read")
		return err
	}
	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) (int64, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ? AND read_at IS NULL", userID).Updates(map[string]interface{}{
		"read_at":   &now,
		"status":    entity.NotificationStatusRead,
		"updated_at": now,
	})
	
	if result.Error != nil {
		r.logger.WithError(result.Error).Error("Failed to mark all notifications as read")
		return 0, result.Error
	}
	
	return result.RowsAffected, nil
}

// MarkAsSent marks a notification as sent
func (r *NotificationRepository) MarkAsSent(ctx context.Context, id string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"sent_at":   &now,
		"status":    entity.NotificationStatusSent,
		"updated_at": now,
	}).Error; err != nil {
		r.logger.WithError(err).Error("Failed to mark notification as sent")
		return err
	}
	return nil
}

// MarkAsDelivered marks a notification as delivered
func (r *NotificationRepository) MarkAsDelivered(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    entity.NotificationStatusDelivered,
		"updated_at": time.Now(),
	}).Error; err != nil {
		r.logger.WithError(err).Error("Failed to mark notification as delivered")
		return err
	}
	return nil
}

// MarkAsFailed marks a notification as failed
func (r *NotificationRepository) MarkAsFailed(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    entity.NotificationStatusFailed,
		"updated_at": time.Now(),
	}).Error; err != nil {
		r.logger.WithError(err).Error("Failed to mark notification as failed")
		return err
	}
	return nil
}

// Delete deletes a notification
func (r *NotificationRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Notification{}, "id = ?", id).Error; err != nil {
		r.logger.WithError(err).Error("Failed to delete notification")
		return err
	}
	return nil
}

// DeleteByUserID deletes all notifications for a user
func (r *NotificationRepository) DeleteByUserID(ctx context.Context, userID string) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Notification{}, "user_id = ?", userID).Error; err != nil {
		r.logger.WithError(err).Error("Failed to delete notifications by user ID")
		return err
	}
	return nil
}

// DeleteExpired deletes expired notifications
func (r *NotificationRepository) DeleteExpired(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&entity.Notification{}, "expires_at IS NOT NULL AND expires_at < ?", time.Now())
	if result.Error != nil {
		r.logger.WithError(result.Error).Error("Failed to delete expired notifications")
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// GetStatsByUserID gets notification statistics for a user
func (r *NotificationRepository) GetStatsByUserID(ctx context.Context, userID string) (*entity.NotificationStats, error) {
	stats := &entity.NotificationStats{}
	
	// Get total notifications
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ?", userID).Count(&stats.TotalNotifications).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get total notifications count")
		return nil, err
	}
	
	// Get unread notifications
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ? AND read_at IS NULL", userID).Count(&stats.UnreadNotifications).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get unread notifications count")
		return nil, err
	}
	
	// Get sent notifications
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ? AND status = ?", userID, entity.NotificationStatusSent).Count(&stats.SentNotifications).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get sent notifications count")
		return nil, err
	}
	
	// Get failed notifications
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ? AND status = ?", userID, entity.NotificationStatusFailed).Count(&stats.FailedNotifications).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get failed notifications count")
		return nil, err
	}
	
	// Get pending notifications
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ? AND status = ?", userID, entity.NotificationStatusPending).Count(&stats.PendingNotifications).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get pending notifications count")
		return nil, err
	}
	
	// Get notifications by type
	stats.ByType = make(map[string]int64)
	var typeStats []struct {
		Type  string
		Count int64
	}
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Select("type, count(*) as count").Where("user_id = ?", userID).Group("type").Scan(&typeStats).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get notifications by type")
		return nil, err
	}
	for _, stat := range typeStats {
		stats.ByType[stat.Type] = stat.Count
	}
	
	// Get notifications by channel
	stats.ByChannel = make(map[string]int64)
	var channelStats []struct {
		Channel string
		Count   int64
	}
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Select("channel, count(*) as count").Where("user_id = ?", userID).Group("channel").Scan(&channelStats).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get notifications by channel")
		return nil, err
	}
	for _, stat := range channelStats {
		stats.ByChannel[stat.Channel] = stat.Count
	}
	
	// Get notifications by status
	stats.ByStatus = make(map[string]int64)
	var statusStats []struct {
		Status string
		Count  int64
	}
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Select("status, count(*) as count").Where("user_id = ?", userID).Group("status").Scan(&statusStats).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get notifications by status")
		return nil, err
	}
	for _, stat := range statusStats {
		stats.ByStatus[stat.Status] = stat.Count
	}
	
	return stats, nil
}

// GetCountByUserID gets notification count by user ID
func (r *NotificationRepository) GetCountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get notification count by user ID")
		return 0, err
	}
	return count, nil
}

// GetUnreadCountByUserID gets unread notification count by user ID
func (r *NotificationRepository) GetUnreadCountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ? AND read_at IS NULL", userID).Count(&count).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get unread notification count by user ID")
		return 0, err
	}
	return count, nil
}

// GetCountByStatus gets notification count by status
func (r *NotificationRepository) GetCountByStatus(ctx context.Context, status entity.NotificationStatus) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("status = ?", status).Count(&count).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get notification count by status")
		return 0, err
	}
	return count, nil
}

// GetCountByType gets notification count by type
func (r *NotificationRepository) GetCountByType(ctx context.Context, notificationType entity.NotificationType) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("type = ?", notificationType).Count(&count).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get notification count by type")
		return 0, err
	}
	return count, nil
}

// GetCountByChannel gets notification count by channel
func (r *NotificationRepository) GetCountByChannel(ctx context.Context, channel entity.NotificationChannel) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("channel = ?", channel).Count(&count).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get notification count by channel")
		return 0, err
	}
	return count, nil
}

// Ping checks database connectivity
func (r *NotificationRepository) Ping(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
