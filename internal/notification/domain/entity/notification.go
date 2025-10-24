package entity

import (
	"time"
)

// Notification represents a notification in the system
type Notification struct {
	ID          string            `json:"id" gorm:"primaryKey"`
	UserID      string            `json:"user_id" gorm:"not null;index"`
	Title       string            `json:"title" gorm:"not null"`
	Message     string            `json:"message" gorm:"not null"`
	Type        NotificationType  `json:"type" gorm:"not null"`
	Status      NotificationStatus `json:"status" gorm:"not null;default:'pending'"`
	Priority    NotificationPriority `json:"priority" gorm:"not null;default:'normal'"`
	Channel     NotificationChannel `json:"channel" gorm:"not null"`
	TemplateID  string            `json:"template_id" gorm:"index"`
	Data        map[string]string `json:"data" gorm:"type:json"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	SentAt      *time.Time        `json:"sent_at"`
	ReadAt      *time.Time        `json:"read_at"`
	ExpiresAt   *time.Time        `json:"expires_at"`
}

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeInfo     NotificationType = "info"
	NotificationTypeWarning  NotificationType = "warning"
	NotificationTypeError    NotificationType = "error"
	NotificationTypeSuccess  NotificationType = "success"
	NotificationTypePayment  NotificationType = "payment"
	NotificationTypeOrder    NotificationType = "order"
	NotificationTypeSystem   NotificationType = "system"
	NotificationTypeMarketing NotificationType = "marketing"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "pending"
	NotificationStatusSent      NotificationStatus = "sent"
	NotificationStatusDelivered NotificationStatus = "delivered"
	NotificationStatusRead      NotificationStatus = "read"
	NotificationStatusFailed    NotificationStatus = "failed"
	NotificationStatusExpired   NotificationStatus = "expired"
)

// NotificationPriority represents the priority of a notification
type NotificationPriority string

const (
	NotificationPriorityLow    NotificationPriority = "low"
	NotificationPriorityNormal NotificationPriority = "normal"
	NotificationPriorityHigh   NotificationPriority = "high"
	NotificationPriorityUrgent NotificationPriority = "urgent"
)

// NotificationChannel represents the delivery channel
type NotificationChannel string

const (
	NotificationChannelEmail    NotificationChannel = "email"
	NotificationChannelSMS      NotificationChannel = "sms"
	NotificationChannelPush     NotificationChannel = "push"
	NotificationChannelInApp    NotificationChannel = "in_app"
	NotificationChannelWebhook NotificationChannel = "webhook"
)

// CreateNotificationRequest represents the request payload for creating a notification
type CreateNotificationRequest struct {
	UserID     string            `json:"user_id" binding:"required"`
	Title      string            `json:"title" binding:"required"`
	Message    string            `json:"message" binding:"required"`
	Type       NotificationType  `json:"type" binding:"required"`
	Priority   NotificationPriority `json:"priority"`
	Channel    NotificationChannel `json:"channel" binding:"required"`
	TemplateID string            `json:"template_id"`
	Data       map[string]string `json:"data"`
	ExpiresAt  *time.Time        `json:"expires_at"`
}

// UpdateNotificationRequest represents the request payload for updating a notification
type UpdateNotificationRequest struct {
	Status  NotificationStatus `json:"status"`
	Title   string             `json:"title"`
	Message string             `json:"message"`
}

// ToDTO converts a Notification entity to a DTO-compatible struct
func (n *Notification) ToDTO() map[string]interface{} {
	return map[string]interface{}{
		"id":          n.ID,
		"user_id":     n.UserID,
		"title":       n.Title,
		"message":     n.Message,
		"type":        n.Type,
		"status":      n.Status,
		"priority":    n.Priority,
		"channel":     n.Channel,
		"template_id": n.TemplateID,
		"data":        n.Data,
		"created_at":  n.CreatedAt,
		"updated_at":  n.UpdatedAt,
		"sent_at":     n.SentAt,
		"read_at":     n.ReadAt,
		"expires_at":  n.ExpiresAt,
	}
}

// ToResponse converts a Notification to response format
func (n *Notification) ToResponse() interface{} {
	return n.ToDTO()
}

// FromCreateRequest converts CreateNotificationRequest to Notification
func (n *Notification) FromCreateRequest(req CreateNotificationRequest) {
	n.UserID = req.UserID
	n.Title = req.Title
	n.Message = req.Message
	n.Type = req.Type
	n.Priority = req.Priority
	n.Channel = req.Channel
	n.TemplateID = req.TemplateID
	n.Data = req.Data
	n.ExpiresAt = req.ExpiresAt
	n.Status = NotificationStatusPending
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
}

// IsExpired checks if the notification is expired
func (n *Notification) IsExpired() bool {
	if n.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*n.ExpiresAt)
}

// IsRead checks if the notification is read
func (n *Notification) IsRead() bool {
	return n.ReadAt != nil
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	now := time.Now()
	n.ReadAt = &now
	n.Status = NotificationStatusRead
	n.UpdatedAt = now
}

// MarkAsSent marks the notification as sent
func (n *Notification) MarkAsSent() {
	now := time.Now()
	n.SentAt = &now
	n.Status = NotificationStatusSent
	n.UpdatedAt = now
}

// MarkAsDelivered marks the notification as delivered
func (n *Notification) MarkAsDelivered() {
	n.Status = NotificationStatusDelivered
	n.UpdatedAt = time.Now()
}

// MarkAsFailed marks the notification as failed
func (n *Notification) MarkAsFailed() {
	n.Status = NotificationStatusFailed
	n.UpdatedAt = time.Now()
}
