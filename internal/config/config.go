package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

const defaultPort = "9080"

type Config struct {
	Organization        string `yaml:"org"`
	PersonalAccessToken string `yaml:"pat"`
	GitHubAppToken      string `yaml:"app_token"`
	Port                string `yaml:"port"`
	IsEnterprise        bool   `yaml:"is_enterprise"`
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
		appToken := os.Getenv("GITHUB_APP_TOKEN")
		port := os.Getenv("PORT")

		if org != "" && (token != "" || appToken != "") {
			if port != "" {
				config.Port = port
			}

			config.Organization = org
			config.PersonalAccessToken = token
			config.GitHubAppToken = appToken
			config.IsEnterprise = isEnterprise == "true"
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
