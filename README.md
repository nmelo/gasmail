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

Or build from source:

```bash
go install github.com/nmelo/gasmail@latest
```

## Usage

Check inbox:
```bash
gm inbox
gm inbox --unread
gm inbox --json
```

Send a message:
```bash
gm send agent-2 -s "Status" -m "How's it going?"
gm send agent-1 -s "Done" -m "Task complete" --reply-to hq-abc123
gm send agent-2 -s "Urgent" -m "Need help" -p 0  # priority 0=urgent
```

Read a message:
```bash
gm read hq-abc123
gm read hq-abc123 --json
```

Delete/archive a message:
```bash
gm delete hq-abc123
gm delete hq-abc123 hq-def456  # multiple messages
```

Check for new mail (for scripts):
```bash
gm check && echo "You have mail"
gm check --json
```

## Identity

Identity is auto-detected from:
1. `--identity` flag
2. `GM_IDENTITY` environment variable
3. Current tmux window name
4. Hostname

Override identity for any command:
```bash
gm inbox --identity agent-1
gm send agent-2 -s "Hello" -i agent-1
```

## Environment Variables

- `GM_IDENTITY` - Override identity detection
- `BEADS_DIR` - Beads database directory (default: `.beads` in current dir or parents)

## Priority Levels

- 0: Urgent
- 1: High
- 2: Normal (default)
- 3: Low

## See also

- [gn](https://github.com/nmelo/gasnudge) - send nudge messages to tmux windows
- [gp](https://github.com/nmelo/gaspeek) - capture output from tmux windows
