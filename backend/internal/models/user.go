package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email           string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,email"`
	PasswordHash    string         `json:"-" gorm:"type:varchar(255);not null"`
	FirstName       *string        `json:"first_name" gorm:"type:varchar(100)"`
	LastName        *string        `json:"last_name" gorm:"type:varchar(100)"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	EmailVerified   bool           `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	LoginCount      int            `json:"login_count" gorm:"default:0"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	var fullName string
	if u.FirstName != nil {
		fullName = *u.FirstName
	}
	if u.LastName != nil {
		if fullName != "" {
			fullName += " "
		}
		fullName += *u.LastName
	}
	return fullName
}

// IsEmailVerified checks if the user's email is verified
func (u *User) IsEmailVerified() bool {
	return u.EmailVerified && u.EmailVerifiedAt != nil
}

// UpdateLoginInfo updates login-related fields
func (u *User) UpdateLoginInfo() {
	now := time.Now()
	u.LastLoginAt = &now
	u.LoginCount++
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	FirstName       *string    `json:"first_name"`
	LastName        *string    `json:"last_name"`
	IsActive        bool       `json:"is_active"`
	EmailVerified   bool       `json:"email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	LoginCount      int        `json:"login_count"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		IsActive:        u.IsActive,
		EmailVerified:   u.EmailVerified,
		EmailVerifiedAt: u.EmailVerifiedAt,
		LastLoginAt:     u.LastLoginAt,
		LoginCount:      u.LoginCount,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}