package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/visualtrecplans/backend/internal/audit"
	"github.com/visualtrecplans/backend/internal/services"
	"github.com/visualtrecplans/backend/pkg/jwt"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// LoginHandler handles user login
func LoginHandler(authService *services.AuthService, jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest

		// Bind JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("Invalid login request format", 
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

		// Authenticate user
		user, err := authService.Login(c.Request.Context(), &req)
		if err != nil {
			logger.Warn("Login failed", 
				logger.String("email", req.Email),
				logger.Error(err),
				logger.String("ip", c.ClientIP()),
			)

			// Log failed login attempt
			audit.GetAuditLogger().LogLoginFailed(
				c.Request.Context(),
				req.Email,
				c.ClientIP(),
				c.GetHeader("User-Agent"),
				err.Error(),
			)

			// Handle specific errors
			switch {
			case errors.Is(err, services.ErrInvalidCredentials):
				c.JSON(http.StatusUnauthorized, ErrorResponse{
					Error:   "Invalid credentials",
					Message: "Email or password is incorrect",
					Code:    "INVALID_CREDENTIALS",
				})
				return

			case errors.Is(err, services.ErrUserInactive):
				c.JSON(http.StatusForbidden, ErrorResponse{
					Error:   "Account inactive",
					Message: "Your account has been deactivated. Please contact support",
					Code:    "ACCOUNT_INACTIVE",
				})
				return

			case errors.Is(err, services.ErrEmailNotVerified):
				c.JSON(http.StatusForbidden, ErrorResponse{
					Error:   "Email not verified",
					Message: "Please verify your email address before logging in",
					Code:    "EMAIL_NOT_VERIFIED",
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
					Error:   "Login failed",
					Message: "Unable to process login. Please try again later",
					Code:    "INTERNAL_ERROR",
				})
				return
			}
		}

		// Generate JWT tokens
		accessToken, refreshToken, accessExpiry, refreshExpiry, err := jwtService.GenerateTokenPair(
			user.ID,
			user.Email,
			user.IsActive,
		)
		if err != nil {
			logger.Error("Failed to generate JWT tokens", 
				logger.String("user_id", user.ID.String()),
				logger.Error(err),
			)
			
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Token generation failed",
				Message: "Unable to generate authentication tokens",
				Code:    "TOKEN_ERROR",
			})
			return
		}

		logger.Info("User logged in successfully", 
			logger.String("user_id", user.ID.String()),
			logger.String("email", user.Email),
			logger.String("ip", c.ClientIP()),
		)

		// Log successful login
		audit.GetAuditLogger().LogLoginSuccess(
			c.Request.Context(),
			user.ID,
			user.Email,
			c.ClientIP(),
			c.GetHeader("User-Agent"),
		)

		// Create user info
		userInfo := UserInfo{
			ID:              user.ID,
			Email:           user.Email,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			IsActive:        user.IsActive,
			EmailVerified:   user.EmailVerified,
			EmailVerifiedAt: user.EmailVerifiedAt,
			LastLoginAt:     user.LastLoginAt,
			LoginCount:      user.LoginCount,
			CreatedAt:       user.CreatedAt,
		}

		// Calculate expires in seconds for frontend convenience
		expiresInSeconds := int(jwtService.GetAccessTokenExpiry().Seconds())

		// Return success response
		response := LoginResponse{
			User:         userInfo,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    expiresInSeconds,
			Message:      "Login successful",
		}

		// Set secure HTTP-only cookie for refresh token (optional)
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(
			"refresh_token",
			refreshToken,
			int(refreshExpiry.Sub(accessExpiry).Seconds()),
			"/api/v1/auth",
			"",
			true,  // secure
			true,  // httpOnly
		)

		c.JSON(http.StatusOK, response)
	}
}

// RefreshTokenHandler handles token refresh
func RefreshTokenHandler(authService *services.AuthService, jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefreshTokenRequest

		// Try to get refresh token from request body first
		if err := c.ShouldBindJSON(&req); err != nil {
			// If no JSON body, try to get from cookie
			if refreshToken, err := c.Cookie("refresh_token"); err == nil {
				req.RefreshToken = refreshToken
			} else {
				logger.Warn("No refresh token provided", 
					logger.String("ip", c.ClientIP()),
				)
				
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error:   "Refresh token required",
					Message: "Please provide a refresh token",
					Code:    "MISSING_REFRESH_TOKEN",
				})
				return
			}
		}

		// Generate new access token from refresh token
		newAccessToken, accessExpiry, err := jwtService.RefreshAccessToken(req.RefreshToken)
		if err != nil {
			logger.Warn("Token refresh failed", 
				logger.Error(err),
				logger.String("ip", c.ClientIP()),
			)

			// Clear refresh token cookie if invalid
			c.SetCookie("refresh_token", "", -1, "/api/v1/auth", "", true, true)

			switch {
			case errors.Is(err, jwt.ErrExpiredToken):
				c.JSON(http.StatusUnauthorized, ErrorResponse{
					Error:   "Refresh token expired",
					Message: "Please login again",
					Code:    "REFRESH_TOKEN_EXPIRED",
				})
				return

			case errors.Is(err, jwt.ErrInvalidToken):
				c.JSON(http.StatusUnauthorized, ErrorResponse{
					Error:   "Invalid refresh token",
					Message: "Please login again",
					Code:    "INVALID_REFRESH_TOKEN",
				})
				return

			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Token refresh failed",
					Message: "Unable to refresh token",
					Code:    "REFRESH_ERROR",
				})
				return
			}
		}

		// Get user info from refresh token claims
		claims, err := jwtService.ValidateRefreshToken(req.RefreshToken)
		if err != nil {
			logger.Error("Failed to validate refresh token claims", logger.Error(err))
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Invalid refresh token",
				Message: "Please login again",
				Code:    "INVALID_REFRESH_TOKEN",
			})
			return
		}

		// Get updated user info
		user, err := authService.GetUserByID(c.Request.Context(), claims.UserID.String())
		if err != nil {
			logger.Error("Failed to get user for token refresh", 
				logger.String("user_id", claims.UserID.String()),
				logger.Error(err),
			)
			
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "User not found",
				Message: "Please login again",
				Code:    "USER_NOT_FOUND",
			})
			return
		}

		logger.Info("Token refreshed successfully", 
			logger.String("user_id", user.ID.String()),
			logger.String("email", user.Email),
		)

		// Calculate expires in seconds
		expiresInSeconds := int(jwtService.GetAccessTokenExpiry().Seconds())

		// Create user info
		userInfo := UserInfo{
			ID:              user.ID,
			Email:           user.Email,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			IsActive:        user.IsActive,
			EmailVerified:   user.EmailVerified,
			EmailVerifiedAt: user.EmailVerifiedAt,
			LastLoginAt:     user.LastLoginAt,
			LoginCount:      user.LoginCount,
			CreatedAt:       user.CreatedAt,
		}

		// Return new access token
		response := LoginResponse{
			User:        userInfo,
			AccessToken: newAccessToken,
			TokenType:   "Bearer",
			ExpiresIn:   expiresInSeconds,
			Message:     "Token refreshed successfully",
		}

		c.JSON(http.StatusOK, response)
	}
}

// LogoutHandler handles user logout
func LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Clear refresh token cookie
		c.SetCookie("refresh_token", "", -1, "/api/v1/auth", "", true, true)

		logger.Info("User logged out", 
			logger.String("ip", c.ClientIP()),
		)

		c.JSON(http.StatusOK, SuccessResponse{
			Success: true,
			Message: "Logged out successfully",
		})
	}
}