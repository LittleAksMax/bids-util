package requests

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ValidateAccessToken validates Forward-Auth headers and injects claims into request context.
func ValidateAccessToken[TClaims jwt.RegisteredClaims](
	sharedSecret []byte,
	maxSkew time.Duration,
	authClaimsKey string,
	authTimestampKey string,
	authSigKey string,
) func(http.Handler) http.Handler {
	secretBytes := []byte(sharedSecret)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claimsB64 := r.Header.Get(authClaimsKey)
			tsStr := r.Header.Get(authTimestampKey)
			sigB64 := r.Header.Get(authSigKey)

			if claimsB64 == "" || tsStr == "" || sigB64 == "" {
				WriteJSON(w, http.StatusUnauthorized, APIResponse{Success: false, Error: "missing auth headers"})
				return
			}

			if err := VerifyAuthSignature(secretBytes, tsStr, claimsB64, sigB64, maxSkew); err != nil {
				WriteJSON(w, http.StatusUnauthorized, APIResponse{Success: false, Error: "invalid auth signature"})
				return
			}

			claimsJSON, err := base64.RawStdEncoding.DecodeString(claimsB64)
			if err != nil {
				WriteJSON(w, http.StatusUnauthorized, APIResponse{Success: false, Error: "invalid auth claims encoding"})
				return
			}

			var claims TClaims
			if err := json.Unmarshal(claimsJSON, &claims); err != nil {
				WriteJSON(w, http.StatusUnauthorized, APIResponse{Success: false, Error: "invalid auth claims"})
				return
			}

			ctx := context.WithValue(r.Context(), authClaimsKey, &claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAPIKey is a middleware that validates API key header
func RequireAPIKey(apiKey, apiKeyKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			providedKey := r.Header.Get(apiKeyKey)
			if providedKey == "" || providedKey != apiKey {
				WriteJSON(w, http.StatusUnauthorized, APIResponse{Success: false, Error: "invalid or missing API key"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
