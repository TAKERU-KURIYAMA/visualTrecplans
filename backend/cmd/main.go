package main

import (
	"github.com/gin-gonic/gin"
	"github.com/visualtrecplans/backend/internal/database"
	"github.com/visualtrecplans/backend/internal/handlers/auth"
	"github.com/visualtrecplans/backend/internal/middleware"
	"github.com/visualtrecplans/backend/internal/services"
	"github.com/visualtrecplans/backend/pkg/config"
	"github.com/visualtrecplans/backend/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", logger.Error(err))
	}

	// Initialize logger
	logger.Info("Starting VisualTrecplans API server",
		logger.String("version", cfg.App.Version),
		logger.String("environment", cfg.App.Environment),
	)

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		logger.Fatal("Failed to connect to database", logger.Error(err))
	}
	defer func() {
		if err := database.Close(); err != nil {
			logger.Error("Failed to close database connection", logger.Error(err))
		}
	}()

	// Run database migrations
	if err := database.Migrate(); err != nil {
		logger.Fatal("Failed to run database migrations", logger.Error(err))
	}

	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// Add security middlewares
	r.Use(middleware.SecurityHeaders(cfg))
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit())
	r.Use(middleware.BruteForceProtection())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		// Check database health
		dbHealth := "ok"
		if err := database.HealthCheck(); err != nil {
			dbHealth = "error"
			logger.Error("Database health check failed", logger.Error(err))
		}

		c.JSON(200, gin.H{
			"status":    "ok",
			"message":   "VisualTrecplans API is running",
			"version":   cfg.App.Version,
			"database":  dbHealth,
		})
	})

	// Initialize services
	authService := services.NewAuthService()
	jwtService := services.NewJWTService(cfg)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// Authentication routes
		auth.RegisterRoutes(v1, authService, jwtService)
	}

	// Start server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	logger.Info("Server starting", 
		logger.String("address", addr),
		logger.Bool("debug", cfg.App.Debug),
	)
	
	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server", logger.Error(err))
	}
}