package product

import (
	"os"
	"strconv"
)

// Config holds the configuration for the product service
type Config struct {
	Port        string
	Environment string
	LogLevel    string
	LogFormat   string
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

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	environment := getEnv("ENVIRONMENT", "development")
	
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: environment,
		LogLevel:    getLogLevelFromEnv(environment),
		LogFormat:   getLogFormatFromEnv(environment),
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

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDatabaseURL returns the complete database connection URL
func (c *Config) GetDatabaseURL() string {
	return "postgres://" + c.Database.User + ":" + c.Database.Password + "@" + c.Database.Host + ":" + c.Database.Port + "/" + c.Database.DBName + "?sslmode=" + c.Database.SSLMode
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
