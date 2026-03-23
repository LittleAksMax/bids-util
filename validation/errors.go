package validation

import (
	"strings"
)

// ValidationError represents validation errors.
type ValidationError struct {
	Fields []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.Fields, ", ") + " required"
}

// EmailValidationError represents email validation errors.
type EmailValidationError struct {
	Fields []string
}

func (e *EmailValidationError) Error() string {
	return strings.Join(e.Fields, ", ") + " must be valid email address(es)"
}

// UUIDValidationError represents UUID validation errors.
type UUIDValidationError struct {
	Fields []string
}

func (e *UUIDValidationError) Error() string {
	return strings.Join(e.Fields, ", ") + " must be valid UUID(s)"
}

// PasswordValidationError represents password validation errors.
type PasswordValidationError struct {
	Fields []string
}

func (e *PasswordValidationError) Error() string {
	return strings.Join(e.Fields, ", ") + " must be at least 8 characters"
}

// RoleValidationError represents role validation errors.
type RoleValidationError struct {
	Fields []string
}

func (e *RoleValidationError) Error() string {
	return strings.Join(e.Fields, ", ") + " must be one of: user, admin"
}
