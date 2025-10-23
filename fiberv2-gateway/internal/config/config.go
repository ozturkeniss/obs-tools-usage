package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sony/gobreaker"
)

// Config holds the configuration for the API Gateway
type Config struct {
	Port        string
	Environment string
	LogLevel    string
	LogFormat   string
	
	// Redis configuration
	Redis RedisConfig
	
	// Services configuration
	Services ServicesConfig
	
	// Circuit breaker configuration
	CircuitBreaker CircuitBreakerConfig
	
	// Load balancer configuration
	LoadBalancer LoadBalancerConfig
	
	// Rate limiting configuration
	RateLimit RateLimitConfig
	
	// Health check configuration
	Health HealthConfig
	
	// Metrics configuration
	Metrics MetricsConfig
}

// ServicesConfig holds configuration for backend services
type ServicesConfig struct {
	Product ProductServiceConfig
	Basket  BasketServiceConfig
	Payment PaymentServiceConfig
}

// ProductServiceConfig holds product service configuration
type ProductServiceConfig struct {
	Name     string
	URLs     []string
	Timeout  int
	Retries  int
	Enabled  bool
}

// BasketServiceConfig holds basket service configuration
type BasketServiceConfig struct {
	Name     string
	URLs     []string
	Timeout  int
	Retries  int
	Enabled  bool
}

// PaymentServiceConfig holds payment service configuration
type PaymentServiceConfig struct {
	Name     string
	URLs     []string
	Timeout  int
	Retries  int
	Enabled  bool
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	Enabled           bool
	MaxRequests       uint32
	Interval          int
	Timeout           int
	ReadyToTrip       func(counts gobreaker.Counts) bool
	OnStateChange     func(name string, from gobreaker.State, to gobreaker.State)
}

// LoadBalancerConfig holds load balancer configuration
type LoadBalancerConfig struct {
	Strategy string // round_robin, least_connections, weighted_round_robin
	Enabled  bool
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled    bool
	Requests   int
	Window     time.Duration
	Burst      int
}

// HealthConfig holds health check configuration
type HealthConfig struct {
	Enabled        bool
	CheckInterval  time.Duration
	Timeout        time.Duration
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Enabled bool
	Path    string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host         string
	Port         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		LogFormat:   getEnv("LOG_FORMAT", "json"),
		
		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnv("REDIS_PORT", "6379"),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
			MaxRetries:   getEnvAsInt("REDIS_MAX_RETRIES", 3),
			DialTimeout:  getEnvAsDuration("REDIS_DIAL_TIMEOUT", "5s"),
			ReadTimeout:  getEnvAsDuration("REDIS_READ_TIMEOUT", "3s"),
			WriteTimeout: getEnvAsDuration("REDIS_WRITE_TIMEOUT", "3s"),
		},
		
		Services: ServicesConfig{
			Product: ProductServiceConfig{
				Name:     getEnv("PRODUCT_SERVICE_NAME", "product-service"),
				URLs:     getEnvSlice("PRODUCT_SERVICE_URLS", []string{"http://localhost:8080"}),
				Timeout:  getEnvAsInt("PRODUCT_SERVICE_TIMEOUT", 30),
				Retries:  getEnvAsInt("PRODUCT_SERVICE_RETRIES", 3),
				Enabled:  getEnvAsBool("PRODUCT_SERVICE_ENABLED", true),
			},
			Basket: BasketServiceConfig{
				Name:     getEnv("BASKET_SERVICE_NAME", "basket-service"),
				URLs:     getEnvSlice("BASKET_SERVICE_URLS", []string{"http://localhost:8081"}),
				Timeout:  getEnvAsInt("BASKET_SERVICE_TIMEOUT", 30),
				Retries:  getEnvAsInt("BASKET_SERVICE_RETRIES", 3),
				Enabled:  getEnvAsBool("BASKET_SERVICE_ENABLED", true),
			},
			Payment: PaymentServiceConfig{
				Name:     getEnv("PAYMENT_SERVICE_NAME", "payment-service"),
				URLs:     getEnvSlice("PAYMENT_SERVICE_URLS", []string{"http://localhost:8082"}),
				Timeout:  getEnvAsInt("PAYMENT_SERVICE_TIMEOUT", 30),
				Retries:  getEnvAsInt("PAYMENT_SERVICE_RETRIES", 3),
				Enabled:  getEnvAsBool("PAYMENT_SERVICE_ENABLED", true),
			},
		},
		
		CircuitBreaker: CircuitBreakerConfig{
			Enabled:     getEnvAsBool("CIRCUIT_BREAKER_ENABLED", true),
			MaxRequests: uint32(getEnvAsInt("CIRCUIT_BREAKER_MAX_REQUESTS", 10)),
			Interval:    getEnvAsInt("CIRCUIT_BREAKER_INTERVAL", 60),
			Timeout:     getEnvAsInt("CIRCUIT_BREAKER_TIMEOUT", 30),
		},
		
		LoadBalancer: LoadBalancerConfig{
			Strategy: getEnv("LOAD_BALANCER_STRATEGY", "round_robin"),
			Enabled:  getEnvAsBool("LOAD_BALANCER_ENABLED", true),
		},
		
		RateLimit: RateLimitConfig{
			Enabled:  getEnvAsBool("RATE_LIMIT_ENABLED", true),
			Requests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
			Window:   getEnvAsDuration("RATE_LIMIT_WINDOW", "1m"),
			Burst:    getEnvAsInt("RATE_LIMIT_BURST", 10),
		},
		
		Health: HealthConfig{
			Enabled:       getEnvAsBool("HEALTH_CHECK_ENABLED", true),
			CheckInterval: getEnvAsDuration("HEALTH_CHECK_INTERVAL", "30s"),
			Timeout:       getEnvAsDuration("HEALTH_CHECK_TIMEOUT", "5s"),
		},
		
		Metrics: MetricsConfig{
			Enabled: getEnvAsBool("METRICS_ENABLED", true),
			Path:    getEnv("METRICS_PATH", "/metrics"),
		},
	}
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
