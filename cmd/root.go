package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var apiKey string
var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "quo",
	Short: "CLI for the OpenPhone messaging API",
	Long:  "quo-cli lets you send and receive SMS messages via the OpenPhone API.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getAPIKey() string {
	if apiKey != "" {
		return apiKey
	}
	return os.Getenv("QUO_API_KEY")
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key (env: QUO_API_KEY)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output raw JSON")
}
