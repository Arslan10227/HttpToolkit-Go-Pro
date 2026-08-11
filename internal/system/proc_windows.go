//go:build windows

package system

import (
	"context"
	"os/exec"
	"syscall"
)

// HideWindow sets the SysProcAttr on Windows to hide the command prompt window.
func HideWindow(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.HideWindow = true
}

// Command creates an exec.Cmd with standard hidden window behavior on Windows.
func Command(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	HideWindow(cmd)
	return cmd
}

// CommandContext creates an exec.Cmd with context and hidden window behavior on Windows.
func CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, arg...)
	HideWindow(cmd)
	return cmd
}
