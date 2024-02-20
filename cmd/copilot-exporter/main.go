package main

import (
	"fmt"
	"github.com/asmild/copilot-metrics-exporter/config"
	"github.com/asmild/copilot-metrics-exporter/internal/prometheus"
	"github.com/spf13/cobra"
	"os"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "copilot-exporter",
	Short: "Github Copilot usage prometheus exporter",
	Long: `GitHub Copilot usage prometheus exporter is a tool that retrieves usage statistics for GitHub Copilot
and exports them in a format suitable for consumption by Prometheus. It collects data such as the total
number of suggestions, the number of accepted suggestions, the total lines suggested and accepted, as well
as language-specific breakdowns. This exporter can help monitor the usage of GitHub Copilot over time
and identify trends in usage patterns.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.GetConfig(configPath)
		if err != nil {
			fmt.Printf("failed to read config file: %v\n", err)
			return
		}
		prometheusexporter.StartExporter(conf)

		//
		//ghc, err := github.NewGitHubClient(*conf)
		//
		//if err != nil {
		//	fmt.Printf("Failed to run GitHub Copilot exporter: %v\n", err)
		//	return
		//}
		//usage, err := ghc.GetCopilotUsage()
		//
		//for i, u := range usage {
		//	fmt.Printf("Entry %d:\n", i+1)
		//	fmt.Println("Day:", u.Day)
		//	fmt.Println("Total Suggestions Count:", u.TotalSuggestionsCount)
		//	fmt.Println("Total Acceptances Count:", u.TotalAcceptancesCount)
		//	fmt.Println("Total Lines Suggested:", u.TotalLinesSuggested)
		//	fmt.Println("Total Lines Accepted:", u.TotalLinesAccepted)
		//	fmt.Println("Total Active Users:", u.TotalActiveUsers)
		//	fmt.Println("Breakdown:")
		//	for _, item := range u.Breakdown {
		//		fmt.Printf("  Language: %s, Editor: %s\n", item.Language, item.Editor)
		//		fmt.Printf("    Suggestions Count: %d, Acceptances Count: %d\n", item.SuggestionsCount, item.AcceptancesCount)
		//		fmt.Printf("    Lines Suggested: %d, Lines Accepted: %d\n", item.LinesSuggested, item.LinesAccepted)
		//		fmt.Printf("    Active Users: %d\n", item.ActiveUsers)
		//	}
		//	fmt.Println()
		//}
		//
		//billing, err := ghc.GetBilling()
		//if err != nil {
		//	fmt.Printf("Failed to get billing data: %v\n", err)
		//	return
		//}
		//
		//fmt.Println("Billing:")
		//fmt.Println("Seat Breakdown:")
		//fmt.Printf("  Total: %d\n", billing.SeatBreakdown.Total)
		//fmt.Printf("  Added This Cycle: %d\n", billing.SeatBreakdown.AddedThisCycle)
		//fmt.Printf("  Pending Invitation: %d\n", billing.SeatBreakdown.PendingInvitation)
		//fmt.Printf("  Pending Cancellation: %d\n", billing.SeatBreakdown.PendingCancellation)
		//fmt.Printf("  Active This Cycle: %d\n", billing.SeatBreakdown.ActiveThisCycle)
		//fmt.Printf("  Inactive This Cycle: %d\n", billing.SeatBreakdown.InactiveThisCycle)
		//fmt.Println("Seat Management Setting:", billing.SeatManagementSetting)
		//fmt.Println("Public Code Suggestions:", billing.PublicCodeSuggestions)

	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file path (default is $HOME/.copilot-exporter/config.yaml)")
}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	if err := Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//var (
//	githubExampleMetric = prometheus.NewGauge(prometheus.GaugeOpts{
//		Name: "github_example_metric",
//		Help: "Example metric exported from GitHub",
//	})
//)

//func init() {
// Set the initial value of the metric
//githubExampleMetric.Set(1)
// Register the metric with Prometheus
//prometheus.MustRegister(githubExampleMetric)
//}

//func main() {
//	configPath := flag.String("c", "", "path to config file")
//	flag.StringVar(configPath, "config", "", "path to config file")
//	flag.Parse()
//
//	usage, err := githubcopilot.GetGitHubCopilotUsage(*configPath)
//
//	if err != nil {
//		fmt.Printf("Failed to read GitHub Copilot config: %v\n", err)
//		return
//	}
//
//	fmt.Println(usage)
//	Start an HTTP server to expose the metrics
//	http.Handle("/metrics", promhttp.Handler())
//	http.ListenAndServe(":8080", nil)
//}
//
