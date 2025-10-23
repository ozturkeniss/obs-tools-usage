package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds the configuration for the product service
type Config struct {
	Port        string
	Environment string
	LogLevel    string
	LogFormat   string
	LogOutput   string
	LogDir      string
	LogFile     string
	LogRotation LogRotationConfig
	Database    DatabaseConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LogRotationConfig holds log rotation configuration
type LogRotationConfig struct {
	Enabled   bool
	MaxSize   int    // Maximum size in MB
	MaxAge    int    // Maximum age in days
	MaxBackups int   // Maximum number of backup files
	Compress  bool   // Whether to compress old log files
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	environment := getEnv("ENVIRONMENT", "development")
	
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: environment,
		LogLevel:    getLogLevelFromEnv(environment),
		LogFormat:   getLogFormatFromEnv(environment),
		LogOutput:   getLogOutputFromEnv(environment),
		LogDir:      getEnv("LOG_DIR", "./logs"),
		LogFile:     getEnv("LOG_FILE", "product-service.log"),
		LogRotation: LogRotationConfig{
			Enabled:    true,
			MaxSize:    100,
			MaxAge:     30,
			MaxBackups: 10,
			Compress:   true,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "obs_tools"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

// GetPort returns the port as an integer
func (c *Config) GetPort() int {
	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return 8080
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


// GetDatabaseURL returns the complete database connection URL
func (c *Config) GetDatabaseURL() string {
	return "postgres://" + c.Database.User + ":" + c.Database.Password + "@" + c.Database.Host + ":" + c.Database.Port + "/" + c.Database.DBName + "?sslmode=" + c.Database.SSLMode
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
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

