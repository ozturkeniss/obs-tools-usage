package handler

import (
	"obs-tools-usage/internal/notification/application/command"
	"obs-tools-usage/internal/notification/application/dto"
	"obs-tools-usage/internal/notification/application/usecase"
)

// CommandHandler handles all commands
type CommandHandler struct {
	notificationUseCase *usecase.NotificationUseCase
}

// NewCommandHandler creates a new command handler
func NewCommandHandler(notificationUseCase *usecase.NotificationUseCase) *CommandHandler {
	return &CommandHandler{
		notificationUseCase: notificationUseCase,
	}
}

// HandleCreateNotification handles CreateNotificationCommand
func (h *CommandHandler) HandleCreateNotification(cmd command.CreateNotificationCommand) (*dto.NotificationResponse, error) {
	return h.notificationUseCase.CreateNotification(
		cmd.UserID,
		cmd.Title,
		cmd.Message,
		cmd.Type,
		cmd.Priority,
		cmd.Channel,
		cmd.TemplateID,
		cmd.Data,
		cmd.ExpiresAt,
	)
}

// HandleUpdateNotification handles UpdateNotificationCommand
func (h *CommandHandler) HandleUpdateNotification(cmd command.UpdateNotificationCommand) (*dto.NotificationResponse, error) {
	return h.notificationUseCase.UpdateNotification(
		cmd.ID,
		cmd.Status,
		cmd.Title,
		cmd.Message,
	)
}

// HandleSendNotification handles SendNotificationCommand
func (h *CommandHandler) HandleSendNotification(cmd command.SendNotificationCommand) (*dto.NotificationResponse, error) {
	return h.notificationUseCase.SendNotification(cmd.ID)
}

// HandleMarkAsRead handles MarkAsReadCommand
func (h *CommandHandler) HandleMarkAsRead(cmd command.MarkAsReadCommand) (*dto.NotificationResponse, error) {
	return h.notificationUseCase.MarkAsRead(cmd.ID)
}

// HandleMarkAllAsRead handles MarkAllAsReadCommand
func (h *CommandHandler) HandleMarkAllAsRead(cmd command.MarkAllAsReadCommand) (*dto.NotificationResponse, error) {
	return h.notificationUseCase.MarkAllAsRead(cmd.UserID)
}

// HandleDeleteNotification handles DeleteNotificationCommand
func (h *CommandHandler) HandleDeleteNotification(cmd command.DeleteNotificationCommand) (*dto.NotificationResponse, error) {
	return h.notificationUseCase.DeleteNotification(cmd.ID)
}

// HandleBulkCreateNotification handles BulkCreateNotificationCommand
func (h *CommandHandler) HandleBulkCreateNotification(cmd command.BulkCreateNotificationCommand) (*dto.NotificationListResponse, error) {
	return h.notificationUseCase.BulkCreateNotification(
		cmd.UserIDs,
		cmd.Title,
		cmd.Message,
		cmd.Type,
		cmd.Priority,
		cmd.Channel,
		cmd.TemplateID,
		cmd.Data,
		cmd.ExpiresAt,
	)
}

// HandleScheduleNotification handles ScheduleNotificationCommand
func (h *CommandHandler) HandleScheduleNotification(cmd command.ScheduleNotificationCommand) (*dto.NotificationResponse, error) {
	return h.notificationUseCase.ScheduleNotification(
		cmd.UserID,
		cmd.Title,
		cmd.Message,
		cmd.Type,
		cmd.Priority,
		cmd.Channel,
		cmd.TemplateID,
		cmd.Data,
		cmd.SendAt,
		cmd.ExpiresAt,
	)
}

// HandleRetryFailedNotification handles RetryFailedNotificationCommand
func (h *CommandHandler) HandleRetryFailedNotification(cmd command.RetryFailedNotificationCommand) (*dto.NotificationResponse, error) {
	return h.notificationUseCase.RetryFailedNotification(cmd.ID)
}

// HandleCleanupExpiredNotifications handles CleanupExpiredNotificationsCommand
func (h *CommandHandler) HandleCleanupExpiredNotifications(cmd command.CleanupExpiredNotificationsCommand) (*dto.NotificationResponse, error) {
	return h.notificationUseCase.CleanupExpiredNotifications()
}
