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
	totalSeatsOccupied    *prometheus.Desc

	githubClient *github.GitHubClient
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
		githubClient: githubClient,
	}
}

func (collector *CopilotMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.totalSuggestionsCount
	ch <- collector.totalAcceptancesCount
	ch <- collector.totalLinesSuggested
	ch <- collector.totalLinesAccepted
	ch <- collector.totalActiveUsers
	ch <- collector.totalSeatsOccupied
}

func (collector *CopilotMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	copilotUsage, err := collector.githubClient.GetCopilotUsage()
	if err != nil {
		fmt.Printf("Failed to get Copilot usage: %v\n", err)
		return
	}

	billing, err := collector.githubClient.GetBilling()
	if err != nil {
		fmt.Printf("Failed to get Copilot billing: %v\n", err)
		return
	}

	lastDayCopilotUsage := helper.GetLastDayData(copilotUsage)
	totalSeatsOccupied := float64(billing.SeatBreakdown.Total)
	totalSuggestionsCount := helper.GetTotalSuggestionsCount(lastDayCopilotUsage)
	totalAcceptancesCount := helper.GetTotalAcceptancesCount(lastDayCopilotUsage)
	totalLinesSuggested := helper.GetTotalLinesSuggested(lastDayCopilotUsage)
	totalLinesAccepted := helper.GetTotalLinesAccepted(lastDayCopilotUsage)
	totalActiveUsers := helper.GetTotalActiveUsers(lastDayCopilotUsage)

	ch <- prometheus.MustNewConstMetric(collector.totalSuggestionsCount, prometheus.GaugeValue, totalSuggestionsCount)
	ch <- prometheus.MustNewConstMetric(collector.totalAcceptancesCount, prometheus.GaugeValue, totalAcceptancesCount)
	ch <- prometheus.MustNewConstMetric(collector.totalLinesSuggested, prometheus.GaugeValue, totalLinesSuggested)
	ch <- prometheus.MustNewConstMetric(collector.totalLinesAccepted, prometheus.GaugeValue, totalLinesAccepted)
	ch <- prometheus.MustNewConstMetric(collector.totalActiveUsers, prometheus.GaugeValue, totalActiveUsers)
	ch <- prometheus.MustNewConstMetric(collector.totalSeatsOccupied, prometheus.GaugeValue, totalSeatsOccupied)

	for _, usage := range lastDayCopilotUsage.Breakdown {
		language := usage.Language
		editor := usage.Editor
		linesAccepted := float64(usage.LinesAccepted)
		linesSuggested := float64(usage.LinesSuggested)

		linesAcceptedDesc := prometheus.NewDesc(
			"github_copilot_lines_accepted_breakdown",
			fmt.Sprintf("Usage breakdown for GitHub Copilot by language and editor."),
			[]string{"language", "editor"},
			nil,
		)

		linesSuggestedDesc := prometheus.NewDesc(
			"github_copilot_lines_suggested_breakdown",
			fmt.Sprintf("Usage breakdown for GitHub Copilot by language and editor."),
			[]string{"language", "editor"},
			nil,
		)

		suggestionsCountDesc := prometheus.NewDesc(
			"github_copilot_suggestions_count_breakdown",
			fmt.Sprintf("Usage breakdown for GitHub Copilot by language and editor."),
			[]string{"language", "editor"},
			nil,
		)

		acceptancesCountDesc := prometheus.NewDesc(
			"github_copilot_acceptances_count_breakdown",
			fmt.Sprintf("Usage breakdown for GitHub Copilot by language and editor."),
			[]string{"language", "editor"},
			nil,
		)
		ch <- prometheus.MustNewConstMetric(linesAcceptedDesc, prometheus.GaugeValue, linesAccepted, language, editor)
		ch <- prometheus.MustNewConstMetric(linesSuggestedDesc, prometheus.GaugeValue, linesSuggested, language, editor)
		ch <- prometheus.MustNewConstMetric(suggestionsCountDesc, prometheus.GaugeValue, float64(usage.SuggestionsCount), language, editor)
		ch <- prometheus.MustNewConstMetric(acceptancesCountDesc, prometheus.GaugeValue, float64(usage.AcceptancesCount), language, editor)
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
	initMetrics(conf)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":"+conf.Port, nil)
	if err != nil {
		fmt.Printf("Failed to start HTTP server: %v\n", err)
	}
}
