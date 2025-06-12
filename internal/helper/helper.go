package helper

import (
	"fmt"
	"github.com/asmild/copilot-metrics-exporter/internal/github"
	"time"
)

func GetLastDayData(usages []github.CopilotMetrics) github.CopilotMetrics {
	var maxDate time.Time
	var lastDayData github.CopilotMetrics

	for _, item := range usages {
		date, err := time.Parse("2006-01-02", item.Date)
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

func GetTotalSuggestionsCount(metrics github.CopilotMetrics) float64 {
	totalSuggestions := 0

	// Sum suggestions from languages
	for _, lang := range metrics.CopilotIDECodeCompletions.Languages {
		totalSuggestions += lang.TotalSuggestions
	}

	// Sum suggestions from editors' models' languages
	for _, editor := range metrics.CopilotIDECodeCompletions.Editors {
		for _, model := range editor.Models {
			for _, lang := range model.Languages {
				totalSuggestions += lang.TotalSuggestions
			}
		}
	}

	return float64(totalSuggestions)
}

func GetTotalAcceptancesCount(metrics github.CopilotMetrics) float64 {
	totalAcceptances := 0

	// Sum accepted suggestions from languages
	for _, lang := range metrics.CopilotIDECodeCompletions.Languages {
		totalAcceptances += lang.TotalAcceptances
	}

	// Sum accepted suggestions from editors' models' languages
	for _, editor := range metrics.CopilotIDECodeCompletions.Editors {
		for _, model := range editor.Models {
			for _, lang := range model.Languages {
				totalAcceptances += lang.TotalAcceptances
			}
		}
	}

	return float64(totalAcceptances)
}

func GetTotalLinesSuggested(metrics github.CopilotMetrics) float64 {
	totalLinesSuggested := 0

	// Sum lines suggested from languages
	for _, lang := range metrics.CopilotIDECodeCompletions.Languages {
		totalLinesSuggested += lang.TotalLinesSuggested
	}

	// Sum lines suggested from editors' models' languages
	for _, editor := range metrics.CopilotIDECodeCompletions.Editors {
		for _, model := range editor.Models {
			for _, lang := range model.Languages {
				totalLinesSuggested += lang.TotalLinesSuggested
			}
		}
	}
	return float64(totalLinesSuggested)
}

func GetTotalLinesAccepted(metrics github.CopilotMetrics) float64 {
	totalLinesAccepted := 0

	// Sum accepted lines from languages
	for _, lang := range metrics.CopilotIDECodeCompletions.Languages {
		totalLinesAccepted += lang.TotalLinesAccepted
	}

	// Sum accepted lines from editors' models' languages
	for _, editor := range metrics.CopilotIDECodeCompletions.Editors {
		for _, model := range editor.Models {
			for _, lang := range model.Languages {
				totalLinesAccepted += lang.TotalLinesAccepted
			}
		}
	}
	return float64(totalLinesAccepted)
}

func GetTotalActiveUsers(metrics github.CopilotMetrics) float64 {
	return float64(metrics.TotalActiveUsers)
}

func GetTotalChats(metrics github.CopilotMetrics) float64 {
	totalChats := 0

	// Sum chats from IDE chat editors' models
	for _, editor := range metrics.CopilotIDEChat.Editors {
		for _, model := range editor.Models {
			totalChats += model.TotalChats
		}
	}

	// Sum chats from Dotcom chat models
	for _, model := range metrics.CopilotDotcomChat.Models {
		totalChats += model.TotalChats
	}
	return float64(totalChats)
}

func GetTotalChatInsertions(metrics github.CopilotMetrics) float64 {
	totalChatAcceptance := 0

	// Sum chat acceptances from IDE chat editors' models
	for _, editor := range metrics.CopilotIDEChat.Editors {
		for _, model := range editor.Models {
			totalChatAcceptance += model.TotalChatInsertions
		}
	}

	// Sum chat acceptances from Dotcom chat models
	for _, model := range metrics.CopilotDotcomChat.Models {
		totalChatAcceptance += model.TotalChatInsertions
	}
	return float64(totalChatAcceptance)
}

func GetTotalChatCopies(metrics github.CopilotMetrics) float64 {
	totalChatAcceptance := 0

	// Sum chat code copying from IDE chat editors' models
	for _, editor := range metrics.CopilotIDEChat.Editors {
		for _, model := range editor.Models {
			totalChatAcceptance += model.TotalChatCopies
		}
	}

	// Sum chat code copying from Dotcom chat models
	for _, model := range metrics.CopilotDotcomChat.Models {
		totalChatAcceptance += model.TotalChatCopies
	}
	return float64(totalChatAcceptance)
}

func GetTotalActiveChatUsers(metrics github.CopilotMetrics) float64 {
	totalActiveChatUsers := metrics.CopilotIDEChat.TotalEngagedUsers +
		metrics.CopilotDotcomChat.TotalEngagedUsers
	return float64(totalActiveChatUsers)
}

func ComputeLanguageBreakdown(metrics github.CopilotMetrics) map[string]map[string]map[string]float64 {
	metricsSum := make(map[string]map[string]map[string]float64)

	// Aggregate language data from direct `languages`
	for _, lang := range metrics.CopilotIDECodeCompletions.Languages {
		if lang.Name == "" { // Skip empty languages
			continue
		}
		editor := "global" // This represents standalone language usage, outside specific editors

		if _, ok := metricsSum[editor]; !ok {
			metricsSum[editor] = make(map[string]map[string]float64)
		}
		if _, ok := metricsSum[editor][lang.Name]; !ok {
			metricsSum[editor][lang.Name] = make(map[string]float64)
		}

		metricsSum[editor][lang.Name]["linesAccepted"] += float64(lang.TotalLinesAccepted)
		metricsSum[editor][lang.Name]["linesSuggested"] += float64(lang.TotalLinesSuggested)
		metricsSum[editor][lang.Name]["suggestionsCount"] += float64(lang.TotalSuggestions)
		metricsSum[editor][lang.Name]["acceptancesCount"] += float64(lang.TotalAcceptances)
		metricsSum[editor][lang.Name]["activeUsers"] += float64(lang.TotalEngagedUsers)
	}

	// Aggregate language data from each editor's models
	for _, editor := range metrics.CopilotIDECodeCompletions.Editors {
		for _, model := range editor.Models {
			for _, lang := range model.Languages {
				if lang.Name == "" { // Skip empty languages
					continue
				}
				if _, ok := metricsSum[editor.Name]; !ok {
					metricsSum[editor.Name] = make(map[string]map[string]float64)
				}
				if _, ok := metricsSum[editor.Name][lang.Name]; !ok {
					metricsSum[editor.Name][lang.Name] = make(map[string]float64)
				}

				metricsSum[editor.Name][lang.Name]["linesAccepted"] += float64(lang.TotalLinesAccepted)
				metricsSum[editor.Name][lang.Name]["linesSuggested"] += float64(lang.TotalLinesSuggested)
				metricsSum[editor.Name][lang.Name]["suggestionsCount"] += float64(lang.TotalSuggestions)
				metricsSum[editor.Name][lang.Name]["acceptancesCount"] += float64(lang.TotalAcceptances)
				metricsSum[editor.Name][lang.Name]["activeUsers"] += float64(lang.TotalEngagedUsers)
			}
		}
	}

	return metricsSum
}
