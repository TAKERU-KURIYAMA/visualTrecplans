package master

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visualtrecplans/backend/internal/middleware"
	"github.com/visualtrecplans/backend/internal/models"
	"github.com/visualtrecplans/backend/internal/services"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// GetMuscleGroups handles GET /api/v1/muscle-groups
func GetMuscleGroups(masterService services.MasterService, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.DefaultQuery("lang", "ja")
		category := c.Query("category")

		var muscleGroups []*models.MuscleGroup
		var err error

		if category != "" {
			muscleGroups, err = masterService.GetMuscleGroupsByCategory(c.Request.Context(), category, lang)
		} else {
			muscleGroups, err = masterService.GetAllMuscleGroups(c.Request.Context(), lang)
		}

		if err != nil {
			logger.Error("Failed to get muscle groups", "error", err, "lang", lang, "category", category)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to fetch muscle groups",
				Code:  "INTERNAL_ERROR",
			})
			return
		}

		response := MuscleGroupListResponse{
			Data: toMuscleGroupResponses(muscleGroups),
		}

		// Long-term cache for master data
		c.Header("Cache-Control", "public, max-age=3600")
		c.JSON(http.StatusOK, response)
	}
}

// GetExercises handles GET /api/v1/exercises
func GetExercises(masterService services.MasterService, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.DefaultQuery("lang", "ja")
		muscleGroupCode := c.Query("muscle_group")

		exercises, err := masterService.GetAllExercises(c.Request.Context(), muscleGroupCode, lang)
		if err != nil {
			logger.Error("Failed to get exercises", "error", err, "lang", lang, "muscleGroupCode", muscleGroupCode)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to fetch exercises",
				Code:  "INTERNAL_ERROR",
			})
			return
		}

		response := ExerciseListResponse{
			Data: toExerciseResponses(exercises),
		}

		// Long-term cache for master data
		c.Header("Cache-Control", "public, max-age=3600")
		c.JSON(http.StatusOK, response)
	}
}

// GetCustomExercises handles GET /api/v1/exercises/custom
func GetCustomExercises(masterService services.MasterService, logger *logger.Logger) gin.HandlerFunc {
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

		exercises, err := masterService.GetUserCustomExercises(c.Request.Context(), userID)
		if err != nil {
			logger.Error("Failed to get custom exercises", "error", err, "userID", userID)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to fetch custom exercises",
				Code:  "INTERNAL_ERROR",
			})
			return
		}

		response := ExerciseListResponse{
			Data: toExerciseResponses(exercises),
		}

		c.JSON(http.StatusOK, response)
	}
}

// CreateCustomExercise handles POST /api/v1/exercises/custom
func CreateCustomExercise(masterService services.MasterService, logger *logger.Logger) gin.HandlerFunc {
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
		var req CreateCustomExerciseRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("Invalid request body", "error", err, "userID", userID)
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid request",
				Code:    "INVALID_REQUEST",
				Details: map[string]interface{}{"validation_error": err.Error()},
			})
			return
		}

		// Create custom exercise
		exercise, err := masterService.CreateCustomExercise(c.Request.Context(), userID, &req)
		if err != nil {
			if errors.Is(err, services.ErrInvalidMuscleGroup) {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error: "Invalid muscle group",
					Code:  "INVALID_MUSCLE_GROUP",
				})
				return
			}

			logger.Error("Failed to create custom exercise", "error", err, "userID", userID)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to create custom exercise",
				Code:  "INTERNAL_ERROR",
			})
			return
		}

		c.JSON(http.StatusCreated, toExerciseResponse(exercise))
	}
}

// GetExerciseIcons handles GET /api/v1/exercise-icons
func GetExerciseIcons(masterService services.MasterService, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		category := c.Query("category")

		icons, err := masterService.GetAllExerciseIcons(c.Request.Context(), category)
		if err != nil {
			logger.Error("Failed to get exercise icons", "error", err, "category", category)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to fetch exercise icons",
				Code:  "INTERNAL_ERROR",
			})
			return
		}

		response := ExerciseIconListResponse{
			Data: toExerciseIconResponses(icons),
		}

		// Very long-term cache for icons
		c.Header("Cache-Control", "public, max-age=86400")
		c.JSON(http.StatusOK, response)
	}
}

// Helper functions to convert models to responses

func toMuscleGroupResponse(mg *models.MuscleGroup) MuscleGroupResponse {
	return MuscleGroupResponse{
		Code:      mg.Code,
		NameJa:    mg.NameJa,
		NameEn:    mg.NameEn,
		Category:  mg.Category,
		ColorCode: mg.ColorCode,
		SortOrder: mg.SortOrder,
	}
}

func toMuscleGroupResponses(muscleGroups []*models.MuscleGroup) []MuscleGroupResponse {
	responses := make([]MuscleGroupResponse, len(muscleGroups))
	for i, mg := range muscleGroups {
		responses[i] = toMuscleGroupResponse(mg)
	}
	return responses
}

func toExerciseResponse(exercise *models.Exercise) ExerciseResponse {
	return ExerciseResponse{
		ID:              exercise.ID,
		MuscleGroupCode: exercise.MuscleGroupCode,
		NameJa:          exercise.NameJa,
		NameEn:          exercise.NameEn,
		IconName:        exercise.IconName,
		IsCustom:        exercise.IsCustom,
		SortOrder:       exercise.SortOrder,
	}
}

func toExerciseResponses(exercises []*models.Exercise) []ExerciseResponse {
	responses := make([]ExerciseResponse, len(exercises))
	for i, exercise := range exercises {
		responses[i] = toExerciseResponse(exercise)
	}
	return responses
}

func toExerciseIconResponse(icon *models.ExerciseIcon) ExerciseIconResponse {
	return ExerciseIconResponse{
		Name:     icon.Name,
		SVGPath:  icon.SVGPath,
		Category: icon.Category,
	}
}

func toExerciseIconResponses(icons []*models.ExerciseIcon) []ExerciseIconResponse {
	responses := make([]ExerciseIconResponse, len(icons))
	for i, icon := range icons {
		responses[i] = toExerciseIconResponse(icon)
	}
	return responses
}