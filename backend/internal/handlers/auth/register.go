package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/visualtrecplans/backend/internal/services"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// RegisterHandler handles user registration
func RegisterHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest

		// Bind JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("Invalid registration request format", 
				logger.Error(err),
				logger.String("ip", c.ClientIP()),
			)
			
			// Handle validation errors
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				fieldErrors := make(map[string]string)
				for _, fieldError := range validationErrors {
					fieldErrors[strings.ToLower(fieldError.Field())] = getValidationErrorMessage(fieldError)
				}
				
				c.JSON(http.StatusBadRequest, ValidationErrorResponse{
					Error:   "Validation failed",
					Message: "Please check the provided data",
					Fields:  fieldErrors,
				})
				return
			}
			
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid request format",
				Message: "Please provide valid JSON data",
			})
			return
		}

		// Sanitize input
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		req.FirstName = strings.TrimSpace(req.FirstName)
		req.LastName = strings.TrimSpace(req.LastName)

		// Register user
		user, err := authService.Register(c.Request.Context(), &req)
		if err != nil {
			logger.Error("Registration failed", 
				logger.String("email", req.Email),
				logger.Error(err),
				logger.String("ip", c.ClientIP()),
			)

			// Handle specific errors
			switch {
			case errors.Is(err, services.ErrUserAlreadyExists):
				c.JSON(http.StatusConflict, ErrorResponse{
					Error:   "User already exists",
					Message: "An account with this email address already exists",
					Code:    "USER_EXISTS",
				})
				return

			case errors.Is(err, services.ErrWeakPassword):
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error:   "Weak password",
					Message: "Password does not meet security requirements",
					Code:    "WEAK_PASSWORD",
					Details: map[string]interface{}{
						"password_requirements": []string{
							"At least 8 characters long",
							"At least 3 of: uppercase, lowercase, numbers, special characters",
							"Not a common password",
							"No sequential patterns",
						},
					},
				})
				return

			case strings.Contains(err.Error(), "validation failed"):
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error:   "Validation error",
					Message: err.Error(),
					Code:    "VALIDATION_ERROR",
				})
				return

			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Registration failed",
					Message: "Unable to create account. Please try again later",
					Code:    "INTERNAL_ERROR",
				})
				return
			}
		}

		logger.Info("User registered successfully", 
			logger.String("user_id", user.ID.String()),
			logger.String("email", user.Email),
			logger.String("ip", c.ClientIP()),
		)

		// Return success response
		response := RegisterResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt,
			Message:   "Account created successfully",
		}

		c.JSON(http.StatusCreated, response)
	}
}

// getValidationErrorMessage returns user-friendly validation error messages
func getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Please provide a valid email address"
	case "min":
		if fe.Field() == "Password" {
			return "Password must be at least 8 characters long"
		}
		return "This field is too short"
	case "max":
		return "This field is too long"
	case "eqfield":
		if fe.Field() == "PasswordConfirm" {
			return "Password confirmation must match password"
		}
		return "This field must match the corresponding field"
	case "password_strength":
		return "Password does not meet strength requirements"
	default:
		return "Invalid value provided"
	}
}