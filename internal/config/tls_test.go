package config

import (
	"testing"
)

func TestValidateTLSConfig(t *testing.T) {
	certFile := "testdata/snakeoil_cert.pem"
	keyFile := "testdata/snakeoil_key.pem"

	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "TLS disabled",
			config: &Config{
				TLS: &TLSConfig{
					Enabled: false,
				},
			},
			expectError: false,
		},
		{
			name: "TLS nil",
			config: &Config{
				TLS: nil,
			},
			expectError: false,
		},
		{
			name: "Valid TLS config",
			config: &Config{
				TLS: &TLSConfig{
					Enabled:  true,
					CertFile: certFile,
					KeyFile:  keyFile,
				},
			},
			expectError: false,
		},
		{
			name: "Missing cert file",
			config: &Config{
				TLS: &TLSConfig{
					Enabled:  true,
					CertFile: "",
					KeyFile:  keyFile,
				},
			},
			expectError: true,
		},
		{
			name: "Missing key file",
			config: &Config{
				TLS: &TLSConfig{
					Enabled:  true,
					CertFile: certFile,
					KeyFile:  "",
				},
			},
			expectError: true,
		},
		{
			name: "Nonexistent cert file",
			config: &Config{
				TLS: &TLSConfig{
					Enabled:  true,
					CertFile: "/path/to/nonexistent/cert.pem",
					KeyFile:  "/path/to/nonexistent/key.pem",
				},
			},
			expectError: true,
		},
		{
			name: "Invalid cert file path",
			config: &Config{
				TLS: &TLSConfig{
					Enabled:  true,
					CertFile: "/path/to/cert.pem",
					KeyFile:  "/path/to/key.pem",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTLSConfig(tt.config)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
