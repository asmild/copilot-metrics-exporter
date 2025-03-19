package client

import (
 "net/http"
 "time"

 "github.com/yourusername/copilot-metrics-exporter/internal/auth"
 "github.com/yourusername/copilot-metrics-exporter/internal/config"
)

type Client struct {
 httpClient   *http.Client
 baseURL      string
 auth         auth.AuthProvider
 orgName      string
 isEnterprise bool
}

func New(cfg *config.Config) (*Client, error) {
 authProvider, err := auth.NewAuthProvider(cfg)
 if err != nil {
  return nil, err
 }

 baseURL := "https://api.github.com"
 if cfg.IsEnterprise {
  baseURL = "https://api.github.com/api/v3" // Adjust if enterprise URL is different
 }

 return &Client{
  httpClient:   &http.Client{Timeout: 30 * time.Second},
  baseURL:      baseURL,
  auth:         authProvider,
  orgName:      cfg.Organization,
  isEnterprise: cfg.IsEnterprise,
 }, nil
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
 token, err := c.auth.Token()
 if err != nil {
  return nil, err
 }

 req.Header.Set("Authorization", "Bearer "+token)
 req.Header.Set("Accept", "application/vnd.github.v3+json")

 return c.httpClient.Do(req)
}

// Other methods that make API calls remain unchanged except for error handling
