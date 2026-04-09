package helper

import (
	"github.com/asmild/copilot-metrics-exporter/internal/github"
)

func GetTotalSuggestionsCount(report *github.UsageReport) float64 {
	if report == nil {
		return 0
	}
	return float64(report.CodeGenerationActivityCount)
}

func GetTotalAcceptancesCount(report *github.UsageReport) float64 {
	if report == nil {
		return 0
	}
	return float64(report.CodeAcceptanceActivityCount)
}

func GetTotalLinesSuggested(report *github.UsageReport) float64 {
	if report == nil {
		return 0
	}
	return float64(report.LocSuggestedToAddSum)
}

func GetTotalLinesAccepted(report *github.UsageReport) float64 {
	if report == nil {
		return 0
	}
	return float64(report.LocAddedSum)
}

func GetTotalActiveUsers(report *github.UsageReport) float64 {
	if report == nil {
		return 0
	}
	return float64(report.DailyActiveUsers)
}

func GetTotalChats(report *github.UsageReport) float64 {
	if report == nil {
		return 0
	}
	return float64(report.UserInitiatedInteractionCount)
}

func GetTotalActiveChatUsers(report *github.UsageReport) float64 {
	if report == nil {
		return 0
	}
	return float64(report.MonthlyActiveChatUsers)
}

// ComputeLanguageBreakdown aggregates metrics by editor and by language.
// Returns map[editor_or_language]map[language_or_"all"]map[metricKey]value.
//
// For IDE breakdown: uses totals_by_ide (editor-level totals).
// For language breakdown: uses totals_by_language_feature (language×feature totals).
func ComputeLanguageBreakdown(report *github.UsageReport) map[string]map[string]map[string]float64 {
	metricsSum := make(map[string]map[string]map[string]float64)

	if report == nil {
		return metricsSum
	}

	// Aggregate language data from totals_by_language_feature.
	// The new format provides language×feature, so we aggregate across features per language
	// and attribute to "all" editors (since IDE-level language breakdown is not available).
	for _, lf := range report.TotalsByLanguageFeature {
		if lf.Language == "" {
			continue
		}
		editor := "all"
		if _, ok := metricsSum[editor]; !ok {
			metricsSum[editor] = make(map[string]map[string]float64)
		}
		if _, ok := metricsSum[editor][lf.Language]; !ok {
			metricsSum[editor][lf.Language] = make(map[string]float64)
		}
		metricsSum[editor][lf.Language]["linesAccepted"] += float64(lf.LocAddedSum)
		metricsSum[editor][lf.Language]["linesSuggested"] += float64(lf.LocSuggestedToAddSum)
		metricsSum[editor][lf.Language]["suggestionsCount"] += float64(lf.CodeGenerationActivityCount)
		metricsSum[editor][lf.Language]["acceptancesCount"] += float64(lf.CodeAcceptanceActivityCount)
	}

	// Aggregate IDE-level data from totals_by_ide.
	// These don't have per-language breakdown, so we use language "all".
	for _, ide := range report.TotalsByIDE {
		if ide.IDE == "" {
			continue
		}
		if _, ok := metricsSum[ide.IDE]; !ok {
			metricsSum[ide.IDE] = make(map[string]map[string]float64)
		}
		if _, ok := metricsSum[ide.IDE]["all"]; !ok {
			metricsSum[ide.IDE]["all"] = make(map[string]float64)
		}
		metricsSum[ide.IDE]["all"]["linesAccepted"] += float64(ide.LocAddedSum)
		metricsSum[ide.IDE]["all"]["linesSuggested"] += float64(ide.LocSuggestedToAddSum)
		metricsSum[ide.IDE]["all"]["suggestionsCount"] += float64(ide.CodeGenerationActivityCount)
		metricsSum[ide.IDE]["all"]["acceptancesCount"] += float64(ide.CodeAcceptanceActivityCount)
	}

	return metricsSum
}