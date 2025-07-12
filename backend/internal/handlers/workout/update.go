package workout

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/visualtrecplans/backend/internal/middleware"
	"github.com/visualtrecplans/backend/internal/services"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// UpdateWorkout handles PUT /api/v1/workouts/:id
func UpdateWorkout(workoutService services.WorkoutService, logger *logger.Logger) gin.HandlerFunc {
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

		// Bind and validate request
		var req UpdateWorkoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("Invalid request body", "error", err, "userID", userID, "workoutID", workoutID)
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid request",
				Code:    "INVALID_REQUEST",
				Details: map[string]interface{}{"validation_error": err.Error()},
			})
			return
		}

		// Update workout
		workout, err := workoutService.Update(c.Request.Context(), workoutID, userID, &req)
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
			if errors.Is(err, services.ErrInvalidMuscleGroup) {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error: "Invalid muscle group",
					Code:  "INVALID_MUSCLE_GROUP",
				})
				return
			}

			logger.Error("Failed to update workout", "error", err, "workoutID", workoutID, "userID", userID)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to update workout",
				Code:  "INTERNAL_ERROR",
			})
			return
		}

		c.JSON(http.StatusOK, toWorkoutResponse(workout))
	}
}

// DeleteWorkout handles DELETE /api/v1/workouts/:id
func DeleteWorkout(workoutService services.WorkoutService, logger *logger.Logger) gin.HandlerFunc {
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

		// Delete workout
		err = workoutService.Delete(c.Request.Context(), workoutID, userID)
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

			logger.Error("Failed to delete workout", "error", err, "workoutID", workoutID, "userID", userID)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to delete workout",
				Code:  "INTERNAL_ERROR",
			})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}