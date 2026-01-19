package mail

import (
	"fmt"
	"os/exec"
	"strconv"
)

// Router handles sending messages via beads
type Router struct {
	Identity string
}

// NewRouter creates a new Router for the given sender identity
func NewRouter(identity string) *Router {
	return &Router{Identity: identity}
}

// Send creates a message in the recipient's mailbox
func (r *Router) Send(msg *Message) error {
	if msg.To == "" {
		return fmt.Errorf("recipient is required")
	}
	if msg.Subject == "" {
		return fmt.Errorf("subject is required")
	}

	// Set defaults
	if msg.From == "" {
		msg.From = r.Identity
	}
	if msg.Priority < 0 || msg.Priority > 3 {
		msg.Priority = PriorityNormal
	}

	// Build labels
	labels := fmt.Sprintf("from:%s,priority:%d", msg.From, msg.Priority)
	if msg.ReplyTo != "" {
		labels += ",reply-to:" + msg.ReplyTo
	}
	if msg.ThreadID != "" {
		labels += ",thread:" + msg.ThreadID
	}

	// Build command
	// bd create <subject> --type=message --assignee=<to> -d <body> --actor=<from> --labels=<labels> --priority=<n>
	args := []string{
		"create",
		msg.Subject,
		"--type=message",
		"--assignee=" + msg.To,
		"--actor=" + msg.From,
		"--labels=" + labels,
		"--priority=" + strconv.Itoa(msg.Priority),
	}

	if msg.Body != "" {
		args = append(args, "-d", msg.Body)
	}

	cmd := exec.Command("bd", args...)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("bd create failed: %s", string(exitErr.Stderr))
		}
		return fmt.Errorf("bd create failed: %w", err)
	}

	return nil
}
