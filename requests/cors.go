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
		AllowedMethods:   allowedMethods,//[]string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   allowedHeaders,//[]string{"Accept", "Authorization", "Content-Type", "X-Auth-Claims", "X-Auth-Ts", "X-Auth-Sig"},
		ExposedHeaders:   exposedHeaders,//[]string{"Set-Cookie"},
		AllowCredentials: allowCredentials,//true,
		MaxAge:           300,
	})

	r.Use(c)
}