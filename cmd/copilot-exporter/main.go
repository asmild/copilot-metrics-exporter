package main

import (
	"flag"
	"fmt"
	"github.com/asmild/copilot-metrics-exporter/internal/config"
	"github.com/asmild/copilot-metrics-exporter/internal/prometheus"
	"os"
)

func execute() error {
	configPath := flag.String("config", "", "Path to the config file")
	flag.StringVar(configPath, "c", "", "Path to the config file (shorthand)")
	flag.Usage = func() {
		fmt.Println("Github Copilot usage prometheus exporter")
		fmt.Println("\nGitHub Copilot usage prometheus exporter is a simple tool that retrieves usage statistics for GitHub Copilot\n" +
			"and exports them in a format suitable for consumption by Prometheus. It collects data such as the total\n" +
			"number of suggestions, the number of accepted suggestions, the total lines suggested and accepted, as well\n" +
			"as language-specific breakdowns. This exporter can help monitor the usage of GitHub Copilot over time\n" +
			"and identify trends in usage patterns.")
		fmt.Println("\nUsage: program [options]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nConfiguration options:")
		fmt.Println("  - config file path (default is $HOME/.copilot-exporter/config.yaml or ./config.yaml)")
		fmt.Println("  - Environment variables:")
		fmt.Println("    GITHUB_ORG, GITHUB_IS_ENTERPRISE, GITHUB_TOKEN, GITHUB_APP_TOKEN, PORT")
	}
	flag.Parse()

	conf, err := config.MustLoad(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v\n", err)
	}
	prometheusexporter.StartExporter(conf)
	return nil
}

func main() {
	if err := execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
