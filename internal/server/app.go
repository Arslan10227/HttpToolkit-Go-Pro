package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/api"
	certmgr "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/interceptors"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/mcp/ctl"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/admin"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/events"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/mitm"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/rules"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/rtc"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/session"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/settings"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/uibridge"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/webextension"
)

type App struct {
	Cfg          *config.Config
	AuthToken    string
	Settings     *settings.Manager
	REST         *api.Server
	Admin        *admin.Server
	Interceptors *interceptors.Registry
	Sessions     *session.Manager
	proxy        *mitm.Server
	webext       *webextension.Manager
}

func Run(cfg *config.Config, sm *settings.Manager) (*App, error) {
	certs, err := certmgr.NewManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("cert: %w", err)
	}

	spki, err := certs.SPKIFingerprint()
	if err != nil {
		return nil, fmt.Errorf("spki: %w", err)
	}

	sessions := session.NewManager()
	bus := events.NewBus()
	eng := rules.NewEngine()
	bpMgr := mitm.NewBreakpointManager()
	proxy := mitm.NewServer(certs, eng, bus)
	proxy.SetSession(sessions)
	proxy.SetBreakpointManager(bpMgr)
	rtcMgr := rtc.NewManager(bus)
	proxy.SetRTC(rtcMgr)
	webext := webextension.NewManager(cfg.AssetsDir)
	adminSrv := admin.New(cfg, proxy, eng, bus, rtcMgr, sessions, webext, certs)

	reg := interceptors.NewRegistry(cfg, spki, certs)
	interceptors.SetWebExtension(webext)
	bridge := uibridge.New(cfg.AuthToken)
	rest := api.New(cfg, certs, spki, sm, reg, bridge, sessions, bpMgr)
	rest.SetShutdownHandler(func() {
		reg.DeactivateAll()
		_ = proxy.Stop()
		_ = adminSrv.Shutdown()
		_ = webext.Cleanup()
		os.Exit(0)
	})

	go func() {
		if err := adminSrv.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "admin: %v\n", err)
		}
	}()

	go func() {
		if err := ctl.New(bridge).Listen(); err != nil {
			fmt.Fprintf(os.Stderr, "ctl: %v\n", err)
		}
	}()

	go func() {
		if err := rest.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "rest: %v\n", err)
		}
	}()

	app := &App{
		Cfg: cfg, AuthToken: cfg.AuthToken, Settings: sm, REST: rest,
		Admin: adminSrv, Interceptors: reg, Sessions: sessions, proxy: proxy,
		webext: webext,
	}

	go waitForSignal(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		reg.DeactivateAll()
		_ = proxy.Stop()
		_ = adminSrv.Shutdown()
		_ = webext.Cleanup()
		_ = rest.Shutdown(ctx)
	})

	return app, nil
}

func (a *App) Shutdown(ctx context.Context) {
	if a.Interceptors != nil {
		a.Interceptors.DeactivateAll()
	}
	if a.proxy != nil {
		_ = a.proxy.Stop()
	}
	if a.Admin != nil {
		_ = a.Admin.Shutdown()
	}
	if a.webext != nil {
		_ = a.webext.Cleanup()
	}
	if a.REST != nil {
		_ = a.REST.Shutdown(ctx)
	}
}

func waitForSignal(fn func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fn()
}
