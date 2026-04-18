package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiKey    string
	serverURL string
)

var rootCmd = &cobra.Command{
	Use:   "openrelay",
	Short: "OpenRelay CLI — local webhook tunnel and management",
	Long:  `OpenRelay CLI forwards webhooks from your OpenRelay server to localhost for local development.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "OpenRelay API key (required)")
	rootCmd.PersistentFlags().StringVar(&serverURL, "server-url", "http://localhost:8081", "OpenRelay server base URL")
}