package validators

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// PasswordStrength validates password strength
func PasswordStrength(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return IsStrongPassword(password)
}

// IsStrongPassword checks if password meets strength requirements
func IsStrongPassword(password string) bool {
	// Minimum length
	if len(password) < 8 {
		return false
	}

	// Maximum length for security reasons
	if len(password) > 128 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	// Check character types
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Require at least 3 of 4 character types
	requirements := 0
	if hasUpper {
		requirements++
	}
	if hasLower {
		requirements++
	}
	if hasNumber {
		requirements++
	}
	if hasSpecial {
		requirements++
	}

	return requirements >= 3
}

// IsCommonPassword checks if password is in common password list
func IsCommonPassword(password string) bool {
	commonPasswords := []string{
		"password", "123456", "123456789", "12345678", "12345",
		"qwerty", "abc123", "password123", "admin", "letmein",
		"welcome", "monkey", "1234567890", "dragon", "123123",
		"football", "iloveyou", "admin123", "welcome123", "password1",
	}

	for _, common := range commonPasswords {
		if password == common {
			return true
		}
	}
	return false
}

// IsSequentialPassword checks for sequential patterns
func IsSequentialPassword(password string) bool {
	// Check for sequential numbers
	sequentialNumbers := regexp.MustCompile(`(012|123|234|345|456|567|678|789|890|987|876|765|654|543|432|321|210)`)
	if sequentialNumbers.MatchString(password) {
		return true
	}

	// Check for sequential letters
	sequentialLetters := regexp.MustCompile(`(?i)(abc|bcd|cde|def|efg|fgh|ghi|hij|ijk|jkl|klm|lmn|mno|nop|opq|pqr|qrs|rst|stu|tuv|uvw|vwx|wxy|xyz|zyx|yxw|xwv|wvu|vut|uts|tsr|srq|rqp|qpo|pon|onm|nml|mlk|lkj|kji|jih|ihg|hgf|gfe|fed|edc|dcb|cba)`)
	if sequentialLetters.MatchString(password) {
		return true
	}

	return false
}

// ValidatePasswordStrength provides detailed password validation
func ValidatePasswordStrength(password string) (bool, []string) {
	var errors []string

	// Length check
	if len(password) < 8 {
		errors = append(errors, "Password must be at least 8 characters long")
	}
	if len(password) > 128 {
		errors = append(errors, "Password must be no more than 128 characters long")
	}

	// Character type checks
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	requirements := 0
	if hasUpper {
		requirements++
	} else {
		errors = append(errors, "Password should contain at least one uppercase letter")
	}

	if hasLower {
		requirements++
	} else {
		errors = append(errors, "Password should contain at least one lowercase letter")
	}

	if hasNumber {
		requirements++
	} else {
		errors = append(errors, "Password should contain at least one number")
	}

	if hasSpecial {
		requirements++
	} else {
		errors = append(errors, "Password should contain at least one special character")
	}

	// Need at least 3 requirements
	if requirements < 3 {
		errors = append(errors, "Password must meet at least 3 of the 4 character type requirements")
	}

	// Common password check
	if IsCommonPassword(password) {
		errors = append(errors, "Password is too common, please choose a more secure password")
	}

	// Sequential pattern check
	if IsSequentialPassword(password) {
		errors = append(errors, "Password contains sequential patterns, please choose a more secure password")
	}

	// Repeated character check
	if hasRepeatedCharacters(password) {
		errors = append(errors, "Password contains too many repeated characters")
	}

	return len(errors) == 0, errors
}

// hasRepeatedCharacters checks for excessive character repetition
func hasRepeatedCharacters(password string) bool {
	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i+1] == password[i+2] {
			return true
		}
	}
	return false
}