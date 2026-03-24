package validation

import (
	"reflect"
)

// ValidateRequiredFields checks if fields marked with validate:"required" are non-empty.
func ValidateRequiredFields(v interface{}) error {
	missingFields := ValidateByTag(v, "required", func(field reflect.StructField, value reflect.Value) bool {
		return value.Kind() == reflect.String && value.String() == ""
	})
	if len(missingFields) > 0 {
		return &ValidationError{Fields: missingFields}
	}
	return nil
}
