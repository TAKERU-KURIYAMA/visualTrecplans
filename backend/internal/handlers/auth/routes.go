package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/visualtrecplans/backend/internal/middleware"
	"github.com/visualtrecplans/backend/internal/services"
	"github.com/visualtrecplans/backend/pkg/jwt"
)

// RegisterRoutes registers all authentication routes
func RegisterRoutes(router *gin.RouterGroup, authService *services.AuthService, jwtService *jwt.JWTService) {
	auth := router.Group("/auth")
	{
		// Public routes (no authentication required)
		auth.POST("/register", RegisterHandler(authService))
		auth.POST("/login", LoginHandler(authService, jwtService))
		auth.POST("/refresh", RefreshTokenHandler(authService, jwtService))
		auth.POST("/logout", LogoutHandler())
		
		// Protected routes (authentication required)
		protected := auth.Group("")
		protected.Use(middleware.AuthRequired(jwtService))
		{
			// User profile endpoints
			protected.GET("/profile", GetProfileHandler(authService))
			protected.PUT("/profile", UpdateProfileHandler(authService))
			
			// Password change endpoint
			protected.PUT("/password", ChangePasswordHandler(authService))
		}
		
		// Email verification endpoint (to be implemented)
		// auth.POST("/verify-email", VerifyEmailHandler(authService))
	}
}