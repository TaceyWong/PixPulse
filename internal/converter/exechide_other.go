//go:build !windows

package converter

import "os/exec"

func hideConsoleWindow(cmd *exec.Cmd) {}
