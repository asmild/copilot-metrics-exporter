package config

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestValidateBasicAuthConfig(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate test hash: %v", err)
	}

	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "Basic auth disabled",
			config: &Config{
				BasicAuth: &BasicAuthConfig{
					Enabled: false,
				},
			},
			expectError: false,
		},
		{
			name: "Basic auth nil",
			config: &Config{
				BasicAuth: nil,
			},
			expectError: false,
		},
		{
			name: "Valid basic auth config",
			config: &Config{
				BasicAuth: &BasicAuthConfig{
					Enabled:  true,
					Username: "admin",
					Password: string(passwordHash),
				},
			},
			expectError: false,
		},
		{
			name: "Missing username",
			config: &Config{
				BasicAuth: &BasicAuthConfig{
					Enabled:  true,
					Username: "",
					Password: string(passwordHash),
				},
			},
			expectError: true,
		},
		{
			name: "Missing password",
			config: &Config{
				BasicAuth: &BasicAuthConfig{
					Enabled:  true,
					Username: "admin",
					Password: "",
				},
			},
			expectError: true,
		},
		{
			name: "Invalid password format",
			config: &Config{
				BasicAuth: &BasicAuthConfig{
					Enabled:  true,
					Username: "admin",
					Password: "plaintext",
				},
			},
			expectError: true,
		},
		{
			name: "Too short bcrypt string",
			config: &Config{
				BasicAuth: &BasicAuthConfig{
					Enabled:  true,
					Username: "admin",
					Password: "$2a$12$short",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBasicAuthConfig(tt.config)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
