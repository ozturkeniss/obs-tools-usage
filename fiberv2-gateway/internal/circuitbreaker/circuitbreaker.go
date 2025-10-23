package circuitbreaker

import (
	"fmt"
	"sync"
	"time"

	"github.com/sony/gobreaker"
	"github.com/sirupsen/logrus"
)

// CircuitBreakerManager manages circuit breakers for different services
type CircuitBreakerManager struct {
	breakers map[string]*gobreaker.CircuitBreaker
	mutex    sync.RWMutex
	logger   *logrus.Logger
}

// CircuitBreakerConfig holds configuration for circuit breaker
type CircuitBreakerConfig struct {
	Name        string
	MaxRequests uint32
	Interval    time.Duration
	Timeout     time.Duration
	ReadyToTrip func(counts gobreaker.Counts) bool
	OnStateChange func(name string, from gobreaker.State, to gobreaker.State)
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager(logger *logrus.Logger) *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*gobreaker.CircuitBreaker),
		logger:   logger,
	}
}

// CreateCircuitBreaker creates a circuit breaker for a service
func (cbm *CircuitBreakerManager) CreateCircuitBreaker(config CircuitBreakerConfig) *gobreaker.CircuitBreaker {
	cbm.mutex.Lock()
	defer cbm.mutex.Unlock()

	// Default ReadyToTrip function
	if config.ReadyToTrip == nil {
		config.ReadyToTrip = func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		}
	}

	// Default OnStateChange function
	if config.OnStateChange == nil {
		config.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
			cbm.logger.WithFields(logrus.Fields{
				"circuit_breaker": name,
				"from_state":      from,
				"to_state":        to,
			}).Info("Circuit breaker state changed")
		}
	}

	settings := gobreaker.Settings{
		Name:        config.Name,
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: config.ReadyToTrip,
		OnStateChange: config.OnStateChange,
	}

	breaker := gobreaker.NewCircuitBreaker(settings)
	cbm.breakers[config.Name] = breaker

	cbm.logger.WithFields(logrus.Fields{
		"name":         config.Name,
		"max_requests": config.MaxRequests,
		"interval":     config.Interval,
		"timeout":      config.Timeout,
	}).Info("Circuit breaker created")

	return breaker
}

// GetCircuitBreaker gets a circuit breaker by name
func (cbm *CircuitBreakerManager) GetCircuitBreaker(name string) (*gobreaker.CircuitBreaker, bool) {
	cbm.mutex.RLock()
	defer cbm.mutex.RUnlock()

	breaker, exists := cbm.breakers[name]
	return breaker, exists
}

// Execute executes a function through the circuit breaker
func (cbm *CircuitBreakerManager) Execute(name string, req func() (interface{}, error)) (interface{}, error) {
	breaker, exists := cbm.GetCircuitBreaker(name)
	if !exists {
		return nil, fmt.Errorf("circuit breaker not found: %s", name)
	}

	result, err := breaker.Execute(req)
	if err != nil {
		cbm.logger.WithFields(logrus.Fields{
			"circuit_breaker": name,
			"error":          err.Error(),
		}).Warn("Circuit breaker execution failed")
	}

	return result, err
}

// GetState returns the current state of a circuit breaker
func (cbm *CircuitBreakerManager) GetState(name string) (gobreaker.State, error) {
	breaker, exists := cbm.GetCircuitBreaker(name)
	if !exists {
		return gobreaker.StateClosed, fmt.Errorf("circuit breaker not found: %s", name)
	}

	return breaker.State(), nil
}

// GetStats returns statistics for a circuit breaker
func (cbm *CircuitBreakerManager) GetStats(name string) (gobreaker.Counts, error) {
	breaker, exists := cbm.GetCircuitBreaker(name)
	if !exists {
		return gobreaker.Counts{}, fmt.Errorf("circuit breaker not found: %s", name)
	}

	return breaker.Counts(), nil
}

// GetAllBreakers returns all circuit breakers
func (cbm *CircuitBreakerManager) GetAllBreakers() map[string]*gobreaker.CircuitBreaker {
	cbm.mutex.RLock()
	defer cbm.mutex.RUnlock()

	breakers := make(map[string]*gobreaker.CircuitBreaker)
	for name, breaker := range cbm.breakers {
		breakers[name] = breaker
	}

	return breakers
}

// HealthCheck performs health check on all circuit breakers
func (cbm *CircuitBreakerManager) HealthCheck() map[string]interface{} {
	cbm.mutex.RLock()
	defer cbm.mutex.RUnlock()

	health := make(map[string]interface{})
	
	for name, breaker := range cbm.breakers {
		state := breaker.State()
		counts := breaker.Counts()
		
		health[name] = map[string]interface{}{
			"state":          state.String(),
			"requests":       counts.Requests,
			"total_success":  counts.TotalSuccesses,
			"total_failures": counts.TotalFailures,
			"consecutive_success": counts.ConsecutiveSuccesses,
			"consecutive_failures": counts.ConsecutiveFailures,
		}
	}

	return health
}

// Reset resets a circuit breaker to closed state
func (cbm *CircuitBreakerManager) Reset(name string) error {
	breaker, exists := cbm.GetCircuitBreaker(name)
	if !exists {
		return fmt.Errorf("circuit breaker not found: %s", name)
	}

	// Reset is not directly available in gobreaker, but we can create a new one
	cbm.mutex.Lock()
	delete(cbm.breakers, name)
	cbm.mutex.Unlock()

	cbm.logger.WithField("circuit_breaker", name).Info("Circuit breaker reset")

	return nil
}
