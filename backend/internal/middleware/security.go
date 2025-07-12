package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"

	"github.com/visualtrecplans/backend/pkg/config"
)

// SecurityHeaders middleware adds security headers to all responses
func SecurityHeaders(cfg *config.Config) gin.HandlerFunc {
	secureConfig := secure.Options{
		// SSL configuration
		SSLRedirect:          cfg.App.Environment == "production",
		SSLTemporaryRedirect: false,
		SSLHost:              "",
		STSSeconds:           31536000, // 1 year
		STSIncludeSubdomains: true,
		STSPreload:           true,

		// Frame protection
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",

		// Content Security Policy
		ContentSecurityPolicy: buildCSP(cfg),

		// HPKP (disabled by default)
		PublicKey: "",

		// Feature Policy / Permissions Policy
		FeaturePolicy: "microphone 'none'; camera 'none'; geolocation 'self'",

		// Custom headers
		CustomFrameOptionsValue: "DENY",
		CustomBrowserXssValue:   "1; mode=block",

		// Development mode settings
		IsDevelopment: cfg.App.Environment == "development",
	}

	// Configure allowed hosts based on environment
	if cfg.App.Environment == "production" {
		secureConfig.AllowedHosts = []string{
			"trecplans.com",
			"www.trecplans.com",
			"api.trecplans.com",
		}
	}

	secureMiddleware := secure.New(secureConfig)

	return func(c *gin.Context) {
		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			// If there was an error, do not continue
			c.AbortWithStatus(400)
			return
		}

		// Add additional security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "microphone=(), camera=(), geolocation=(self)")
		
		// Add cache control for sensitive endpoints
		if isAuthEndpoint(c.Request.URL.Path) {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	}
}

// buildCSP constructs Content Security Policy based on environment
func buildCSP(cfg *config.Config) string {
	if cfg.App.Environment == "development" {
		// More relaxed CSP for development
		return "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval' localhost:* 127.0.0.1:*; " +
			"style-src 'self' 'unsafe-inline' localhost:* 127.0.0.1:*; " +
			"img-src 'self' data: blob: localhost:* 127.0.0.1:*; " +
			"connect-src 'self' localhost:* 127.0.0.1:* ws: wss:; " +
			"font-src 'self' data:; " +
			"object-src 'none'; " +
			"media-src 'self'; " +
			"frame-src 'none'"
	}

	// Strict CSP for production
	return "default-src 'self'; " +
		"script-src 'self'; " +
		"style-src 'self' 'unsafe-inline'; " +
		"img-src 'self' data: blob:; " +
		"connect-src 'self' https://api.trecplans.com; " +
		"font-src 'self' data:; " +
		"object-src 'none'; " +
		"media-src 'self'; " +
		"frame-src 'none'; " +
		"base-uri 'self'; " +
		"form-action 'self'"
}

// isAuthEndpoint checks if the path is an authentication endpoint
func isAuthEndpoint(path string) bool {
	authPaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/refresh",
		"/api/v1/auth/logout",
		"/api/v1/auth/password",
	}

	for _, authPath := range authPaths {
		if path == authPath {
			return true
		}
	}
	return false
}

// RequestID middleware adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.Header("X-Request-ID", generateRequestID())
		c.AbortWithStatus(500)
	})
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	// Simple timestamp-based ID for now
	// In production, you might want to use UUID or similar
	return time.Now().Format("20060102150405.000000")
}