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

// ValidateNonNegativeFields checks if fields marked with validate:"nonnegative" are non-negative numbers.
func ValidateNonNegativeFields(v interface{}) error {
	nonNegativeFields := ValidateByTag(v, "nonnegative", func(field reflect.StructField, value reflect.Value) bool {
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return value.Int() < 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return false // unsigned types are always non-negative
		case reflect.Float32, reflect.Float64:
			return value.Float() < 0
		default:
			return false
		}
	})
	if len(nonNegativeFields) > 0 {
		return &ValidationError{Fields: nonNegativeFields}
	}
	return nil
}
