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
