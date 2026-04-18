//go:build windows

package converter

import (
	"os/exec"
	"syscall"
)

// hideConsoleWindow prevents a brief console window flash when a GUI app spawns CLI tools.
func hideConsoleWindow(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.HideWindow = true
}
