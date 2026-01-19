package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var identityFlag string

var rootCmd = &cobra.Command{
	Use:   "gm",
	Short: "Agent-to-agent messaging CLI",
	Long:  `gm is a CLI tool for agent-to-agent messaging using beads as the storage backend.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&identityFlag, "identity", "i", "", "Your identity (default: auto-detect)")
}
