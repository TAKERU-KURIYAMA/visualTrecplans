package master

import (
	"github.com/gin-gonic/gin"
	"github.com/visualtrecplans/backend/internal/middleware"
	"github.com/visualtrecplans/backend/internal/services"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// RegisterRoutes registers all master data-related routes
func RegisterRoutes(router *gin.RouterGroup, masterService services.MasterService, logger *logger.Logger) {
	// Public routes (no authentication required)
	router.GET("/muscle-groups", GetMuscleGroups(masterService, logger))
	router.GET("/exercises", GetExercises(masterService, logger))
	router.GET("/exercise-icons", GetExerciseIcons(masterService, logger))

	// Authenticated routes for custom exercises
	authenticated := router.Group("/exercises")
	authenticated.Use(middleware.AuthRequired())
	{
		authenticated.GET("/custom", GetCustomExercises(masterService, logger))
		authenticated.POST("/custom", CreateCustomExercise(masterService, logger))
	}
}