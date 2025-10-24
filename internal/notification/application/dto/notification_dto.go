package dto

import (
	"obs-tools-usage/internal/notification/domain/entity"
	"time"
)

// NotificationResponse represents the response for notification operations
type NotificationResponse struct {
	Success      bool                    `json:"success"`
	Message      string                  `json:"message"`
	Notification *entity.Notification    `json:"notification,omitempty"`
}

// NotificationListResponse represents the response for notification list operations
type NotificationListResponse struct {
	Success       bool                    `json:"success"`
	Message       string                  `json:"message"`
	Notifications []*entity.Notification  `json:"notifications"`
	Total         int64                   `json:"total"`
	UnreadCount   int64                   `json:"unread_count"`
}

// NotificationStatsResponse represents the response for notification statistics
type NotificationStatsResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Stats   *entity.NotificationStats `json:"stats"`
}

// CreateNotificationRequest represents the request to create a notification
type CreateNotificationRequest struct {
	UserID     string                        `json:"user_id" binding:"required"`
	Title      string                        `json:"title" binding:"required"`
	Message    string                        `json:"message" binding:"required"`
	Type       entity.NotificationType       `json:"type" binding:"required"`
	Priority   entity.NotificationPriority   `json:"priority"`
	Channel    entity.NotificationChannel    `json:"channel" binding:"required"`
	TemplateID string                        `json:"template_id"`
	Data       map[string]string             `json:"data"`
	ExpiresAt  *time.Time                    `json:"expires_at"`
}

// UpdateNotificationRequest represents the request to update a notification
type UpdateNotificationRequest struct {
	Status  entity.NotificationStatus `json:"status"`
	Title   string                    `json:"title"`
	Message string                    `json:"message"`
}

// SendNotificationRequest represents the request to send a notification
type SendNotificationRequest struct {
	ID string `json:"id" binding:"required"`
}

// MarkAsReadRequest represents the request to mark notification as read
type MarkAsReadRequest struct {
	ID string `json:"id" binding:"required"`
}

// MarkAllAsReadRequest represents the request to mark all notifications as read
type MarkAllAsReadRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// GetNotificationsRequest represents the request to get notifications
type GetNotificationsRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// GetNotificationStatsRequest represents the request to get notification statistics
type GetNotificationStatsRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// BulkCreateNotificationRequest represents the request to create multiple notifications
type BulkCreateNotificationRequest struct {
	UserIDs    []string                      `json:"user_ids" binding:"required"`
	Title      string                        `json:"title" binding:"required"`
	Message    string                        `json:"message" binding:"required"`
	Type       entity.NotificationType       `json:"type" binding:"required"`
	Priority   entity.NotificationPriority   `json:"priority"`
	Channel    entity.NotificationChannel    `json:"channel" binding:"required"`
	TemplateID string                        `json:"template_id"`
	Data       map[string]string             `json:"data"`
	ExpiresAt  *time.Time                    `json:"expires_at"`
}

// ScheduleNotificationRequest represents the request to schedule a notification
type ScheduleNotificationRequest struct {
	UserID     string                        `json:"user_id" binding:"required"`
	Title      string                        `json:"title" binding:"required"`
	Message    string                        `json:"message" binding:"required"`
	Type       entity.NotificationType       `json:"type" binding:"required"`
	Priority   entity.NotificationPriority   `json:"priority"`
	Channel    entity.NotificationChannel    `json:"channel" binding:"required"`
	TemplateID string                        `json:"template_id"`
	Data       map[string]string             `json:"data"`
	SendAt     time.Time                     `json:"send_at" binding:"required"`
	ExpiresAt  *time.Time                    `json:"expires_at"`
}
