package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/backup"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/logger"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/mcp"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/mcp/ctlclient"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/server"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/settings"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/system"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var frontendAssets embed.FS

//go:embed all:assets/overrides
var overridesAssets embed.FS

func handleStartupArg(arg string) {
	lower := strings.ToLower(arg)
	if strings.HasSuffix(lower, ".har") {
		if content, err := os.ReadFile(arg); err == nil {
			logger.Info("Emitting open-har-file event", map[string]any{"path": arg})
			application.Get().Event.Emit("open-har-file", map[string]any{
				"path":    arg,
				"content": string(content),
			})
		} else {
			logger.Info("Failed to read startup HAR file", map[string]any{"path": arg, "error": err.Error()})
		}
	} else if strings.HasPrefix(arg, "httptoolkitpro://") {
		logger.Info("Emitting auth-callback event", map[string]any{"url": arg})
		application.Get().Event.Emit("auth-callback", arg)
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

	if err := fs.WalkDir(overridesAssets, "assets/overrides", func(path string, d fs.DirEntry, err error) error {
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
		data, err := fs.ReadFile(overridesAssets, path)
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

// ---------------------------------------------------------------------------
// Subcommand: server (headless mode)
// ---------------------------------------------------------------------------

func runServerCmd(args []string) {
	cfg, err := config.LoadDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}
	extractEmbeddedOverrides(cfg)
	sm := settings.NewManager(cfg.ConfigDir)
	sm.SetBackup(backup.New(cfg))
	if err := sm.Load(); err != nil {
		logger.Error(err, map[string]any{"msg": "settings load failed"})
	}
	if os.Getenv("HTK_SERVER_TOKEN") == "" {
		os.Setenv("HTK_SERVER_TOKEN", cfg.AuthToken)
	}
	app, err := server.Run(cfg, sm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "server: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "httptoolkit-go listening REST=%d admin=%d token=%s\n",
		cfg.ServerPort, cfg.AdminPort, app.AuthToken)
	select {} // block forever
}

// ---------------------------------------------------------------------------
// Subcommand: mcp (MCP stdio server / status)
// ---------------------------------------------------------------------------

func runMCPCmd(args []string) {
	if len(args) > 0 && (args[0] == "stdio" || args[0] == "--stdio") {
		if err := mcp.RunStdio(os.Stderr); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	client, err := ctlclient.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ready, err := client.Status()
	if err != nil || ready["ready"] != true {
		fmt.Fprintln(os.Stderr, "UI bridge not ready")
		os.Exit(1)
	}
	ops, err := client.Operations()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	out := map[string]any{
		"protocol": "httptoolkit-mcp",
		"version":  "1.0.0",
		"tools":    ops,
	}
	data, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(data))
}

// ---------------------------------------------------------------------------
// Subcommand: ctl (CLI client)
// ---------------------------------------------------------------------------

func runCTLCmd(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: HttpToolkit-Pro ctl <status|operations|execute> [args-json]")
		os.Exit(2)
	}
	client, err := ctlclient.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	switch args[0] {
	case "status":
		out, err := client.Status()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		printJSON(out)
	case "operations":
		out, err := client.Operations()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		printJSON(out)
	case "execute":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: HttpToolkit-Pro ctl execute <name> <args-json>")
			os.Exit(2)
		}
		var jargs map[string]any
		if err := json.Unmarshal([]byte(args[2]), &jargs); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		out, err := client.Execute(args[1], jargs)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		printJSON(out)
	default:
		fmt.Fprintln(os.Stderr, "unknown command:", args[0])
		os.Exit(2)
	}
}

func printJSON(v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

// ---------------------------------------------------------------------------
// Default: launch GUI desktop app
// ---------------------------------------------------------------------------

func runDesktopApp() {
	// Parse command line arguments manually to check for verbose flags.
	verbose := false
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--verbose" {
			verbose = true
			break
		}
	}
	logger.SetVerbose(verbose)
	logger.Info("HttpToolkit Pro starting", map[string]any{"version": config.AppVersion()})

	// Disable WebView2 web security & allow mixed content so that Vercel (HTTPS)
	// UI can make REST/WebSocket requests to local (HTTP) server ports.
	// In Wails v3, AdditionalBrowserArgs on WindowsOptions achieves this.
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

	// Server startup/shutdown callbacks for the ShellService
	var serverApp *server.App
	serverStartup := func() {
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
		serverApp = app
	}
	serverShutdown := func() {
		if serverApp != nil {
			serverApp.Shutdown(context.Background())
		}
	}

	shellService := NewShellService(cfg, sm, serverStartup, serverShutdown)

	// Sub FS: strip the "frontend/dist/" prefix so Wails serves index.html at /
	subFS, err := fs.Sub(frontendAssets, "frontend/dist")
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

	logger.Info("Starting Wails v3 desktop app", map[string]any{
		"adminPort": cfg.AdminPort,
		"restPort":  cfg.ServerPort,
	})

	app := application.New(application.Options{
		Name:        "HttpToolkit Pro",
		Description: "Native Go MITM proxy + desktop shell for HttpToolkit Pro",
		Services: []application.Service{
			application.NewService(shellService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(subFS),
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "tech.httptoolkit.pro",
			OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
				window := application.Get().Window.Current()
				if window != nil {
					window.UnMinimise()
					window.Show()
					window.Focus()
				}
				for _, arg := range data.Args {
					handleStartupArg(arg)
				}
			},
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Windows: application.WindowsOptions{
			AdditionalBrowserArgs: []string{
				"--disable-web-security",
				"--allow-running-insecure-content",
				"--remote-debugging-port=9222",
			},
		},
		OnShutdown: func() {
			logger.Info("Application shutting down", nil)
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:           config.AppTitle(),
		Width:           1024,
		Height:          700,
		MinWidth:        800,
		MinHeight:       600,
		BackgroundColour: application.NewRGBA(15, 17, 23, 255),
	})

	// Center the window after creation (v3 doesn't have a WindowCentered option)
	window.Center()

	// Handle startup args after DOM is ready — use a short delay since v3
	// doesn't have an OnDomReady callback on the app options.
	if len(startupArgs) > 0 {
		go func() {
			time.Sleep(1500 * time.Millisecond)
			for _, arg := range startupArgs {
				handleStartupArg(arg)
			}
		}()
	}

	logger.Info("Running Wails v3 application", nil)
	if err := app.Run(); err != nil {
		logger.Info("Wails run failed", map[string]any{"error": err.Error()})
		fmt.Fprintf(os.Stderr, "wails: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	// Subcommand dispatch: if the first arg matches a known subcommand,
	// run that mode instead of the GUI.
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "server":
			runServerCmd(os.Args[2:])
			return
		case "mcp":
			runMCPCmd(os.Args[2:])
			return
		case "ctl":
			runCTLCmd(os.Args[2:])
			return
		case "-h", "--help", "help":
			printHelp()
			return
		}
	}
	runDesktopApp()
}

func printHelp() {
	fmt.Println("HttpToolkit Pro — Native Go MITM proxy + desktop shell")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  HttpToolkit-Pro              Launch the desktop GUI (default)")
	fmt.Println("  HttpToolkit-Pro server       Run as headless server (no GUI)")
	fmt.Println("  HttpToolkit-Pro mcp stdio    Run MCP stdio server for AI assistants")
	fmt.Println("  HttpToolkit-Pro mcp          Show MCP status and available tools")
	fmt.Println("  HttpToolkit-Pro ctl status   Query running instance status")
	fmt.Println("  HttpToolkit-Pro ctl operations  List available operations")
	fmt.Println("  HttpToolkit-Pro ctl execute <name> <args-json>  Execute an operation")
	fmt.Println("  HttpToolkit-Pro -v           Launch GUI with verbose logging")
	fmt.Println("  HttpToolkit-Pro help         Show this help message")
}
