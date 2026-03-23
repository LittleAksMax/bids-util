package validation

import (
	"net/mail"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// ValidateEmails checks if fields marked with validate:"email" have valid email format.
func ValidateEmails(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	var invalidFields []string

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Check if field has validated tag with "email"
		validateTag := field.Tag.Get("validate")
		if !strings.Contains(validateTag, "email") {
			continue
		}

		// Only validate string fields
		if fieldValue.Kind() == reflect.String {
			emailStr := fieldValue.String()
			// Skip empty strings - let required validation handle that
			if emailStr == "" {
				continue
			}

			// Validate email format using net/mail
			if _, err := mail.ParseAddress(emailStr); err != nil {
				jsonTag := field.Tag.Get("json")
				fieldName := strings.Split(jsonTag, ",")[0]
				if fieldName == "" {
					fieldName = field.Name
				}
				invalidFields = append(invalidFields, fieldName)
			}
		}
	}

	if len(invalidFields) > 0 {
		return &EmailValidationError{Fields: invalidFields}
	}

	return nil
}

// ValidateUUIDs checks if fields marked with validate:"uuid" have valid UUID format.
func ValidateUUIDs(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	var invalidFields []string

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Check if field has validate tag with "uuid"
		validateTag := field.Tag.Get("validate")
		if !strings.Contains(validateTag, "uuid") {
			continue
		}

		// Only validate string fields
		if fieldValue.Kind() == reflect.String {
			uuidStr := fieldValue.String()
			// Skip empty strings - let required validation handle that
			if uuidStr == "" {
				continue
			}

			// Validate UUID format using google/uuid
			if err := uuid.Validate(uuidStr); err != nil {
				jsonTag := field.Tag.Get("json")
				fieldName := strings.Split(jsonTag, ",")[0]
				if fieldName == "" {
					fieldName = field.Name
				}
				invalidFields = append(invalidFields, fieldName)
			}
		}
	}

	if len(invalidFields) > 0 {
		return &UUIDValidationError{Fields: invalidFields}
	}

	return nil
}

// ValidatePasswords checks if fields marked with validate:"password" meet minimum strength requirements.
func ValidatePasswords(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	var invalidFields []string

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Check if field has validate tag with "password"
		validateTag := field.Tag.Get("validate")
		if !strings.Contains(validateTag, "password") {
			continue
		}

		// Only validate string fields
		if fieldValue.Kind() == reflect.String {
			passwordStr := fieldValue.String()
			// Skip empty strings - let required validation handle that
			if passwordStr == "" {
				continue
			}

			// Validate password strength (minimum 8 characters)
			if len(passwordStr) < 8 {
				jsonTag := field.Tag.Get("json")
				fieldName := strings.Split(jsonTag, ",")[0]
				if fieldName == "" {
					fieldName = field.Name
				}
				invalidFields = append(invalidFields, fieldName)
			}
		}
	}

	if len(invalidFields) > 0 {
		return &PasswordValidationError{Fields: invalidFields}
	}

	return nil
}

// ValidateRoles checks if fields marked with validate:"role" contain allowed role values.
func ValidateRoles(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	var invalidFields []string
	allowed := map[string]struct{}{
		"user":  {},
		"admin": {},
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Check if field has validate tag with "role"
		validateTag := field.Tag.Get("validate")
		if !strings.Contains(validateTag, "role") {
			continue
		}

		// Only validate string fields
		if fieldValue.Kind() == reflect.String {
			roleStr := strings.ToLower(strings.TrimSpace(fieldValue.String()))
			// Skip empty strings - let required validation handle that
			if roleStr == "" {
				continue
			}

			if _, ok := allowed[roleStr]; !ok {
				jsonTag := field.Tag.Get("json")
				fieldName := strings.Split(jsonTag, ",")[0]
				if fieldName == "" {
					fieldName = field.Name
				}
				invalidFields = append(invalidFields, fieldName)
			}
		}
	}

	if len(invalidFields) > 0 {
		return &RoleValidationError{Fields: invalidFields}
	}

	return nil
}
