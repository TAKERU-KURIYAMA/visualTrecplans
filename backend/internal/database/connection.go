package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/visualtrecplans/backend/internal/models"
	"github.com/visualtrecplans/backend/pkg/config"
	appLogger "github.com/visualtrecplans/backend/pkg/logger"
)

var db *gorm.DB

// Connect establishes database connection with configuration
func Connect(cfg *config.Config) error {
	dsn := cfg.GetDatabaseDSN()
	
	// Configure GORM logger
	var gormLogger logger.Interface
	switch cfg.App.Environment {
	case "production":
		gormLogger = logger.Default.LogMode(logger.Error)
	case "development":
		gormLogger = logger.Default.LogMode(logger.Info)
	default:
		gormLogger = logger.Default.LogMode(logger.Warn)
	}

	// Open database connection
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	appLogger.Info("Database connection established successfully",
		appLogger.String("host", cfg.Database.Host),
		appLogger.Int("port", cfg.Database.Port),
		appLogger.String("database", cfg.Database.DBName),
	)

	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	if db == nil {
		panic("Database not initialized. Call Connect() first")
	}
	return db
}

// Close closes the database connection
func Close() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying sql.DB: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}

// Migrate runs database migrations
func Migrate() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	appLogger.Info("Running database migrations...")

	// Auto-migrate models
	err := db.AutoMigrate(
		&models.User{},
		&models.Workout{},
		&models.MuscleGroup{},
		&models.Exercise{},
		&models.ExerciseIcon{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	appLogger.Info("Database migrations completed successfully")
	return nil
}

// HealthCheck checks if the database is healthy
func HealthCheck() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// GetStats returns database connection statistics
func GetStats() sql.DBStats {
	if db == nil {
		return sql.DBStats{}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return sql.DBStats{}
	}

	return sqlDB.Stats()
}

// Transaction executes a function within a database transaction
func Transaction(fn func(tx *gorm.DB) error) error {
	return db.Transaction(fn)
}

// IsConnected checks if database connection is active
func IsConnected() bool {
	if db == nil {
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx) == nil
}