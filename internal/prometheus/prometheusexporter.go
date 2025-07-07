package prometheusexporter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/asmild/copilot-metrics-exporter/internal/config"
	"github.com/asmild/copilot-metrics-exporter/internal/github"
	"github.com/asmild/copilot-metrics-exporter/internal/helper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CopilotMetricsCollector struct {
	totalSuggestionsCount *prometheus.Desc
	totalAcceptancesCount *prometheus.Desc
	totalLinesSuggested   *prometheus.Desc
	totalLinesAccepted    *prometheus.Desc
	totalActiveUsers      *prometheus.Desc
	totalChats            *prometheus.Desc
	totalChatInsertions   *prometheus.Desc
	totalChatCopies       *prometheus.Desc
	totalActiveChatUsers  *prometheus.Desc
	totalSeatsOccupied    *prometheus.Desc
	linesAcceptedDesc     *prometheus.Desc
	linesSuggestedDesc    *prometheus.Desc
	suggestionsCountDesc  *prometheus.Desc
	acceptancesCountDesc  *prometheus.Desc
	activeUsers           *prometheus.Desc
	githubClient          *github.Client
}

func NewCopilotMetricsCollector(githubClient *github.Client) *CopilotMetricsCollector {
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
		totalChats: prometheus.NewDesc("github_copilot_total_chats",
			"The total number of chats initiated by users last day.",
			nil, nil,
		),
		totalChatInsertions: prometheus.NewDesc("github_copilot_total_chat_insertions",
			"Total number of chat acceptances made by GitHub Copilot last day.",
			nil, nil,
		),
		totalChatCopies: prometheus.NewDesc("github_copilot_total_chat_copies",
			"The number of times users copied a code suggestion from Copilot Chat using the keyboard, or the 'Copy' UI element last day.",
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
			"Lines suggested breakdown for GitHub Copilot by language and editor.",
			[]string{"language", "editor"},
			nil,
		),

		suggestionsCountDesc: prometheus.NewDesc(
			"github_copilot_suggestions_count_breakdown",
			"Suggestions count breakdown for GitHub Copilot by language and editor.",
			[]string{"language", "editor"},
			nil,
		),

		acceptancesCountDesc: prometheus.NewDesc(
			"github_copilot_acceptances_count_breakdown",
			"Acceptance count breakdown for GitHub Copilot by language and editor.",
			[]string{"language", "editor"},
			nil,
		),

		activeUsers: prometheus.NewDesc(
			"github_copilot_active_users_breakdown",
			"Active users breakdown for GitHub Copilot by language and editor.",
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
	ch <- collector.totalChatInsertions
	ch <- collector.totalChatCopies
	ch <- collector.totalActiveChatUsers
	ch <- collector.totalSeatsOccupied
	ch <- collector.linesAcceptedDesc
	ch <- collector.linesSuggestedDesc
	ch <- collector.suggestionsCountDesc
	ch <- collector.acceptancesCountDesc
	ch <- collector.activeUsers
}

func (collector *CopilotMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	since := time.Now().AddDate(0, 0, -1).Format("2006-01-02T15:04:05Z")
	copilotUsage, err := collector.githubClient.GetCopilotMetrics(&since)

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
	ch <- prometheus.MustNewConstMetric(collector.totalChatInsertions, prometheus.GaugeValue, helper.GetTotalChatInsertions(lastDayCopilotUsage))
	ch <- prometheus.MustNewConstMetric(collector.totalChats, prometheus.GaugeValue, helper.GetTotalChats(lastDayCopilotUsage))
	ch <- prometheus.MustNewConstMetric(collector.totalChatCopies, prometheus.GaugeValue, helper.GetTotalChatCopies(lastDayCopilotUsage))
	ch <- prometheus.MustNewConstMetric(collector.totalActiveChatUsers, prometheus.GaugeValue, helper.GetTotalActiveChatUsers(lastDayCopilotUsage))
	ch <- prometheus.MustNewConstMetric(collector.totalSeatsOccupied, prometheus.GaugeValue, float64(billing.TotalSeats))

	metricsSum := helper.ComputeLanguageBreakdown(lastDayCopilotUsage)
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

func collectMetrics(conf *config.Config) {
	ghc, err := github.NewClient(conf)
	if err != nil {
		fmt.Printf("failed to run GitHub Copilot exporter: %v\n", err)
		return
	}

	collector := NewCopilotMetricsCollector(ghc)
	prometheus.MustRegister(collector)
}

func initMetrics(conf *config.Config) {
	collectMetrics(conf)
}

func StartExporter(conf *config.Config) {
	fmt.Println("Starting exporter on port", conf.Port)
	if conf.IsEnterprise {
		fmt.Println("Enterprise:", conf.Organization)
	} else {
		fmt.Println("Organization:", conf.Organization)
	}

	initMetrics(conf)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	addr := fmt.Sprintf(":%s", conf.Port)

	if conf.TLS != nil && conf.TLS.Enabled {
		fmt.Println("TLS enabled - using HTTPS")
		err := http.ListenAndServeTLS(addr, conf.TLS.CertFile, conf.TLS.KeyFile, mux)
		if err != nil {
			fmt.Printf("Failed to start HTTPS server: %v\n", err)
		}
	} else {
		fmt.Println("TLS disabled - using HTTP")
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			fmt.Printf("Failed to start HTTP server: %v\n", err)
		}
	}
	fmt.Println("Exporter started successfully")
}
