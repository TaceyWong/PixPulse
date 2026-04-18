package converter

import (
	"bytes"
	"fmt"
	"os/exec"
)

// RunCommand executes a command and returns combined output or error
func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	hideConsoleWindow(cmd)

	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("command failed: %w\nOutput: %s", err, out.String())
	}

	return out.String(), nil
}
