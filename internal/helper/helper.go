package helper

import (
	"fmt"
	"github.com/asmild/copilot-metrics-exporter/internal/github"
	"time"
)

func GetLastDayData(usages []github.CopilotUsage) github.CopilotUsage {

	var maxDate time.Time
	var lastDayData github.CopilotUsage

	for _, item := range usages {
		date, err := time.Parse("2006-01-02", item.Day)
		if err != nil {
			fmt.Println("Failed to parse date: ", err)
			continue
		}
		if date.After(maxDate) {
			maxDate = date
			lastDayData = item
		}
	}
	return lastDayData
}

func GetTotalSuggestionsCount(usage github.CopilotUsage) float64 {
	return float64(usage.TotalSuggestionsCount)
}

func GetTotalAcceptancesCount(usage github.CopilotUsage) float64 {
	return float64(usage.TotalAcceptancesCount)
}

func GetTotalLinesSuggested(usage github.CopilotUsage) float64 {
	return float64(usage.TotalLinesSuggested)
}

func GetTotalLinesAccepted(usage github.CopilotUsage) float64 {
	return float64(usage.TotalLinesAccepted)
}

func GetTotalActiveUsers(usage github.CopilotUsage) float64 {
	return float64(usage.TotalActiveUsers)
}
