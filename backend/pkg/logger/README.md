# Logger Package

A structured logging package for the Visual Training Record Plans backend, built on top of [uber-go/zap](https://github.com/uber-go/zap).

## Features

- Structured logging with zap
- Multiple log levels (debug, info, warn, error, fatal)
- Context fields support for request tracing
- Environment-based configuration
- Both structured and formatted logging APIs
- Global logger instance with helper functions

## Configuration

The logger is configured using environment variables:

- `LOG_LEVEL`: Sets the minimum log level (default: "info")
  - Supported values: "debug", "info", "warn", "error", "fatal"
- `LOG_FORMAT`: Sets the output format (default: "json")
  - "json": Production-ready JSON format
  - "console": Human-readable format with colors (for development)

## Usage

### Basic Logging

```go
import "github.com/trecplans/backend/pkg/logger"

// Simple messages
logger.Debug("Debug message")
logger.Info("Info message")
logger.Warn("Warning message")
logger.Error("Error message")

// Formatted messages
logger.Infof("User %s logged in", username)
logger.Errorf("Failed to connect to database: %v", err)
```

### Structured Logging

```go
// Log with structured fields
logger.Info("User action",
    logger.String("user_id", "123"),
    logger.String("action", "login"),
    logger.Bool("success", true),
    logger.Int("duration_ms", 150),
)

// Log with error
logger.Error("Operation failed",
    logger.String("operation", "create_user"),
    logger.ErrorField(err),
)
```

### Context Logging

```go
// Create a logger with persistent fields
requestLogger := logger.WithContext(
    logger.String("request_id", "req-12345"),
    logger.String("user_id", "user-789"),
)

// All logs from this logger include the context fields
requestLogger.Info("Processing request")
requestLogger.Debug("Validating input")
```

### Sugared Logger

```go
// Get sugared logger for more convenient API
sugar := logger.GetSugaredLogger()
sugar.Infow("Failed to fetch URL",
    "url", url,
    "attempt", 3,
    "backoff", time.Second,
)

// Create sugared logger with context
contextSugar := logger.WithSugaredContext(
    "request_id", "xyz789",
    "service", "api",
)
contextSugar.Infof("API call completed in %dms", 150)
```

## Field Helpers

The package provides convenience functions for common field types:

```go
logger.String("key", "value")      // string field
logger.Int("count", 42)            // int field
logger.Int64("id", 12345)          // int64 field
logger.Float64("rate", 0.95)       // float64 field
logger.Bool("enabled", true)       // bool field
logger.ErrorField(err)             // error field
logger.Field("data", customStruct) // any type field
```

## Best Practices

1. **Use structured logging**: Prefer structured fields over formatted strings for better searchability
2. **Add context**: Use `WithContext` to add request IDs and other tracing information
3. **Log levels**: Use appropriate log levels (debug for development, info for general information, warn for warnings, error for errors)
4. **Performance**: The structured logger is faster than formatted logging
5. **Defer Sync**: In your main function, defer `logger.Sync()` to ensure all logs are flushed

## Example in HTTP Handler

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    requestID := generateRequestID()
    
    // Create request-scoped logger
    reqLogger := logger.WithContext(
        logger.String("request_id", requestID),
        logger.String("method", r.Method),
        logger.String("path", r.URL.Path),
        logger.String("remote_addr", r.RemoteAddr),
    )
    
    reqLogger.Info("Request started")
    
    // Use the logger throughout the request
    if err := processRequest(r); err != nil {
        reqLogger.Error("Request failed", logger.ErrorField(err))
        http.Error(w, "Internal Server Error", 500)
        return
    }
    
    reqLogger.Info("Request completed successfully")
}
```