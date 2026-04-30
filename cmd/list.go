package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/jneb802/quo-cli/client"
	"github.com/spf13/cobra"
)

var listPhoneNumberID string
var listParticipant string
var listLimit int

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List messages",
	Long:  "List messages for a phone number and participant.",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("API key required: set QUO_API_KEY or use --api-key")
		}

		pnID := listPhoneNumberID
		if pnID == "" {
			pnID = getPhoneNumberID()
		}
		if pnID == "" {
			return fmt.Errorf("--phone-number-id is required (or set QUO_PHONE_NUMBER_ID)")
		}

		if listParticipant == "" {
			return fmt.Errorf("--participant is required")
		}

		c := client.New(apiKey)
		resp, err := c.ListMessages(pnID, listParticipant, listLimit, "")
		if err != nil {
			return err
		}

		if jsonOutput {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(resp)
		}

		if len(resp.Data) == 0 {
			fmt.Println("No messages found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIME\tDIR\tFROM\tTO\tTEXT")
		for _, m := range resp.Data {
			t, _ := time.Parse(time.RFC3339, m.CreatedAt)
			text := m.Text
			if len(text) > 50 {
				text = text[:47] + "..."
			}
			to := ""
			if len(m.To) > 0 {
				to = m.To[0]
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				t.Local().Format("Jan 02 15:04"),
				m.Direction,
				m.From,
				to,
				text,
			)
		}
		w.Flush()
		return nil
	},
}

func init() {
	listCmd.Flags().StringVar(&listPhoneNumberID, "phone-number-id", "", "OpenPhone number ID (env: QUO_PHONE_NUMBER_ID)")
	listCmd.Flags().StringVar(&listParticipant, "participant", "", "participant phone number (required)")
	listCmd.Flags().IntVar(&listLimit, "limit", 20, "max results (1-100)")
	rootCmd.AddCommand(listCmd)
}
