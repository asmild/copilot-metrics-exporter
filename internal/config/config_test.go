package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// TLS Configuration Tests

func TestGetConfigFromEnv_WithTLS(t *testing.T) {
	certFile := "testdata/snakeoil_cert.pem"
	keyFile := "testdata/snakeoil_key.pem"

	// Set environment variables for the test
	os.Setenv("GITHUB_ORG", TEST_ORG)
	os.Setenv("GITHUB_TOKEN", TEST_TOKEN)
	os.Setenv("TLS_ENABLED", "true")
	os.Setenv("TLS_CERT_FILE", certFile)
	os.Setenv("TLS_KEY_FILE", keyFile)

	defer func() {
		os.Unsetenv("GITHUB_ORG")
		os.Unsetenv("GITHUB_TOKEN")
		os.Unsetenv("TLS_ENABLED")
		os.Unsetenv("TLS_CERT_FILE")
		os.Unsetenv("TLS_KEY_FILE")
	}()

	var emptyString string
	config, err := MustLoad(&emptyString)
	require.NoError(t, err)

	assert.Equal(t, TEST_ORG, config.Organization)
	assert.Equal(t, TEST_TOKEN, config.PersonalAccessToken)
	assert.NotNil(t, config.TLS)
	assert.True(t, config.TLS.Enabled)
	assert.Equal(t, certFile, config.TLS.CertFile)
	assert.Equal(t, keyFile, config.TLS.KeyFile)
}

func TestGetConfigFromFile_WithTLS(t *testing.T) {
	certFile := "testdata/snakeoil_cert.pem"
	keyFile := "testdata/snakeoil_key.pem"

	// Create a configuration file with TLS settings
	yamlContent := fmt.Sprintf(`
org: %s
pat: %s
port: "%s"
is_enterprise: %s
tls:
  enabled: true
  cert_file: %s
  key_file: %s
`, TEST_ORG, TEST_TOKEN, TEST_PORT, TEST_IS_ENT, certFile, keyFile)

	tmpFile, err := os.CreateTemp("", "config*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Write YAML to the file
	_, err = tmpFile.Write([]byte(yamlContent))
	require.NoError(t, err)
	tmpFile.Close()

	// Check that the configuration is correctly retrieved from the file
	tmpFilePath := tmpFile.Name()
	config, err := MustLoad(&tmpFilePath)
	require.NoError(t, err)

	assert.Equal(t, TEST_ORG, config.Organization)
	assert.Equal(t, TEST_TOKEN, config.PersonalAccessToken)
	assert.Equal(t, TEST_PORT, config.Port)
	assert.True(t, config.IsEnterprise)
	assert.NotNil(t, config.TLS)
	assert.True(t, config.TLS.Enabled)
	assert.Equal(t, certFile, config.TLS.CertFile)
	assert.Equal(t, keyFile, config.TLS.KeyFile)
}

func TestTLSValidation_MissingCertFile(t *testing.T) {
	yamlContent := fmt.Sprintf(`
org: %s
pat: %s
tls:
  enabled: true
  key_file: /path/to/key.pem
`, TEST_ORG, TEST_TOKEN)

	tmpFile, err := os.CreateTemp("", "config*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(yamlContent))
	require.NoError(t, err)
	tmpFile.Close()

	tmpFilePath := tmpFile.Name()
	_, err = MustLoad(&tmpFilePath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load TLS certificate and key")
}

func TestTLSValidation_MissingKeyFile(t *testing.T) {
	yamlContent := fmt.Sprintf(`
org: %s
pat: %s
tls:
  enabled: true
  cert_file: /path/to/cert.pem
`, TEST_ORG, TEST_TOKEN)

	tmpFile, err := os.CreateTemp("", "config*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(yamlContent))
	require.NoError(t, err)
	tmpFile.Close()

	tmpFilePath := tmpFile.Name()
	_, err = MustLoad(&tmpFilePath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load TLS certificate and key")
}

func TestTLSValidation_NonexistentCertFile(t *testing.T) {
	yamlContent := fmt.Sprintf(`
org: %s
pat: %s
tls:
  enabled: true
  cert_file: /path/to/nonexistent/cert.pem
  key_file: /path/to/nonexistent/key.pem
`, TEST_ORG, TEST_TOKEN)

	tmpFile, err := os.CreateTemp("", "config*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(yamlContent))
	require.NoError(t, err)
	tmpFile.Close()

	tmpFilePath := tmpFile.Name()
	_, err = MustLoad(&tmpFilePath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load TLS certificate and key")
}

func TestTLSValidation_EmptyCertFile(t *testing.T) {
	// Create a temporary empty certificate file
	// This simulates a case where the cert file exists but is empty
	certFile, err := os.CreateTemp("", "empty_cert*.pem")
	require.NoError(t, err)
	defer os.Remove(certFile.Name())
	keyFile := "testdata/snakeoil_key.pem"

	yamlContent := fmt.Sprintf(`
org: %s
pat: %s
tls:
  enabled: true
  cert_file: %s
  key_file: %s
`, TEST_ORG, TEST_TOKEN, certFile.Name(), keyFile)

	tmpFile, err := os.CreateTemp("", "config*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(yamlContent))
	require.NoError(t, err)
	tmpFile.Close()

	tmpFilePath := tmpFile.Name()
	_, err = MustLoad(&tmpFilePath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find any PEM data in certificate input")
}

func TestTLSValidation_DisabledTLS(t *testing.T) {
	yamlContent := fmt.Sprintf(`
org: %s
pat: %s
tls:
  enabled: false
`, TEST_ORG, TEST_TOKEN)

	tmpFile, err := os.CreateTemp("", "config*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(yamlContent))
	require.NoError(t, err)
	tmpFile.Close()

	tmpFilePath := tmpFile.Name()
	config, err := MustLoad(&tmpFilePath)
	require.NoError(t, err)

	assert.NotNil(t, config.TLS)
	assert.False(t, config.TLS.Enabled)
}

func TestTLSValidation_NoTLSConfig(t *testing.T) {
	yamlContent := fmt.Sprintf(`
org: %s
pat: %s
`, TEST_ORG, TEST_TOKEN)

	tmpFile, err := os.CreateTemp("", "config*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(yamlContent))
	require.NoError(t, err)
	tmpFile.Close()

	tmpFilePath := tmpFile.Name()
	config, err := MustLoad(&tmpFilePath)
	require.NoError(t, err)

	assert.Nil(t, config.TLS)
}
