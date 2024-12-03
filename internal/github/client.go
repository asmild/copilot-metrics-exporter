package github

import (
	"encoding/json"
	"fmt"
	"github.com/asmild/copilot-metrics-exporter/config"
	"github.com/asmild/copilot-metrics-exporter/internal/requests"
	"net/http"
)

type GitHubClient struct {
	token        string
	appToken     string
	Organization string
	client       *http.Client
	baseApiUrl   string
}

func NewGitHubClient(conf config.ExporterConfig) (*GitHubClient, error) {
	client := &http.Client{}
	endpoint := "orgs"
	if conf.IsEnterprise {
		endpoint = "enterprises"
	}
	return &GitHubClient{
		token:        conf.PersonalAccessToken,
		Organization: conf.Organization,
		client:       client,
		baseApiUrl:   fmt.Sprintf("https://api.github.com/%s/%s", endpoint, conf.Organization),
	}, nil
}

func (c *GitHubClient) makeRequest(method, endpoint string, data interface{}) (*http.Response, error) {
	headers := map[string]string{
		"Accept":               "application/vnd.github+json",
		"X-GitHub-Api-Version": "2022-11-28",
	}

	if c.token != "" {
		headers["Authorization"] = "Bearer " + c.token
	} else if c.appToken != "" {
		headers["Authorization"] = "Bearer " + c.appToken
	}

	url := fmt.Sprintf("%s/%s", c.baseApiUrl, endpoint)
	res, err := requests.HttpRequester(c.client, url, headers, method, data)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var errorMessage struct {
			Message          string `json:"message"`
			DocumentationURL string `json:"documentation_url"`
		}
		err = json.NewDecoder(res.Body).Decode(&errorMessage)
		if err != nil {
			return nil, fmt.Errorf("HTTP request failed with status %d: %s", res.StatusCode, res.Status)
		}
		return nil, fmt.Errorf("HTTP request to %s failed with status %d: %s (%s)", url, res.StatusCode, errorMessage.Message, errorMessage.DocumentationURL)
	}

	return res, nil
}

func (c *GitHubClient) get(endpoint string) (*http.Response, error) {
	return c.makeRequest("GET", endpoint, nil)
}

func (c *GitHubClient) post(endpoint string, data interface{}) (*http.Response, error) {
	return c.makeRequest("POST", endpoint, data)
}
