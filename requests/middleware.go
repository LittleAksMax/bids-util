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
	authClaimsKey string,
	maxSkew time.Duration,
) func(http.Handler) http.Handler {
	secretBytes := []byte(sharedSecret)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claimsB64 := r.Header.Get("X-Auth-Claims")
			tsStr := r.Header.Get("X-Auth-Ts")
			sigB64 := r.Header.Get("X-Auth-Sig")

			if claimsB64 == "" || tsStr == "" || sigB64 == "" {
				writeJSON(w, http.StatusUnauthorized, ApiResponse{Success: false, Error: "missing auth headers"})
				return
			}

			if err := VerifyAuthSignature(secretBytes, tsStr, claimsB64, sigB64, maxSkew); err != nil {
				writeJSON(w, http.StatusUnauthorized, ApiResponse{Success: false, Error: "invalid auth signature"})
				return
			}

			claimsJSON, err := base64.RawStdEncoding.DecodeString(claimsB64)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, ApiResponse{Success: false, Error: "invalid auth claims encoding"})
				return
			}

			var claims TClaims
			if err := json.Unmarshal(claimsJSON, &claims); err != nil {
				writeJSON(w, http.StatusUnauthorized, ApiResponse{Success: false, Error: "invalid auth claims"})
				return
			}

			ctx := context.WithValue(r.Context(), authClaimsKey, &claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}