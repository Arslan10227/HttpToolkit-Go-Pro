package main

import (
	"context"
	"fmt"
	"os/exec"
	goruntime "runtime"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/logger"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/server"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/system"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type ShellApp struct {
	ctx   context.Context
	cfg   *config.Config
	app   *server.App
	token string
}

func NewShellApp(cfg *config.Config, app *server.App) *ShellApp {
	return &ShellApp{cfg: cfg, app: app, token: cfg.AuthToken}
}

func (s *ShellApp) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *ShellApp) GetServerAuthToken() string {
	logger.Info("Wails bound method GetServerAuthToken called", nil)
	return s.token
}

func (s *ShellApp) GetServerPort() int {
	logger.Info("Wails bound method GetServerPort called", map[string]any{"port": s.cfg.ServerPort})
	return s.cfg.ServerPort
}

func (s *ShellApp) GetMockttpPort() int {
	logger.Info("Wails bound method GetMockttpPort called", map[string]any{"port": s.cfg.AdminPort})
	return s.cfg.AdminPort
}

func (s *ShellApp) FrontendLog(level string, msg string) {
	if level == "error" {
		logger.Error(fmt.Errorf("[Frontend] %s", msg), nil)
	} else {
		logger.Info(fmt.Sprintf("[Frontend] %s", msg), nil)
	}
}

func (s *ShellApp) GetDesktopVersion() string {
	logger.Info("Wails bound method GetDesktopVersion called", map[string]any{"version": config.ServerVersion()})
	return config.ServerVersion()
}

func (s *ShellApp) OpenURL(url string) {
	if s.ctx != nil {
		wailsruntime.BrowserOpenURL(s.ctx, url)
	}
}

func (s *ShellApp) OpenExternalAuth(url string) {
	if url == "" {
		url = "http://127.0.0.1:45457/auth/callback"
	}
	s.OpenURL(url)
}

func (s *ShellApp) SelectFilePath() string {
	if s.ctx == nil {
		return ""
	}
	path, err := wailsruntime.OpenFileDialog(s.ctx, wailsruntime.OpenDialogOptions{
		Title: "Select file",
	})
	if err != nil {
		return ""
	}
	return path
}

func (s *ShellApp) SelectApplication() string {
	if s.ctx == nil {
		return ""
	}
	path, err := wailsruntime.OpenFileDialog(s.ctx, wailsruntime.OpenDialogOptions{
		Title: "Select application",
		Filters: []wailsruntime.FileFilter{
			{DisplayName: "Executables", Pattern: "*"},
		},
	})
	if err != nil {
		return ""
	}
	return path
}

func (s *ShellApp) ShutdownServer() {
	if s.app != nil && s.app.REST != nil {
		_ = s.app.REST.Shutdown(s.ctx)
	}
}

func (s *ShellApp) OpenHarWithDefaultApp(path string) error {
	var cmd *exec.Cmd
	switch goruntime.GOOS {
	case "windows":
		cmd = system.Command("cmd", "/c", "start", "", path)
	case "darwin":
		cmd = system.Command("open", path)
	default:
		cmd = system.Command("xdg-open", path)
	}
	return cmd.Start()
}

func (s *ShellApp) GetLogPath() string {
	return logger.GetLogPath()
}

func (s *ShellApp) OpenLogFile() {
	logPath := logger.GetLogPath()
	if logPath == "" {
		return
	}
	var cmd *exec.Cmd
	switch goruntime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", logPath)
	case "darwin":
		cmd = exec.Command("open", logPath)
	default:
		cmd = exec.Command("xdg-open", logPath)
	}
	_ = cmd.Start()
}

func (s *ShellApp) IsHarAssociated() bool {
	return system.IsHarAssociated()
}

func (s *ShellApp) RegisterHarAssociation() bool {
	return system.RegisterHarAssociation()
}

func (s *ShellApp) UnregisterHarAssociation() bool {
	return system.UnregisterHarAssociation()
}
