package main

import (
	"errors"
	"os"
	"time"

	"github.com/visualtrecplans/backend/pkg/logger"
)

func main() {
	// The logger is automatically initialized when the package is imported
	// It reads LOG_LEVEL and LOG_FORMAT from environment variables

	// Basic logging
	logger.Info("Application started")
	logger.Debug("Debug mode enabled")

	// Structured logging with fields
	logger.Info("Server configuration",
		logger.String("host", "localhost"),
		logger.Int("port", 8080),
		logger.Bool("tls_enabled", false),
	)

	// Simulate user login
	simulateUserLogin("john.doe@example.com", 12345)

	// Simulate an error
	if err := simulateOperation(); err != nil {
		logger.Error("Operation failed",
			logger.ErrorField(err),
			logger.String("operation", "data_sync"),
			logger.Int("retry_count", 3),
		)
	}

	// Using formatted logging (sugared)
	logger.Infof("Processing %d items in queue", 42)
	logger.Warnf("Cache miss rate: %.2f%%", 15.5)

	// Context logging for a request
	handleRequest("req-123", "user-456")

	// Ensure all logs are flushed before exit
	if err := logger.Sync(); err != nil {
		// Handle sync error (usually can be ignored)
		os.Exit(1)
	}

	logger.Info("Application shutdown complete")
}

func simulateUserLogin(email string, userID int) {
	start := time.Now()
	
	// Create a logger with user context
	userLogger := logger.WithContext(
		logger.String("email", email),
		logger.Int("user_id", userID),
		logger.String("action", "login"),
	)

	userLogger.Debug("Validating credentials")
	userLogger.Info("User authenticated successfully")
	userLogger.Info("Login completed",
		logger.Duration("duration", time.Since(start)),
		logger.String("ip_address", "192.168.1.100"),
	)
}

func simulateOperation() error {
	logger.Debug("Starting operation")
	// Simulate an error
	return errors.New("connection timeout")
}

func handleRequest(requestID, userID string) {
	// Create request-scoped logger
	reqLogger := logger.WithContext(
		logger.String("request_id", requestID),
		logger.String("user_id", userID),
		logger.String("endpoint", "/api/v1/training-records"),
		logger.String("method", "POST"),
	)

	reqLogger.Info("Request received")
	reqLogger.Debug("Parsing request body")
	
	// Simulate processing
	time.Sleep(100 * time.Millisecond)
	
	reqLogger.Info("Request processed successfully",
		logger.Int("status_code", 200),
		logger.Int64("response_time_ms", 105),
		logger.Int("response_size_bytes", 2048),
	)
}