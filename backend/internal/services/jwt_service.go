package services

import (
	"time"

	"github.com/visualtrecplans/backend/pkg/config"
	"github.com/visualtrecplans/backend/pkg/jwt"
)

// NewJWTService creates a JWT service instance from configuration
func NewJWTService(cfg *config.Config) *jwt.JWTService {
	// Parse expiry durations
	accessExpiry := 24 * time.Hour  // Default 24 hours
	refreshExpiry := 7 * 24 * time.Hour // Default 7 days

	// Parse from config if available
	if cfg.JWT.ExpiresIn != "" {
		if duration, err := time.ParseDuration(cfg.JWT.ExpiresIn); err == nil {
			accessExpiry = duration
		}
	}

	if cfg.JWT.RefreshExpiry != "" {
		if duration, err := time.ParseDuration(cfg.JWT.RefreshExpiry); err == nil {
			refreshExpiry = duration
		}
	}

	return jwt.NewJWTService(
		cfg.JWT.Secret,
		accessExpiry,
		refreshExpiry,
		cfg.App.Name,
	)
}