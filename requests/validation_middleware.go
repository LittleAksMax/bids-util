package requests

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ContextKey type for context keys to avoid collisions.
type ContextKey string

const RequestBodyKey ContextKey = "requestBody"

// RegisterMiddleware attaches common middleware to the router.
func RegisterMiddleware(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
}

// ValidateRequest is a generic middleware that validates request body fields based on a list of validation functions.
// Decodes the JSON body into the provided type T and validates all marked fields.
func ValidateRequest[T any](validationFuncs []func(T any) error) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var reqValue T
			if err := json.NewDecoder(r.Body).Decode(&reqValue); err != nil {
				WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request body"})
				return
			}

			if err := validateRequestBody(&reqValue, validationFuncs); err != nil {
				WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
				return
			}

			ctx := context.WithValue(r.Context(), RequestBodyKey, &reqValue)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validateRequestBody(value any, validationFuncs []func(any) error) error {
	reflectValue := reflect.ValueOf(value)
	for reflectValue.Kind() == reflect.Ptr {
		if reflectValue.IsNil() {
			return applyValidations(value, validationFuncs)
		}
		reflectValue = reflectValue.Elem()
	}

	if reflectValue.Kind() != reflect.Slice && reflectValue.Kind() != reflect.Array {
		return applyValidations(value, validationFuncs)
	}

	for i := 0; i < reflectValue.Len(); i++ {
		if err := applyValidations(reflectValue.Index(i).Interface(), validationFuncs); err != nil {
			return err
		}
	}

	return nil
}

func applyValidations(value any, validationFuncs []func(any) error) error {
	for _, validate := range validationFuncs {
		if err := validate(value); err != nil {
			return err
		}
	}

	return nil
}

// GetRequestBody retrieves the validated request body from the context.
func GetRequestBody[T any](r *http.Request) *T {
	if body := r.Context().Value(RequestBodyKey); body != nil {
		if typedBody, ok := body.(*T); ok {
			return typedBody
		}
	}
	return nil
}
