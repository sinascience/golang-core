package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Global validator instance
var validate = validator.New()

// ValidateStruct performs validation on a struct using its 'validate' tags.
func ValidateStruct(payload interface{}) map[string]string {
	// Perform validation
	err := validate.Struct(payload)
	if err == nil {
		return nil // No errors
	}

	// Type assert the error to ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// This should not happen with standard validation, but it's good practice
		return map[string]string{"error": "Invalid validation error"}
	}

	// Create a map to hold our custom error messages
	errorMessages := make(map[string]string)

	for _, fieldErr := range validationErrors {
		// Use a switch to create user-friendly messages
		fieldName := strings.ToLower(fieldErr.Field())
		switch fieldErr.Tag() {
		case "required":
			errorMessages[fieldName] = fieldName + " is required"
		case "email":
			errorMessages[fieldName] = "must be a valid email address"
		case "min":
			errorMessages[fieldName] = fieldName + " must be at least " + fieldErr.Param() + " characters long"
		case "max":
			errorMessages[fieldName] = fieldName + " must be at most " + fieldErr.Param() + " characters long"
		default:
			errorMessages[fieldName] = "invalid value"
		}
	}

	return errorMessages
}
