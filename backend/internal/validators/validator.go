package validators

import (
	"github.com/gin-gonic/gin/binding"
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

// ValidateMuscleGroup validates muscle group codes
func ValidateMuscleGroup(fl validator.FieldLevel) bool {
	validGroups := []string{"chest", "back", "shoulders", "arms", "core", "legs", "glutes", "full_body"}
	value := fl.Field().String()
	
	for _, group := range validGroups {
		if group == value {
			return true
		}
	}
	return false
}

// SetupCustomValidators registers all custom validators
func SetupCustomValidators() {
	// Register validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", PasswordStrength)
		v.RegisterValidation("muscle_group", ValidateMuscleGroup)
	}
}