package auth

import (
	"crypto/subtle"
	"net/http"

	"github.com/asmild/copilot-metrics-exporter/internal/config"
	"golang.org/x/crypto/bcrypt"
)

// BasicAuthMiddleware returns a middleware that performs basic authentication
func BasicAuthMiddleware(cfg *config.BasicAuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If basic auth is not enabled, skip authentication
			if cfg == nil || !cfg.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Get basic auth credentials from request
			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Copilot Metrics Exporter"`)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}

			// Check username using constant time comparison
			if subtle.ConstantTimeCompare([]byte(username), []byte(cfg.Username)) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="Copilot Metrics Exporter"`)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}

			// Check password using bcrypt
			if err := bcrypt.CompareHashAndPassword([]byte(cfg.Password), []byte(password)); err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="Copilot Metrics Exporter"`)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}

			// Authentication successful, continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}
