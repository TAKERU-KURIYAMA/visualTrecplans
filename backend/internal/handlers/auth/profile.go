package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/visualtrecplans/backend/internal/middleware"
	"github.com/visualtrecplans/backend/internal/services"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// GetProfileHandler handles fetching user profile
func GetProfileHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := middleware.GetUserID(c)
		if err != nil {
			logger.Error("Failed to get user ID from context", logger.Error(err))
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Authentication required",
				Message: "User not authenticated",
				Code:    "AUTH_REQUIRED",
			})
			return
		}

		// Get user from database
		user, err := authService.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			logger.Error("Failed to get user profile", 
				logger.String("user_id", userID),
				logger.Error(err),
			)

			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "User not found",
				Message: "Unable to retrieve user profile",
				Code:    "USER_NOT_FOUND",
			})
			return
		}

		logger.Info("Profile retrieved successfully", 
			logger.String("user_id", user.ID.String()),
		)

		// Return user profile
		response := user.ToResponse()
		c.JSON(http.StatusOK, response)
	}
}

// UpdateProfileHandler handles updating user profile
func UpdateProfileHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := middleware.GetUserID(c)
		if err != nil {
			logger.Error("Failed to get user ID from context", logger.Error(err))
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Authentication required",
				Message: "User not authenticated",
				Code:    "AUTH_REQUIRED",
			})
			return
		}

		var req UpdateProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("Invalid profile update request", 
				logger.String("user_id", userID),
				logger.Error(err),
			)

			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid request",
				Message: "Please provide valid profile data",
				Code:    "INVALID_REQUEST",
			})
			return
		}

		// Get current user
		user, err := authService.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			logger.Error("Failed to get user for profile update", 
				logger.String("user_id", userID),
				logger.Error(err),
			)

			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "User not found",
				Message: "Unable to update profile",
				Code:    "USER_NOT_FOUND",
			})
			return
		}

		// Update fields
		if req.FirstName != nil {
			user.FirstName = req.FirstName
		}
		if req.LastName != nil {
			user.LastName = req.LastName
		}

		// Save updated user
		if err := authService.UpdateUser(c.Request.Context(), user); err != nil {
			logger.Error("Failed to update user profile", 
				logger.String("user_id", userID),
				logger.Error(err),
			)

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Update failed",
				Message: "Unable to update profile",
				Code:    "UPDATE_ERROR",
			})
			return
		}

		logger.Info("Profile updated successfully", 
			logger.String("user_id", user.ID.String()),
		)

		// Return updated user profile
		response := user.ToResponse()
		c.JSON(http.StatusOK, response)
	}
}

// ChangePasswordHandler handles password changes
func ChangePasswordHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := middleware.GetUserID(c)
		if err != nil {
			logger.Error("Failed to get user ID from context", logger.Error(err))
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Authentication required",
				Message: "User not authenticated",
				Code:    "AUTH_REQUIRED",
			})
			return
		}

		var req ChangePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("Invalid password change request", 
				logger.String("user_id", userID),
				logger.Error(err),
			)

			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid request",
				Message: "Please provide valid password data",
				Code:    "INVALID_REQUEST",
			})
			return
		}

		// Change password
		err = authService.ChangePassword(c.Request.Context(), userID, &req)
		if err != nil {
			logger.Error("Failed to change password", 
				logger.String("user_id", userID),
				logger.Error(err),
			)

			switch err {
			case services.ErrInvalidCredentials:
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error:   "Invalid current password",
					Message: "The current password is incorrect",
					Code:    "INVALID_PASSWORD",
				})
				return

			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Password change failed",
					Message: "Unable to change password",
					Code:    "PASSWORD_CHANGE_ERROR",
				})
				return
			}
		}

		logger.Info("Password changed successfully", 
			logger.String("user_id", userID),
		)

		c.JSON(http.StatusOK, SuccessResponse{
			Success: true,
			Message: "Password changed successfully",
		})
	}
}