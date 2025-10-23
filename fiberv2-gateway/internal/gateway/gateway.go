package gateway

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"fiberv2-gateway/internal/circuitbreaker"
	"fiberv2-gateway/internal/config"
	"fiberv2-gateway/internal/loadbalancer"
	"fiberv2-gateway/internal/proxy"
)

// Gateway manages the API Gateway functionality
type Gateway struct {
	config           *config.Config
	logger           *logrus.Logger
	circuitBreaker   *circuitbreaker.CircuitBreakerManager
	loadBalancers    map[string]*loadbalancer.LoadBalancer
	reverseProxy     *proxy.ReverseProxy
	mutex            sync.RWMutex
}

// NewGateway creates a new API Gateway
func NewGateway(cfg *config.Config, logger *logrus.Logger) *Gateway {
	return &Gateway{
		config:         cfg,
		logger:         logger,
		circuitBreaker: circuitbreaker.NewCircuitBreakerManager(logger),
		loadBalancers:  make(map[string]*loadbalancer.LoadBalancer),
		reverseProxy:   proxy.NewReverseProxy(proxy.ProxyConfig{
			Timeout:   30 * time.Second,
			Retries:   3,
			RetryDelay: 1 * time.Second,
			StripPath: false,
			AddHeaders: map[string]string{
				"X-Gateway": "FiberV2-Gateway",
			},
		}, logger),
	}
}

// SetupRoutes sets up all the gateway routes
func SetupRoutes(app *fiber.App, cfg *config.Config, logger *logrus.Logger) {
	gateway := NewGateway(cfg, logger)
	
	// Initialize services
	gateway.initializeServices()
	
	// Setup service routes
	gateway.setupServiceRoutes(app)
	
	// Setup admin routes
	gateway.setupAdminRoutes(app)
}

// initializeServices initializes all backend services
func (g *Gateway) initializeServices() {
	// Initialize Product Service
	if g.config.Services.Product.Enabled {
		g.initializeService("product", g.config.Services.Product.URLs, g.config.Services.Product.Timeout)
	}

	// Initialize Basket Service
	if g.config.Services.Basket.Enabled {
		g.initializeService("basket", g.config.Services.Basket.URLs, g.config.Services.Basket.Timeout)
	}

	// Initialize Payment Service
	if g.config.Services.Payment.Enabled {
		g.initializeService("payment", g.config.Services.Payment.URLs, g.config.Services.Payment.Timeout)
	}
}

// initializeService initializes a single service with load balancer and circuit breaker
func (g *Gateway) initializeService(serviceName string, urls []string, timeout int) {
	// Create load balancer for the service
	lb := loadbalancer.NewLoadBalancer(
		loadbalancer.Strategy(g.config.LoadBalancer.Strategy),
		g.logger,
	)

	// Add backends to load balancer
	for i, url := range urls {
		weight := 1 // Default weight
		if err := lb.AddBackend(url, weight); err != nil {
			g.logger.WithError(err).WithField("service", serviceName).Error("Failed to add backend")
		} else {
			g.logger.WithFields(logrus.Fields{
				"service": serviceName,
				"backend": url,
				"weight":  weight,
			}).Info("Backend added")
		}
	}

	// Store load balancer
	g.mutex.Lock()
	g.loadBalancers[serviceName] = lb
	g.mutex.Unlock()

	// Create circuit breaker for the service
	if g.config.CircuitBreaker.Enabled {
		cbConfig := circuitbreaker.CircuitBreakerConfig{
			Name:        serviceName,
			MaxRequests: g.config.CircuitBreaker.MaxRequests,
			Interval:    time.Duration(g.config.CircuitBreaker.Interval) * time.Second,
			Timeout:     time.Duration(g.config.CircuitBreaker.Timeout) * time.Second,
		}

		g.circuitBreaker.CreateCircuitBreaker(cbConfig)
	}

	g.logger.WithField("service", serviceName).Info("Service initialized")
}

// setupServiceRoutes sets up routes for backend services
func (g *Gateway) setupServiceRoutes(app *fiber.App) {
	// Product Service Routes
	if g.config.Services.Product.Enabled {
		productGroup := app.Group("/api/products")
		g.setupServiceGroup(productGroup, "product")
	}

	// Basket Service Routes
	if g.config.Services.Basket.Enabled {
		basketGroup := app.Group("/api/baskets")
		g.setupServiceGroup(basketGroup, "basket")
	}

	// Payment Service Routes
	if g.config.Services.Payment.Enabled {
		paymentGroup := app.Group("/api/payments")
		g.setupServiceGroup(paymentGroup, "payment")
	}
}

// setupServiceGroup sets up routes for a service group
func (g *Gateway) setupServiceGroup(group *fiber.Router, serviceName string) {
	// Catch-all route for the service
	group.All("/*", g.createServiceHandler(serviceName))
}

// createServiceHandler creates a handler for a service
func (g *Gateway) createServiceHandler(serviceName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get load balancer for the service
		lb, exists := g.loadBalancers[serviceName]
		if !exists {
			g.logger.WithField("service", serviceName).Error("Load balancer not found")
			return c.Status(503).JSON(fiber.Map{
				"error": "Service not available",
			})
		}

		// Get backend from load balancer
		backend, err := lb.GetBackend()
		if err != nil {
			g.logger.WithFields(logrus.Fields{
				"service": serviceName,
				"error":   err.Error(),
			}).Error("No healthy backends available")
			return c.Status(503).JSON(fiber.Map{
				"error": "No healthy backends available",
			})
		}

		// Increment connection count
		lb.IncrementConnection(backend)

		// Decrement connection count when done
		defer lb.DecrementConnection(backend)

		// Execute through circuit breaker if enabled
		if g.config.CircuitBreaker.Enabled {
			return g.executeWithCircuitBreaker(c, serviceName, backend)
		}

		// Execute directly
		return g.executeRequest(c, backend)
	}
}

// executeWithCircuitBreaker executes request through circuit breaker
func (g *Gateway) executeWithCircuitBreaker(c *fiber.Ctx, serviceName string, backend *loadbalancer.Backend) error {
	result, err := g.circuitBreaker.Execute(serviceName, func() (interface{}, error) {
		// Create a copy of the context for the circuit breaker
		ctx := c.Context()
		
		// Execute the request
		err := g.reverseProxy.FastHTTPProxy(c, backend.URL.String())
		if err != nil {
			// Increment failed request count
			g.loadBalancers[serviceName].IncrementFailedRequest(backend)
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		g.logger.WithFields(logrus.Fields{
			"service": serviceName,
			"backend": backend.URL.String(),
			"error":   err.Error(),
		}).Error("Circuit breaker execution failed")

		// Increment failed request count
		g.loadBalancers[serviceName].IncrementFailedRequest(backend)

		return c.Status(503).JSON(fiber.Map{
			"error": "Service temporarily unavailable",
		})
	}

	_ = result // Result is not used in this context
	return nil
}

// executeRequest executes request directly
func (g *Gateway) executeRequest(c *fiber.Ctx, backend *loadbalancer.Backend) error {
	err := g.reverseProxy.FastHTTPProxy(c, backend.URL.String())
	if err != nil {
		// Find the service name for this backend
		serviceName := g.findServiceNameByBackend(backend.URL.String())
		if serviceName != "" {
			g.loadBalancers[serviceName].IncrementFailedRequest(backend)
		}

		g.logger.WithFields(logrus.Fields{
			"backend": backend.URL.String(),
			"error":   err.Error(),
		}).Error("Request execution failed")

		return c.Status(502).JSON(fiber.Map{
			"error": "Backend service error",
		})
	}

	return nil
}

// findServiceNameByBackend finds the service name for a backend URL
func (g *Gateway) findServiceNameByBackend(backendURL string) string {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	for serviceName, lb := range g.loadBalancers {
		stats := lb.GetStats()
		for _, stat := range stats {
			if stat["url"] == backendURL {
				return serviceName
			}
		}
	}

	return ""
}

// setupAdminRoutes sets up administrative routes
func (g *Gateway) setupAdminRoutes(app *fiber.App) {
	admin := app.Group("/admin")

	// Gateway status
	admin.Get("/status", g.getGatewayStatus)

	// Service status
	admin.Get("/services", g.getServicesStatus)

	// Load balancer stats
	admin.Get("/loadbalancer/:service", g.getLoadBalancerStats)

	// Circuit breaker stats
	admin.Get("/circuitbreaker/:service", g.getCircuitBreakerStats)

	// Health check
	admin.Get("/health", g.getHealthCheck)
}

// getGatewayStatus returns the overall gateway status
func (g *Gateway) getGatewayStatus(c *fiber.Ctx) error {
	status := fiber.Map{
		"status":    "healthy",
		"timestamp": time.Now(),
		"services":  make(map[string]interface{}),
	}

	g.mutex.RLock()
	defer g.mutex.RUnlock()

	for serviceName, lb := range g.loadBalancers {
		healthy := lb.GetHealthyBackends()
		total := lb.GetTotalBackends()
		
		status["services"].(map[string]interface{})[serviceName] = fiber.Map{
			"healthy_backends": healthy,
			"total_backends":   total,
			"status":           func() string {
				if healthy == 0 {
					return "unhealthy"
				} else if healthy < total {
					return "degraded"
				}
				return "healthy"
			}(),
		}
	}

	return c.JSON(status)
}

// getServicesStatus returns the status of all services
func (g *Gateway) getServicesStatus(c *fiber.Ctx) error {
	services := make(map[string]interface{})

	g.mutex.RLock()
	defer g.mutex.RUnlock()

	for serviceName, lb := range g.loadBalancers {
		services[serviceName] = lb.GetStats()
	}

	return c.JSON(services)
}

// getLoadBalancerStats returns load balancer statistics for a service
func (g *Gateway) getLoadBalancerStats(c *fiber.Ctx) error {
	serviceName := c.Params("service")

	lb, exists := g.loadBalancers[serviceName]
	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "Service not found",
		})
	}

	return c.JSON(lb.GetStats())
}

// getCircuitBreakerStats returns circuit breaker statistics for a service
func (g *Gateway) getCircuitBreakerStats(c *fiber.Ctx) error {
	serviceName := c.Params("service")

	state, err := g.circuitBreaker.GetState(serviceName)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Circuit breaker not found",
		})
	}

	stats, err := g.circuitBreaker.GetStats(serviceName)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Circuit breaker not found",
		})
	}

	return c.JSON(fiber.Map{
		"state": state.String(),
		"stats": stats,
	})
}

// getHealthCheck returns the health check status
func (g *Gateway) getHealthCheck(c *fiber.Ctx) error {
	health := fiber.Map{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	}

	// Check if all services are healthy
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	for serviceName, lb := range g.loadBalancers {
		if lb.GetHealthyBackends() == 0 {
			health["status"] = "unhealthy"
			health["unhealthy_services"] = append(
				health["unhealthy_services"].([]string),
				serviceName,
			)
		}
	}

	statusCode := 200
	if health["status"] == "unhealthy" {
		statusCode = 503
	}

	return c.Status(statusCode).JSON(health)
}
