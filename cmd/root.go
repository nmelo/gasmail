package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var identityFlag string

var rootCmd = &cobra.Command{
	Use:   "gm",
	Short: "Agent-to-agent messaging CLI",
	Long: `gasmail (gm) provides persistent agent-to-agent messaging using beads as storage.

WHEN TO USE gm vs gn/ga:
  gm  - Durable messages that persist across sessions and restarts
  gn  - Immediate interrupts (requires both agents running in tmux)
  ga  - Queued messages (requires both agents running in tmux)

BEHAVIOR:
  - Messages stored in beads database (survives restarts)
  - Agents poll inbox to receive messages
  - Supports priorities, threading (reply-to), and JSON output
  - Works across tmux sessions and even across machines (shared beads)

IDENTITY RESOLUTION (in priority order):
  1. --identity / -i flag
  2. GM_IDENTITY environment variable
  3. Current tmux window name
  4. Hostname

PRIORITY LEVELS:
  0 = Urgent   (process immediately)
  1 = High     (prioritize over normal work)
  2 = Normal   (default)
  3 = Low      (handle when convenient)

USE CASES FOR AGENT COORDINATION:
  - Request work from another agent and get results later
  - Coordinate agents that aren't always running simultaneously
  - Build task queues for worker agents
  - Send status updates that won't be lost if recipient is busy
  - Thread conversations with reply-to for context

SUBCOMMANDS:
  inbox   - List messages in your inbox
  send    - Send a message to another agent
  read    - Read a specific message by ID
  delete  - Close/archive messages
  check   - Exit 0 if new mail, 1 if none (for scripts)

EXAMPLES:
  gm inbox                           # List all messages
  gm inbox --unread                  # Only unread messages
  gm inbox --json                    # JSON output for parsing
  gm send worker-1 -s "Task" -m "Run tests"
  gm send worker-1 -s "Urgent" -m "Stop" -p 0
  gm send worker-1 -s "Re: Task" -m "Done" --reply-to hq-abc123
  gm read hq-abc123                  # Read specific message
  gm read hq-abc123 --json           # JSON format
  gm delete hq-abc123                # Archive message
  gm delete hq-abc123 hq-def456      # Multiple messages
  gm check && echo "You have mail"   # Script usage

ENVIRONMENT VARIABLES:
  GM_IDENTITY  - Override identity detection
  BEADS_DIR    - Beads database directory (default: .beads in cwd or parents)

RELATED TOOLS:
  gn (gasnudge) - Interrupt agents urgently (requires tmux)
  ga (gasadd)   - Queue messages without interrupting (requires tmux)
  gp (gaspeek)  - Read output from agent windows (requires tmux)`,
}

func Execute(version string) {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&identityFlag, "identity", "i", "", "Your identity (default: auto-detect)")
}
