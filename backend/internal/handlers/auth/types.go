package auth

import (
	"time"

	"github.com/google/uuid"
)

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email" validate:"required,email"`
	Password        string `json:"password" binding:"required,min=8" validate:"required,min=8,password_strength"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password" validate:"required,eqfield=Password"`
	FirstName       string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName        string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
}

// RegisterResponse represents the response for successful user registration
type RegisterResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName *string   `json:"first_name,omitempty"`
	LastName  *string   `json:"last_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"message"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required"`
}

// LoginResponse represents the response for successful user login
type LoginResponse struct {
	User         UserInfo `json:"user"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token,omitempty"`
	TokenType    string   `json:"token_type"`
	ExpiresIn    int      `json:"expires_in"`
	Message      string   `json:"message"`
}

// UserInfo represents user information in responses
type UserInfo struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	FirstName       *string    `json:"first_name,omitempty"`
	LastName        *string    `json:"last_name,omitempty"`
	IsActive        bool       `json:"is_active"`
	EmailVerified   bool       `json:"email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	LoginCount      int        `json:"login_count"`
	CreatedAt       time.Time  `json:"created_at"`
}

// RefreshTokenRequest represents the request for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" validate:"required"`
}

// ChangePasswordRequest represents the request for password change
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" validate:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8" validate:"required,min=8,password_strength"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=NewPassword" validate:"required,eqfield=NewPassword"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Code    string                 `json:"code,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ValidationErrorResponse represents validation error response
type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UpdateProfileRequest represents the request for updating user profile
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	TokenType string    `json:"token_type"` // "access" or "refresh"
	IssuedAt  int64     `json:"iat"`
	ExpiresAt int64     `json:"exp"`
}