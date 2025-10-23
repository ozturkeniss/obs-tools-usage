package loadbalancer

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

// Strategy defines the load balancing strategy
type Strategy string

const (
	RoundRobin        Strategy = "round_robin"
	LeastConnections  Strategy = "least_connections"
	WeightedRoundRobin Strategy = "weighted_round_robin"
	Random            Strategy = "random"
)

// Backend represents a backend server
type Backend struct {
	URL            *url.URL
	Weight         int
	ActiveConns    int64
	TotalRequests  int64
	FailedRequests int64
	LastHealthCheck time.Time
	Healthy        bool
	mutex          sync.RWMutex
}

// LoadBalancer manages backend servers and load balancing
type LoadBalancer struct {
	backends  []*Backend
	strategy  Strategy
	current   int64
	mutex     sync.RWMutex
	logger    *logrus.Logger
	rand      *rand.Rand
}

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(strategy Strategy, logger *logrus.Logger) *LoadBalancer {
	return &LoadBalancer{
		backends: make([]*Backend, 0),
		strategy: strategy,
		logger:   logger,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// AddBackend adds a backend server to the load balancer
func (lb *LoadBalancer) AddBackend(backendURL string, weight int) error {
	parsedURL, err := url.Parse(backendURL)
	if err != nil {
		return fmt.Errorf("invalid backend URL: %w", err)
	}

	backend := &Backend{
		URL:     parsedURL,
		Weight:  weight,
		Healthy: true,
	}

	lb.mutex.Lock()
	lb.backends = append(lb.backends, backend)
	lb.mutex.Unlock()

	lb.logger.WithFields(logrus.Fields{
		"backend": backendURL,
		"weight":  weight,
	}).Info("Backend added to load balancer")

	return nil
}

// RemoveBackend removes a backend server from the load balancer
func (lb *LoadBalancer) RemoveBackend(backendURL string) error {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	for i, backend := range lb.backends {
		if backend.URL.String() == backendURL {
			lb.backends = append(lb.backends[:i], lb.backends[i+1:]...)
			lb.logger.WithField("backend", backendURL).Info("Backend removed from load balancer")
			return nil
		}
	}

	return fmt.Errorf("backend not found: %s", backendURL)
}

// GetBackend returns the next backend server based on the strategy
func (lb *LoadBalancer) GetBackend() (*Backend, error) {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	if len(lb.backends) == 0 {
		return nil, fmt.Errorf("no backends available")
	}

	// Filter healthy backends
	healthyBackends := make([]*Backend, 0)
	for _, backend := range lb.backends {
		if backend.Healthy {
			healthyBackends = append(healthyBackends, backend)
		}
	}

	if len(healthyBackends) == 0 {
		return nil, fmt.Errorf("no healthy backends available")
	}

	switch lb.strategy {
	case RoundRobin:
		return lb.roundRobin(healthyBackends)
	case LeastConnections:
		return lb.leastConnections(healthyBackends)
	case WeightedRoundRobin:
		return lb.weightedRoundRobin(healthyBackends)
	case Random:
		return lb.random(healthyBackends)
	default:
		return lb.roundRobin(healthyBackends)
	}
}

// roundRobin implements round-robin load balancing
func (lb *LoadBalancer) roundRobin(backends []*Backend) (*Backend, error) {
	if len(backends) == 0 {
		return nil, fmt.Errorf("no backends available")
	}

	index := atomic.AddInt64(&lb.current, 1) % int64(len(backends))
	backend := backends[index]
	
	atomic.AddInt64(&backend.TotalRequests, 1)
	return backend, nil
}

// leastConnections implements least connections load balancing
func (lb *LoadBalancer) leastConnections(backends []*Backend) (*Backend, error) {
	if len(backends) == 0 {
		return nil, fmt.Errorf("no backends available")
	}

	var selected *Backend
	minConns := int64(^uint64(0) >> 1) // Max int64

	for _, backend := range backends {
		conns := atomic.LoadInt64(&backend.ActiveConns)
		if conns < minConns {
			minConns = conns
			selected = backend
		}
	}

	if selected != nil {
		atomic.AddInt64(&selected.TotalRequests, 1)
	}

	return selected, nil
}

// weightedRoundRobin implements weighted round-robin load balancing
func (lb *LoadBalancer) weightedRoundRobin(backends []*Backend) (*Backend, error) {
	if len(backends) == 0 {
		return nil, fmt.Errorf("no backends available")
	}

	// Calculate total weight
	totalWeight := 0
	for _, backend := range backends {
		totalWeight += backend.Weight
	}

	if totalWeight == 0 {
		return lb.roundRobin(backends)
	}

	// Get current index and increment
	index := atomic.AddInt64(&lb.current, 1)
	
	// Find backend based on weight
	currentWeight := int(index % int64(totalWeight))
	weightSum := 0

	for _, backend := range backends {
		weightSum += backend.Weight
		if currentWeight < weightSum {
			atomic.AddInt64(&backend.TotalRequests, 1)
			return backend, nil
		}
	}

	// Fallback to first backend
	atomic.AddInt64(&backends[0].TotalRequests, 1)
	return backends[0], nil
}

// random implements random load balancing
func (lb *LoadBalancer) random(backends []*Backend) (*Backend, error) {
	if len(backends) == 0 {
		return nil, fmt.Errorf("no backends available")
	}

	index := lb.rand.Intn(len(backends))
	backend := backends[index]
	
	atomic.AddInt64(&backend.TotalRequests, 1)
	return backend, nil
}

// IncrementConnection increments the active connection count for a backend
func (lb *LoadBalancer) IncrementConnection(backend *Backend) {
	atomic.AddInt64(&backend.ActiveConns, 1)
}

// DecrementConnection decrements the active connection count for a backend
func (lb *LoadBalancer) DecrementConnection(backend *Backend) {
	atomic.AddInt64(&backend.ActiveConns, -1)
}

// IncrementFailedRequest increments the failed request count for a backend
func (lb *LoadBalancer) IncrementFailedRequest(backend *Backend) {
	atomic.AddInt64(&backend.FailedRequests, 1)
}

// SetBackendHealth sets the health status of a backend
func (lb *LoadBalancer) SetBackendHealth(backendURL string, healthy bool) error {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	for _, backend := range lb.backends {
		if backend.URL.String() == backendURL {
			backend.mutex.Lock()
			backend.Healthy = healthy
			backend.LastHealthCheck = time.Now()
			backend.mutex.Unlock()

			lb.logger.WithFields(logrus.Fields{
				"backend": backendURL,
				"healthy": healthy,
			}).Info("Backend health status updated")

			return nil
		}
	}

	return fmt.Errorf("backend not found: %s", backendURL)
}

// GetStats returns statistics for all backends
func (lb *LoadBalancer) GetStats() []map[string]interface{} {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	stats := make([]map[string]interface{}, len(lb.backends))
	
	for i, backend := range lb.backends {
		backend.mutex.RLock()
		stats[i] = map[string]interface{}{
			"url":               backend.URL.String(),
			"weight":            backend.Weight,
			"active_connections": atomic.LoadInt64(&backend.ActiveConns),
			"total_requests":    atomic.LoadInt64(&backend.TotalRequests),
			"failed_requests":   atomic.LoadInt64(&backend.FailedRequests),
			"healthy":           backend.Healthy,
			"last_health_check": backend.LastHealthCheck,
		}
		backend.mutex.RUnlock()
	}

	return stats
}

// GetHealthyBackends returns the count of healthy backends
func (lb *LoadBalancer) GetHealthyBackends() int {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	healthy := 0
	for _, backend := range lb.backends {
		if backend.Healthy {
			healthy++
		}
	}

	return healthy
}

// GetTotalBackends returns the total count of backends
func (lb *LoadBalancer) GetTotalBackends() int {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	return len(lb.backends)
}
