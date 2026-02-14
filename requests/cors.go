package requests

import (
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// ApplyCORS configures CORS with an allow-all default or a restricted list of origins.
// Pass the ALLOWED_ORIGINS values from config. Use "*" to allow all.
func ApplyCORS(r chi.Router, allowedOrigins []string, allowedMethods []string, allowedHeaders []string, exposedHeaders []string, allowCredentials bool, maxAge int) {
	allowAll := false
	for _, o := range allowedOrigins {
		if strings.TrimSpace(o) == "*" {
			allowAll = true
			break
		}
	}

	c := cors.Handler(cors.Options{
		AllowedOrigins: func() []string {
			if allowAll {
				return []string{"*"}
			}
			return allowedOrigins
		}(),
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		ExposedHeaders:   exposedHeaders,
		AllowCredentials: allowCredentials,
		MaxAge:           300,
	})

	r.Use(c)
}