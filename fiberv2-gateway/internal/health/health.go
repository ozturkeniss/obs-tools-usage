package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// HealthChecker represents a health check function
type HealthChecker func(ctx context.Context) error

// HealthManager manages health checks for the gateway
type HealthManager struct {
	checkers map[string]HealthChecker
	logger   *logrus.Logger
}

// NewHealthManager creates a new health manager
func NewHealthManager(logger *logrus.Logger) *HealthManager {
	return &HealthManager{
		checkers: make(map[string]HealthChecker),
		logger:   logger,
	}
}

// AddHealthChecker adds a health checker
func (hm *HealthManager) AddHealthChecker(name string, checker HealthChecker) {
	hm.checkers[name] = checker
	hm.logger.WithField("checker", name).Info("Health checker added")
}

// RemoveHealthChecker removes a health checker
func (hm *HealthManager) RemoveHealthChecker(name string) {
	delete(hm.checkers, name)
	hm.logger.WithField("checker", name).Info("Health checker removed")
}

// CheckHealth performs all health checks
func (hm *HealthManager) CheckHealth(ctx context.Context) map[string]interface{} {
	results := make(map[string]interface{})
	overallHealthy := true

	for name, checker := range hm.checkers {
		start := time.Now()
		err := checker(ctx)
		duration := time.Since(start)

		healthy := err == nil
		if !healthy {
			overallHealthy = false
		}

		results[name] = map[string]interface{}{
			"healthy":  healthy,
			"duration": duration.String(),
			"error":    func() string {
				if err != nil {
					return err.Error()
				}
				return ""
			}(),
		}
	}

	return map[string]interface{}{
		"healthy": overallHealthy,
		"checks":  results,
		"timestamp": time.Now(),
	}
}

// SetupHealthRoutes sets up health check routes
func SetupHealthRoutes(app *fiber.App) {
	health := app.Group("/health")

	// Basic health check
	health.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now(),
			"version":   "1.0.0",
		})
	})

	// Detailed health check
	health.Get("/detailed", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		// Create a simple health manager for basic checks
		hm := NewHealthManager(logrus.New())
		
		// Add basic health checkers
		hm.AddHealthChecker("gateway", func(ctx context.Context) error {
			// Basic gateway health check
			return nil
		})

		results := hm.CheckHealth(ctx)
		
		statusCode := 200
		if !results["healthy"].(bool) {
			statusCode = 503
		}

		return c.Status(statusCode).JSON(results)
	})

	// Readiness check
	health.Get("/ready", func(c *fiber.Ctx) error {
		// Check if the gateway is ready to serve requests
		return c.JSON(fiber.Map{
			"status":    "ready",
			"timestamp": time.Now(),
		})
	})

	// Liveness check
	health.Get("/live", func(c *fiber.Ctx) error {
		// Check if the gateway is alive
		return c.JSON(fiber.Map{
			"status":    "alive",
			"timestamp": time.Now(),
		})
	})
}
