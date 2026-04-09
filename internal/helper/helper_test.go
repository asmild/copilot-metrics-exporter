package helper

import (
	"testing"

	"github.com/asmild/copilot-metrics-exporter/internal/github"
)

func getExampleReport() *github.UsageReport {
	return &github.UsageReport{
		Day:                           "2026-04-08",
		DailyActiveUsers:              116,
		MonthlyActiveChatUsers:        72,
		UserInitiatedInteractionCount: 145,
		CodeGenerationActivityCount:   2164,
		CodeAcceptanceActivityCount:   546,
		LocSuggestedToAddSum:          4170,
		LocAddedSum:                   1353,
		TotalsByIDE: []github.IDETotals{
			{
				IDE:                         "intellij",
				CodeGenerationActivityCount: 1210,
				CodeAcceptanceActivityCount: 275,
				LocSuggestedToAddSum:        2207,
				LocAddedSum:                 707,
			},
			{
				IDE:                         "vscode",
				CodeGenerationActivityCount: 864,
				CodeAcceptanceActivityCount: 232,
				LocSuggestedToAddSum:        1347,
				LocAddedSum:                 516,
			},
		},
		TotalsByLanguageFeature: []github.LanguageFeatureTotals{
			{
				Language:                    "java",
				Feature:                     "code_completion",
				CodeGenerationActivityCount: 416,
				CodeAcceptanceActivityCount: 126,
				LocSuggestedToAddSum:        893,
				LocAddedSum:                 283,
			},
			{
				Language:                    "python",
				Feature:                     "code_completion",
				CodeGenerationActivityCount: 203,
				CodeAcceptanceActivityCount: 62,
				LocSuggestedToAddSum:        319,
				LocAddedSum:                 133,
			},
		},
	}
}

func TestGetTotalSuggestionsCount(t *testing.T) {
	report := getExampleReport()
	count := GetTotalSuggestionsCount(report)
	if count != 2164 {
		t.Errorf("Expected 2164, got %f", count)
	}
}

func TestGetTotalAcceptancesCount(t *testing.T) {
	report := getExampleReport()
	count := GetTotalAcceptancesCount(report)
	if count != 546 {
		t.Errorf("Expected 546, got %f", count)
	}
}

func TestGetTotalLinesSuggested(t *testing.T) {
	report := getExampleReport()
	count := GetTotalLinesSuggested(report)
	if count != 4170 {
		t.Errorf("Expected 4170, got %f", count)
	}
}

func TestGetTotalLinesAccepted(t *testing.T) {
	report := getExampleReport()
	count := GetTotalLinesAccepted(report)
	if count != 1353 {
		t.Errorf("Expected 1353, got %f", count)
	}
}

func TestGetTotalActiveUsers(t *testing.T) {
	report := getExampleReport()
	count := GetTotalActiveUsers(report)
	if count != 116 {
		t.Errorf("Expected 116, got %f", count)
	}
}

func TestGetTotalChats(t *testing.T) {
	report := getExampleReport()
	count := GetTotalChats(report)
	if count != 145 {
		t.Errorf("Expected 145, got %f", count)
	}
}

func TestGetTotalActiveChatUsers(t *testing.T) {
	report := getExampleReport()
	count := GetTotalActiveChatUsers(report)
	if count != 72 {
		t.Errorf("Expected 72, got %f", count)
	}
}

func TestNilReportSafety(t *testing.T) {
	nilChecks := []struct {
		name string
		fn   func(*github.UsageReport) float64
	}{
		{"GetTotalSuggestionsCount", GetTotalSuggestionsCount},
		{"GetTotalAcceptancesCount", GetTotalAcceptancesCount},
		{"GetTotalLinesSuggested", GetTotalLinesSuggested},
		{"GetTotalLinesAccepted", GetTotalLinesAccepted},
		{"GetTotalActiveUsers", GetTotalActiveUsers},
		{"GetTotalChats", GetTotalChats},
		{"GetTotalActiveChatUsers", GetTotalActiveChatUsers},
	}
	for _, tc := range nilChecks {
		if tc.fn(nil) != 0 {
			t.Errorf("%s: expected 0 for nil report", tc.name)
		}
	}
	if len(ComputeLanguageBreakdown(nil)) != 0 {
		t.Error("ComputeLanguageBreakdown: expected empty breakdown for nil report")
	}
}

func TestComputeLanguageBreakdown(t *testing.T) {
	report := getExampleReport()
	breakdown := ComputeLanguageBreakdown(report)

	// Should have IDE entries
	intellij, ok := breakdown["intellij"]
	if !ok {
		t.Fatal("Expected 'intellij' in breakdown")
	}
	allLang := intellij["all"]
	if allLang["suggestionsCount"] != 1210 {
		t.Errorf("Expected 1210 suggestions for intellij, got %f", allLang["suggestionsCount"])
	}

	// Should have language entries under "all" editor
	allEditor, ok := breakdown["all"]
	if !ok {
		t.Fatal("Expected 'all' editor in breakdown")
	}
	java := allEditor["java"]
	if java["suggestionsCount"] != 416 {
		t.Errorf("Expected 416 suggestions for java, got %f", java["suggestionsCount"])
	}
	if java["linesAccepted"] != 283 {
		t.Errorf("Expected 283 lines accepted for java, got %f", java["linesAccepted"])
	}
}