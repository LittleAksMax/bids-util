package validation

import (
	"net/mail"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// ValidateEmails checks if fields marked with validate:"email" have valid email format.
func ValidateEmails(v interface{}) error {
	invalidFields := ValidateByTag(v, "email", func(field reflect.StructField, value reflect.Value) bool {
		if value.Kind() != reflect.String {
			return false
		}
		emailStr := value.String()
		if emailStr == "" {
			return false // let required validation handle empty
		}
		_, err := mail.ParseAddress(emailStr)
		return err != nil
	})
	if len(invalidFields) > 0 {
		return &EmailValidationError{Fields: invalidFields}
	}
	return nil
}

// ValidateUUIDs checks if fields marked with validate:"uuid" have valid UUID format.
func ValidateUUIDs(v interface{}) error {
	invalidFields := ValidateByTag(v, "uuid", func(field reflect.StructField, value reflect.Value) bool {
		if value.Kind() != reflect.String {
			return false
		}
		uuidStr := value.String()
		if uuidStr == "" {
			return false // let required validation handle empty
		}
		return uuid.Validate(uuidStr) != nil
	})
	if len(invalidFields) > 0 {
		return &UUIDValidationError{Fields: invalidFields}
	}
	return nil
}

// ValidatePasswords checks if fields marked with validate:"password" meet minimum strength requirements.
func ValidatePasswords(v interface{}) error {
	invalidFields := ValidateByTag(v, "password", func(field reflect.StructField, value reflect.Value) bool {
		if value.Kind() != reflect.String {
			return false
		}
		passwordStr := value.String()
		if passwordStr == "" {
			return false // let required validation handle empty
		}
		return len(passwordStr) < 8
	})
	if len(invalidFields) > 0 {
		return &PasswordValidationError{Fields: invalidFields}
	}
	return nil
}

// ValidateRoles checks if fields marked with validate:"role" contain allowed role values.
func ValidateRoles(v interface{}) error {
	allowed := map[string]struct{}{"user": {}, "admin": {}}
	invalidFields := ValidateByTag(v, "role", func(field reflect.StructField, value reflect.Value) bool {
		if value.Kind() != reflect.String {
			return false
		}
		roleStr := strings.ToLower(strings.TrimSpace(value.String()))
		if roleStr == "" {
			return false // let required validation handle empty
		}
		_, ok := allowed[roleStr]
		return !ok
	})
	if len(invalidFields) > 0 {
		return &RoleValidationError{Fields: invalidFields}
	}
	return nil
}
