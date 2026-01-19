package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nmelo/gasmail/internal/identity"
	"github.com/nmelo/gasmail/internal/mail"
	"github.com/spf13/cobra"
)

var inboxUnread bool
var inboxJSON bool

var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "List messages in your inbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := identity.GetIdentity(identityFlag)
		if err != nil {
			return fmt.Errorf("failed to determine identity: %w", err)
		}

		mailbox := mail.NewMailbox(id)

		var messages []*mail.Message
		if inboxUnread {
			messages, err = mailbox.ListUnread()
		} else {
			messages, err = mailbox.List()
		}
		if err != nil {
			return fmt.Errorf("failed to list messages: %w", err)
		}

		if inboxJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(messages)
		}

		if len(messages) == 0 {
			fmt.Println("No messages")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tFROM\tSUBJECT\tTIME")
		for _, msg := range messages {
			timeStr := mail.FormatTimeAgo(msg.Timestamp)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", msg.ID, msg.From, msg.Subject, timeStr)
		}
		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(inboxCmd)
	inboxCmd.Flags().BoolVarP(&inboxUnread, "unread", "u", false, "Show only unread messages")
	inboxCmd.Flags().BoolVarP(&inboxJSON, "json", "j", false, "Output as JSON")
}
