//go:build wire
// +build wire

package main

import (
	"github.com/google/wire"
	"obs-tools-usage/internal/notification/application/handler"
	"obs-tools-usage/internal/notification/application/usecase"
	"obs-tools-usage/internal/notification/infrastructure/config"
	"obs-tools-usage/internal/notification/infrastructure/metrics"
	"obs-tools-usage/internal/notification/infrastructure/persistence"
	"obs-tools-usage/internal/notification/interfaces/http"
)

// WireSet is the wire provider set for notification service
var WireSet = wire.NewSet(
	// Config
	config.LoadConfig,
	
	// Repository
	persistence.NewNotificationRepository,
	
	// Use case
	usecase.NewNotificationUseCase,
	
	// Handlers
	handler.NewCommandHandler,
	handler.NewQueryHandler,
	
	// Metrics
	metrics.NewNotificationMetrics,
	
	// HTTP Handler
	http.NewNotificationHandler,
)
