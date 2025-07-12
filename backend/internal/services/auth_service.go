package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/visualtrecplans/backend/internal/database"
	"github.com/visualtrecplans/backend/internal/handlers/auth"
	"github.com/visualtrecplans/backend/internal/models"
	"github.com/visualtrecplans/backend/internal/validators"
	"github.com/visualtrecplans/backend/pkg/logger"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound = errors.New("user not found")
	ErrUserInactive = errors.New("user account is inactive")
	ErrEmailNotVerified = errors.New("email not verified")
	ErrWeakPassword = errors.New("password does not meet strength requirements")
)

// AuthService handles authentication-related business logic
type AuthService struct {
	db        *gorm.DB
	validator *validators.CustomValidator
}

// NewAuthService creates a new AuthService instance
func NewAuthService() *AuthService {
	return &AuthService{
		db:        database.GetDB(),
		validator: validators.NewValidator(),
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *auth.RegisterRequest) (*models.User, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		logger.Warn("Registration validation failed", 
			logger.String("email", req.Email),
			logger.Error(err),
		)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Additional password strength validation
	isStrong, passwordErrors := validators.ValidatePasswordStrength(req.Password)
	if !isStrong {
		logger.Warn("Weak password provided", 
			logger.String("email", req.Email),
			logger.Any("errors", passwordErrors),
		)
		return nil, fmt.Errorf("%w: %v", ErrWeakPassword, passwordErrors)
	}

	// Check if user already exists
	var existingUser models.User
	result := s.db.WithContext(ctx).Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		logger.Warn("Registration attempt with existing email", 
			logger.String("email", req.Email),
		)
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Error("Database error during user lookup", 
			logger.String("email", req.Email),
			logger.Error(result.Error),
		)
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password", 
			logger.String("email", req.Email),
			logger.Error(err),
		)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsActive:     true,
		EmailVerified: false,
	}

	// Set optional fields
	if req.FirstName != "" {
		user.FirstName = &req.FirstName
	}
	if req.LastName != "" {
		user.LastName = &req.LastName
	}

	// Save user in transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logger.Error("Failed to create user", 
			logger.String("email", req.Email),
			logger.Error(err),
		)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Info("User registered successfully", 
		logger.String("user_id", user.ID.String()),
		logger.String("email", user.Email),
	)

	return user, nil
}

// Login authenticates a user and returns user info
func (s *AuthService) Login(ctx context.Context, req *auth.LoginRequest) (*models.User, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		logger.Warn("Login validation failed", 
			logger.String("email", req.Email),
			logger.Error(err),
		)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Find user by email
	var user models.User
	result := s.db.WithContext(ctx).Where("email = ?", req.Email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Warn("Login attempt with non-existent email", 
			logger.String("email", req.Email),
		)
		return nil, ErrInvalidCredentials
	}
	if result.Error != nil {
		logger.Error("Database error during login", 
			logger.String("email", req.Email),
			logger.Error(result.Error),
		)
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	// Check if user is active
	if !user.IsActive {
		logger.Warn("Login attempt with inactive account", 
			logger.String("email", req.Email),
			logger.String("user_id", user.ID.String()),
		)
		return nil, ErrUserInactive
	}

	// Verify password
	if err := s.verifyPassword(req.Password, user.PasswordHash); err != nil {
		logger.Warn("Login attempt with invalid password", 
			logger.String("email", req.Email),
			logger.String("user_id", user.ID.String()),
		)
		return nil, ErrInvalidCredentials
	}

	// Update login information
	user.UpdateLoginInfo()
	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		logger.Error("Failed to update login info", 
			logger.String("user_id", user.ID.String()),
			logger.Error(err),
		)
		// Don't fail login for this error
	}

	logger.Info("User logged in successfully", 
		logger.String("user_id", user.ID.String()),
		logger.String("email", user.Email),
	)

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	result := s.db.WithContext(ctx).Where("id = ? AND is_active = ?", userID, true).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if result.Error != nil {
		logger.Error("Database error getting user by ID", 
			logger.String("user_id", userID),
			logger.Error(result.Error),
		)
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	return &user, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *auth.ChangePasswordRequest) error {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Additional password strength validation
	isStrong, passwordErrors := validators.ValidatePasswordStrength(req.NewPassword)
	if !isStrong {
		return fmt.Errorf("%w: %v", ErrWeakPassword, passwordErrors)
	}

	// Get user
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := s.verifyPassword(req.CurrentPassword, user.PasswordHash); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := s.hashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.PasswordHash = hashedPassword
	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		logger.Error("Failed to update password", 
			logger.String("user_id", userID),
			logger.Error(err),
		)
		return fmt.Errorf("failed to update password: %w", err)
	}

	logger.Info("Password changed successfully", 
		logger.String("user_id", userID),
	)

	return nil
}

// hashPassword hashes a password using bcrypt
func (s *AuthService) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// verifyPassword verifies a password against its hash
func (s *AuthService) verifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// DeactivateUser deactivates a user account (soft delete)
func (s *AuthService) DeactivateUser(ctx context.Context, userID string) error {
	result := s.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_active", false)
	
	if result.Error != nil {
		logger.Error("Failed to deactivate user", 
			logger.String("user_id", userID),
			logger.Error(result.Error),
		)
		return fmt.Errorf("failed to deactivate user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	logger.Info("User deactivated successfully", 
		logger.String("user_id", userID),
	)

	return nil
}

// VerifyEmail marks a user's email as verified
func (s *AuthService) VerifyEmail(ctx context.Context, userID string) error {
	now := time.Now()
	result := s.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"email_verified":    true,
			"email_verified_at": now,
		})
	
	if result.Error != nil {
		logger.Error("Failed to verify email", 
			logger.String("user_id", userID),
			logger.Error(result.Error),
		)
		return fmt.Errorf("failed to verify email: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	logger.Info("Email verified successfully", 
		logger.String("user_id", userID),
	)

	return nil
}

// UpdateUser updates user information
func (s *AuthService) UpdateUser(ctx context.Context, user *models.User) error {
	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		logger.Error("Failed to update user", 
			logger.String("user_id", user.ID.String()),
			logger.Error(err),
		)
		return fmt.Errorf("failed to update user: %w", err)
	}

	logger.Info("User updated successfully", 
		logger.String("user_id", user.ID.String()),
	)

	return nil
}