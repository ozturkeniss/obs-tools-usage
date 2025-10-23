package config

import (
	"os"
	"strconv"
)

// Config holds the configuration for the basket service
type Config struct {
	Port        string
	Environment string
	LogLevel    string
	LogFormat   string
	LogOutput   string
	LogDir      string
	LogFile     string
	Redis       RedisConfig
	Product     ProductConfig
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

// ProductConfig holds product service configuration
type ProductConfig struct {
	ServiceURL string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	environment := getEnv("ENVIRONMENT", "development")
	
	return &Config{
		Port:        getEnv("PORT", "8081"),
		Environment: environment,
		LogLevel:    getLogLevelFromEnv(environment),
		LogFormat:   getLogFormatFromEnv(environment),
		LogOutput:   getLogOutputFromEnv(environment),
		LogDir:      getEnv("LOG_DIR", "./logs"),
		LogFile:     getEnv("LOG_FILE", "basket-service.log"),
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			PoolSize: getEnvAsInt("REDIS_POOL_SIZE", 10),
		},
		Product: ProductConfig{
			ServiceURL: getEnv("PRODUCT_SERVICE_URL", "localhost:50050"),
		},
	}
}

// GetPort returns the port as an integer
func (c *Config) GetPort() int {
	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return 8081
	}
	return port
}

// IsDevelopment returns true if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
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

// getLogLevelFromEnv determines log level from environment
func getLogLevelFromEnv(environment string) string {
	// First check LOG_LEVEL environment variable
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		return logLevel
	}
	
	// Default log levels based on environment
	switch environment {
	case "production":
		return "warn"
	case "staging":
		return "info"
	case "development", "dev":
		return "debug"
	default:
		return "info"
	}
}

// getLogFormatFromEnv determines log format from environment
func getLogFormatFromEnv(environment string) string {
	// First check LOG_FORMAT environment variable
	if logFormat := os.Getenv("LOG_FORMAT"); logFormat != "" {
		return logFormat
	}
	
	// Default formats based on environment
	switch environment {
	case "production":
		return "json"
	case "staging", "development", "dev":
		return "text"
	default:
		return "text"
	}
}

// getLogOutputFromEnv determines log output from environment
func getLogOutputFromEnv(environment string) string {
	// First check LOG_OUTPUT environment variable
	if logOutput := os.Getenv("LOG_OUTPUT"); logOutput != "" {
		return logOutput
	}
	
	// Default outputs based on environment
	switch environment {
	case "production":
		return "file"
	case "staging":
		return "both"
	case "development", "dev":
		return "console"
	default:
		return "console"
	}
}
