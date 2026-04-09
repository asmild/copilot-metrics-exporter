package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asmild/copilot-metrics-exporter/internal/config"
	gogithub "github.com/google/go-github/v84/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	conf := config.Config{
		PersonalAccessToken: "test-token",
		Organization:        "test-org",
		IsEnterprise:        false,
	}
	client, err := NewClient(&conf)
	require.NoError(t, err)
	assert.Equal(t, "test-org", client.org)
	assert.False(t, client.isEnterprise)
	assert.NotNil(t, client.gh)
}

func TestNewClientEnterprise(t *testing.T) {
	conf := config.Config{
		PersonalAccessToken: "test-token",
		Organization:        "test-enterprise",
		IsEnterprise:        true,
	}
	client, err := NewClient(&conf)
	require.NoError(t, err)
	assert.Equal(t, "test-enterprise", client.org)
	assert.True(t, client.isEnterprise)
}

func TestNewClientNoAuth(t *testing.T) {
	conf := config.Config{
		Organization: "test-org",
	}
	_, err := NewClient(&conf)
	assert.Error(t, err)
}

func TestGetMetrics(t *testing.T) {
	reportData := UsageReport{
		Day:                         "2026-04-08",
		DailyActiveUsers:            116,
		CodeGenerationActivityCount: 2164,
		CodeAcceptanceActivityCount: 546,
		LocSuggestedToAddSum:        4170,
		LocAddedSum:                 1353,
	}

	// Test server for report download (pre-signed URL)
	reportServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(reportData)
	}))
	defer reportServer.Close()

	// Test server for GitHub API (report endpoint returns download links)
	reportResponse := gogithub.CopilotDailyMetricsReport{
		DownloadLinks: []string{reportServer.URL + "/report.json"},
		ReportDay:     "2026-04-08",
	}

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/orgs/test-org/copilot/metrics/reports/organization-1-day")
		json.NewEncoder(w).Encode(reportResponse)
	}))
	defer apiServer.Close()

	conf := config.Config{
		PersonalAccessToken: "test-token",
		Organization:        "test-org",
		IsEnterprise:        false,
	}
	client, err := NewClient(&conf)
	require.NoError(t, err)

	client.gh, err = client.gh.WithEnterpriseURLs(apiServer.URL+"/", apiServer.URL+"/")
	require.NoError(t, err)

	report, err := client.GetMetrics(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "2026-04-08", report.Day)
	assert.Equal(t, 116, report.DailyActiveUsers)
	assert.Equal(t, 2164, report.CodeGenerationActivityCount)
	assert.Equal(t, 1353, report.LocAddedSum)
}

func TestGetTotalSeats(t *testing.T) {
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/orgs/test-org/copilot/billing/seats")
		resp := gogithub.ListCopilotSeatsResponse{
			TotalSeats: 424,
			Seats: []*gogithub.CopilotSeatDetails{
				{CreatedAt: &gogithub.Timestamp{}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer apiServer.Close()

	conf := config.Config{
		PersonalAccessToken: "test-token",
		Organization:        "test-org",
		IsEnterprise:        false,
	}
	client, err := NewClient(&conf)
	require.NoError(t, err)

	client.gh, err = client.gh.WithEnterpriseURLs(apiServer.URL+"/", apiServer.URL+"/")
	require.NoError(t, err)

	total, err := client.GetTotalSeats(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(424), total)
}