package http

import (
	"github.com/gin-gonic/gin"
	"obs-tools-usage/internal/notification/application/handler"
)

// SetupRoutes configures all notification routes
func SetupRoutes(
	r *gin.Engine,
	commandHandler *handler.CommandHandler,
	queryHandler *handler.QueryHandler,
) {
	// Create notification handler
	notificationHandler := NewNotificationHandler(
		commandHandler,
		queryHandler,
		nil, // metrics will be injected later
		nil, // logger will be injected later
	)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Notification routes
		notifications := v1.Group("/notifications")
		{
			// CRUD operations
			notifications.POST("", notificationHandler.CreateNotification)
			notifications.GET("/:id", notificationHandler.GetNotification)
			notifications.PUT("/:id", notificationHandler.UpdateNotification)
			notifications.DELETE("/:id", notificationHandler.DeleteNotification)
			
			// Notification actions
			notifications.POST("/:id/send", notificationHandler.SendNotification)
			notifications.POST("/:id/read", notificationHandler.MarkAsRead)
			notifications.POST("/:id/retry", notificationHandler.RetryFailedNotification)
			
			// Bulk operations
			notifications.POST("/read-all", notificationHandler.MarkAllAsRead)
			notifications.POST("/bulk", notificationHandler.BulkCreateNotification)
			notifications.POST("/schedule", notificationHandler.ScheduleNotification)
			notifications.POST("/cleanup", notificationHandler.CleanupExpiredNotifications)
			
			// Query operations
			notifications.GET("", notificationHandler.GetNotifications)
			notifications.GET("/unread", notificationHandler.GetUnreadNotifications)
			notifications.GET("/stats", notificationHandler.GetNotificationStats)
		}
		
		// Health check
		v1.GET("/health", notificationHandler.HealthCheck)
	}
	
	// Root health check
	r.GET("/health", notificationHandler.HealthCheck)
}
