package validation

import (
	"reflect"
	"strings"
)

// ValidateRequiredFields checks if fields marked with validate:"required" are non-empty.
func ValidateRequiredFields(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	var missingFields []string

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Check if field has validate:"required" tag
		validateTag := field.Tag.Get("validate")
		if !strings.Contains(validateTag, "required") {
			continue
		}

		// Check if string field is empty
		if fieldValue.Kind() == reflect.String && fieldValue.String() == "" {
			jsonTag := field.Tag.Get("json")
			fieldName := strings.Split(jsonTag, ",")[0]
			if fieldName == "" {
				fieldName = field.Name
			}
			missingFields = append(missingFields, fieldName)
		}
	}

	if len(missingFields) > 0 {
		return &ValidationError{Fields: missingFields}
	}

	return nil
}
