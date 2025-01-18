package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	TEST_ORG    = "test-org"
	TEST_TOKEN  = "test-token"
	TEST_PORT   = "9090"
	TEST_IS_ENT = "true"
)

func TestGetConfigFromEnv(t *testing.T) {
	// Set environment variables for the test

	os.Setenv("GITHUB_ORG", TEST_ORG)
	os.Setenv("GITHUB_TOKEN", TEST_TOKEN)
	os.Setenv("GITHUB_IS_ENTERPRISE", TEST_IS_ENT)
	os.Setenv("PORT", TEST_PORT)

	// Check that the configuration is correctly retrieved from the environment
	var emptyString string
	config, err := MustLoad(&emptyString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Reset environment variables
	os.Setenv("GITHUB_ORG", "")
	os.Setenv("GITHUB_TOKEN", "")
	os.Setenv("GITHUB_IS_ENTERPRISE", "")
	os.Setenv("PORT", "")

	assert.Equal(t, TEST_ORG, config.Organization)
	assert.Equal(t, TEST_TOKEN, config.PersonalAccessToken)
	assert.Equal(t, TEST_PORT, config.Port)
	assert.True(t, config.IsEnterprise)
}

func TestGetConfigFromFile(t *testing.T) {
	// Create a temporary configuration file
	yamlContent := fmt.Sprintf(`
org: %s
pat: %s
port: "%s"
is_enterprise: %s
`, TEST_ORG, TEST_TOKEN, TEST_PORT, TEST_IS_ENT)

	tmpFile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Delete the file after the test

	// Write YAML to the file
	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Check that the configuration is correctly retrieved from the file
	tmpFilePath := tmpFile.Name()
	config, err := MustLoad(&tmpFilePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, TEST_ORG, config.Organization)
	assert.Equal(t, TEST_TOKEN, config.PersonalAccessToken)
	assert.Equal(t, TEST_PORT, config.Port)
	assert.True(t, config.IsEnterprise)
}

func TestMustLoad_MissingConfigFile(t *testing.T) {
	// Try to get the configuration from a non-existent file
	nonexistentFile := "/path/to/nonexistent/config.yaml"
	_, err := MustLoad(&nonexistentFile)
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	assert.Contains(t, err.Error(), "config file not found")
}

func TestMustLoad_DefaultPort(t *testing.T) {
	// Create a configuration file without specifying the port
	yamlContent := fmt.Sprintf(`
org: %s
pat: %s
is_enterprise: %s
`, TEST_ORG, TEST_TOKEN, TEST_IS_ENT)

	tmpFile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Delete the file after the test

	// Write YAML to the file
	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Check that the configuration is correctly retrieved from the file
	tmpFilePath := tmpFile.Name()
	config, err := MustLoad(&tmpFilePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, TEST_ORG, config.Organization)
	assert.Equal(t, TEST_TOKEN, config.PersonalAccessToken)
	assert.Equal(t, "9080", config.Port) // Default port
}

func TestMustLoad_MissingEnvVars(t *testing.T) {
	os.Clearenv()

	configPath := ""
	config, err := MustLoad(&configPath)

	require.Error(t, err)
	assert.Nil(t, config)
}

func TestMustLoad_FallbackToDefaultConfigFile(t *testing.T) {
	// Test loading configuration using default paths
	yamlContent := fmt.Sprintf(`
org: %s
pat: %s
`, TEST_ORG, TEST_TOKEN)

	tmpFile, err := os.Create(defaultConfigPaths[0])
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Delete the file after the test

	// Write YAML to the file
	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Check that the configuration is correctly retrieved from the file
	emptyString := ""
	config, err := MustLoad(&emptyString)

	require.NoError(t, err)
	assert.Equal(t, TEST_ORG, config.Organization)
	assert.Equal(t, TEST_TOKEN, config.PersonalAccessToken)
	assert.Equal(t, defaultPort, config.Port)
	assert.False(t, config.IsEnterprise)
}
