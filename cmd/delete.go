package cmd

import (
	"fmt"

	"github.com/nmelo/gasmail/internal/identity"
	"github.com/nmelo/gasmail/internal/mail"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <message-id> [message-id...]",
	Short: "Delete (close) messages",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := identity.GetIdentity(identityFlag)
		if err != nil {
			return fmt.Errorf("failed to determine identity: %w", err)
		}

		mailbox := mail.NewMailbox(id)

		for _, messageID := range args {
			if err := mailbox.Delete(messageID); err != nil {
				return fmt.Errorf("failed to delete message %s: %w", messageID, err)
			}
			fmt.Printf("Deleted %s\n", messageID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
