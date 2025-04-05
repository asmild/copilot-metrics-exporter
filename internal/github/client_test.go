package github

import (
	"github.com/asmild/copilot-metrics-exporter/internal/config"
	"testing"
)

func TestNewGitHubClient(t *testing.T) {
	conf := config.Config{
		PersonalAccessToken: "test-token",
		Organization:        "test-org",
		IsEnterprise:        false,
	}
	client, err := NewClient(&conf)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	token, err := client.auth.GetToken()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if token != "test-token" {
		t.Errorf("Expected token 'test-token', but got: %s", token)
	}

	if client.orgName != "test-org" {
		t.Errorf("Expected organization 'test-org', but got: %s", client.orgName)
	}

	if client.baseURL != "https://api.github.com/orgs/test-org" {
		t.Errorf("Expected baseApiUrl 'https://api.github.com/orgs/test-org', but got: %s", client.baseURL)
	}
}
