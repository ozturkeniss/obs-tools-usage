package product

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
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
	
	// Set output based on environment and configuration
	output := getLogOutput(environment)
	Logger.SetOutput(output)
	
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

// getLogOutput determines where to write logs based on environment and configuration
func getLogOutput(environment string) io.Writer {
	// Check LOG_OUTPUT environment variable
	logOutput := strings.ToLower(os.Getenv("LOG_OUTPUT"))
	
	switch logOutput {
	case "file":
		return getFileOutput()
	case "both", "file+console":
		return io.MultiWriter(os.Stdout, getFileOutput())
	case "console", "stdout":
		return os.Stdout
	default:
		// Default behavior based on environment
		switch environment {
		case "production":
			return getFileOutput()
		case "staging", "development", "dev":
			return os.Stdout
		default:
			return os.Stdout
		}
	}
}

// getFileOutput creates a file output for logging with rotation
func getFileOutput() io.Writer {
	// Get log directory from environment or use default
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "./logs"
	}
	
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// If we can't create the directory, fall back to stdout
		return os.Stdout
	}
	
	// Get log filename from environment or use default
	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		logFile = "product-service.log"
	}
	
	// Create full path
	logPath := filepath.Join(logDir, logFile)
	
	// Check if log rotation is enabled
	rotationEnabled := getEnv("LOG_ROTATION_ENABLED", "true")
	if strings.ToLower(rotationEnabled) == "false" {
		// Simple file output without rotation
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return os.Stdout
		}
		return file
	}
	
	// Configure log rotation
	rotation := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    getLogRotationMaxSize(),
		MaxAge:     getLogRotationMaxAge(),
		MaxBackups: getLogRotationMaxBackups(),
		Compress:   getLogRotationCompress(),
		LocalTime:  true,
	}
	
	return rotation
}


// GetLogger returns the configured logger instance
func GetLogger() *logrus.Logger {
	return Logger
}
