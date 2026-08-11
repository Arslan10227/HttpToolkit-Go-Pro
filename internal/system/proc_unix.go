//go:build !windows

package system

import (
	"context"
	"os/exec"
)

// HideWindow is a no-op on non-Windows platforms.
func HideWindow(cmd *exec.Cmd) {
	// No-op
}

// Command creates a standard exec.Cmd on non-Windows platforms.
func Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

// CommandContext creates a standard exec.Cmd with context on non-Windows platforms.
func CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, arg...)
}
