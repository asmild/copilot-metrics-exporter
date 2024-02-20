package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type ExporterConfig struct {
	Organization        string `yaml:"org"`
	PersonalAccessToken string `yaml:"pat"`
	Port                string `yaml:"port"`
}

var defaultConfigPaths = []string{
	"./config.yaml",
	"./config.yml",
	"~/.copilot-exporter/config.yaml",
	"~/.copilot-exporter/config.yml",
}

var defaultPort = 9080

func GetConfig(configPath string) (*ExporterConfig, error) {
	if configPath == "" {
		for _, p := range defaultConfigPaths {
			if _, err := os.Stat(p); err == nil {
				configPath = p
				break
			}
		}
		if configPath == "" {
			return nil, fmt.Errorf("config file not found in default paths: %v", defaultConfigPaths)
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ExporterConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML config data: %w", err)
	}

	if config.Port == "" {
		config.Port = fmt.Sprintf("%d", defaultPort)
	}

	return &config, nil
}
