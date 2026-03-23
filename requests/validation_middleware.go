package requests

import (
	"context"
	"encoding/json"
	"net/http"

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

			for _, validate := range validationFuncs {
				if err := validate(&reqValue); err != nil {
					WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
					return
				}
			}

			ctx := context.WithValue(r.Context(), RequestBodyKey, &reqValue)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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
