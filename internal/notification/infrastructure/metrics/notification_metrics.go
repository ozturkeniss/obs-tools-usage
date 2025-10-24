package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// NotificationMetrics holds all notification-related metrics
type NotificationMetrics struct {
	// Counter metrics
	NotificationsCreatedTotal    prometheus.Counter
	NotificationsSentTotal       prometheus.Counter
	NotificationsDeliveredTotal  prometheus.Counter
	NotificationsFailedTotal     prometheus.Counter
	NotificationsReadTotal       prometheus.Counter
	NotificationsDeletedTotal    prometheus.Counter
	
	// Counter metrics by type
	NotificationsByTypeTotal     *prometheus.CounterVec
	NotificationsByChannelTotal  *prometheus.CounterVec
	NotificationsByPriorityTotal *prometheus.CounterVec
	NotificationsByStatusTotal   *prometheus.CounterVec
	
	// Counter metrics for events
	EventsProcessedTotal         *prometheus.CounterVec
	EventsFailedTotal            *prometheus.CounterVec
	
	// Histogram metrics
	NotificationProcessingDuration prometheus.Histogram
	EventProcessingDuration       prometheus.Histogram
	DatabaseOperationDuration     *prometheus.HistogramVec
	
	// Gauge metrics
	ActiveNotifications          prometheus.Gauge
	PendingNotifications         prometheus.Gauge
	FailedNotifications          prometheus.Gauge
	UnreadNotifications          prometheus.Gauge
	
	// Gauge metrics by user
	NotificationsPerUser         *prometheus.GaugeVec
	UnreadNotificationsPerUser   *prometheus.GaugeVec
}

// NewNotificationMetrics creates a new notification metrics instance
func NewNotificationMetrics() *NotificationMetrics {
	return &NotificationMetrics{
		// Counter metrics
		NotificationsCreatedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "notification_created_total",
			Help: "Total number of notifications created",
		}),
		
		NotificationsSentTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "notification_sent_total",
			Help: "Total number of notifications sent",
		}),
		
		NotificationsDeliveredTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "notification_delivered_total",
			Help: "Total number of notifications delivered",
		}),
		
		NotificationsFailedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "notification_failed_total",
			Help: "Total number of notifications that failed to send",
		}),
		
		NotificationsReadTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "notification_read_total",
			Help: "Total number of notifications marked as read",
		}),
		
		NotificationsDeletedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "notification_deleted_total",
			Help: "Total number of notifications deleted",
		}),
		
		// Counter metrics by type
		NotificationsByTypeTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "notification_by_type_total",
			Help: "Total number of notifications by type",
		}, []string{"type"}),
		
		NotificationsByChannelTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "notification_by_channel_total",
			Help: "Total number of notifications by channel",
		}, []string{"channel"}),
		
		NotificationsByPriorityTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "notification_by_priority_total",
			Help: "Total number of notifications by priority",
		}, []string{"priority"}),
		
		NotificationsByStatusTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "notification_by_status_total",
			Help: "Total number of notifications by status",
		}, []string{"status"}),
		
		// Counter metrics for events
		EventsProcessedTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "notification_event_processed_total",
			Help: "Total number of events processed",
		}, []string{"event_type", "status"}),
		
		EventsFailedTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "notification_event_failed_total",
			Help: "Total number of events that failed to process",
		}, []string{"event_type", "error_type"}),
		
		// Histogram metrics
		NotificationProcessingDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "notification_processing_duration_seconds",
			Help:    "Time taken to process notifications",
			Buckets: prometheus.DefBuckets,
		}),
		
		EventProcessingDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "notification_event_processing_duration_seconds",
			Help:    "Time taken to process events",
			Buckets: prometheus.DefBuckets,
		}),
		
		DatabaseOperationDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "notification_database_operation_duration_seconds",
			Help:    "Time taken for database operations",
			Buckets: prometheus.DefBuckets,
		}, []string{"operation"}),
		
		// Gauge metrics
		ActiveNotifications: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "notification_active_total",
			Help: "Current number of active notifications",
		}),
		
		PendingNotifications: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "notification_pending_total",
			Help: "Current number of pending notifications",
		}),
		
		FailedNotifications: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "notification_failed_total",
			Help: "Current number of failed notifications",
		}),
		
		UnreadNotifications: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "notification_unread_total",
			Help: "Current number of unread notifications",
		}),
		
		// Gauge metrics by user
		NotificationsPerUser: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "notification_per_user_total",
			Help: "Current number of notifications per user",
		}, []string{"user_id"}),
		
		UnreadNotificationsPerUser: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "notification_unread_per_user_total",
			Help: "Current number of unread notifications per user",
		}, []string{"user_id"}),
	}
}

// IncrementNotificationCreated increments the notification created counter
func (m *NotificationMetrics) IncrementNotificationCreated(notificationType, channel, priority string) {
	m.NotificationsCreatedTotal.Inc()
	m.NotificationsByTypeTotal.WithLabelValues(notificationType).Inc()
	m.NotificationsByChannelTotal.WithLabelValues(channel).Inc()
	m.NotificationsByPriorityTotal.WithLabelValues(priority).Inc()
}

// IncrementNotificationSent increments the notification sent counter
func (m *NotificationMetrics) IncrementNotificationSent(notificationType, channel string) {
	m.NotificationsSentTotal.Inc()
	m.NotificationsByTypeTotal.WithLabelValues(notificationType).Inc()
	m.NotificationsByChannelTotal.WithLabelValues(channel).Inc()
}

// IncrementNotificationDelivered increments the notification delivered counter
func (m *NotificationMetrics) IncrementNotificationDelivered(notificationType, channel string) {
	m.NotificationsDeliveredTotal.Inc()
	m.NotificationsByTypeTotal.WithLabelValues(notificationType).Inc()
	m.NotificationsByChannelTotal.WithLabelValues(channel).Inc()
}

// IncrementNotificationFailed increments the notification failed counter
func (m *NotificationMetrics) IncrementNotificationFailed(notificationType, channel, errorType string) {
	m.NotificationsFailedTotal.Inc()
	m.NotificationsByTypeTotal.WithLabelValues(notificationType).Inc()
	m.NotificationsByChannelTotal.WithLabelValues(channel).Inc()
}

// IncrementNotificationRead increments the notification read counter
func (m *NotificationMetrics) IncrementNotificationRead(userID string) {
	m.NotificationsReadTotal.Inc()
}

// IncrementNotificationDeleted increments the notification deleted counter
func (m *NotificationMetrics) IncrementNotificationDeleted(notificationType string) {
	m.NotificationsDeletedTotal.Inc()
	m.NotificationsByTypeTotal.WithLabelValues(notificationType).Inc()
}

// IncrementEventProcessed increments the event processed counter
func (m *NotificationMetrics) IncrementEventProcessed(eventType, status string) {
	m.EventsProcessedTotal.WithLabelValues(eventType, status).Inc()
}

// IncrementEventFailed increments the event failed counter
func (m *NotificationMetrics) IncrementEventFailed(eventType, errorType string) {
	m.EventsFailedTotal.WithLabelValues(eventType, errorType).Inc()
}

// UpdateNotificationStatus updates notification status gauges
func (m *NotificationMetrics) UpdateNotificationStatus(status string, count float64) {
	switch status {
	case "pending":
		m.PendingNotifications.Set(count)
	case "failed":
		m.FailedNotifications.Set(count)
	case "unread":
		m.UnreadNotifications.Set(count)
	}
}

// UpdateNotificationCounts updates notification count gauges
func (m *NotificationMetrics) UpdateNotificationCounts(active, pending, failed, unread float64) {
	m.ActiveNotifications.Set(active)
	m.PendingNotifications.Set(pending)
	m.FailedNotifications.Set(failed)
	m.UnreadNotifications.Set(unread)
}

// UpdateUserNotificationCounts updates per-user notification count gauges
func (m *NotificationMetrics) UpdateUserNotificationCounts(userID string, total, unread float64) {
	m.NotificationsPerUser.WithLabelValues(userID).Set(total)
	m.UnreadNotificationsPerUser.WithLabelValues(userID).Set(unread)
}

// RecordNotificationProcessingDuration records notification processing duration
func (m *NotificationMetrics) RecordNotificationProcessingDuration(duration float64) {
	m.NotificationProcessingDuration.Observe(duration)
}

// RecordEventProcessingDuration records event processing duration
func (m *NotificationMetrics) RecordEventProcessingDuration(duration float64) {
	m.EventProcessingDuration.Observe(duration)
}

// RecordDatabaseOperationDuration records database operation duration
func (m *NotificationMetrics) RecordDatabaseOperationDuration(operation string, duration float64) {
	m.DatabaseOperationDuration.WithLabelValues(operation).Observe(duration)
}
