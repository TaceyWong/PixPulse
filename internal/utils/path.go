package utils

import (
	"os"
	"path/filepath"
)

// GetExeDir returns the directory where the executable or the project is located
func GetExeDir() string {
	// 1. Check for an environment variable or a specific development flag if needed
	// 2. Default to executable directory
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	dir := filepath.Dir(exe)

	// If we are in 'wails dev', the executable is in a temp folder.
	// We check if "bin" exists in the current directory as a fallback.
	if _, err := os.Stat(filepath.Join(dir, "bin")); os.IsNotExist(err) {
		// Try current working directory
		cwd, _ := os.Getwd()
		if _, err := os.Stat(filepath.Join(cwd, "bin")); err == nil {
			return cwd
		}
	}

	return dir
}

// GetBinPath returns the absolute path to a binary in the bin folder
func GetBinPath(name string) string {
	// On Windows, append .exe if not present
	if filepath.Ext(name) == "" {
		name += ".exe"
	}
	return filepath.Join(GetExeDir(), "bin", name)
}
