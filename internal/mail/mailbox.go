package mail

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Mailbox handles inbox operations via beads CLI
type Mailbox struct {
	Identity string
}

// NewMailbox creates a new Mailbox for the given identity
func NewMailbox(identity string) *Mailbox {
	return &Mailbox{Identity: identity}
}

// List returns all messages in the inbox
func (m *Mailbox) List() ([]*Message, error) {
	cmd := exec.Command("bd", "list",
		"--type=message",
		"--assignee="+m.Identity,
		"--status=open",
		"--json",
	)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("bd list failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("bd list failed: %w", err)
	}

	if len(output) == 0 || string(output) == "null" || string(output) == "[]\n" {
		return []*Message{}, nil
	}

	return ParseBeadsOutput(output)
}

// ListUnread returns only unread messages
func (m *Mailbox) ListUnread() ([]*Message, error) {
	messages, err := m.List()
	if err != nil {
		return nil, err
	}

	unread := make([]*Message, 0)
	for _, msg := range messages {
		if !msg.Read {
			unread = append(unread, msg)
		}
	}
	return unread, nil
}

// Get retrieves a single message by ID
func (m *Mailbox) Get(id string) (*Message, error) {
	cmd := exec.Command("bd", "show", id, "--json")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("bd show failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("bd show failed: %w", err)
	}

	var issue BeadsIssue
	if err := json.Unmarshal(output, &issue); err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}

	return ParseBeadsIssue(&issue), nil
}

// MarkRead marks a message as read by adding the "read" label
func (m *Mailbox) MarkRead(id string) error {
	cmd := exec.Command("bd", "label", "add", id, "read")
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("bd label add failed: %s", string(exitErr.Stderr))
		}
		return fmt.Errorf("bd label add failed: %w", err)
	}
	return nil
}

// Delete closes a message in beads (marks as handled)
func (m *Mailbox) Delete(id string) error {
	cmd := exec.Command("bd", "close", id)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("bd close failed: %s", string(exitErr.Stderr))
		}
		return fmt.Errorf("bd close failed: %w", err)
	}
	return nil
}

// CountUnread returns the number of unread messages
func (m *Mailbox) CountUnread() (int, error) {
	messages, err := m.ListUnread()
	if err != nil {
		return 0, err
	}
	return len(messages), nil
}
