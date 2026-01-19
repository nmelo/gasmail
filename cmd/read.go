package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nmelo/gasmail/internal/identity"
	"github.com/nmelo/gasmail/internal/mail"
	"github.com/spf13/cobra"
)

var readJSON bool

var readCmd = &cobra.Command{
	Use:   "read <message-id>",
	Short: "Read a message",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		messageID := args[0]

		id, err := identity.GetIdentity(identityFlag)
		if err != nil {
			return fmt.Errorf("failed to determine identity: %w", err)
		}

		mailbox := mail.NewMailbox(id)

		msg, err := mailbox.Get(messageID)
		if err != nil {
			return fmt.Errorf("failed to get message: %w", err)
		}

		// Mark as read
		if err := mailbox.MarkRead(messageID); err != nil {
			// Non-fatal, just warn
			fmt.Fprintf(os.Stderr, "Warning: could not mark message as read: %v\n", err)
		}

		if readJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(msg)
		}

		fmt.Printf("From:    %s\n", msg.From)
		fmt.Printf("To:      %s\n", msg.To)
		fmt.Printf("Subject: %s\n", msg.Subject)
		fmt.Printf("Date:    %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Println()
		if msg.Body != "" {
			fmt.Println(msg.Body)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().BoolVarP(&readJSON, "json", "j", false, "Output as JSON")
}
