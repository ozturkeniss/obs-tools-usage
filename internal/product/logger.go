package product

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	
	// Get environment
	environment := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if environment == "" {
		environment = "development"
	}
	
	// Set log level based on environment and LOG_LEVEL variable
	logLevel := getLogLevel(environment)
	Logger.SetLevel(logLevel)
	
	// Set formatter based on environment
	if environment == "production" {
		// Production: JSON format for log aggregation
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		// Development/Staging: Pretty format for readability
		Logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceColors:     true,
		})
	}
	
	// Set output to stdout
	Logger.SetOutput(os.Stdout)
	
	// Log configuration
	Logger.WithFields(logrus.Fields{
		"environment": environment,
		"log_level":   logLevel.String(),
		"formatter":   getFormatterName(environment),
	}).Info("Logger initialized")
}

// getLogLevel determines log level based on environment and LOG_LEVEL variable
func getLogLevel(environment string) logrus.Level {
	// First check LOG_LEVEL environment variable
	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if logLevel != "" {
		switch logLevel {
		case "debug":
			return logrus.DebugLevel
		case "info":
			return logrus.InfoLevel
		case "warn", "warning":
			return logrus.WarnLevel
		case "error":
			return logrus.ErrorLevel
		case "fatal":
			return logrus.FatalLevel
		case "panic":
			return logrus.PanicLevel
		}
	}
	
	// Default log levels based on environment
	switch environment {
	case "production":
		return logrus.WarnLevel // Only warnings and errors in production
	case "staging":
		return logrus.InfoLevel // Info and above in staging
	case "development", "dev":
		return logrus.DebugLevel // All logs in development
	default:
		return logrus.InfoLevel // Default to info level
	}
}

// getFormatterName returns the formatter name for logging
func getFormatterName(environment string) string {
	if environment == "production" {
		return "json"
	}
	return "text"
}

// GetLogger returns the configured logger instance
func GetLogger() *logrus.Logger {
	return Logger
}
