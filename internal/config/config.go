package config

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultPort = "9080"

// GitHubApp holds GitHub App authentication configuration
type GitHubApp struct {
	AppID          int64  `yaml:"app_id"`
	InstallationID int64  `yaml:"installation_id"`
	PrivateKeyPath string `yaml:"private_key_path"`
	PrivateKey     string `yaml:"private_key"`
}

// TLSConfig holds TLS certificate configuration
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// BasicAuthConfig holds basic authentication configuration
type BasicAuthConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"` // bcrypt hashed password
}

// Config holds the application configuration
type Config struct {
	Organization        string           `yaml:"org"`
	PersonalAccessToken string           `yaml:"pat"`
	GitHubApp           *GitHubApp       `yaml:"github_app"`
	Port                string           `yaml:"port"`
	IsEnterprise        bool             `yaml:"is_enterprise"`
	TLS                 *TLSConfig       `yaml:"tls"`
	BasicAuth           *BasicAuthConfig `yaml:"basic_auth"`
}

var defaultConfigPaths = []string{
	"./config.yaml",
	"./config.yml",
	"~/.copilot-exporter/config.yaml",
	"~/.copilot-exporter/config.yml",
}

func MustLoad(configPath *string) (*Config, error) {
	var config Config

	if *configPath == "" {
		org := os.Getenv("GITHUB_ORG")
		isEnterprise := os.Getenv("GITHUB_IS_ENTERPRISE")
		token := os.Getenv("GITHUB_TOKEN")
		port := os.Getenv("PORT")

		// TLS configuration from environment variables
		tlsEnabled := os.Getenv("TLS_ENABLED")
		tlsCertFile := os.Getenv("TLS_CERT_FILE")
		tlsKeyFile := os.Getenv("TLS_KEY_FILE")

		// Basic Auth configuration from environment variables
		basicAuthEnabled := os.Getenv("BASIC_AUTH_ENABLED")
		basicAuthUsername := os.Getenv("BASIC_AUTH_USERNAME")
		basicAuthPassword := os.Getenv("BASIC_AUTH_PASSWORD")

		// Check for GitHub App environment variables
		appIDStr := os.Getenv("GITHUB_APP_ID")
		installIDStr := os.Getenv("GITHUB_APP_INSTALLATION_ID")

		if org != "" && (token != "" || (appIDStr != "" && installIDStr != "")) {
			if port != "" {
				config.Port = port
			}

			config.Organization = org
			config.PersonalAccessToken = token
			config.IsEnterprise = isEnterprise == "true"

			// Configure TLS if environment variables are set
			if tlsEnabled == "true" {
				config.TLS = &TLSConfig{
					Enabled:  true,
					CertFile: tlsCertFile,
					KeyFile:  tlsKeyFile,
				}
			}

			// Configure Basic Auth if environment variables are set
			if basicAuthEnabled == "true" && basicAuthUsername != "" && basicAuthPassword != "" {
				config.BasicAuth = &BasicAuthConfig{
					Enabled:  true,
					Username: basicAuthUsername,
					Password: basicAuthPassword, // Expect this to be bcrypt hashed
				}
			}

			// If GitHub App environment variables are provided
			if appIDStr != "" && installIDStr != "" {
				appID, err := strconv.ParseInt(appIDStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid GITHUB_APP_ID: %v", err)
				}

				instID, err := strconv.ParseInt(installIDStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid GITHUB_APP_INSTALLATION_ID: %v", err)
				}

				config.GitHubApp = &GitHubApp{
					AppID:          appID,
					InstallationID: instID,
					PrivateKeyPath: os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH"),
					PrivateKey:     os.Getenv("GITHUB_APP_PRIVATE_KEY"),
				}
			}

			return &config, nil
		}

		for _, defaultConfigPath := range defaultConfigPaths {
			if _, err := os.Stat(defaultConfigPath); err == nil {
				*configPath = defaultConfigPath
				break
			}
		}
	}

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", *configPath)
	}

	data, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to load config from file: %w", err)
	}

	if config.Port == "" {
		config.Port = defaultPort
	}

	// Validate TLS configuration
	if err := validateTLSConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid TLS configuration: %v", err)
	}

	// Validate Basic Auth configuration
	if err := validateBasicAuthConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid Basic Auth configuration: %v", err)
	}

	return &config, nil
}

// validateTLSConfig validates the TLS configuration
func validateTLSConfig(config *Config) error {
	if config.TLS == nil || !config.TLS.Enabled {
		return nil
	}

	// Optionally, you can add more checks here, such as validating the certificate format
	_, err := tls.LoadX509KeyPair(config.TLS.CertFile, config.TLS.KeyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate and key: %v", err)
	}

	return nil
}

// validateBasicAuthConfig validates the basic auth configuration
func validateBasicAuthConfig(config *Config) error {
	if config.BasicAuth == nil || !config.BasicAuth.Enabled {
		return nil
	}

	if config.BasicAuth.Username == "" {
		return fmt.Errorf("username is required when basic auth is enabled")
	}

	if config.BasicAuth.Password == "" {
		return fmt.Errorf("password is required when basic auth is enabled")
	}

	// Check if password looks like a bcrypt hash
	if len(config.BasicAuth.Password) < 60 || !strings.HasPrefix(config.BasicAuth.Password, "$2") {
		return fmt.Errorf("password must be a bcrypt hash (starts with $2 and at least 60 characters)")
	}

	return nil
}
