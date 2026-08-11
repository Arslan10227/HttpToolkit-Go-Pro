package main

import (
	"context"
	"fmt"
	"os/exec"
	goruntime "runtime"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/logger"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/server"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/settings"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/system"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// ShellService is the Wails v3 service that exposes Go methods to the frontend.
// It replaces the v2 ShellApp bound struct. Methods with exported names are
// automatically available to the frontend via generated bindings.
type ShellService struct {
	cfg      *config.Config
	sm       *settings.Manager
	app      *server.App
	token    string
	startup  func() // called during ServiceStartup to launch the server
	shutdown func() // called during ServiceShutdown to stop the server
}

// NewShellService creates a new ShellService. The startup/shutdown callbacks
// are wired by main.go to start/stop the embedded server.
func NewShellService(cfg *config.Config, sm *settings.Manager, startup, shutdown func()) *ShellService {
	return &ShellService{
		cfg:      cfg,
		sm:       sm,
		token:    cfg.AuthToken,
		startup:  startup,
		shutdown: shutdown,
	}
}

// ServiceStartup is called by Wails v3 when the application starts.
// The context is valid for the application's lifetime.
func (s *ShellService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	logger.Info("ShellService ServiceStartup — launching server", nil)
	if s.startup != nil {
		s.startup()
	}
	return nil
}

// ServiceShutdown is called by Wails v3 when the application is shutting down.
func (s *ShellService) ServiceShutdown() error {
	logger.Info("ShellService ServiceShutdown — stopping server", nil)
	if s.shutdown != nil {
		s.shutdown()
	}
	return nil
}

func (s *ShellService) GetServerAuthToken() string {
	logger.Info("Wails bound method GetServerAuthToken called", nil)
	return s.token
}

func (s *ShellService) GetServerPort() int {
	logger.Info("Wails bound method GetServerPort called", map[string]any{"port": s.cfg.ServerPort})
	return s.cfg.ServerPort
}

func (s *ShellService) GetMockttpPort() int {
	logger.Info("Wails bound method GetMockttpPort called", map[string]any{"port": s.cfg.AdminPort})
	return s.cfg.AdminPort
}

func (s *ShellService) FrontendLog(level string, msg string) {
	if level == "error" {
		logger.Error(fmt.Errorf("[Frontend] %s", msg), nil)
	} else {
		logger.Info(fmt.Sprintf("[Frontend] %s", msg), nil)
	}
}

func (s *ShellService) GetDesktopVersion() string {
	logger.Info("Wails bound method GetDesktopVersion called", map[string]any{"version": config.ServerVersion()})
	return config.ServerVersion()
}

func (s *ShellService) OpenURL(url string) {
	if err := application.Get().Browser.OpenURL(url); err != nil {
		logger.Info("Failed to open URL", map[string]any{"url": url, "error": err.Error()})
	}
}

func (s *ShellService) OpenExternalAuth(url string) {
	if url == "" {
		url = "http://127.0.0.1:45457/auth/callback"
	}
	s.OpenURL(url)
}

func (s *ShellService) SelectFilePath() string {
	dialog := application.Get().Dialog.OpenFile()
	dialog.SetTitle("Select file")
	path, err := dialog.PromptForSingleSelection()
	if err != nil {
		return ""
	}
	return path
}

func (s *ShellService) SelectApplication() string {
	dialog := application.Get().Dialog.OpenFile()
	dialog.SetTitle("Select application")
	path, err := dialog.PromptForSingleSelection()
	if err != nil {
		return ""
	}
	return path
}

func (s *ShellService) ShutdownServer() {
	if s.app != nil && s.app.REST != nil {
		_ = s.app.REST.Shutdown(context.Background())
	}
}

func (s *ShellService) OpenHarWithDefaultApp(path string) error {
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

func (s *ShellService) GetLogPath() string {
	return logger.GetLogPath()
}

func (s *ShellService) OpenLogFile() {
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

func (s *ShellService) IsHarAssociated() bool {
	return system.IsHarAssociated()
}

func (s *ShellService) RegisterHarAssociation() bool {
	return system.RegisterHarAssociation()
}

func (s *ShellService) UnregisterHarAssociation() bool {
	return system.UnregisterHarAssociation()
}
