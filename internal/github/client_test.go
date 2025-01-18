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
	client, err := NewGitHubClient(conf)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if client.token != "test-token" {
		t.Errorf("Expected token 'test-token', but got: %s", client.token)
	}
	if client.Organization != "test-org" {
		t.Errorf("Expected organization 'test-org', but got: %s", client.Organization)
	}
	if client.baseApiUrl != "https://api.github.com/orgs/test-org" {
		t.Errorf("Expected baseApiUrl 'https://api.github.com/orgs/test-org', but got: %s", client.baseApiUrl)
	}
}
