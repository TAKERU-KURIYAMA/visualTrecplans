package main

import (
	"github.com/gin-gonic/gin"
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

	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "VisualTrecplans API is running",
			"version": cfg.App.Version,
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
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