package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/visualtrecplans/backend/internal/handlers/auth"
	"github.com/visualtrecplans/backend/pkg/jwt"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// AuthRequired middleware enforces authentication for protected routes
func AuthRequired(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractTokenFromHeader(c)
		if err != nil {
			logger.Warn("Authentication failed - no token provided", 
				logger.String("path", c.Request.URL.Path),
				logger.String("ip", c.ClientIP()),
				logger.Error(err),
			)

			c.JSON(http.StatusUnauthorized, auth.ErrorResponse{
				Error:   "Authentication required",
				Message: "Please provide a valid authentication token",
				Code:    "AUTH_REQUIRED",
			})
			c.Abort()
			return
		}

		// Validate access token
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			logger.Warn("Authentication failed - invalid token", 
				logger.String("path", c.Request.URL.Path),
				logger.String("ip", c.ClientIP()),
				logger.Error(err),
			)

			var errorResponse auth.ErrorResponse
			switch {
			case errors.Is(err, jwt.ErrExpiredToken):
				errorResponse = auth.ErrorResponse{
					Error:   "Token expired",
					Message: "Your session has expired. Please login again",
					Code:    "TOKEN_EXPIRED",
				}
			case errors.Is(err, jwt.ErrInvalidToken):
				errorResponse = auth.ErrorResponse{
					Error:   "Invalid token",
					Message: "The provided token is invalid",
					Code:    "INVALID_TOKEN",
				}
			default:
				errorResponse = auth.ErrorResponse{
					Error:   "Authentication failed",
					Message: "Unable to verify authentication token",
					Code:    "AUTH_FAILED",
				}
			}

			c.JSON(http.StatusUnauthorized, errorResponse)
			c.Abort()
			return
		}

		// Check if user is active
		if !claims.IsActive {
			logger.Warn("Authentication failed - inactive user", 
				logger.String("user_id", claims.UserID.String()),
				logger.String("email", claims.Email),
				logger.String("path", c.Request.URL.Path),
			)

			c.JSON(http.StatusForbidden, auth.ErrorResponse{
				Error:   "Account inactive",
				Message: "Your account has been deactivated",
				Code:    "ACCOUNT_INACTIVE",
			})
			c.Abort()
			return
		}

		// Set user context
		setUserContext(c, claims)

		logger.Debug("Authentication successful", 
			logger.String("user_id", claims.UserID.String()),
			logger.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}

// AuthOptional middleware provides optional authentication for public routes
func AuthOptional(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractTokenFromHeader(c)
		if err != nil {
			// No token provided, continue without authentication
			c.Next()
			return
		}

		// Try to validate token, but don't fail if invalid
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			logger.Debug("Optional authentication failed", 
				logger.String("path", c.Request.URL.Path),
				logger.Error(err),
			)
			// Continue without authentication
			c.Next()
			return
		}

		// Check if user is active
		if !claims.IsActive {
			logger.Debug("Optional authentication - inactive user", 
				logger.String("user_id", claims.UserID.String()),
			)
			// Continue without authentication
			c.Next()
			return
		}

		// Set user context if token is valid
		setUserContext(c, claims)

		logger.Debug("Optional authentication successful", 
			logger.String("user_id", claims.UserID.String()),
		)

		c.Next()
	}
}

// extractTokenFromHeader extracts Bearer token from Authorization header
func extractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	// Check if it starts with Bearer
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	// Extract token
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token == "" {
		return "", errors.New("bearer token is empty")
	}

	return token, nil
}

// setUserContext sets user information in Gin context
func setUserContext(c *gin.Context, claims *jwt.CustomClaims) {
	c.Set("user_id", claims.UserID.String())
	c.Set("user_uuid", claims.UserID)
	c.Set("email", claims.Email)
	c.Set("is_active", claims.IsActive)
	c.Set("user_claims", claims)
}

// GetUserID retrieves user ID from context
func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", errors.New("user not authenticated")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", errors.New("invalid user ID in context")
	}

	return userIDStr, nil
}

// GetUserUUID retrieves user UUID from context
func GetUserUUID(c *gin.Context) (uuid.UUID, error) {
	userUUID, exists := c.Get("user_uuid")
	if !exists {
		return uuid.Nil, errors.New("user not authenticated")
	}

	userUUIDVal, ok := userUUID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid user UUID in context")
	}

	return userUUIDVal, nil
}

// GetUserEmail retrieves user email from context
func GetUserEmail(c *gin.Context) (string, error) {
	email, exists := c.Get("email")
	if !exists {
		return "", errors.New("user not authenticated")
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", errors.New("invalid email in context")
	}

	return emailStr, nil
}

// GetUserClaims retrieves full JWT claims from context
func GetUserClaims(c *gin.Context) (*jwt.CustomClaims, error) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return nil, errors.New("user not authenticated")
	}

	claimsVal, ok := claims.(*jwt.CustomClaims)
	if !ok {
		return nil, errors.New("invalid claims in context")
	}

	return claimsVal, nil
}

// IsAuthenticated checks if the current request is authenticated
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// RequireAdmin middleware checks if user has admin privileges
func RequireAdmin(jwtService *jwt.JWTService) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// First ensure user is authenticated
		AuthRequired(jwtService)(c)
		
		// If authentication failed, AuthRequired already handled the response
		if c.IsAborted() {
			return
		}

		// TODO: Implement admin check when role system is added
		// For now, we'll add a placeholder
		userEmail, err := GetUserEmail(c)
		if err != nil {
			c.JSON(http.StatusForbidden, auth.ErrorResponse{
				Error:   "Admin access required",
				Message: "This endpoint requires administrator privileges",
				Code:    "ADMIN_REQUIRED",
			})
			c.Abort()
			return
		}

		// Temporary admin check - in production this should be role-based
		if !strings.HasSuffix(userEmail, "@trecplans.com") {
			logger.Warn("Admin access denied", 
				logger.String("email", userEmail),
				logger.String("path", c.Request.URL.Path),
			)

			c.JSON(http.StatusForbidden, auth.ErrorResponse{
				Error:   "Admin access required",
				Message: "This endpoint requires administrator privileges",
				Code:    "ADMIN_REQUIRED",
			})
			c.Abort()
			return
		}

		logger.Info("Admin access granted", 
			logger.String("email", userEmail),
			logger.String("path", c.Request.URL.Path),
		)

		c.Next()
	})
}

// CORS middleware for handling Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// In development, allow localhost origins
		if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}