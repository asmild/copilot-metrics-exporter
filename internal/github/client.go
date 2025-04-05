package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/asmild/copilot-metrics-exporter/internal/auth"
	"github.com/asmild/copilot-metrics-exporter/internal/config"
	"net/http"
	"time"
)

type Client struct {
	httpClient   *http.Client
	baseURL      string
	auth         auth.Provider
	orgName      string
	isEnterprise bool
}

func NewClient(cfg *config.Config) (*Client, error) {
	authProvider, err := auth.NewAuthProvider(cfg)
	if err != nil {
		return nil, err
	}

	org := "orgs"
	if cfg.IsEnterprise { // Adjust if enterprise URL is different
		org = "enterprises"
	}

	baseURL := fmt.Sprintf("https://api.github.com/%s/%s", org, cfg.Organization)

	return &Client{
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		baseURL:      baseURL,
		auth:         authProvider,
		orgName:      cfg.Organization,
		isEnterprise: cfg.IsEnterprise,
	}, nil
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	token, err := c.auth.GetToken()
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	return c.httpClient.Do(req)
}

// Other methods that make API calls remain unchanged except for error handling
func (c *Client) get(endpoint string) (*http.Response, error) {
	// Create a GET request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

func (c *Client) post(endpoint string, data interface{}) (*http.Response, error) {
	// Marshal data into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Create a POST request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}
