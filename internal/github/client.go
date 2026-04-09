package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/asmild/copilot-metrics-exporter/internal/auth"
	"github.com/asmild/copilot-metrics-exporter/internal/config"
	gogithub "github.com/google/go-github/v84/github"
)

// authTransport is an http.RoundTripper that injects the auth token into requests.
type authTransport struct {
	provider auth.Provider
	base     http.RoundTripper
}

// RoundTrip injects the auth token into every outgoing request.
// This allows go-github to work with both PAT and GitHub App auth transparently,
// since GitHub App tokens expire and need to be refreshed via auth.Provider.GetToken().
func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := t.provider.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	// The new report endpoints require API version 2026-03-10;
	// billing endpoints use go-github's default (2022-11-28).
	// The new Copilot endpoints require API version 2026-03-10
	if strings.Contains(req.URL.Path, "/copilot/") {
		req.Header.Set("X-GitHub-Api-Version", "2026-03-10")
	}
	return t.base.RoundTrip(req)
}

// Client wraps go-github's client for Copilot metrics collection.
type Client struct {
	gh           *gogithub.Client
	reportClient *http.Client // plain client for downloading pre-signed report URLs (no auth)
	org          string
	isEnterprise bool
}

// NewClient creates a new GitHub client using go-github with the configured auth provider.
func NewClient(cfg *config.Config) (*Client, error) {
	authProvider, err := auth.NewAuthProvider(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth provider: %w", err)
	}

	transport := &authTransport{
		provider: authProvider,
		base:     http.DefaultTransport,
	}

	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	ghClient := gogithub.NewClient(httpClient)

	// Separate plain HTTP client for downloading pre-signed report URLs from Azure Blob Storage.
	// These URLs use SAS tokens for auth and must NOT include the GitHub Bearer token.
	reportClient := &http.Client{Timeout: 60 * time.Second}

	return &Client{
		gh:           ghClient,
		reportClient: reportClient,
		org:          cfg.Organization,
		isEnterprise: cfg.IsEnterprise,
	}, nil
}

// GetMetrics fetches Copilot usage metrics via the new report endpoints (apiVersion 2026-03-10).
func (c *Client) GetMetrics(ctx context.Context) (*UsageReport, error) {
	downloadLinks, err := c.getReportDownloadLinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics report: %w", err)
	}

	// The 1-day report typically has a single download link with a single report object.
	for _, link := range downloadLinks {
		report, err := c.downloadAndParseReport(ctx, link)
		if err != nil {
			return nil, fmt.Errorf("failed to download report from %s: %w", link, err)
		}
		if report != nil {
			return report, nil
		}
	}

	return nil, fmt.Errorf("no report data found in download links")
}

// getReportDownloadLinks calls the 1-day report endpoint for yesterday (UTC)
// and returns the download links. Falls back to the day before if yesterday's
// data is not yet available.
// Fix for https://github.com/asmild/copilot-metrics-exporter/issues/19:
// Uses UTC to avoid timezone mismatches, and falls back to day-before-yesterday
// in case GitHub hasn't finished processing yesterday's metrics yet.
func (c *Client) getReportDownloadLinks(ctx context.Context) ([]string, error) {
	yesterday := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02")

	links, err := c.fetchDailyReport(ctx, yesterday)
	if err != nil {
		// Yesterday's data may not be ready yet; try the day before
		dayBefore := time.Now().UTC().AddDate(0, 0, -2).Format("2006-01-02")
		links, err = c.fetchDailyReport(ctx, dayBefore)
		if err != nil {
			return nil, fmt.Errorf("failed to get daily report: %w", err)
		}
	}
	return links, nil
}

func (c *Client) fetchDailyReport(ctx context.Context, day string) ([]string, error) {
	opts := &gogithub.CopilotMetricsReportOptions{Day: day}

	if c.isEnterprise {
		report, _, err := c.gh.Copilot.GetEnterpriseDailyMetricsReport(ctx, c.org, opts)
		if err != nil {
			return nil, err
		}
		return report.DownloadLinks, nil
	}

	report, _, err := c.gh.Copilot.GetOrganizationDailyMetricsReport(ctx, c.org, opts)
	if err != nil {
		return nil, err
	}
	return report.DownloadLinks, nil
}

// downloadAndParseReport downloads a report JSON file and parses it into UsageReport.
func (c *Client) downloadAndParseReport(ctx context.Context, url string) (*UsageReport, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.reportClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download report: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d downloading report", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read report body: %w", err)
	}

	var report UsageReport
	if err := json.Unmarshal(body, &report); err != nil {
		return nil, fmt.Errorf("failed to decode report JSON: %w", err)
	}

	return &report, nil
}

// GetTotalSeats returns the total number of Copilot seats assigned.
func (c *Client) GetTotalSeats(ctx context.Context) (int64, error) {
	opts := &gogithub.ListOptions{Page: 1, PerPage: 1}

	if c.isEnterprise {
		resp, _, err := c.gh.Copilot.ListCopilotEnterpriseSeats(ctx, c.org, opts)
		if err != nil {
			return 0, fmt.Errorf("failed to get enterprise seats: %w", err)
		}
		return resp.TotalSeats, nil
	}

	resp, _, err := c.gh.Copilot.ListCopilotSeats(ctx, c.org, opts)
	if err != nil {
		return 0, fmt.Errorf("failed to get org seats: %w", err)
	}
	return resp.TotalSeats, nil
}
