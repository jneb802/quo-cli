package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jneb802/quo-cli/client"
	"github.com/spf13/cobra"
)

var sendFrom string
var sendTo string

var sendCmd = &cobra.Command{
	Use:   "send [message]",
	Short: "Send a text message",
	Long:  "Send a text message from your OpenPhone number to a recipient.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("API key required: set QUO_API_KEY or use --api-key")
		}

		from := sendFrom
		if from == "" {
			from = os.Getenv("QUO_FROM")
		}
		if from == "" {
			return fmt.Errorf("--from is required (or set QUO_FROM)")
		}

		if sendTo == "" {
			return fmt.Errorf("--to is required")
		}

		content := strings.Join(args, " ")
		c := client.New(apiKey)
		msg, err := c.SendMessage(from, sendTo, content)
		if err != nil {
			return err
		}

		if jsonOutput {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(msg)
		}

		fmt.Printf("Message sent: %s (status: %s)\n", msg.ID, msg.Status)
		return nil
	},
}

func init() {
	sendCmd.Flags().StringVar(&sendFrom, "from", "", "sender phone number or ID (env: QUO_FROM)")
	sendCmd.Flags().StringVar(&sendTo, "to", "", "recipient phone number (required)")
	rootCmd.AddCommand(sendCmd)
}
