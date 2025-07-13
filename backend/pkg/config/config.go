package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	CORS     CORSConfig
	Redis    RedisConfig
	Email    EmailConfig
	Storage  StorageConfig
}

// AppConfig holds application configuration
type AppConfig struct {
	Name        string
	Environment string
	Version     string
	Debug       bool
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret        string
	ExpiresIn     string
	RefreshExpiry string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port    string
	Host    string
	Timeout int
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	From         string
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	UploadPath    string
	MaxUploadSize int64
}

var cfg *Config

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		logger.Debug("No .env file found, using system environment variables")
	}

	// Initialize viper
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg = &Config{
		App: AppConfig{
			Name:        getEnvOrDefault("APP_NAME", "Trecplans"),
			Environment: getEnvOrDefault("APP_ENV", "development"),
			Version:     getEnvOrDefault("APP_VERSION", "1.0.0"),
			Debug:       getBoolEnvOrDefault("APP_DEBUG", true),
		},
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getIntEnvOrDefault("DB_PORT", 5432),
			User:     getEnvOrDefault("DB_USER", "trecplans"),
			Password: getEnvOrDefault("DB_PASSWORD", "password"),
			DBName:   getEnvOrDefault("DB_NAME", "trecplans_dev"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:        getEnvOrDefault("JWT_SECRET", "your-super-secret-jwt-key"),
			ExpiresIn:     getEnvOrDefault("JWT_EXPIRY", "24h"),
			RefreshExpiry: getEnvOrDefault("JWT_REFRESH_EXPIRY", "168h"),
		},
		Server: ServerConfig{
			Port:    getEnvOrDefault("APP_PORT", "8080"),
			Host:    getEnvOrDefault("APP_HOST", ""),
			Timeout: getIntEnvOrDefault("SERVER_TIMEOUT", 30),
		},
		CORS: CORSConfig{
			AllowedOrigins: getStringSliceEnvOrDefault("CORS_ALLOWED_ORIGINS", []string{"http://localhost:5173", "http://localhost:3000"}),
			AllowedMethods: getStringSliceEnvOrDefault("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders: getStringSliceEnvOrDefault("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"}),
		},
		Redis: RedisConfig{
			Host:     getEnvOrDefault("REDIS_HOST", "localhost"),
			Port:     getIntEnvOrDefault("REDIS_PORT", 6379),
			Password: getEnvOrDefault("REDIS_PASSWORD", ""),
			DB:       getIntEnvOrDefault("REDIS_DB", 0),
		},
		Email: EmailConfig{
			SMTPHost:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getIntEnvOrDefault("SMTP_PORT", 587),
			SMTPUser:     getEnvOrDefault("SMTP_USER", ""),
			SMTPPassword: getEnvOrDefault("SMTP_PASSWORD", ""),
			From:         getEnvOrDefault("SMTP_FROM", "noreply@trecplans.com"),
		},
		Storage: StorageConfig{
			UploadPath:    getEnvOrDefault("UPLOAD_PATH", "./uploads"),
			MaxUploadSize: getInt64EnvOrDefault("MAX_UPLOAD_SIZE", 10485760), // 10MB
		},
	}

	// Validate required configuration
	if err := validate(cfg); err != nil {
		return nil, err
	}

	logger.Info("Configuration loaded successfully",
		logger.String("environment", cfg.App.Environment),
		logger.String("app_name", cfg.App.Name),
		logger.String("version", cfg.App.Version),
	)

	return cfg, nil
}

// Get returns the loaded configuration
func Get() *Config {
	if cfg == nil {
		panic("Configuration not loaded. Call Load() first")
	}
	return cfg
}

// validate validates required configuration
func validate(cfg *Config) error {
	// Validate database configuration
	if cfg.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if cfg.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if cfg.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	// Validate JWT configuration
	if cfg.JWT.Secret == "" || cfg.JWT.Secret == "your-super-secret-jwt-key" {
		if cfg.App.Environment == "production" {
			return fmt.Errorf("JWT_SECRET must be set in production")
		}
		logger.Warn("Using default JWT secret. This is insecure for production use")
	}

	// Validate server configuration
	if cfg.Server.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}

	return nil
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getInt64EnvOrDefault(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnvOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getStringSliceEnvOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}