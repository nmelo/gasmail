# gasmail Implementation Plan

## Overview

Create a CLI tool (`gm`) for agent-to-agent messaging using beads as the storage backend. This is a standalone extraction of gastown's mail system.

Companions:
- [gn](https://github.com/nmelo/gasnudge) - send nudge messages to tmux windows
- [gp](https://github.com/nmelo/gaspeek) - capture output from tmux windows

## Architecture

```
gasmail/
├── main.go
├── go.mod
├── internal/
│   ├── mail/
│   │   ├── message.go      # Message struct and helpers
│   │   ├── mailbox.go      # Inbox operations via beads
│   │   └── router.go       # Send routing
│   └── identity/
│       └── identity.go     # Identity resolution
└── cmd/
    ├── root.go
    ├── inbox.go
    ├── send.go
    ├── read.go
    ├── delete.go
    └── check.go
```

## Dependencies

- **beads** (`bd` CLI) - must be installed and available in PATH
- Messages stored as beads issues with `--type=message`

## Core Components

### 1. Message Struct (`internal/mail/message.go`)

```go
type Message struct {
    ID        string    `json:"id"`
    From      string    `json:"from"`
    To        string    `json:"to"`
    Subject   string    `json:"subject"`
    Body      string    `json:"body"`
    Timestamp time.Time `json:"timestamp"`
    Priority  int       `json:"priority"`  // 0=urgent, 1=high, 2=normal, 3=low
    Read      bool      `json:"read"`
    ThreadID  string    `json:"thread_id,omitempty"`
    ReplyTo   string    `json:"reply_to,omitempty"`
}
```

### 2. Identity Resolution (`internal/identity/identity.go`)

Identity is determined by:
1. `--identity` flag (explicit)
2. `GM_IDENTITY` environment variable
3. Tmux session name (if inside tmux)
4. Hostname (fallback)

```go
func GetIdentity(explicit string) (string, error)
func IsInsideTmux() bool
func GetTmuxSession() (string, error)
```

### 3. Mailbox Operations (`internal/mail/mailbox.go`)

Uses beads CLI commands:

```go
// List inbox messages
// bd list --type=message --assignee=<identity> --status=open --json
func (m *Mailbox) List() ([]*Message, error)

// Get single message
// bd show <id> --json
func (m *Mailbox) Get(id string) (*Message, error)

// Mark as read (close in beads)
// bd close <id>
func (m *Mailbox) MarkRead(id string) error

// Delete (close in beads)
// bd close <id>
func (m *Mailbox) Delete(id string) error

// Count unread
func (m *Mailbox) CountUnread() (int, error)
```

### 4. Router (`internal/mail/router.go`)

Send messages via beads:

```go
// Send creates a message in recipient's mailbox
// bd create <subject> --type=message --assignee=<to> -d <body> --actor=<from> --labels=from:<from>,priority:<n>
func (r *Router) Send(msg *Message) error
```

## CLI Commands

### `gm inbox`

```
gm inbox [flags]

Flags:
  -i, --identity STRING    Your identity (default: auto-detect)
  -u, --unread             Show only unread messages
  -j, --json               Output as JSON
```

Output format:
```
ID          FROM        SUBJECT                    TIME
hq-abc123   agent-1     Status update              2m ago
hq-def456   agent-2     Task complete              1h ago
```

### `gm send`

```
gm send <recipient> [flags]

Flags:
  -s, --subject STRING     Message subject (required)
  -m, --message STRING     Message body
  -i, --identity STRING    Your identity (default: auto-detect)
  -p, --priority INT       Priority 0-3 (default: 2)
  -r, --reply-to STRING    Message ID this replies to
```

Examples:
```bash
gm send agent-2 -s "Status check" -m "How's the build going?"
gm send agent-1 -s "Done" -m "Task complete" --reply-to hq-abc123
```

### `gm read`

```
gm read <message-id> [flags]

Flags:
  -j, --json               Output as JSON
```

Output format:
```
From:    agent-1
To:      agent-2
Subject: Status update
Date:    2024-01-15 10:30:00

How's the build going?
```

### `gm delete`

```
gm delete <message-id> [message-id...]
```

Closes messages in beads (marks as handled).

### `gm check`

```
gm check [flags]

Flags:
  -i, --identity STRING    Your identity (default: auto-detect)
  -j, --json               Output as JSON

Exit codes:
  0 - New mail available
  1 - No new mail
```

For use in scripts/hooks to check for mail.

## Beads Commands Reference

```bash
# Create message
bd create "Subject" --type=message --assignee=recipient -d "Body" --actor=sender --labels="from:sender,priority:2"

# List inbox
bd list --type=message --assignee=identity --status=open --json

# Show message
bd show hq-abc123 --json

# Close/delete message
bd close hq-abc123

# Add label (mark read without closing)
bd label add hq-abc123 read
```

## Beads JSON Format

When querying with `--json`, beads returns:
```json
[
  {
    "id": "hq-abc123",
    "title": "Subject line",
    "description": "Message body",
    "type": "message",
    "status": "open",
    "assignee": "recipient-identity",
    "priority": 2,
    "labels": ["from:sender", "priority:2"],
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

Parse `from:` label to extract sender. Parse `read` label to determine read status.

## Implementation Steps

1. Initialize Go module: `go mod init github.com/nmelo/gasmail`
2. Create identity package with auto-detection
3. Create mail package with Message struct
4. Implement Mailbox (list, get, delete via bd commands)
5. Implement Router (send via bd create)
6. Create CLI commands with Cobra
7. Create main.go entry point
8. Build and test
9. Create README
10. Create Homebrew formula in `~/Desktop/Projects/homebrew-tap/Formula/gm.rb`
11. Create GitHub repo and push

## Environment Variables

- `GM_IDENTITY` - Override identity detection
- `BEADS_DIR` - Beads database directory (default: `.beads` in current dir or parents)

## README Template

```markdown
# gm

CLI tool for agent-to-agent messaging using beads.

Messaging extracted from [gastown](https://github.com/steveyegge/gastown).

## Prerequisites

Requires [beads](https://github.com/steveyegge/beads) (`bd` CLI) to be installed.

## Installation

```bash
brew tap nmelo/tap
brew install gm
```

## Usage

Check inbox:
```bash
gm inbox
gm inbox --unread
```

Send a message:
```bash
gm send agent-2 -s "Status" -m "How's it going?"
```

Read a message:
```bash
gm read hq-abc123
```

Delete/archive a message:
```bash
gm delete hq-abc123
```

Check for new mail (for scripts):
```bash
gm check && echo "You have mail"
```

## Identity

Identity is auto-detected from:
1. `--identity` flag
2. `GM_IDENTITY` environment variable
3. Current tmux session name
4. Hostname

## See also

- [gn](https://github.com/nmelo/gasnudge) - send nudge messages to tmux windows
- [gp](https://github.com/nmelo/gaspeek) - capture output from tmux windows
```

## Reference Code

From gastown (`~/gastown/`):
- `internal/mail/mailbox.go` - Mailbox implementation
- `internal/mail/router.go` - Send routing (lines 580-626)
- `internal/mail/types.go` - Message struct and BeadsMessage parsing
- `internal/cmd/mail.go` - CLI structure
- `internal/cmd/mail_send.go` - Send implementation
- `internal/cmd/mail_inbox.go` - Inbox implementation
