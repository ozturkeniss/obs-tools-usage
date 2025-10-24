package service

import (
	"errors"
	"obs-tools-usage/internal/notification/domain/entity"
	"time"
)

// NotificationDomainService handles domain-specific business logic
type NotificationDomainService struct{}

// NewNotificationDomainService creates a new domain service
func NewNotificationDomainService() *NotificationDomainService {
	return &NotificationDomainService{}
}

// ValidateNotification performs domain validation on notification data
func (s *NotificationDomainService) ValidateNotification(notification entity.Notification) error {
	if notification.UserID == "" {
		return errors.New("user ID cannot be empty")
	}
	if notification.Title == "" {
		return errors.New("notification title cannot be empty")
	}
	if notification.Message == "" {
		return errors.New("notification message cannot be empty")
	}
	if notification.Type == "" {
		return errors.New("notification type cannot be empty")
	}
	if notification.Channel == "" {
		return errors.New("notification channel cannot be empty")
	}
	return nil
}

// ValidateCreateRequest validates a create notification request
func (s *NotificationDomainService) ValidateCreateRequest(req entity.CreateNotificationRequest) error {
	if req.UserID == "" {
		return errors.New("user ID is required")
	}
	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.Message == "" {
		return errors.New("message is required")
	}
	if req.Type == "" {
		return errors.New("type is required")
	}
	if req.Channel == "" {
		return errors.New("channel is required")
	}
	
	// Validate type
	if !s.IsValidNotificationType(req.Type) {
		return errors.New("invalid notification type")
	}
	
	// Validate channel
	if !s.IsValidNotificationChannel(req.Channel) {
		return errors.New("invalid notification channel")
	}
	
	// Validate priority
	if req.Priority != "" && !s.IsValidNotificationPriority(req.Priority) {
		return errors.New("invalid notification priority")
	}
	
	return nil
}

// IsValidNotificationType checks if the notification type is valid
func (s *NotificationDomainService) IsValidNotificationType(notificationType entity.NotificationType) bool {
	validTypes := []entity.NotificationType{
		entity.NotificationTypeInfo,
		entity.NotificationTypeWarning,
		entity.NotificationTypeError,
		entity.NotificationTypeSuccess,
		entity.NotificationTypePayment,
		entity.NotificationTypeOrder,
		entity.NotificationTypeSystem,
		entity.NotificationTypeMarketing,
	}
	
	for _, validType := range validTypes {
		if notificationType == validType {
			return true
		}
	}
	return false
}

// IsValidNotificationChannel checks if the notification channel is valid
func (s *NotificationDomainService) IsValidNotificationChannel(channel entity.NotificationChannel) bool {
	validChannels := []entity.NotificationChannel{
		entity.NotificationChannelEmail,
		entity.NotificationChannelSMS,
		entity.NotificationChannelPush,
		entity.NotificationChannelInApp,
		entity.NotificationChannelWebhook,
	}
	
	for _, validChannel := range validChannels {
		if channel == validChannel {
			return true
		}
	}
	return false
}

// IsValidNotificationPriority checks if the notification priority is valid
func (s *NotificationDomainService) IsValidNotificationPriority(priority entity.NotificationPriority) bool {
	validPriorities := []entity.NotificationPriority{
		entity.NotificationPriorityLow,
		entity.NotificationPriorityNormal,
		entity.NotificationPriorityHigh,
		entity.NotificationPriorityUrgent,
	}
	
	for _, validPriority := range validPriorities {
		if priority == validPriority {
			return true
		}
	}
	return false
}

// ShouldSendImmediately determines if a notification should be sent immediately
func (s *NotificationDomainService) ShouldSendImmediately(notification entity.Notification) bool {
	// High priority notifications should be sent immediately
	if notification.Priority == entity.NotificationPriorityHigh || 
	   notification.Priority == entity.NotificationPriorityUrgent {
		return true
	}
	
	// System notifications should be sent immediately
	if notification.Type == entity.NotificationTypeSystem {
		return true
	}
	
	// Payment and order notifications should be sent immediately
	if notification.Type == entity.NotificationTypePayment || 
	   notification.Type == entity.NotificationTypeOrder {
		return true
	}
	
	return false
}

// GetDefaultPriority returns the default priority for a notification type
func (s *NotificationDomainService) GetDefaultPriority(notificationType entity.NotificationType) entity.NotificationPriority {
	switch notificationType {
	case entity.NotificationTypeError:
		return entity.NotificationPriorityHigh
	case entity.NotificationTypePayment, entity.NotificationTypeOrder:
		return entity.NotificationPriorityHigh
	case entity.NotificationTypeWarning:
		return entity.NotificationPriorityNormal
	case entity.NotificationTypeMarketing:
		return entity.NotificationPriorityLow
	default:
		return entity.NotificationPriorityNormal
	}
}

// IsExpired checks if a notification is expired
func (s *NotificationDomainService) IsExpired(notification entity.Notification) bool {
	if notification.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*notification.ExpiresAt)
}

// ShouldRetry determines if a failed notification should be retried
func (s *NotificationDomainService) ShouldRetry(notification entity.Notification, retryCount int) bool {
	// Don't retry if already successful
	if notification.Status == entity.NotificationStatusSent || 
	   notification.Status == entity.NotificationStatusDelivered {
		return false
	}
	
	// Don't retry if expired
	if s.IsExpired(notification) {
		return false
	}
	
	// Don't retry if too many attempts
	maxRetries := 3
	if retryCount >= maxRetries {
		return false
	}
	
	// Retry based on priority
	switch notification.Priority {
	case entity.NotificationPriorityUrgent:
		return retryCount < 5
	case entity.NotificationPriorityHigh:
		return retryCount < 3
	case entity.NotificationPriorityNormal:
		return retryCount < 2
	case entity.NotificationPriorityLow:
		return retryCount < 1
	default:
		return retryCount < 2
	}
}
