package requests

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents access token payload.
type Claims struct {
	Role string `json:"role"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

// ValidateAccessToken validates Forward-Auth headers and injects claims into request context.
func ValidateAccessToken(
	sharedSecret []byte,
	maxSkew time.Duration,
	authClaimsKey string,
	authTimestampKey string,
	authSigKey string,
) func(http.Handler) http.Handler {
	secretBytes := sharedSecret
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

			// claimsB64 is actually a full JWT (header.payload.signature), not just base64-encoded JSON
			// We need to parse and verify the JWT signature to extract the payload
			var claims Claims
			if strings.Count(claimsB64, ".") == 2 {
				// Parse and verify JWT
				token, err := jwt.ParseWithClaims(claimsB64, &claims, func(t *jwt.Token) (interface{}, error) {
					// Verify it's using HMAC signing method
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
					}
					return secretBytes, nil
				})

				if err != nil || !token.Valid {
					WriteJSON(w, http.StatusUnauthorized, APIResponse{Success: false, Error: "invalid JWT token"})
					return
				}
			} else {
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
