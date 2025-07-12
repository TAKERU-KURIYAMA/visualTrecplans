package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Global logger instance
	log        *zap.Logger
	sugar      *zap.SugaredLogger
	atomicLevel zap.AtomicLevel
)

// Field represents a logger field
type Field = zap.Field

// Common field constructors
var (
	String  = zap.String
	Int     = zap.Int
	Int64   = zap.Int64
	Float64 = zap.Float64
	Bool    = zap.Bool
	Error   = zap.Error
	Any     = zap.Any
	Time    = zap.Time
	Duration = zap.Duration
)

func init() {
	// Initialize logger with environment configuration
	logLevel := getEnvOrDefault("LOG_LEVEL", "info")
	logFormat := getEnvOrDefault("LOG_FORMAT", "json")

	// Set log level
	atomicLevel = zap.NewAtomicLevel()
	switch strings.ToLower(logLevel) {
	case "debug":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "warn", "warning":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	default:
		atomicLevel.SetLevel(zapcore.InfoLevel)
	}

	// Configure encoder
	var encoderConfig zapcore.EncoderConfig
	if logFormat == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Build encoder
	var encoder zapcore.Encoder
	if logFormat == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Build logger
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		atomicLevel,
	)

	log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar = log.Sugar()
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Debug logs a message at debug level
func Debug(msg string, fields ...Field) {
	log.Debug(msg, fields...)
}

// Info logs a message at info level
func Info(msg string, fields ...Field) {
	log.Info(msg, fields...)
}

// Warn logs a message at warn level
func Warn(msg string, fields ...Field) {
	log.Warn(msg, fields...)
}

// Error logs a message at error level
func Error(msg string, fields ...Field) {
	log.Error(msg, fields...)
}

// Fatal logs a message at fatal level and exits
func Fatal(msg string, fields ...Field) {
	log.Fatal(msg, fields...)
}

// Debugf logs a formatted message at debug level
func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

// Infof logs a formatted message at info level
func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

// Warnf logs a formatted message at warn level
func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

// Errorf logs a formatted message at error level
func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

// Fatalf logs a formatted message at fatal level and exits
func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}

// WithContext creates a new logger with additional context fields
func WithContext(fields ...Field) *zap.Logger {
	return log.With(fields...)
}

// WithSugaredContext creates a new sugared logger with additional context
func WithSugaredContext(keysAndValues ...interface{}) *zap.SugaredLogger {
	return sugar.With(keysAndValues...)
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	return log
}

// GetSugaredLogger returns the global sugared logger instance
func GetSugaredLogger() *zap.SugaredLogger {
	return sugar
}

// SetLevel dynamically changes the log level
func SetLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "warn", "warning":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	}
}

// Sync flushes any buffered log entries
func Sync() error {
	return log.Sync()
}