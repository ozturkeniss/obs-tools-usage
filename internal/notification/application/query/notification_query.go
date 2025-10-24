package query

import (
	"obs-tools-usage/internal/notification/domain/entity"
)

// GetNotificationQuery represents a query to get a notification by ID
type GetNotificationQuery struct {
	ID string `json:"id" binding:"required"`
}

// GetNotificationsByUserQuery represents a query to get notifications by user ID
type GetNotificationsByUserQuery struct {
	UserID string `json:"user_id" binding:"required"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// GetUnreadNotificationsQuery represents a query to get unread notifications for a user
type GetUnreadNotificationsQuery struct {
	UserID string `json:"user_id" binding:"required"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// GetNotificationStatsQuery represents a query to get notification statistics
type GetNotificationStatsQuery struct {
	UserID string `json:"user_id" binding:"required"`
}

// GetNotificationsByTypeQuery represents a query to get notifications by type
type GetNotificationsByTypeQuery struct {
	UserID string                      `json:"user_id" binding:"required"`
	Type   entity.NotificationType     `json:"type" binding:"required"`
	Limit  int                         `json:"limit"`
	Offset int                         `json:"offset"`
}

// GetNotificationsByChannelQuery represents a query to get notifications by channel
type GetNotificationsByChannelQuery struct {
	UserID  string                      `json:"user_id" binding:"required"`
	Channel entity.NotificationChannel `json:"channel" binding:"required"`
	Limit   int                         `json:"limit"`
	Offset  int                         `json:"offset"`
}

// GetNotificationsByPriorityQuery represents a query to get notifications by priority
type GetNotificationsByPriorityQuery struct {
	UserID  string                        `json:"user_id" binding:"required"`
	Priority entity.NotificationPriority  `json:"priority" binding:"required"`
	Limit   int                           `json:"limit"`
	Offset  int                           `json:"offset"`
}

// SearchNotificationsQuery represents a query to search notifications
type SearchNotificationsQuery struct {
	UserID    string `json:"user_id" binding:"required"`
	Query     string `json:"query"`
	Type      string `json:"type"`
	Channel   string `json:"channel"`
	Status    string `json:"status"`
	Priority  string `json:"priority"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

// GetNotificationCountQuery represents a query to get notification count
type GetNotificationCountQuery struct {
	UserID string `json:"user_id" binding:"required"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// GetRecentNotificationsQuery represents a query to get recent notifications
type GetRecentNotificationsQuery struct {
	UserID string `json:"user_id" binding:"required"`
	Hours  int    `json:"hours"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
