package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/backup"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/logger"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/server"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/settings"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/system"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:assets
var assets embed.FS

func handleStartupArg(ctx context.Context, arg string) {
	if ctx == nil {
		return
	}
	lower := strings.ToLower(arg)
	if strings.HasSuffix(lower, ".har") {
		if content, err := os.ReadFile(arg); err == nil {
			logger.Info("Emitting open-har-file event", map[string]any{"path": arg})
			wailsruntime.EventsEmit(ctx, "open-har-file", map[string]any{
				"path":    arg,
				"content": string(content),
			})
		} else {
			logger.Info("Failed to read startup HAR file", map[string]any{"path": arg, "error": err.Error()})
		}
	} else if strings.HasPrefix(arg, "httptoolkitpro://") {
		logger.Info("Emitting auth-callback event", map[string]any{"url": arg})
		wailsruntime.EventsEmit(ctx, "auth-callback", arg)
	}
}

func extractEmbeddedOverrides(cfg *config.Config) {
	overridesJar := filepath.Join(cfg.AssetsDir, "overrides", "java-agent.jar")
	if _, err := os.Stat(overridesJar); err == nil {
		// Real overrides directory is already available (dev/build with external assets).
		return
	}

	// The bundled Wails executable has no on-disk assets directory, so extract
	// the overrides needed by interceptors (java-agent.jar, frida, etc.) into
	// the user's config directory and use that as the runtime AssetsDir.
	targetDir := cfg.ConfigDir
	overridesDir := filepath.Join(targetDir, "overrides")

	if err := fs.WalkDir(assets, "assets/overrides", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(path, "assets/overrides/")
		if rel == "" || rel == "." {
			return nil
		}
		targetPath := filepath.Join(overridesDir, filepath.FromSlash(rel))
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0o755)
		}
		data, err := fs.ReadFile(assets, path)
		if err != nil {
			return err
		}
		return os.WriteFile(targetPath, data, 0o644)
	}); err != nil {
		logger.Info("Failed to extract embedded overrides", map[string]any{"error": err.Error()})
		return
	}

	cfg.AssetsDir = targetDir
}

func main() {
	// Parse command line arguments manually to check for verbose flags (e.g. -v, --verbose)
	// We avoid using standard flag.Parse() to prevent WebView2 spawned subprocesses (which receive unrecognized arguments) from crashing the engine.
	verbose := false
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--verbose" {
			verbose = true
			break
		}
	}
	logger.SetVerbose(verbose)
	logger.Info("HttpToolkit Pro starting", map[string]any{"version": config.AppVersion()})

	// Disable WebView2 web security & allow mixed content so that Vercel (HTTPS) UI can make REST/WebSocket requests to local (HTTP) server ports
	os.Setenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", "--disable-web-security --allow-running-insecure-content --remote-debugging-port=9222")

	cfg, err := config.LoadDefault()
	if err != nil {
		logger.Info("Config load failed", map[string]any{"error": err.Error()})
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}

	extractEmbeddedOverrides(cfg)

	logger.Info("Config loaded", map[string]any{
		"serverPort": cfg.ServerPort,
		"adminPort":  cfg.AdminPort,
		"assetsDir":  cfg.AssetsDir,
		"devMode":    cfg.DevMode,
	})

	if os.Getenv("HTK_SERVER_TOKEN") == "" {
		os.Setenv("HTK_SERVER_TOKEN", cfg.AuthToken)
	}
	sm := settings.NewManager(cfg.ConfigDir)
	sm.SetBackup(backup.New(cfg))
	if err := sm.Load(); err != nil {
		logger.Error(err, map[string]any{"msg": "settings load failed"})
	}

	logger.Info("Auth token set", map[string]any{"token": cfg.AuthToken[:8] + "..."})

	// Register deep link custom protocol in Windows registry
	if err := system.RegisterProtocol(); err != nil {
		logger.Info("Protocol registration failed", map[string]any{"error": err.Error()})
	}

	// Capture startup arguments if any (.har or deep link url)
	var startupArgs []string
	for _, arg := range os.Args {
		lower := strings.ToLower(arg)
		if strings.HasSuffix(lower, ".har") || strings.HasPrefix(lower, "httptoolkitpro://") {
			startupArgs = append(startupArgs, arg)
		}
	}

	shell := NewShellApp(cfg, nil)

	// Sub FS: strip the "assets/" prefix so Wails serves index.html at /
	// (Go embed preserves directory structure, so assets/index.html needs
	//  to be mapped to / for Wails to find it)
	subFS, err := fs.Sub(assets, "assets")
	if err != nil {
		logger.Info("fs.Sub failed", map[string]any{"error": err.Error()})
		fmt.Fprintf(os.Stderr, "fs.Sub: %v\n", err)
		os.Exit(1)
	}

	// Verify index.html is in the embedded FS
	if f, openErr := subFS.Open("index.html"); openErr != nil {
		logger.Info("WARN: index.html not found in embedded assets", map[string]any{"error": openErr.Error()})
	} else {
		f.Close()
		logger.Info("index.html found in embedded assets — offline-first mode active", nil)
	}

	logger.Info("Starting Wails WebView2 shell", map[string]any{
		"adminPort": cfg.AdminPort,
		"restPort":  cfg.ServerPort,
	})

	err = wails.Run(&options.App{
		Title:     config.AppTitle(),
		Width:     1024,
		Height:    700,
		MinWidth:  800,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: subFS,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 17, B: 23, A: 255},
		OnStartup: func(ctx context.Context) {
			logger.Info("Wails OnStartup called — WebView2 ready", nil)
			shell.SetContext(ctx)
			go func() {
				logger.Info("Starting Go server services asynchronously in background", nil)
				app, err := server.Run(cfg, sm)
				if err != nil {
					logger.Info("Server startup failed in background", map[string]any{"error": err.Error()})
					fmt.Fprintf(os.Stderr, "server error: %v\n", err)
					os.Exit(1)
				}
				logger.Info("Server services started in background", map[string]any{
					"restPort":  cfg.ServerPort,
					"adminPort": cfg.AdminPort,
				})
				shell.app = app
			}()
		},
		OnDomReady: func(ctx context.Context) {
			logger.Info("Wails OnDomReady called — DOM loaded", nil)
			wailsruntime.WindowCenter(ctx)
			if len(startupArgs) > 0 {
				go func() {
					time.Sleep(1000 * time.Millisecond)
					for _, arg := range startupArgs {
						handleStartupArg(ctx, arg)
					}
				}()
			}
		},
		OnBeforeClose: func(ctx context.Context) bool {
			if sm != nil && sm.ConfirmBeforeClose() {
				selection, _ := wailsruntime.MessageDialog(ctx, wailsruntime.MessageDialogOptions{
					Type:          wailsruntime.QuestionDialog,
					Title:         "Close HTTP Toolkit?",
					Message:       "Are you sure you want to close HTTP Toolkit?",
					Buttons:       []string{"Yes", "No"},
					DefaultButton: "No",
					CancelButton:  "No",
				})
				return selection != "Yes"
			}
			return false
		},
		OnShutdown: func(ctx context.Context) {
			logger.Info("Wails OnShutdown called — cleaning up", nil)
			if shell.app != nil {
				shell.app.Shutdown(ctx)
			}
		},
		Bind: []interface{}{shell},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "tech.httptoolkit.pro",
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				if shell.ctx != nil {
					wailsruntime.WindowUnminimise(shell.ctx)
					wailsruntime.WindowShow(shell.ctx)
					for _, arg := range secondInstanceData.Args {
						handleStartupArg(shell.ctx, arg)
					}
				}
			},
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisablePinchZoom:     true,
		},
	})
	if err != nil {
		logger.Info("Wails run failed", map[string]any{"error": err.Error()})
		fmt.Fprintf(os.Stderr, "wails: %v\n", err)
		os.Exit(1)
	}
}
