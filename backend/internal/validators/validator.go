package validators

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps the validator instance
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator creates a new custom validator instance
func NewValidator() *CustomValidator {
	v := validator.New()
	
	// Register custom validation functions
	v.RegisterValidation("password_strength", PasswordStrength)
	
	return &CustomValidator{
		validator: v,
	}
}

// Validate validates a struct using the custom validator
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// Engine returns the underlying validator engine
func (cv *CustomValidator) Engine() interface{} {
	return cv.validator
}