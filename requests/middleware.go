package requests

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
	sharedSecret string,
	accessTokenSecret string,
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

			// Decode claimsB64 from base64 to get the actual string
			claimsBytes, err := base64.RawStdEncoding.DecodeString(claimsB64)
			if err != nil {
				WriteJSON(w, http.StatusUnauthorized, APIResponse{Success: false, Error: "invalid auth claims encoding"})
				return
			}
			claimsStr := string(claimsBytes)

			// Check if claimsStr is a JWT (header.payload.signature) or plain JSON
			var claims Claims
			if strings.Count(claimsStr, ".") == 2 {
				// Parse and verify JWT
				token, err := jwt.ParseWithClaims(claimsStr, &claims, func(t *jwt.Token) (interface{}, error) {
					// Verify we are using HMAC signing method
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
					}
					return accessTokenSecret, nil
				})

				if err != nil || !token.Valid {
					WriteJSON(w, http.StatusUnauthorized, APIResponse{Success: false, Error: "invalid JWT token"})
					return
				}
			} else {
				// Plain JSON - unmarshal directly
				if err := json.Unmarshal(claimsBytes, &claims); err != nil {
					WriteJSON(w, http.StatusUnauthorized, APIResponse{Success: false, Error: "invalid auth claims"})
					return
				}
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
