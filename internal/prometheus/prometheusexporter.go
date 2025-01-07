package prometheusexporter

import (
	"fmt"
	"github.com/asmild/copilot-metrics-exporter/config"
	"github.com/asmild/copilot-metrics-exporter/internal/github"
	"github.com/asmild/copilot-metrics-exporter/internal/helper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type CopilotMetricsCollector struct {
	totalSuggestionsCount *prometheus.Desc
	totalAcceptancesCount *prometheus.Desc
	totalLinesSuggested   *prometheus.Desc
	totalLinesAccepted    *prometheus.Desc
	totalActiveUsers      *prometheus.Desc
	totalChatAcceptances  *prometheus.Desc
	totalChatTurns        *prometheus.Desc
	totalActiveChatUsers  *prometheus.Desc
	totalSeatsOccupied    *prometheus.Desc
	linesAcceptedDesc     *prometheus.Desc
	linesSuggestedDesc    *prometheus.Desc
	suggestionsCountDesc  *prometheus.Desc
	acceptancesCountDesc  *prometheus.Desc
	activeUsers           *prometheus.Desc
	githubClient          *github.GitHubClient
}

func NewCopilotMetricsCollector(githubClient *github.GitHubClient) *CopilotMetricsCollector {
	return &CopilotMetricsCollector{
		totalSeatsOccupied: prometheus.NewDesc("github_copilot_total_seats_occupied",
			"Total number of seats occupied by users.",
			nil, nil,
		),
		totalSuggestionsCount: prometheus.NewDesc("github_copilot_total_suggestions_count",
			"Total number of suggestions made by GitHub Copilot last day.",
			nil, nil,
		),
		totalAcceptancesCount: prometheus.NewDesc("github_copilot_total_acceptances_count",
			"Total number of suggestions accepted by users last day.",
			nil, nil,
		),
		totalLinesSuggested: prometheus.NewDesc("github_copilot_total_lines_suggested",
			"Total number of lines suggested by GitHub Copilot last day.",
			nil, nil,
		),
		totalLinesAccepted: prometheus.NewDesc("github_copilot_total_lines_accepted",
			"Total number of lines accepted by users last day.",
			nil, nil,
		),
		totalActiveUsers: prometheus.NewDesc("github_copilot_total_active_users",
			"Total number of active users utilizing GitHub Copilot last day.",
			nil, nil,
		),
        totalChatAcceptances: prometheus.NewDesc("github_copilot_total_chat_acceptances",
            "Total number of chat acceptances made by GitHub Copilot last day.",
            nil, nil,
        ),
        totalChatTurns: prometheus.NewDesc("github_copilot_total_chat_turns",
            "Total number of chat turns made by GitHub Copilot last day.",
            nil, nil,
        ),
        totalActiveChatUsers: prometheus.NewDesc("github_copilot_total_active_chat_users",
            "Total number of active chat users utilizing GitHub Copilot last day.",
            nil, nil,
        ),
		linesAcceptedDesc: prometheus.NewDesc("github_copilot_lines_accepted_breakdown",
			"Lines accepted breakdown for GitHub Copilot by language and editor.",
			[]string{"language", "editor"}, nil,
		),

		linesSuggestedDesc: prometheus.NewDesc(
			"github_copilot_lines_suggested_breakdown",
			fmt.Sprintf("Lines suggested breakdown for GitHub Copilot by language and editor."),
			[]string{"language", "editor"},
			nil,
		),

		suggestionsCountDesc: prometheus.NewDesc(
			"github_copilot_suggestions_count_breakdown",
			fmt.Sprintf("Suggestions count breakdown for GitHub Copilot by language and editor."),
			[]string{"language", "editor"},
			nil,
		),

		acceptancesCountDesc: prometheus.NewDesc(
			"github_copilot_acceptances_count_breakdown",
			fmt.Sprintf("Acceptanse count breakdown for GitHub Copilot by language and editor."),
			[]string{"language", "editor"},
			nil,
		),

		activeUsers: prometheus.NewDesc(
			"github_copilot_active_users_breakdown",
			fmt.Sprintf("Active users breakdown for GitHub Copilot by language and editor."),
			[]string{"language", "editor"},
			nil,
		),
		githubClient: githubClient,
	}
}

func (collector *CopilotMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.totalSuggestionsCount
	ch <- collector.totalAcceptancesCount
	ch <- collector.totalLinesSuggested
	ch <- collector.totalLinesAccepted
	ch <- collector.totalActiveUsers
	ch <- collector.totalChatAcceptances
	ch <- collector.totalChatTurns
	ch <- collector.totalActiveChatUsers
	ch <- collector.totalSeatsOccupied
	ch <- collector.linesAcceptedDesc
	ch <- collector.linesSuggestedDesc
	ch <- collector.suggestionsCountDesc
	ch <- collector.acceptancesCountDesc
	ch <- collector.activeUsers
}

func (collector *CopilotMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	copilotUsage, err := collector.githubClient.GetCopilotUsage()
	if err != nil {
		fmt.Printf("Failed to get Copilot usage: %v\n", err)
		return
	}

	billing, err := collector.githubClient.GetBillingSeats()
	if err != nil {
		fmt.Printf("Failed to get Copilot billing: %v\n", err)
		return
	}

	lastDayCopilotUsage := helper.GetLastDayData(copilotUsage)
	ch <- prometheus.MustNewConstMetric(collector.totalSuggestionsCount, prometheus.GaugeValue, helper.GetTotalSuggestionsCount(lastDayCopilotUsage))
	ch <- prometheus.MustNewConstMetric(collector.totalAcceptancesCount, prometheus.GaugeValue, helper.GetTotalAcceptancesCount(lastDayCopilotUsage))
	ch <- prometheus.MustNewConstMetric(collector.totalLinesSuggested, prometheus.GaugeValue, helper.GetTotalLinesSuggested(lastDayCopilotUsage))
	ch <- prometheus.MustNewConstMetric(collector.totalLinesAccepted, prometheus.GaugeValue, helper.GetTotalLinesAccepted(lastDayCopilotUsage))
	ch <- prometheus.MustNewConstMetric(collector.totalActiveUsers, prometheus.GaugeValue, helper.GetTotalActiveUsers(lastDayCopilotUsage))
	ch <- prometheus.MustNewConstMetric(collector.totalChatAcceptances, prometheus.GaugeValue, helper.GetTotalChatAcceptances(lastDayCopilotUsage))
	ch <- prometheus.MustNewConstMetric(collector.totalChatTurns, prometheus.GaugeValue, helper.GetTotalChatTurns(lastDayCopilotUsage))
	ch <- prometheus.MustNewConstMetric(collector.totalActiveChatUsers, prometheus.GaugeValue, helper.GetTotalActiveChatUsers(lastDayCopilotUsage))
	ch <- prometheus.MustNewConstMetric(collector.totalSeatsOccupied, prometheus.GaugeValue, float64(billing.TotalSeats))

	metricsSum := make(map[string]map[string]map[string]float64)
	for _, usage := range lastDayCopilotUsage.Breakdown {
		language := usage.Language
		editor := usage.Editor

		if _, ok := metricsSum[editor]; !ok {
			metricsSum[editor] = make(map[string]map[string]float64)
		}

		if _, ok := metricsSum[editor][language]; !ok {
			metricsSum[editor][language] = make(map[string]float64)
		}

		metricsSum[editor][language]["linesAccepted"] += float64(usage.LinesAccepted)
		metricsSum[editor][language]["linesSuggested"] += float64(usage.LinesSuggested)
		metricsSum[editor][language]["suggestionsCount"] += float64(usage.SuggestionsCount)
		metricsSum[editor][language]["acceptancesCount"] += float64(usage.AcceptancesCount)
		metricsSum[editor][language]["activeUsers"] += float64(usage.ActiveUsers)
	}
	for editor, languages := range metricsSum {
		for language, metrics := range languages {

			ch <- prometheus.MustNewConstMetric(collector.linesAcceptedDesc, prometheus.GaugeValue, metrics["linesAccepted"], language, editor)
			ch <- prometheus.MustNewConstMetric(collector.linesSuggestedDesc, prometheus.GaugeValue, metrics["linesSuggested"], language, editor)
			ch <- prometheus.MustNewConstMetric(collector.suggestionsCountDesc, prometheus.GaugeValue, metrics["suggestionsCount"], language, editor)
			ch <- prometheus.MustNewConstMetric(collector.acceptancesCountDesc, prometheus.GaugeValue, metrics["acceptancesCount"], language, editor)
			ch <- prometheus.MustNewConstMetric(collector.activeUsers, prometheus.GaugeValue, metrics["activeUsers"], language, editor)
		}
	}
}

func collectMetrics(conf *config.ExporterConfig) {
	ghc, err := github.NewGitHubClient(*conf)
	if err != nil {
		fmt.Printf("Failed to run GitHub Copilot exporter: %v\n", err)
		return
	}

	if err != nil {
		fmt.Printf("Failed to get Copilot usage: %v\n", err)
		return
	}

	collector := NewCopilotMetricsCollector(ghc)
	prometheus.MustRegister(collector)
}

func initMetrics(conf *config.ExporterConfig) {
	collectMetrics(conf)
}

func StartExporter(conf *config.ExporterConfig) {
	fmt.Println("Starting exporter on port", conf.Port)
	if conf.IsEnterprise {
		fmt.Println("Enterprise:", conf.Organization)
	} else {
		fmt.Println("Organization:", conf.Organization)
	}
	initMetrics(conf)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":"+conf.Port, nil)
	if err != nil {
		fmt.Printf("Failed to start HTTP server: %v\n", err)
	}
}
