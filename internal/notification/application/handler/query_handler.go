package handler

import (
	"obs-tools-usage/internal/notification/application/query"
	"obs-tools-usage/internal/notification/application/usecase"
)

// QueryHandler handles all queries
type QueryHandler struct {
	notificationUseCase *usecase.NotificationUseCase
}

// NewQueryHandler creates a new query handler
func NewQueryHandler(notificationUseCase *usecase.NotificationUseCase) *QueryHandler {
	return &QueryHandler{
		notificationUseCase: notificationUseCase,
	}
}

// HandleGetNotification handles GetNotificationQuery
func (h *QueryHandler) HandleGetNotification(q query.GetNotificationQuery) (*dto.NotificationResponse, error) {
	return h.notificationUseCase.GetNotification(q.ID)
}

// HandleGetNotificationsByUser handles GetNotificationsByUserQuery
func (h *QueryHandler) HandleGetNotificationsByUser(q query.GetNotificationsByUserQuery) (*dto.NotificationListResponse, error) {
	return h.notificationUseCase.GetNotificationsByUser(
		q.UserID,
		q.Limit,
		q.Offset,
		q.Status,
		q.Type,
	)
}

// HandleGetUnreadNotifications handles GetUnreadNotificationsQuery
func (h *QueryHandler) HandleGetUnreadNotifications(q query.GetUnreadNotificationsQuery) (*dto.NotificationListResponse, error) {
	return h.notificationUseCase.GetUnreadNotifications(
		q.UserID,
		q.Limit,
		q.Offset,
	)
}

// HandleGetNotificationStats handles GetNotificationStatsQuery
func (h *QueryHandler) HandleGetNotificationStats(q query.GetNotificationStatsQuery) (*dto.NotificationStatsResponse, error) {
	return h.notificationUseCase.GetNotificationStats(q.UserID)
}

// HandleGetNotificationsByType handles GetNotificationsByTypeQuery
func (h *QueryHandler) HandleGetNotificationsByType(q query.GetNotificationsByTypeQuery) (*dto.NotificationListResponse, error) {
	return h.notificationUseCase.GetNotificationsByType(
		q.UserID,
		q.Type,
		q.Limit,
		q.Offset,
	)
}

// HandleGetNotificationsByChannel handles GetNotificationsByChannelQuery
func (h *QueryHandler) HandleGetNotificationsByChannel(q query.GetNotificationsByChannelQuery) (*dto.NotificationListResponse, error) {
	return h.notificationUseCase.GetNotificationsByChannel(
		q.UserID,
		q.Channel,
		q.Limit,
		q.Offset,
	)
}

// HandleGetNotificationsByPriority handles GetNotificationsByPriorityQuery
func (h *QueryHandler) HandleGetNotificationsByPriority(q query.GetNotificationsByPriorityQuery) (*dto.NotificationListResponse, error) {
	return h.notificationUseCase.GetNotificationsByPriority(
		q.UserID,
		q.Priority,
		q.Limit,
		q.Offset,
	)
}

// HandleSearchNotifications handles SearchNotificationsQuery
func (h *QueryHandler) HandleSearchNotifications(q query.SearchNotificationsQuery) (*dto.NotificationListResponse, error) {
	return h.notificationUseCase.SearchNotifications(
		q.UserID,
		q.Query,
		q.Type,
		q.Channel,
		q.Status,
		q.Priority,
		q.StartDate,
		q.EndDate,
		q.Limit,
		q.Offset,
	)
}

// HandleGetNotificationCount handles GetNotificationCountQuery
func (h *QueryHandler) HandleGetNotificationCount(q query.GetNotificationCountQuery) (*dto.NotificationStatsResponse, error) {
	return h.notificationUseCase.GetNotificationCount(
		q.UserID,
		q.Status,
		q.Type,
	)
}

// HandleGetRecentNotifications handles GetRecentNotificationsQuery
func (h *QueryHandler) HandleGetRecentNotifications(q query.GetRecentNotificationsQuery) (*dto.NotificationListResponse, error) {
	return h.notificationUseCase.GetRecentNotifications(
		q.UserID,
		q.Hours,
		q.Limit,
		q.Offset,
	)
}
