package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jneb802/quo-cli/client"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <message-id>",
	Short: "Get a message by ID",
	Long:  "Retrieve a single message by its unique identifier.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("API key required: set QUO_API_KEY or use --api-key")
		}

		c := client.New(apiKey)
		msg, err := c.GetMessage(args[0])
		if err != nil {
			return err
		}

		if jsonOutput {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(msg)
		}

		t, _ := time.Parse(time.RFC3339, msg.CreatedAt)
		fmt.Printf("ID:        %s\n", msg.ID)
		fmt.Printf("Direction: %s\n", msg.Direction)
		fmt.Printf("Status:    %s\n", msg.Status)
		fmt.Printf("From:      %s\n", msg.From)
		fmt.Printf("To:        %s\n", strings.Join(msg.To, ", "))
		fmt.Printf("Time:      %s\n", t.Local().Format(time.RFC1123))
		fmt.Printf("Text:      %s\n", msg.Text)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
