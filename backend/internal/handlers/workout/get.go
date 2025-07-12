package workout

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/visualtrecplans/backend/internal/middleware"
	"github.com/visualtrecplans/backend/internal/services"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// GetWorkouts handles GET /api/v1/workouts
func GetWorkouts(workoutService services.WorkoutService, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, err := middleware.GetUserUUID(c)
		if err != nil {
			logger.Error("Failed to get user ID from context", "error", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Unauthorized",
				Code:  "AUTH_REQUIRED",
			})
			return
		}

		// Bind query parameters
		var filter WorkoutFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			logger.Warn("Invalid query parameters", "error", err, "userID", userID)
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid query parameters",
				Code:    "INVALID_QUERY",
				Details: map[string]interface{}{"validation_error": err.Error()},
			})
			return
		}

		// Set defaults
		if filter.Page <= 0 {
			filter.Page = 1
		}
		if filter.PerPage <= 0 {
			filter.PerPage = 20
		}
		if filter.PerPage > 100 {
			filter.PerPage = 100
		}
		if filter.OrderBy == "" {
			filter.OrderBy = "performed_at"
		}
		if filter.Order == "" {
			filter.Order = "desc"
		}

		// Get workouts
		workouts, total, err := workoutService.GetUserWorkouts(c.Request.Context(), userID, &filter)
		if err != nil {
			logger.Error("Failed to get workouts", "error", err, "userID", userID)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to get workouts",
				Code:  "INTERNAL_ERROR",
			})
			return
		}

		// Calculate pagination
		totalPages := int(math.Ceil(float64(total) / float64(filter.PerPage)))

		// Convert to response
		response := WorkoutListResponse{
			Workouts:   toWorkoutResponses(workouts),
			Total:      total,
			Page:       filter.Page,
			PerPage:    filter.PerPage,
			TotalPages: totalPages,
		}

		c.JSON(http.StatusOK, response)
	}
}

// GetWorkout handles GET /api/v1/workouts/:id
func GetWorkout(workoutService services.WorkoutService, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, err := middleware.GetUserUUID(c)
		if err != nil {
			logger.Error("Failed to get user ID from context", "error", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Unauthorized",
				Code:  "AUTH_REQUIRED",
			})
			return
		}

		// Parse workout ID
		workoutIDStr := c.Param("id")
		workoutID, err := uuid.Parse(workoutIDStr)
		if err != nil {
			logger.Warn("Invalid workout ID", "workoutID", workoutIDStr, "error", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid workout ID",
				Code:  "INVALID_ID",
			})
			return
		}

		// Get workout
		workout, err := workoutService.GetByID(c.Request.Context(), workoutID, userID)
		if err != nil {
			if errors.Is(err, services.ErrWorkoutNotFound) {
				c.JSON(http.StatusNotFound, ErrorResponse{
					Error: "Workout not found",
					Code:  "NOT_FOUND",
				})
				return
			}
			if errors.Is(err, services.ErrUnauthorized) {
				c.JSON(http.StatusForbidden, ErrorResponse{
					Error: "Access denied",
					Code:  "ACCESS_DENIED",
				})
				return
			}

			logger.Error("Failed to get workout", "error", err, "workoutID", workoutID, "userID", userID)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to get workout",
				Code:  "INTERNAL_ERROR",
			})
			return
		}

		c.JSON(http.StatusOK, toWorkoutResponse(workout))
	}
}

// GetWorkoutStats handles GET /api/v1/workouts/stats
func GetWorkoutStats(workoutService services.WorkoutService, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, err := middleware.GetUserUUID(c)
		if err != nil {
			logger.Error("Failed to get user ID from context", "error", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Unauthorized",
				Code:  "AUTH_REQUIRED",
			})
			return
		}

		// Get period from query (default: "month")
		period := c.DefaultQuery("period", "month")
		validPeriods := map[string]bool{
			"week":  true,
			"month": true,
			"year":  true,
		}

		if !validPeriods[period] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid period. Must be one of: week, month, year",
				Code:  "INVALID_PERIOD",
			})
			return
		}

		// Get stats
		stats, err := workoutService.GetStats(c.Request.Context(), userID, period)
		if err != nil {
			logger.Error("Failed to get workout stats", "error", err, "userID", userID, "period", period)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to get workout stats",
				Code:  "INTERNAL_ERROR",
			})
			return
		}

		// Convert to response format
		response := WorkoutStatsResponse{
			TotalWorkouts:     stats.TotalWorkouts,
			TotalSets:         stats.TotalSets,
			TotalReps:         stats.TotalReps,
			TotalWeightLifted: stats.TotalWeightLifted,
			WorkoutsByMuscle:  stats.WorkoutsByMuscle,
			MostUsedExercises: make([]ExerciseCount, len(stats.MostUsedExercises)),
			RecentWorkouts:    toWorkoutResponses(stats.RecentWorkouts),
			WeeklyProgress:    make([]WeeklyProgress, len(stats.WeeklyProgress)),
		}

		// Convert exercise counts
		for i, ec := range stats.MostUsedExercises {
			response.MostUsedExercises[i] = ExerciseCount{
				ExerciseName: ec.ExerciseName,
				Count:        ec.Count,
			}
		}

		// Convert weekly progress
		for i, wp := range stats.WeeklyProgress {
			response.WeeklyProgress[i] = WeeklyProgress{
				Week:         wp.Week,
				WorkoutCount: wp.WorkoutCount,
				TotalSets:    wp.TotalSets,
				TotalReps:    wp.TotalReps,
				TotalWeight:  wp.TotalWeight,
			}
		}

		c.JSON(http.StatusOK, response)
	}
}