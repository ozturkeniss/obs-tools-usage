package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the configuration for the notification service
type Config struct {
	// Server configuration
	Port         string
	Environment  string
	
	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	
	// Kafka configuration
	KafkaBrokers string
	
	// Logging configuration
	LogLevel  string
	LogFormat string
	LogOutput string
	
	// Notification configuration
	DefaultRetryAttempts int
	NotificationTTL      time.Duration
	CleanupInterval      time.Duration
	
	// Rate limiting
	RateLimitEnabled bool
	RateLimitRPS     int
	
	// Metrics configuration
	MetricsEnabled bool
	MetricsPath    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		// Server configuration
		Port:        getEnv("PORT", "8084"),
		Environment: getEnv("ENVIRONMENT", "development"),
		
		// Database configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "notification_service"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
		
		// Kafka configuration
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
		
		// Logging configuration
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
		LogOutput: getEnv("LOG_OUTPUT", "console"),
		
		// Notification configuration
		DefaultRetryAttempts: getEnvAsInt("DEFAULT_RETRY_ATTEMPTS", 3),
		NotificationTTL:      getEnvAsDuration("NOTIFICATION_TTL", 24*time.Hour),
		CleanupInterval:      getEnvAsDuration("CLEANUP_INTERVAL", 1*time.Hour),
		
		// Rate limiting
		RateLimitEnabled: getEnvAsBool("RATE_LIMIT_ENABLED", true),
		RateLimitRPS:     getEnvAsInt("RATE_LIMIT_RPS", 100),
		
		// Metrics configuration
		MetricsEnabled: getEnvAsBool("METRICS_ENABLED", true),
		MetricsPath:    getEnv("METRICS_PATH", "/metrics"),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool gets an environment variable as boolean with a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvAsDuration gets an environment variable as duration with a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
