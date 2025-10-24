package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"obs-tools-usage/internal/notification/application/command"
	"obs-tools-usage/internal/notification/application/dto"
	"obs-tools-usage/internal/notification/application/handler"
	"obs-tools-usage/internal/notification/application/query"
	"obs-tools-usage/internal/notification/domain/entity"
	"obs-tools-usage/internal/notification/infrastructure/metrics"
)

// NotificationHandler handles HTTP requests for notifications
type NotificationHandler struct {
	commandHandler *handler.CommandHandler
	queryHandler   *handler.QueryHandler
	metrics        *metrics.NotificationMetrics
	logger         *logrus.Logger
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(
	commandHandler *handler.CommandHandler,
	queryHandler *handler.QueryHandler,
	metrics *metrics.NotificationMetrics,
	logger *logrus.Logger,
) *NotificationHandler {
	return &NotificationHandler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
		metrics:        metrics,
		logger:         logger,
	}
}

// CreateNotification handles POST /notifications
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	start := time.Now()
	defer func() {
		h.metrics.RecordNotificationProcessingDuration(time.Since(start).Seconds())
	}()

	var req dto.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind create notification request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to command
	cmd := command.CreateNotificationCommand{
		UserID:     req.UserID,
		Title:      req.Title,
		Message:    req.Message,
		Type:       req.Type,
		Priority:   req.Priority,
		Channel:    req.Channel,
		TemplateID: req.TemplateID,
		Data:       req.Data,
		ExpiresAt:  req.ExpiresAt,
	}

	// Handle command
	response, err := h.commandHandler.HandleCreateNotification(cmd)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create notification")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
		return
	}

	// Update metrics
	if response.Success {
		h.metrics.IncrementNotificationCreated(
			string(req.Type),
			string(req.Channel),
			string(req.Priority),
		)
	}

	c.JSON(http.StatusCreated, response)
}

// GetNotification handles GET /notifications/:id
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID is required"})
		return
	}

	// Convert to query
	q := query.GetNotificationQuery{ID: id}

	// Handle query
	response, err := h.queryHandler.HandleGetNotification(q)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get notification")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notification"})
		return
	}

	if !response.Success {
		c.JSON(http.StatusNotFound, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateNotification handles PUT /notifications/:id
func (h *NotificationHandler) UpdateNotification(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID is required"})
		return
	}

	var req dto.UpdateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind update notification request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to command
	cmd := command.UpdateNotificationCommand{
		ID:      id,
		Status:  req.Status,
		Title:   req.Title,
		Message: req.Message,
	}

	// Handle command
	response, err := h.commandHandler.HandleUpdateNotification(cmd)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update notification")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SendNotification handles POST /notifications/:id/send
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID is required"})
		return
	}

	// Convert to command
	cmd := command.SendNotificationCommand{ID: id}

	// Handle command
	response, err := h.commandHandler.HandleSendNotification(cmd)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send notification")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
		return
	}

	// Update metrics
	if response.Success {
		h.metrics.IncrementNotificationSent(
			string(response.Notification.Type),
			string(response.Notification.Channel),
		)
	}

	c.JSON(http.StatusOK, response)
}

// MarkAsRead handles POST /notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID is required"})
		return
	}

	// Convert to command
	cmd := command.MarkAsReadCommand{ID: id}

	// Handle command
	response, err := h.commandHandler.HandleMarkAsRead(cmd)
	if err != nil {
		h.logger.WithError(err).Error("Failed to mark notification as read")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	// Update metrics
	if response.Success {
		h.metrics.IncrementNotificationRead(response.Notification.UserID)
	}

	c.JSON(http.StatusOK, response)
}

// MarkAllAsRead handles POST /notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	var req dto.MarkAllAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind mark all as read request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to command
	cmd := command.MarkAllAsReadCommand{UserID: req.UserID}

	// Handle command
	response, err := h.commandHandler.HandleMarkAllAsRead(cmd)
	if err != nil {
		h.logger.WithError(err).Error("Failed to mark all notifications as read")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark all notifications as read"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteNotification handles DELETE /notifications/:id
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID is required"})
		return
	}

	// Convert to command
	cmd := command.DeleteNotificationCommand{ID: id}

	// Handle command
	response, err := h.commandHandler.HandleDeleteNotification(cmd)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete notification")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
		return
	}

	// Update metrics
	if response.Success {
		h.metrics.IncrementNotificationDeleted("unknown") // We don't have type info here
	}

	c.JSON(http.StatusOK, response)
}

// GetNotifications handles GET /notifications
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.Query("status")
	notificationType := c.Query("type")

	// Convert to query
	q := query.GetNotificationsByUserQuery{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
		Status: status,
		Type:   notificationType,
	}

	// Handle query
	response, err := h.queryHandler.HandleGetNotificationsByUser(q)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get notifications")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notifications"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUnreadNotifications handles GET /notifications/unread
func (h *NotificationHandler) GetUnreadNotifications(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Convert to query
	q := query.GetUnreadNotificationsQuery{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}

	// Handle query
	response, err := h.queryHandler.HandleGetUnreadNotifications(q)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get unread notifications")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get unread notifications"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetNotificationStats handles GET /notifications/stats
func (h *NotificationHandler) GetNotificationStats(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Convert to query
	q := query.GetNotificationStatsQuery{UserID: userID}

	// Handle query
	response, err := h.queryHandler.HandleGetNotificationStats(q)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get notification stats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notification stats"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// BulkCreateNotification handles POST /notifications/bulk
func (h *NotificationHandler) BulkCreateNotification(c *gin.Context) {
	var req dto.BulkCreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind bulk create notification request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to command
	cmd := command.BulkCreateNotificationCommand{
		UserIDs:    req.UserIDs,
		Title:      req.Title,
		Message:    req.Message,
		Type:       req.Type,
		Priority:   req.Priority,
		Channel:    req.Channel,
		TemplateID: req.TemplateID,
		Data:       req.Data,
		ExpiresAt:  req.ExpiresAt,
	}

	// Handle command
	response, err := h.commandHandler.HandleBulkCreateNotification(cmd)
	if err != nil {
		h.logger.WithError(err).Error("Failed to bulk create notifications")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bulk create notifications"})
		return
	}

	// Update metrics
	if response.Success {
		for _, notification := range response.Notifications {
			h.metrics.IncrementNotificationCreated(
				string(notification.Type),
				string(notification.Channel),
				string(notification.Priority),
			)
		}
	}

	c.JSON(http.StatusCreated, response)
}

// ScheduleNotification handles POST /notifications/schedule
func (h *NotificationHandler) ScheduleNotification(c *gin.Context) {
	var req dto.ScheduleNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind schedule notification request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to command
	cmd := command.ScheduleNotificationCommand{
		UserID:     req.UserID,
		Title:      req.Title,
		Message:    req.Message,
		Type:       req.Type,
		Priority:   req.Priority,
		Channel:    req.Channel,
		TemplateID: req.TemplateID,
		Data:       req.Data,
		SendAt:     req.SendAt,
		ExpiresAt:  req.ExpiresAt,
	}

	// Handle command
	response, err := h.commandHandler.HandleScheduleNotification(cmd)
	if err != nil {
		h.logger.WithError(err).Error("Failed to schedule notification")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to schedule notification"})
		return
	}

	// Update metrics
	if response.Success {
		h.metrics.IncrementNotificationCreated(
			string(req.Type),
			string(req.Channel),
			string(req.Priority),
		)
	}

	c.JSON(http.StatusCreated, response)
}

// RetryFailedNotification handles POST /notifications/:id/retry
func (h *NotificationHandler) RetryFailedNotification(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID is required"})
		return
	}

	// Convert to command
	cmd := command.RetryFailedNotificationCommand{ID: id}

	// Handle command
	response, err := h.commandHandler.HandleRetryFailedNotification(cmd)
	if err != nil {
		h.logger.WithError(err).Error("Failed to retry notification")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retry notification"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CleanupExpiredNotifications handles POST /notifications/cleanup
func (h *NotificationHandler) CleanupExpiredNotifications(c *gin.Context) {
	// Convert to command
	cmd := command.CleanupExpiredNotificationsCommand{}

	// Handle command
	response, err := h.commandHandler.HandleCleanupExpiredNotifications(cmd)
	if err != nil {
		h.logger.WithError(err).Error("Failed to cleanup expired notifications")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup expired notifications"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheck handles GET /health
func (h *NotificationHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "notification-service",
	})
}
