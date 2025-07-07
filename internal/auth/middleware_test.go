package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asmild/copilot-metrics-exporter/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func TestBasicAuthMiddleware(t *testing.T) {
	password := "testpass"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate test hash: %v", err)
	}

	tests := []struct {
		name           string
		config         *config.BasicAuthConfig
		username       string
		password       string
		expectedStatus int
	}{
		{
			name:           "No auth configured",
			config:         nil,
			username:       "",
			password:       "",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Auth disabled",
			config: &config.BasicAuthConfig{
				Enabled:  false,
				Username: "",
				Password: "",
			},
			username:       "",
			password:       "",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Valid credentials",
			config: &config.BasicAuthConfig{
				Enabled:  true,
				Username: "admin",
				Password: string(passwordHash),
			},
			username:       "admin",
			password:       password,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid username",
			config: &config.BasicAuthConfig{
				Enabled:  true,
				Username: "admin",
				Password: string(passwordHash),
			},
			username:       "wrong",
			password:       string(passwordHash),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid password",
			config: &config.BasicAuthConfig{
				Enabled:  true,
				Username: "admin",
				Password: string(passwordHash),
			},
			username:       "admin",
			password:       "wrongpass",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "No credentials provided",
			config: &config.BasicAuthConfig{
				Enabled:  true,
				Username: "admin",
				Password: string(passwordHash),
			},
			username:       "",
			password:       "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			// Wrap with basic auth middleware
			middleware := BasicAuthMiddleware(tt.config)
			wrappedHandler := middleware(handler)

			// Create a test request
			req := httptest.NewRequest("GET", "/metrics", nil)
			if tt.username != "" || tt.password != "" {
				req.SetBasicAuth(tt.username, tt.password)
			}

			// Create a response recorder
			w := httptest.NewRecorder()

			// Call the handler
			wrappedHandler.ServeHTTP(w, req)

			// Check the response
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check WWW-Authenticate header is set for unauthorized requests
			if tt.expectedStatus == http.StatusUnauthorized {
				if w.Header().Get("WWW-Authenticate") == "" {
					t.Error("Expected WWW-Authenticate header to be set")
				}
			}
		})
	}
}
