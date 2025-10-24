package command

import (
	"obs-tools-usage/internal/notification/domain/entity"
	"time"
)

// CreateNotificationCommand represents a command to create a notification
type CreateNotificationCommand struct {
	UserID     string            `json:"user_id" binding:"required"`
	Title      string            `json:"title" binding:"required"`
	Message    string            `json:"message" binding:"required"`
	Type       entity.NotificationType `json:"type" binding:"required"`
	Priority   entity.NotificationPriority `json:"priority"`
	Channel    entity.NotificationChannel `json:"channel" binding:"required"`
	TemplateID string            `json:"template_id"`
	Data       map[string]string `json:"data"`
	ExpiresAt  *time.Time        `json:"expires_at"`
}

// ToDTO converts CreateNotificationCommand to CreateNotificationRequest
func (c CreateNotificationCommand) ToDTO() entity.CreateNotificationRequest {
	return entity.CreateNotificationRequest{
		UserID:     c.UserID,
		Title:      c.Title,
		Message:    c.Message,
		Type:       c.Type,
		Priority:   c.Priority,
		Channel:    c.Channel,
		TemplateID: c.TemplateID,
		Data:       c.Data,
		ExpiresAt:  c.ExpiresAt,
	}
}

// UpdateNotificationCommand represents a command to update a notification
type UpdateNotificationCommand struct {
	ID      string                      `json:"id" binding:"required"`
	Status  entity.NotificationStatus   `json:"status"`
	Title   string                      `json:"title"`
	Message string                      `json:"message"`
}

// SendNotificationCommand represents a command to send a notification
type SendNotificationCommand struct {
	ID string `json:"id" binding:"required"`
}

// MarkAsReadCommand represents a command to mark a notification as read
type MarkAsReadCommand struct {
	ID string `json:"id" binding:"required"`
}

// MarkAllAsReadCommand represents a command to mark all notifications as read for a user
type MarkAllAsReadCommand struct {
	UserID string `json:"user_id" binding:"required"`
}

// DeleteNotificationCommand represents a command to delete a notification
type DeleteNotificationCommand struct {
	ID string `json:"id" binding:"required"`
}

// BulkCreateNotificationCommand represents a command to create multiple notifications
type BulkCreateNotificationCommand struct {
	UserIDs    []string                  `json:"user_ids" binding:"required"`
	Title      string                    `json:"title" binding:"required"`
	Message    string                    `json:"message" binding:"required"`
	Type       entity.NotificationType   `json:"type" binding:"required"`
	Priority   entity.NotificationPriority `json:"priority"`
	Channel    entity.NotificationChannel `json:"channel" binding:"required"`
	TemplateID string                   `json:"template_id"`
	Data       map[string]string        `json:"data"`
	ExpiresAt  *time.Time               `json:"expires_at"`
}

// ScheduleNotificationCommand represents a command to schedule a notification
type ScheduleNotificationCommand struct {
	UserID     string            `json:"user_id" binding:"required"`
	Title      string            `json:"title" binding:"required"`
	Message    string            `json:"message" binding:"required"`
	Type       entity.NotificationType `json:"type" binding:"required"`
	Priority   entity.NotificationPriority `json:"priority"`
	Channel    entity.NotificationChannel `json:"channel" binding:"required"`
	TemplateID string            `json:"template_id"`
	Data       map[string]string  `json:"data"`
	SendAt     time.Time          `json:"send_at" binding:"required"`
	ExpiresAt  *time.Time         `json:"expires_at"`
}

// RetryFailedNotificationCommand represents a command to retry a failed notification
type RetryFailedNotificationCommand struct {
	ID string `json:"id" binding:"required"`
}

// CleanupExpiredNotificationsCommand represents a command to cleanup expired notifications
type CleanupExpiredNotificationsCommand struct {
	// No fields needed - cleanup all expired notifications
}
