package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nmelo/gasmail/internal/identity"
	"github.com/nmelo/gasmail/internal/mail"
	"github.com/spf13/cobra"
)

var checkJSON bool

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for new mail",
	Long: `Check for new mail. Exits with code 0 if new mail is available, 1 if not.
Useful for scripts and hooks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := identity.GetIdentity(identityFlag)
		if err != nil {
			return fmt.Errorf("failed to determine identity: %w", err)
		}

		mailbox := mail.NewMailbox(id)

		count, err := mailbox.CountUnread()
		if err != nil {
			return fmt.Errorf("failed to check mail: %w", err)
		}

		if checkJSON {
			result := map[string]interface{}{
				"identity": id,
				"unread":   count,
				"has_mail": count > 0,
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		if count == 0 {
			os.Exit(1)
		}

		fmt.Printf("You have %d unread message(s)\n", count)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVarP(&checkJSON, "json", "j", false, "Output as JSON")
}
