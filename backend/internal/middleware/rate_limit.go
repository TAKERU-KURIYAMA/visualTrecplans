package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/visualtrecplans/backend/internal/handlers/auth"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// RateLimiter represents a simple rate limiter
type RateLimiter struct {
	requests map[string]*ClientRequests
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// ClientRequests tracks requests for a specific client
type ClientRequests struct {
	count     int
	window    time.Time
	lastReset time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*ClientRequests),
		limit:    limit,
		window:   window,
	}

	// Start cleanup routine
	go rl.cleanup()

	return rl
}

// Allow checks if a request should be allowed
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	client, exists := rl.requests[clientID]

	if !exists {
		rl.requests[clientID] = &ClientRequests{
			count:     1,
			window:    now,
			lastReset: now,
		}
		return true
	}

	// Reset window if expired
	if now.Sub(client.window) >= rl.window {
		client.count = 1
		client.window = now
		client.lastReset = now
		return true
	}

	// Check if limit exceeded
	if client.count >= rl.limit {
		return false
	}

	client.count++
	return true
}

// cleanup removes old entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for clientID, client := range rl.requests {
			if now.Sub(client.lastReset) > rl.window*2 {
				delete(rl.requests, clientID)
			}
		}
		rl.mu.Unlock()
	}
}

// Global rate limiters
var (
	generalLimiter = NewRateLimiter(1000, time.Hour)        // 1000 requests per hour for general API
	authLimiter    = NewRateLimiter(20, time.Minute*15)     // 20 requests per 15 minutes for auth
	loginLimiter   = NewRateLimiter(5, time.Minute*15)      // 5 login attempts per 15 minutes
	registerLimiter = NewRateLimiter(3, time.Hour)          // 3 registrations per hour
)

// RateLimit middleware applies rate limiting based on endpoint
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := getClientID(c)
		path := c.Request.URL.Path
		method := c.Request.Method

		var limiter *RateLimiter
		var limitType string

		// Choose appropriate limiter based on endpoint
		switch {
		case path == "/api/v1/auth/login" && method == "POST":
			limiter = loginLimiter
			limitType = "login"
		case path == "/api/v1/auth/register" && method == "POST":
			limiter = registerLimiter
			limitType = "register"
		case strings.HasPrefix(path, "/api/v1/auth"):
			limiter = authLimiter
			limitType = "auth"
		default:
			limiter = generalLimiter
			limitType = "general"
		}

		if !limiter.Allow(clientID) {
			logger.Warn("Rate limit exceeded", 
				logger.String("client_id", clientID),
				logger.String("path", path),
				logger.String("method", method),
				logger.String("limit_type", limitType),
				logger.String("ip", c.ClientIP()),
			)

			c.JSON(http.StatusTooManyRequests, auth.ErrorResponse{
				Error:   "Rate limit exceeded",
				Message: getRateLimitMessage(limitType),
				Code:    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getClientID generates a client ID for rate limiting
func getClientID(c *gin.Context) string {
	// Try to get user ID from context first (for authenticated requests)
	if userID, exists := c.Get("user_id"); exists {
		if userIDStr, ok := userID.(string); ok {
			return "user:" + userIDStr
		}
	}

	// Fall back to IP address for unauthenticated requests
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	
	// Create a more specific identifier including user agent hash
	// This helps prevent simple IP spoofing while not storing full user agent
	return "ip:" + ip + ":" + hashString(userAgent)[:8]
}

// hashString creates a simple hash of a string
func hashString(s string) string {
	h := uint32(2166136261)
	for _, c := range s {
		h = (h ^ uint32(c)) * 16777619
	}
	return string(rune(h))
}

// getRateLimitMessage returns appropriate message for rate limit type
func getRateLimitMessage(limitType string) string {
	switch limitType {
	case "login":
		return "Too many login attempts. Please wait 15 minutes before trying again."
	case "register":
		return "Too many registration attempts. Please wait 1 hour before trying again."
	case "auth":
		return "Too many authentication requests. Please wait 15 minutes before trying again."
	default:
		return "Too many requests. Please slow down your request rate."
	}
}

// IPWhitelist middleware allows certain IPs to bypass rate limiting
func IPWhitelist(whitelist []string) gin.HandlerFunc {
	whitelistMap := make(map[string]bool)
	for _, ip := range whitelist {
		whitelistMap[ip] = true
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if whitelistMap[clientIP] {
			logger.Debug("IP whitelisted, bypassing rate limit", 
				logger.String("ip", clientIP),
			)
			c.Set("rate_limit_bypass", true)
		}
		c.Next()
	}
}

// BruteForceProtection provides additional protection against brute force attacks
func BruteForceProtection() gin.HandlerFunc {
	// Track failed attempts per IP/user combination
	failedAttempts := make(map[string]*FailedAttempts)
	var mu sync.RWMutex

	return func(c *gin.Context) {
		// Only apply to login endpoint
		if c.Request.URL.Path != "/api/v1/auth/login" {
			c.Next()
			return
		}

		clientID := getClientID(c)
		
		mu.RLock()
		attempts, exists := failedAttempts[clientID]
		mu.RUnlock()

		// Check if client is temporarily blocked
		if exists && attempts.IsBlocked() {
			logger.Warn("Brute force protection triggered", 
				logger.String("client_id", clientID),
				logger.Int("failed_attempts", attempts.Count),
				logger.String("ip", c.ClientIP()),
			)

			c.JSON(http.StatusTooManyRequests, auth.ErrorResponse{
				Error:   "Account temporarily locked",
				Message: "Too many failed login attempts. Please try again later.",
				Code:    "ACCOUNT_LOCKED",
			})
			c.Abort()
			return
		}

		c.Next()

		// Check response status to track failed attempts
		if c.Writer.Status() == http.StatusUnauthorized {
			mu.Lock()
			if attempts == nil {
				failedAttempts[clientID] = &FailedAttempts{
					Count:     1,
					FirstFail: time.Now(),
					LastFail:  time.Now(),
				}
			} else {
				attempts.Count++
				attempts.LastFail = time.Now()
			}
			mu.Unlock()

			logger.Info("Failed login attempt recorded", 
				logger.String("client_id", clientID),
				logger.Int("attempt_count", failedAttempts[clientID].Count),
			)
		} else if c.Writer.Status() == http.StatusOK {
			// Clear failed attempts on successful login
			mu.Lock()
			delete(failedAttempts, clientID)
			mu.Unlock()
		}
	}
}

// FailedAttempts tracks failed login attempts
type FailedAttempts struct {
	Count     int
	FirstFail time.Time
	LastFail  time.Time
}

// IsBlocked checks if the client should be blocked
func (fa *FailedAttempts) IsBlocked() bool {
	now := time.Now()
	
	// Reset if it's been more than 1 hour since first failed attempt
	if now.Sub(fa.FirstFail) > time.Hour {
		fa.Count = 0
		return false
	}

	// Block if more than 5 failed attempts in the last hour
	if fa.Count >= 5 {
		// Block for increasing durations: 5min, 15min, 30min, 1hour
		var blockDuration time.Duration
		switch {
		case fa.Count < 10:
			blockDuration = time.Minute * 5
		case fa.Count < 15:
			blockDuration = time.Minute * 15
		case fa.Count < 20:
			blockDuration = time.Minute * 30
		default:
			blockDuration = time.Hour
		}

		return now.Sub(fa.LastFail) < blockDuration
	}

	return false
}