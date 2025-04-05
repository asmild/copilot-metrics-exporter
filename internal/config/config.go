package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

const defaultPort = "9080"

// GitHubApp holds GitHub App authentication configuration
type GitHubApp struct {
	AppID          int64  `yaml:"app_id"`
	InstallationID int64  `yaml:"installation_id"`
	PrivateKeyPath string `yaml:"private_key_path"`
	PrivateKey     string `yaml:"private_key"`
}

// Config holds the application configuration
type Config struct {
	Organization        string     `yaml:"org"`
	PersonalAccessToken string     `yaml:"pat"`
	GitHubApp           *GitHubApp `yaml:"github_app"`
	Port                string     `yaml:"port"`
	IsEnterprise        bool       `yaml:"is_enterprise"`
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

	return &config, nil
}
