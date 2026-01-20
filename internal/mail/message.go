package mail

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Priority levels for messages
const (
	PriorityUrgent = 0
	PriorityHigh   = 1
	PriorityNormal = 2
	PriorityLow    = 3
)

// Message represents a mail message
type Message struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Timestamp time.Time `json:"timestamp"`
	Priority  int       `json:"priority"`
	Read      bool      `json:"read"`
	ThreadID  string    `json:"thread_id,omitempty"`
	ReplyTo   string    `json:"reply_to,omitempty"`
}

// BeadsIssue represents the JSON structure returned by beads
type BeadsIssue struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Type        string   `json:"issue_type"`
	Status      string   `json:"status"`
	Assignee    string   `json:"assignee"`
	Owner       string   `json:"owner"`
	Priority    int      `json:"priority"`
	Labels      []string `json:"labels"`
	CreatedAt   string   `json:"created_at"`
	CreatedBy   string   `json:"created_by"`
	UpdatedAt   string   `json:"updated_at"`
}

// ParseBeadsIssue converts a BeadsIssue to a Message
func ParseBeadsIssue(issue *BeadsIssue) *Message {
	msg := &Message{
		ID:       issue.ID,
		Subject:  issue.Title,
		Body:     issue.Description,
		To:       issue.Assignee,
		Priority: issue.Priority,
	}

	// Parse timestamp
	if t, err := time.Parse(time.RFC3339, issue.CreatedAt); err == nil {
		msg.Timestamp = t
	}

	// Parse labels for sender and read status
	for _, label := range issue.Labels {
		if strings.HasPrefix(label, "from:") {
			msg.From = strings.TrimPrefix(label, "from:")
		}
		if label == "read" {
			msg.Read = true
		}
		if strings.HasPrefix(label, "thread:") {
			msg.ThreadID = strings.TrimPrefix(label, "thread:")
		}
		if strings.HasPrefix(label, "reply-to:") {
			msg.ReplyTo = strings.TrimPrefix(label, "reply-to:")
		}
	}

	return msg
}

// ParseBeadsOutput parses JSON output from beads list command
func ParseBeadsOutput(data []byte) ([]*Message, error) {
	var issues []BeadsIssue
	if err := json.Unmarshal(data, &issues); err != nil {
		return nil, err
	}

	messages := make([]*Message, 0, len(issues))
	for i := range issues {
		messages = append(messages, ParseBeadsIssue(&issues[i]))
	}

	return messages, nil
}

// FormatTimeAgo returns a human-readable relative time
func FormatTimeAgo(t time.Time) string {
	diff := time.Since(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1m ago"
		}
		return fmt.Sprintf("%dm ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1h ago"
		}
		return fmt.Sprintf("%dh ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1d ago"
		}
		return fmt.Sprintf("%dd ago", days)
	default:
		return t.Format("Jan 2")
	}
}
