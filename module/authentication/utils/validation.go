package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validator is a singleton instance of the validator.
var validate = validator.New()

// ValidationError represents a single validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

// ValidateStruct validates a struct based on the 'validate' tags.
// It returns a slice of ValidationErrors if any.
func ValidateStruct(s interface{}) []*ValidationError {
	var errors []*ValidationError
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Message = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", err.Field(), err.Tag())
			errors = append(errors, &element)
		}
	}
	return errors
}
