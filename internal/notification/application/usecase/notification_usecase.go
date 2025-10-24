package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"obs-tools-usage/internal/notification/application/dto"
	"obs-tools-usage/internal/notification/domain/entity"
	"obs-tools-usage/internal/notification/domain/repository"
	"obs-tools-usage/internal/notification/domain/service"
)

// NotificationUseCase handles notification business logic
type NotificationUseCase struct {
	notificationRepo     repository.NotificationRepository
	domainService        *service.NotificationDomainService
	logger               *logrus.Logger
}

// NewNotificationUseCase creates a new notification use case
func NewNotificationUseCase(
	notificationRepo repository.NotificationRepository,
	logger *logrus.Logger,
) *NotificationUseCase {
	return &NotificationUseCase{
		notificationRepo: notificationRepo,
		domainService:    service.NewNotificationDomainService(),
		logger:           logger,
	}
}

// CreateNotification creates a new notification
func (u *NotificationUseCase) CreateNotification(
	userID, title, message string,
	notificationType entity.NotificationType,
	priority entity.NotificationPriority,
	channel entity.NotificationChannel,
	templateID string,
	data map[string]string,
	expiresAt *time.Time,
) (*dto.NotificationResponse, error) {
	// Set default priority if not provided
	if priority == "" {
		priority = u.domainService.GetDefaultPriority(notificationType)
	}

	// Create notification entity
	notification := &entity.Notification{
		ID:         uuid.New().String(),
		UserID:     userID,
		Title:      title,
		Message:    message,
		Type:       notificationType,
		Priority:   priority,
		Channel:    channel,
		TemplateID: templateID,
		Data:       data,
		Status:     entity.NotificationStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
	}

	// Validate notification
	if err := u.domainService.ValidateNotification(*notification); err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	// Save to database
	ctx := context.Background()
	if err := u.notificationRepo.Create(ctx, notification); err != nil {
		u.logger.WithError(err).Error("Failed to create notification")
		return &dto.NotificationResponse{
			Success: false,
			Message: "Failed to create notification",
		}, err
	}

	// Send notification if should be sent immediately
	if u.domainService.ShouldSendImmediately(*notification) {
		go u.sendNotification(notification)
	}

	u.logger.WithFields(logrus.Fields{
		"notification_id": notification.ID,
		"user_id":         userID,
		"type":            notificationType,
		"channel":         channel,
	}).Info("Notification created")

	return &dto.NotificationResponse{
		Success:      true,
		Message:      "Notification created successfully",
		Notification: notification,
	}, nil
}

// UpdateNotification updates an existing notification
func (u *NotificationUseCase) UpdateNotification(
	id string,
	status entity.NotificationStatus,
	title, message string,
) (*dto.NotificationResponse, error) {
	ctx := context.Background()

	// Get existing notification
	notification, err := u.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Message: "Notification not found",
		}, err
	}

	// Update fields
	if status != "" {
		notification.Status = status
	}
	if title != "" {
		notification.Title = title
	}
	if message != "" {
		notification.Message = message
	}
	notification.UpdatedAt = time.Now()

	// Save changes
	if err := u.notificationRepo.Update(ctx, notification); err != nil {
		u.logger.WithError(err).Error("Failed to update notification")
		return &dto.NotificationResponse{
			Success: false,
			Message: "Failed to update notification",
		}, err
	}

	return &dto.NotificationResponse{
		Success:      true,
		Message:      "Notification updated successfully",
		Notification: notification,
	}, nil
}

// SendNotification sends a notification
func (u *NotificationUseCase) SendNotification(id string) (*dto.NotificationResponse, error) {
	ctx := context.Background()

	// Get notification
	notification, err := u.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Message: "Notification not found",
		}, err
	}

	// Send notification
	if err := u.sendNotification(notification); err != nil {
		// Mark as failed
		notification.MarkAsFailed()
		u.notificationRepo.Update(ctx, notification)

		return &dto.NotificationResponse{
			Success: false,
			Message: "Failed to send notification",
		}, err
	}

	// Mark as sent
	notification.MarkAsSent()
	u.notificationRepo.Update(ctx, notification)

	return &dto.NotificationResponse{
		Success:      true,
		Message:      "Notification sent successfully",
		Notification: notification,
	}, nil
}

// MarkAsRead marks a notification as read
func (u *NotificationUseCase) MarkAsRead(id string) (*dto.NotificationResponse, error) {
	ctx := context.Background()

	if err := u.notificationRepo.MarkAsRead(ctx, id); err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Message: "Failed to mark notification as read",
		}, err
	}

	// Get updated notification
	notification, err := u.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Message: "Notification not found",
		}, err
	}

	return &dto.NotificationResponse{
		Success:      true,
		Message:      "Notification marked as read",
		Notification: notification,
	}, nil
}

// MarkAllAsRead marks all notifications as read for a user
func (u *NotificationUseCase) MarkAllAsRead(userID string) (*dto.NotificationResponse, error) {
	ctx := context.Background()

	count, err := u.notificationRepo.MarkAllAsRead(ctx, userID)
	if err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Message: "Failed to mark notifications as read",
		}, err
	}

	return &dto.NotificationResponse{
		Success: true,
		Message: fmt.Sprintf("Marked %d notifications as read", count),
	}, nil
}

// DeleteNotification deletes a notification
func (u *NotificationUseCase) DeleteNotification(id string) (*dto.NotificationResponse, error) {
	ctx := context.Background()

	if err := u.notificationRepo.Delete(ctx, id); err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Message: "Failed to delete notification",
		}, err
	}

	return &dto.NotificationResponse{
		Success: true,
		Message: "Notification deleted successfully",
	}, nil
}

// GetNotification gets a notification by ID
func (u *NotificationUseCase) GetNotification(id string) (*dto.NotificationResponse, error) {
	ctx := context.Background()

	notification, err := u.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Message: "Notification not found",
		}, err
	}

	return &dto.NotificationResponse{
		Success:      true,
		Message:      "Notification retrieved successfully",
		Notification: notification,
	}, nil
}

// GetNotificationsByUser gets notifications for a user
func (u *NotificationUseCase) GetNotificationsByUser(
	userID, status, notificationType string,
	limit, offset int,
) (*dto.NotificationListResponse, error) {
	ctx := context.Background()

	var notifications []*entity.Notification
	var err error

	if status != "" {
		notifications, err = u.notificationRepo.GetByUserIDAndStatus(
			ctx, userID, entity.NotificationStatus(status), limit, offset,
		)
	} else if notificationType != "" {
		notifications, err = u.notificationRepo.GetByUserIDAndType(
			ctx, userID, entity.NotificationType(notificationType), limit, offset,
		)
	} else {
		notifications, err = u.notificationRepo.GetByUserID(ctx, userID, limit, offset)
	}

	if err != nil {
		return &dto.NotificationListResponse{
			Success: false,
			Message: "Failed to get notifications",
		}, err
	}

	// Get total count
	total, _ := u.notificationRepo.GetCountByUserID(ctx, userID)
	unreadCount, _ := u.notificationRepo.GetUnreadCountByUserID(ctx, userID)

	return &dto.NotificationListResponse{
		Success:       true,
		Message:       "Notifications retrieved successfully",
		Notifications: notifications,
		Total:         total,
		UnreadCount:   unreadCount,
	}, nil
}

// GetUnreadNotifications gets unread notifications for a user
func (u *NotificationUseCase) GetUnreadNotifications(
	userID string,
	limit, offset int,
) (*dto.NotificationListResponse, error) {
	ctx := context.Background()

	notifications, err := u.notificationRepo.GetUnreadByUserID(ctx, userID)
	if err != nil {
		return &dto.NotificationListResponse{
			Success: false,
			Message: "Failed to get unread notifications",
		}, err
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start >= len(notifications) {
		notifications = []*entity.Notification{}
	} else if end > len(notifications) {
		end = len(notifications)
	}

	if start < len(notifications) {
		notifications = notifications[start:end]
	}

	unreadCount := int64(len(notifications))

	return &dto.NotificationListResponse{
		Success:       true,
		Message:       "Unread notifications retrieved successfully",
		Notifications: notifications,
		Total:         unreadCount,
		UnreadCount:   unreadCount,
	}, nil
}

// GetNotificationStats gets notification statistics for a user
func (u *NotificationUseCase) GetNotificationStats(userID string) (*dto.NotificationStatsResponse, error) {
	ctx := context.Background()

	stats, err := u.notificationRepo.GetStatsByUserID(ctx, userID)
	if err != nil {
		return &dto.NotificationStatsResponse{
			Success: false,
			Message: "Failed to get notification statistics",
		}, err
	}

	return &dto.NotificationStatsResponse{
		Success: true,
		Message: "Notification statistics retrieved successfully",
		Stats:   stats,
	}, nil
}

// BulkCreateNotification creates multiple notifications
func (u *NotificationUseCase) BulkCreateNotification(
	userIDs []string,
	title, message string,
	notificationType entity.NotificationType,
	priority entity.NotificationPriority,
	channel entity.NotificationChannel,
	templateID string,
	data map[string]string,
	expiresAt *time.Time,
) (*dto.NotificationListResponse, error) {
	var notifications []*entity.Notification
	var errors []error

	for _, userID := range userIDs {
		response, err := u.CreateNotification(
			userID, title, message, notificationType,
			priority, channel, templateID, data, expiresAt,
		)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		notifications = append(notifications, response.Notification)
	}

	if len(errors) > 0 {
		u.logger.WithField("error_count", len(errors)).Warn("Some notifications failed to create")
	}

	return &dto.NotificationListResponse{
		Success:       true,
		Message:       fmt.Sprintf("Created %d notifications", len(notifications)),
		Notifications: notifications,
		Total:         int64(len(notifications)),
	}, nil
}

// ScheduleNotification schedules a notification for later sending
func (u *NotificationUseCase) ScheduleNotification(
	userID, title, message string,
	notificationType entity.NotificationType,
	priority entity.NotificationPriority,
	channel entity.NotificationChannel,
	templateID string,
	data map[string]string,
	sendAt time.Time,
	expiresAt *time.Time,
) (*dto.NotificationResponse, error) {
	// Create notification with scheduled send time
	notification := &entity.Notification{
		ID:         uuid.New().String(),
		UserID:     userID,
		Title:      title,
		Message:    message,
		Type:       notificationType,
		Priority:   priority,
		Channel:    channel,
		TemplateID: templateID,
		Data:       data,
		Status:     entity.NotificationStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
	}

	// Add send time to data
	if notification.Data == nil {
		notification.Data = make(map[string]string)
	}
	notification.Data["scheduled_send_at"] = sendAt.Format(time.RFC3339)

	ctx := context.Background()
	if err := u.notificationRepo.Create(ctx, notification); err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Message: "Failed to schedule notification",
		}, err
	}

	// Schedule sending (in production, use a job queue like Redis or RabbitMQ)
	go u.scheduleNotification(notification, sendAt)

	return &dto.NotificationResponse{
		Success:      true,
		Message:      "Notification scheduled successfully",
		Notification: notification,
	}, nil
}

// RetryFailedNotification retries a failed notification
func (u *NotificationUseCase) RetryFailedNotification(id string) (*dto.NotificationResponse, error) {
	ctx := context.Background()

	notification, err := u.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Message: "Notification not found",
		}, err
	}

	if notification.Status != entity.NotificationStatusFailed {
		return &dto.NotificationResponse{
			Success: false,
			Message: "Notification is not in failed status",
		}, nil
	}

	// Reset status and retry
	notification.Status = entity.NotificationStatusPending
	notification.UpdatedAt = time.Now()
	u.notificationRepo.Update(ctx, notification)

	// Retry sending
	if err := u.sendNotification(notification); err != nil {
		notification.MarkAsFailed()
		u.notificationRepo.Update(ctx, notification)

		return &dto.NotificationResponse{
			Success: false,
			Message: "Failed to retry notification",
		}, err
	}

	notification.MarkAsSent()
	u.notificationRepo.Update(ctx, notification)

	return &dto.NotificationResponse{
		Success:      true,
		Message:      "Notification retried successfully",
		Notification: notification,
	}, nil
}

// CleanupExpiredNotifications removes expired notifications
func (u *NotificationUseCase) CleanupExpiredNotifications() (*dto.NotificationResponse, error) {
	ctx := context.Background()

	count, err := u.notificationRepo.DeleteExpired(ctx)
	if err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Message: "Failed to cleanup expired notifications",
		}, err
	}

	return &dto.NotificationResponse{
		Success: true,
		Message: fmt.Sprintf("Cleaned up %d expired notifications", count),
	}, nil
}

// sendNotification sends a notification through the appropriate channel
func (u *NotificationUseCase) sendNotification(notification *entity.Notification) error {
	u.logger.WithFields(logrus.Fields{
		"notification_id": notification.ID,
		"user_id":         notification.UserID,
		"channel":         notification.Channel,
		"type":            notification.Type,
	}).Info("Sending notification")

	// Simulate sending notification
	// In production, implement actual sending logic for each channel
	switch notification.Channel {
	case entity.NotificationChannelEmail:
		return u.sendEmailNotification(notification)
	case entity.NotificationChannelSMS:
		return u.sendSMSNotification(notification)
	case entity.NotificationChannelPush:
		return u.sendPushNotification(notification)
	case entity.NotificationChannelInApp:
		return u.sendInAppNotification(notification)
	case entity.NotificationChannelWebhook:
		return u.sendWebhookNotification(notification)
	default:
		return fmt.Errorf("unsupported notification channel: %s", notification.Channel)
	}
}

// sendEmailNotification sends email notification
func (u *NotificationUseCase) sendEmailNotification(notification *entity.Notification) error {
	// Implement email sending logic
	u.logger.WithField("notification_id", notification.ID).Info("Sending email notification")
	return nil
}

// sendSMSNotification sends SMS notification
func (u *NotificationUseCase) sendSMSNotification(notification *entity.Notification) error {
	// Implement SMS sending logic
	u.logger.WithField("notification_id", notification.ID).Info("Sending SMS notification")
	return nil
}

// sendPushNotification sends push notification
func (u *NotificationUseCase) sendPushNotification(notification *entity.Notification) error {
	// Implement push notification logic
	u.logger.WithField("notification_id", notification.ID).Info("Sending push notification")
	return nil
}

// sendInAppNotification sends in-app notification
func (u *NotificationUseCase) sendInAppNotification(notification *entity.Notification) error {
	// In-app notifications are already stored in database
	u.logger.WithField("notification_id", notification.ID).Info("In-app notification ready")
	return nil
}

// sendWebhookNotification sends webhook notification
func (u *NotificationUseCase) sendWebhookNotification(notification *entity.Notification) error {
	// Implement webhook sending logic
	u.logger.WithField("notification_id", notification.ID).Info("Sending webhook notification")
	return nil
}

// scheduleNotification schedules a notification for later sending
func (u *NotificationUseCase) scheduleNotification(notification *entity.Notification, sendAt time.Time) {
	time.Sleep(time.Until(sendAt))
	u.sendNotification(notification)
}
