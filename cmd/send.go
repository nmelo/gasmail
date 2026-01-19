package cmd

import (
	"fmt"

	"github.com/nmelo/gasmail/internal/identity"
	"github.com/nmelo/gasmail/internal/mail"
	"github.com/spf13/cobra"
)

var sendSubject string
var sendMessage string
var sendPriority int
var sendReplyTo string

var sendCmd = &cobra.Command{
	Use:   "send <recipient>",
	Short: "Send a message to another agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		recipient := args[0]

		id, err := identity.GetIdentity(identityFlag)
		if err != nil {
			return fmt.Errorf("failed to determine identity: %w", err)
		}

		if sendSubject == "" {
			return fmt.Errorf("subject is required (use -s or --subject)")
		}

		router := mail.NewRouter(id)

		msg := &mail.Message{
			To:       recipient,
			Subject:  sendSubject,
			Body:     sendMessage,
			Priority: sendPriority,
			ReplyTo:  sendReplyTo,
		}

		if err := router.Send(msg); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		fmt.Printf("Message sent to %s\n", recipient)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().StringVarP(&sendSubject, "subject", "s", "", "Message subject (required)")
	sendCmd.Flags().StringVarP(&sendMessage, "message", "m", "", "Message body")
	sendCmd.Flags().IntVarP(&sendPriority, "priority", "p", mail.PriorityNormal, "Priority 0-3 (0=urgent, 1=high, 2=normal, 3=low)")
	sendCmd.Flags().StringVarP(&sendReplyTo, "reply-to", "r", "", "Message ID this replies to")
}
