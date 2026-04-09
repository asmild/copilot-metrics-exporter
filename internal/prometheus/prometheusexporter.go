package prometheusexporter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asmild/copilot-metrics-exporter/internal/auth"
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
	totalActiveChatUsers  *prometheus.Desc
	totalSeatsOccupied        *prometheus.Desc
	weeklyActiveUsers         *prometheus.Desc
	monthlyActiveUsers        *prometheus.Desc
	featureSuggestionsCount   *prometheus.Desc
	featureAcceptancesCount   *prometheus.Desc
	featureInteractionsCount  *prometheus.Desc
	linesAcceptedDesc         *prometheus.Desc
	linesSuggestedDesc    *prometheus.Desc
	suggestionsCountDesc  *prometheus.Desc
	acceptancesCountDesc  *prometheus.Desc
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
		totalActiveChatUsers: prometheus.NewDesc("github_copilot_total_active_chat_users",
			"Total number of active chat users utilizing GitHub Copilot last day.",
			nil, nil,
		),
		weeklyActiveUsers: prometheus.NewDesc("github_copilot_weekly_active_users",
			"Number of active Copilot users in the last 7 days.",
			nil, nil,
		),
		monthlyActiveUsers: prometheus.NewDesc("github_copilot_monthly_active_users",
			"Number of active Copilot users in the last 30 days.",
			nil, nil,
		),
		featureSuggestionsCount: prometheus.NewDesc("github_copilot_feature_suggestions_count",
			"Code generation activity count by Copilot feature.",
			[]string{"feature"}, nil,
		),
		featureAcceptancesCount: prometheus.NewDesc("github_copilot_feature_acceptances_count",
			"Code acceptance activity count by Copilot feature.",
			[]string{"feature"}, nil,
		),
		featureInteractionsCount: prometheus.NewDesc("github_copilot_feature_interactions_count",
			"User-initiated interaction count by Copilot feature.",
			[]string{"feature"}, nil,
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

		githubClient: githubClient,
	}
}

func (collector *CopilotMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.totalSuggestionsCount
	ch <- collector.totalAcceptancesCount
	ch <- collector.totalLinesSuggested
	ch <- collector.totalLinesAccepted
	ch <- collector.totalActiveUsers
	ch <- collector.totalActiveChatUsers
	ch <- collector.totalSeatsOccupied
	ch <- collector.weeklyActiveUsers
	ch <- collector.monthlyActiveUsers
	ch <- collector.featureSuggestionsCount
	ch <- collector.featureAcceptancesCount
	ch <- collector.featureInteractionsCount
	ch <- collector.linesAcceptedDesc
	ch <- collector.linesSuggestedDesc
	ch <- collector.suggestionsCountDesc
	ch <- collector.acceptancesCountDesc
}

func (collector *CopilotMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	report, err := collector.githubClient.GetMetrics(ctx)
	if err != nil {
		fmt.Printf("Failed to get Copilot usage: %v\n", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(collector.totalSuggestionsCount, prometheus.GaugeValue, helper.GetTotalSuggestionsCount(report))
	ch <- prometheus.MustNewConstMetric(collector.totalAcceptancesCount, prometheus.GaugeValue, helper.GetTotalAcceptancesCount(report))
	ch <- prometheus.MustNewConstMetric(collector.totalLinesSuggested, prometheus.GaugeValue, helper.GetTotalLinesSuggested(report))
	ch <- prometheus.MustNewConstMetric(collector.totalLinesAccepted, prometheus.GaugeValue, helper.GetTotalLinesAccepted(report))
	ch <- prometheus.MustNewConstMetric(collector.totalActiveUsers, prometheus.GaugeValue, helper.GetTotalActiveUsers(report))
	ch <- prometheus.MustNewConstMetric(collector.totalChats, prometheus.GaugeValue, helper.GetTotalChats(report))
	ch <- prometheus.MustNewConstMetric(collector.totalActiveChatUsers, prometheus.GaugeValue, helper.GetTotalActiveChatUsers(report))
	ch <- prometheus.MustNewConstMetric(collector.weeklyActiveUsers, prometheus.GaugeValue, float64(report.WeeklyActiveUsers))
	ch <- prometheus.MustNewConstMetric(collector.monthlyActiveUsers, prometheus.GaugeValue, float64(report.MonthlyActiveUsers))

	for _, f := range report.TotalsByFeature {
		ch <- prometheus.MustNewConstMetric(collector.featureSuggestionsCount, prometheus.GaugeValue, float64(f.CodeGenerationActivityCount), f.Feature)
		ch <- prometheus.MustNewConstMetric(collector.featureAcceptancesCount, prometheus.GaugeValue, float64(f.CodeAcceptanceActivityCount), f.Feature)
		ch <- prometheus.MustNewConstMetric(collector.featureInteractionsCount, prometheus.GaugeValue, float64(f.UserInitiatedInteractionCount), f.Feature)
	}

	totalSeats, err := collector.githubClient.GetTotalSeats(ctx)
	if err != nil {
		fmt.Printf("Warning: Copilot seats unavailable: %v\n", err)
	}
	ch <- prometheus.MustNewConstMetric(collector.totalSeatsOccupied, prometheus.GaugeValue, float64(totalSeats))

	metricsSum := helper.ComputeLanguageBreakdown(report)
	for editor, languages := range metricsSum {
		for language, metrics := range languages {

			ch <- prometheus.MustNewConstMetric(collector.linesAcceptedDesc, prometheus.GaugeValue, metrics["linesAccepted"], language, editor)
			ch <- prometheus.MustNewConstMetric(collector.linesSuggestedDesc, prometheus.GaugeValue, metrics["linesSuggested"], language, editor)
			ch <- prometheus.MustNewConstMetric(collector.suggestionsCountDesc, prometheus.GaugeValue, metrics["suggestionsCount"], language, editor)
			ch <- prometheus.MustNewConstMetric(collector.acceptancesCountDesc, prometheus.GaugeValue, metrics["acceptancesCount"], language, editor)
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

	// Print basic auth status
	if conf.BasicAuth != nil && conf.BasicAuth.Enabled {
		fmt.Println("Basic Auth: enabled")
	} else {
		fmt.Println("Basic Auth: disabled")
	}

	initMetrics(conf)

	mux := http.NewServeMux()

	// Create the metrics handler
	metricsHandler := promhttp.Handler()

	// Apply basic auth middleware if configured
	if conf.BasicAuth != nil && conf.BasicAuth.Enabled {
		basicAuthMiddleware := auth.BasicAuthMiddleware(conf.BasicAuth)
		metricsHandler = basicAuthMiddleware(metricsHandler)
	}

	mux.Handle("/metrics", metricsHandler)
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
}