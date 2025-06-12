package helper

import (
	"github.com/asmild/copilot-metrics-exporter/internal/github"
	"testing"
)

func getExampleMetrics() github.CopilotMetrics {
	return github.CopilotMetrics{
		TotalActiveUsers: 53,
		Date:             "2025-06-13",
		CopilotIDECodeCompletions: github.IDECodeCompletions{
			Languages: []github.Language{
				{
					Name:                "python",
					TotalSuggestions:    200,
					TotalAcceptances:    25,
					TotalLinesSuggested: 500,
					TotalLinesAccepted:  123,
				},
				{
					Name:                "ruby",
					TotalSuggestions:    150,
					TotalAcceptances:    75,
					TotalLinesSuggested: 167,
					TotalLinesAccepted:  100,
				},
			},
			Editors: []github.Editor{
				{
					Name: "vscode",
					Models: []github.Model{
						{
							Name: "default",
							Languages: []github.Language{
								{
									Name:                "python",
									TotalSuggestions:    100,
									TotalAcceptances:    25,
									TotalLinesSuggested: 233,
									TotalLinesAccepted:  283,
								},
								{
									Name:                "ruby",
									TotalSuggestions:    50,
									TotalAcceptances:    25,
									TotalLinesSuggested: 100,
									TotalLinesAccepted:  294,
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestGetLastDayData(t *testing.T) {
	usages := []github.CopilotMetrics{
		{Date: "2025-06-11", TotalActiveUsers: 100},
		{Date: "2025-06-12", TotalActiveUsers: 200},
		{Date: "2025-06-13", TotalActiveUsers: 300},
	}

	lastDayData := GetLastDayData(usages)
	if lastDayData.Date != "2025-06-13" {
		t.Errorf("Expected last day to be 2025-06-13, but got %s", lastDayData.Date)
	}
}

func TestGetTotalSuggestionsCount(t *testing.T) {
	exampleMetrics := getExampleMetrics()
	count := GetTotalSuggestionsCount(exampleMetrics)
	if count != 500 {
		t.Errorf("Expected total suggestions count to be 500, but got %f", count)
	}
}

func TestGetTotalAcceptancesCount(t *testing.T) {
	exampleMetrics := getExampleMetrics()
	count := GetTotalAcceptancesCount(exampleMetrics)
	if count != 150 {
		t.Errorf("Expected total acceptances count to be 150, but got %f", count)
	}
}

func TestGetTotalLinesSuggested(t *testing.T) {
	exampleMetrics := getExampleMetrics()
	count := GetTotalLinesSuggested(exampleMetrics)
	if count != 1000 {
		t.Errorf("Expected total lines suggested to be 1000, but got %f", count)
	}
}

func TestGetTotalLinesAccepted(t *testing.T) {
	exampleMetrics := getExampleMetrics()
	count := GetTotalLinesAccepted(exampleMetrics)
	if count != 800 {
		t.Errorf("Expected total lines accepted to be 800, but got %f", count)
	}
}

func TestGetTotalActiveUsers(t *testing.T) {
	exampleMetrics := getExampleMetrics()
	count := GetTotalActiveUsers(exampleMetrics)
	if count != 53 {
		t.Errorf("Expected total active users to be 50, but got %f", count)
	}
}

func TestGetLastDayDataWithInvalidDate(t *testing.T) {
	usages := []github.CopilotMetrics{
		{Date: "2025-01-01", TotalActiveUsers: 100},
		{Date: "invalid-date", TotalActiveUsers: 200},
	}

	lastDayData := GetLastDayData(usages)
	if lastDayData.Date != "2025-01-01" {
		t.Errorf("Expected last day to be 2025-01-01, but got %s", lastDayData.Date)
	}
}
