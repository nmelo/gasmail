package identity

import (
	"os"
	"os/exec"
	"strings"
)

// GetIdentity returns the current identity based on priority:
// 1. Explicit identity (from --identity flag)
// 2. GM_IDENTITY environment variable
// 3. Tmux window name (if inside tmux)
// 4. Hostname (fallback)
func GetIdentity(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}

	if envID := os.Getenv("GM_IDENTITY"); envID != "" {
		return envID, nil
	}

	if IsInsideTmux() {
		if window, err := GetTmuxWindow(); err == nil && window != "" {
			return window, nil
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "unknown", err
	}
	return hostname, nil
}

// IsInsideTmux returns true if running inside a tmux session
func IsInsideTmux() bool {
	return os.Getenv("TMUX") != ""
}

// GetTmuxWindow returns the current tmux window name
func GetTmuxWindow() (string, error) {
	cmd := exec.Command("tmux", "display-message", "-p", "#W")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
