package helper

import (
	"github.com/asmild/copilot-metrics-exporter/internal/github"
	"testing"
)

func TestGetLastDayData(t *testing.T) {
	usages := []github.CopilotUsage{
		{Day: "2025-01-01", TotalSuggestionsCount: 100},
		{Day: "2025-01-02", TotalSuggestionsCount: 200},
		{Day: "2025-01-03", TotalSuggestionsCount: 300},
	}

	lastDayData := GetLastDayData(usages)
	if lastDayData.Day != "2025-01-03" {
		t.Errorf("Expected last day to be 2025-01-03, but got %s", lastDayData.Day)
	}
}

func TestGetTotalSuggestionsCount(t *testing.T) {
	usage := github.CopilotUsage{TotalSuggestionsCount: 500}
	count := GetTotalSuggestionsCount(usage)
	if count != 500 {
		t.Errorf("Expected total suggestions count to be 500, but got %f", count)
	}
}

func TestGetTotalAcceptancesCount(t *testing.T) {
	usage := github.CopilotUsage{TotalAcceptancesCount: 150}
	count := GetTotalAcceptancesCount(usage)
	if count != 150 {
		t.Errorf("Expected total acceptances count to be 150, but got %f", count)
	}
}

func TestGetTotalLinesSuggested(t *testing.T) {
	usage := github.CopilotUsage{TotalLinesSuggested: 1000}
	count := GetTotalLinesSuggested(usage)
	if count != 1000 {
		t.Errorf("Expected total lines suggested to be 1000, but got %f", count)
	}
}

func TestGetTotalLinesAccepted(t *testing.T) {
	usage := github.CopilotUsage{TotalLinesAccepted: 800}
	count := GetTotalLinesAccepted(usage)
	if count != 800 {
		t.Errorf("Expected total lines accepted to be 800, but got %f", count)
	}
}

func TestGetTotalActiveUsers(t *testing.T) {
	usage := github.CopilotUsage{TotalActiveUsers: 50}
	count := GetTotalActiveUsers(usage)
	if count != 50 {
		t.Errorf("Expected total active users to be 50, but got %f", count)
	}
}

func TestGetLastDayDataWithInvalidDate(t *testing.T) {
	usages := []github.CopilotUsage{
		{Day: "2025-01-01", TotalSuggestionsCount: 100},
		{Day: "invalid-date", TotalSuggestionsCount: 200},
	}

	lastDayData := GetLastDayData(usages)
	if lastDayData.Day != "2025-01-01" {
		t.Errorf("Expected last day to be 2025-01-01, but got %s", lastDayData.Day)
	}
}
