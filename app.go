package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"PixPulse/internal/converter"

	"github.com/wailsapp/wails/v2/pkg/options"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	converter *converter.Converter
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		converter: converter.NewConverter(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// onSecondInstanceLaunch runs when a second process is started; the first instance stays running.
// See https://wails.golang.ac.cn/docs/guides/single-instance-lock/
func (a *App) onSecondInstanceLaunch(data options.SecondInstanceData) {
	if a.ctx == nil {
		return
	}
	wailsRuntime.WindowUnminimise(a.ctx)
	wailsRuntime.Show(a.ctx)
	go wailsRuntime.EventsEmit(a.ctx, "secondInstance", data.Args, data.WorkingDirectory)
}

// SelectFile opens a native file dialog to select an image
func (a *App) SelectFile() string {
	selection, err := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "选择图片",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Images (*.png;*.jpg;*.svg)", Pattern: "*.png;*.jpg;*.jpeg;*.svg"},
		},
	})
	if err != nil {
		return ""
	}
	return selection
}

// GetFileBase64 reads a file and returns its content as a base64 encoded data URI
func (a *App) GetFileBase64(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	
	ext := filepath.Ext(path)
	var mimeType string
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".svg":
		mimeType = "image/svg+xml"
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	default:
		mimeType = "application/octet-stream"
	}

	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data)), nil
}

// OpenDirectory opens the system file explorer and selects the specified file
func (a *App) OpenDirectory(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	// Clean path for Windows
	if runtime.GOOS == "windows" {
		absPath = filepath.Clean(absPath)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// The comma after /select is crucial and should be part of the same argument or handled carefully
		cmd = exec.Command("explorer.exe", "/select,"+absPath)
	case "darwin":
		cmd = exec.Command("open", "-R", absPath)
	default:
		dir := filepath.Dir(absPath)
		cmd = exec.Command("xdg-open", dir)
	}
	
	// Execute without blocking
	go cmd.Run()
}

// Convert processes an image conversion request
func (a *App) Convert(req converter.ConvertRequest) (*converter.ConvertResponse, error) {
	return a.converter.Convert(req)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
