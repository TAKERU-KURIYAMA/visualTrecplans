package workout

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/visualtrecplans/backend/internal/middleware"
	"github.com/visualtrecplans/backend/internal/models"
	"github.com/visualtrecplans/backend/internal/services"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// CreateWorkout handles POST /api/v1/workouts
func CreateWorkout(workoutService services.WorkoutService, logger *logger.Logger) gin.HandlerFunc {
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

		// Bind and validate request
		var req CreateWorkoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("Invalid request body", "error", err, "userID", userID)
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid request",
				Code:    "INVALID_REQUEST",
				Details: map[string]interface{}{"validation_error": err.Error()},
			})
			return
		}

		// Set default performed_at if not provided
		if req.PerformedAt.IsZero() {
			req.PerformedAt = time.Now()
		}

		// Create workout
		workout, err := workoutService.Create(c.Request.Context(), userID, &req)
		if err != nil {
			if errors.Is(err, services.ErrInvalidMuscleGroup) {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error: "Invalid muscle group",
					Code:  "INVALID_MUSCLE_GROUP",
				})
				return
			}

			logger.Error("Failed to create workout", "error", err, "userID", userID)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to create workout",
				Code:  "INTERNAL_ERROR",
			})
			return
		}

		// Return response
		c.JSON(http.StatusCreated, toWorkoutResponse(workout))
	}
}

// toWorkoutResponse converts a workout model to API response
func toWorkoutResponse(w *models.Workout) WorkoutResponse {
	return WorkoutResponse{
		ID:           w.ID.String(),
		MuscleGroup:  w.MuscleGroup,
		ExerciseName: w.ExerciseName,
		ExerciseIcon: w.ExerciseIcon,
		WeightKg:     w.WeightKg,
		Reps:         w.Reps,
		Sets:         w.Sets,
		Notes:        w.Notes,
		PerformedAt:  w.PerformedAt,
		CreatedAt:    w.CreatedAt,
		UpdatedAt:    w.UpdatedAt,
	}
}

// toWorkoutResponses converts multiple workout models to API responses
func toWorkoutResponses(workouts []*models.Workout) []WorkoutResponse {
	responses := make([]WorkoutResponse, len(workouts))
	for i, w := range workouts {
		responses[i] = toWorkoutResponse(w)
	}
	return responses
}