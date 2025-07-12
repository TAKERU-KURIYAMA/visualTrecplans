package workout

import (
	"github.com/gin-gonic/gin"
	"github.com/visualtrecplans/backend/internal/middleware"
	"github.com/visualtrecplans/backend/internal/services"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// RegisterRoutes registers all workout-related routes
func RegisterRoutes(router *gin.RouterGroup, workoutService services.WorkoutService, logger *logger.Logger) {
	workouts := router.Group("/workouts")
	{
		// Authenticated routes
		workouts.Use(middleware.AuthRequired())
		
		// Create workout
		workouts.POST("", CreateWorkout(workoutService, logger))
		
		// Get user workouts with filtering
		workouts.GET("", GetWorkouts(workoutService, logger))
		
		// Get single workout
		workouts.GET("/:id", GetWorkout(workoutService, logger))
		
		// Update workout
		workouts.PUT("/:id", UpdateWorkout(workoutService, logger))
		
		// Delete workout
		workouts.DELETE("/:id", DeleteWorkout(workoutService, logger))
		
		// Get workout statistics
		workouts.GET("/stats", GetWorkoutStats(workoutService, logger))
	}
}