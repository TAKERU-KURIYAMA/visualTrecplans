package logger_test

import (
	"os"
	"testing"

	"github.com/visualtrecplans/backend/pkg/logger"
	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	// Test basic logging functions
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	// Test formatted logging
	logger.Debugf("Debug: %s", "formatted message")
	logger.Infof("Info: %s", "formatted message")
	logger.Warnf("Warning: %s", "formatted message")
	logger.Errorf("Error: %s", "formatted message")

	// Test structured logging with fields
	logger.Info("User login",
		logger.String("username", "john_doe"),
		logger.Int("user_id", 123),
		logger.Bool("success", true),
	)

	// Test context logging
	contextLogger := logger.WithContext(
		logger.String("request_id", "abc123"),
		logger.String("service", "auth"),
	)
	contextLogger.Info("Processing request")

	// Test sugared context logging
	sugaredLogger := logger.WithSugaredContext(
		"request_id", "xyz789",
		"service", "api",
	)
	sugaredLogger.Infof("API call completed in %dms", 150)

	// Test error field
	err := os.ErrNotExist
	logger.Error("File operation failed", logger.ErrorField(err))

	// Ensure logs are flushed
	if err := logger.Sync(); err != nil {
		t.Logf("Failed to sync logger: %v", err)
	}
}

func TestLogLevelFromEnv(t *testing.T) {
	// This test demonstrates that the logger respects LOG_LEVEL env var
	// Set LOG_LEVEL=debug before running to see debug messages
	logger.Debug("This message only appears if LOG_LEVEL=debug")
	logger.Info("This message appears at info level and above")
}

func ExampleWithContext() {
	// Create a logger with context for a specific request
	requestLogger := logger.WithContext(
		logger.String("request_id", "req-12345"),
		logger.String("user_id", "user-789"),
		logger.String("endpoint", "/api/v1/users"),
	)

	// All logs from this logger will include the context fields
	requestLogger.Info("Request started")
	requestLogger.Debug("Validating input")
	requestLogger.Info("Request completed successfully")
}