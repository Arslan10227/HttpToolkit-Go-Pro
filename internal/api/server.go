package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/api/gql"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/interceptors"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/logger"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/origins"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/mitm"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/session"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/settings"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/uibridge"
)

type Server struct {
	cfg          *config.Config
	certs        *cert.Manager
	spki         string
	settings     *settings.Manager
	interceptors *interceptors.Registry
	bridge       *uibridge.Bridge
	sessions     *session.Manager
	bpMgr        *mitm.BreakpointManager
	mux          *http.ServeMux
	httpServer   *http.Server
	shutdownOnce sync.Once
	onShutdown   func()
}

func New(cfg *config.Config, certs *cert.Manager, spki string, sm *settings.Manager, reg *interceptors.Registry, bridge *uibridge.Bridge, sessions *session.Manager, bpMgr *mitm.BreakpointManager) *Server {
	s := &Server{cfg: cfg, certs: certs, spki: spki, settings: sm, interceptors: reg, bridge: bridge, sessions: sessions, bpMgr: bpMgr}
	s.mux = http.NewServeMux()
	s.routes()
	return s
}

func (s *Server) SetShutdownHandler(fn func()) {
	s.onShutdown = fn
}

func (s *Server) routes() {
	// Routes are registered using Go 1.22+ ServeMux method-qualified patterns
	// ("METHOD /path") wherever a route only ever supports one HTTP method,
	// so each handler no longer has to hand-roll its own method check.
	// Note: because "/" is also registered below as a catch-all, ServeMux's
	// built-in "wrong method -> 405" behavior is masked for unmatched
	// method+path combos here (see golang/go#65648) — they fall through to
	// handleRoot's 404 instead. That's an acceptable trade-off: no client
	// (webui's proxyAdminClient.ts) ever sends the wrong method to these
	// routes, so this only affects the status code for hypothetical
	// malformed requests, not real traffic.
	s.mux.HandleFunc("/", s.handleRoot)
	s.mux.Handle("POST /graphql", gql.Handler(s)) // GraphQL API (mirrors httptoolkit-server's GraphQL endpoint)
	s.mux.HandleFunc("GET /auth/client-token", s.handleClientToken)
	s.mux.HandleFunc("GET /auth/oauth", s.handleOAuthStart)
	s.mux.HandleFunc("GET /auth/desktop-session", s.handleDesktopSessionGet)
	s.mux.HandleFunc("POST /auth/desktop-session", s.handleDesktopSessionPost)
	s.mux.Handle("DELETE /auth/desktop-session", http.HandlerFunc(s.handleDesktopSessionDelete))
	s.mux.HandleFunc("GET /auth/callback", s.handleAuthCallback)
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /ready", s.handleHealth)
	s.mux.HandleFunc("GET /version", s.handleVersion)
	s.mux.HandleFunc("GET /config", s.handleConfig)
	s.mux.HandleFunc("POST /config/backup", s.handleConfigBackupPost)
	s.mux.HandleFunc("GET /config/restore", s.handleConfigRestoreGet)
	s.mux.HandleFunc("GET /config/sync-status", s.handleConfigSyncStatusGet)
	s.mux.HandleFunc("GET /config/network-interfaces", s.handleNetworkInterfaces)
	s.mux.HandleFunc("GET /certificate/status", s.handleCertStatus)
	s.mux.HandleFunc("GET /certificate/export", s.handleCertExport)
	s.mux.HandleFunc("POST /certificate/install", s.handleCertInstall)
	s.mux.HandleFunc("GET /java/versions", s.handleJavaVersions)
	s.mux.HandleFunc("GET /java/certificate/status", s.handleJavaCertStatus)
	s.mux.HandleFunc("POST /java/certificate/install", s.handleJavaCertInstall)
	s.mux.HandleFunc("GET /interceptors", s.handleInterceptors)
	// Android ADB certificate endpoints (match before the /interceptors/ catch-all)
	s.mux.HandleFunc("GET /interceptors/android-adb/certificate/status", s.handleAndroidAdbCertStatus)
	s.mux.HandleFunc("POST /interceptors/android-adb/certificate/install", s.handleAndroidAdbCertInstall)
	// Sub-paths have variable depth (/interceptors/{id}/metadata[/{subId}],
	// /activate/{port}, /deactivate/{port}) and support GET+POST depending on
	// the segment, so they keep their own internal method/path parsing.
	s.mux.HandleFunc("/interceptors/", s.handleInterceptorSub)
	s.mux.HandleFunc("/session/breakpoint/", s.handleBreakpointSub)
	s.mux.HandleFunc("POST /client/send", s.handleClientSend)
	s.mux.HandleFunc("POST /snippets/generate", s.handleSnippets)
	s.mux.HandleFunc("POST /rules/bulk-create", s.handleRulesBulkCreate)
	s.mux.HandleFunc("POST /update", s.handleUpdate)
	s.mux.HandleFunc("POST /shutdown", s.handleShutdown)
	s.mux.HandleFunc("GET /mcp/status", s.handleMCPStatus)
	s.mux.HandleFunc("GET /mcp/tools", s.handleMCPTools)
	s.mux.HandleFunc("POST /webhooks/capture", s.handleWebhooksCapture)
	s.mux.HandleFunc("/ui-operations", s.bridge.HandleWebSocket)
	s.mux.HandleFunc("GET /settings", s.handleGetSettings)
	s.mux.HandleFunc("POST /settings", s.handlePostSettings)
	s.mux.HandleFunc("POST /backup/export", s.handleBackupExport)
	s.mux.HandleFunc("POST /backup/import", s.handleBackupImport)
	s.mux.HandleFunc("POST /notifications/register", s.handleNotificationRegister)
	s.mux.HandleFunc("GET /notifications/status", s.handleNotificationStatus)
	s.mux.HandleFunc("GET /cloud/status", s.handleCloudStatus)
	s.mux.HandleFunc("POST /cloud/sync", s.handleCloudSync)
}

func (s *Server) ListenAndServe() error {
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", s.cfg.ServerPort),
		Handler: s.middleware(s.mux),
	}
	logger.Info("REST API server listening", map[string]any{"addr": s.httpServer.Addr})
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

func (s *Server) HandleUpgrade(w http.ResponseWriter, r *http.Request) {
	s.bridge.HandleWebSocket(w, r)
}

func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("REST API request", map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
			"origin": r.Header.Get("Origin"),
			"remote": r.RemoteAddr,
		})
		if r.Header.Get("Access-Control-Request-Private-Network") != "" {
			w.Header().Set("Access-Control-Allow-Private-Network", "true")
		}
		origin := r.Header.Get("Origin")
		allowed := origins.IsAllowed(origin, s.cfg.DevMode)
		logger.Debug("REST CORS check", map[string]any{"origin": origin, "allowed": allowed})
		if allowed && origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Max-Age", "86400")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		if r.URL.Path == "/ui-operations" {
			next.ServeHTTP(w, r)
			return
		}

		if r.URL.Path == "/auth/desktop-session" {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		expected := "Bearer " + s.cfg.AuthToken
		if auth != expected {
			logger.Info("REST auth rejected", map[string]any{"path": r.URL.Path, "got": auth[:min(len(auth), 16)] + "..."})
			writeAPIError(w, http.StatusForbidden, "Invalid or missing auth token")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isPublicPath(path string) bool {
	switch path {
	case "/", "/health", "/ready", "/auth/client-token", "/auth/oauth", "/auth/callback",
		"/config/backup", "/config/restore", "/config/sync-status":
		return true
	}
	if path == "/auth/desktop-session" {
		return true
	}
	return false
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		writeAPIError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, map[string]any{
		"ok": true, "version": config.ServerVersion(),
		"message": "HttpToolkit Server REST+GraphQL API",
		"endpoints": map[string]string{
			"health": "/health", "ready": "/ready", "config": "/config",
			"interceptors": "/interceptors", "graphql": "/graphql",
		},
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{"ok": true, "version": config.ServerVersion()})
}

func (s *Server) handleVersion(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{"version": config.ServerVersion()})
}

func (s *Server) handleClientToken(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{"token": s.cfg.AuthToken})
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]any{"config": s.buildConfig()})
}

func (s *Server) buildConfig() map[string]any {
	dnsServers := []string{}
	ruleKeys := []string{}
	if s.sessions != nil {
		dnsServers = s.sessions.DNSServers()
		ruleKeys = s.sessions.RuleParameterKeys()
	}
	cfg := map[string]any{
		"authToken":              s.cfg.AuthToken,
		"certificatePath":        s.certs.CertPath(),
		"certificateContent":     s.certs.CertPEM(),
		"certificateFingerprint": s.spki,
		"certificateFiles":       s.certs.CertFiles(),
		"dnsServers":             dnsServers,
		"ruleParameterKeys":      ruleKeys,
		"toolPaths":              toolPaths(s.cfg),
	}
	if sp := systemProxy(); sp != nil {
		cfg["systemProxy"] = sp
	}
	return cfg
}

func toolPaths(cfg *config.Config) map[string][]string {
	return map[string][]string{
		"ctl": {resolveToolPath("htk-ctl", cfg)},
		"mcp": {resolveToolPath("htk-mcp", cfg)},
	}
}

func resolveToolPath(name string, cfg *config.Config) string {
	if p, err := execLookPath(name); err == nil {
		return p
	}
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	if exe, err := os.Executable(); err == nil {
		candidate := filepath.Join(filepath.Dir(exe), name+ext)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	candidate := filepath.Join(cfg.AssetsDir, "bin", name+ext)
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	return name + ext
}

func execLookPath(name string) (string, error) {
	return exec.LookPath(name)
}

func (s *Server) handleNetworkInterfaces(w http.ResponseWriter, _ *http.Request) {
	ifaces, _ := net.Interfaces()
	out := make(map[string][]map[string]any)
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		entries := make([]map[string]any, 0, len(addrs))
		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil {
				continue
			}
			family := "IPv6"
			if v4 := ip.To4(); v4 != nil {
				ip = v4
				family = "IPv4"
			}
			internal := ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
			entries = append(entries, map[string]any{
				"address":  ip.String(),
				"family":   family,
				"internal": internal,
			})
		}
		if len(entries) > 0 {
			out[iface.Name] = entries
		}
	}
	writeJSON(w, map[string]any{"networkInterfaces": out})
}

func (s *Server) handleShutdown(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{"success": true})
	s.shutdownOnce.Do(func() {
		if s.onShutdown != nil {
			go s.onShutdown()
		}
	})
}

func (s *Server) handleUpdate(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{"success": true})
}

func (s *Server) handleRulesBulkCreate(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{
			"message": "Use Mockttp admin rule APIs (setHttpRules / setWebSocketRules). This REST stub does not persist rules.",
			"code":    "NOT_IMPLEMENTED",
		},
	})
}

func (s *Server) handleWebhooksCapture(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{
			"message": "Webhook capture forwarding is not configured on this server build",
		},
	})
}

func (s *Server) handleNotImplemented(w http.ResponseWriter, _ *http.Request) {
	writeAPIError(w, http.StatusNotImplemented, "Not implemented")
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func writeAPIError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"message": msg},
	})
}

func (s *Server) handleBreakpointSub(w http.ResponseWriter, r *http.Request) {
	if s.bpMgr == nil {
		writeAPIError(w, http.StatusNotImplemented, "Breakpoint manager not configured")
		return
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/session/breakpoint/"), "/")
	if len(parts) < 2 {
		writeAPIError(w, http.StatusBadRequest, "Invalid breakpoint path")
		return
	}
	id := parts[0]
	action := parts[1]

	pr := s.bpMgr.Get(id)
	if pr == nil {
		writeAPIError(w, http.StatusNotFound, "Breakpoint not found or already resumed")
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch action {
	case "resume":
		var body struct {
			Method  string            `json:"method"`
			URL     string            `json:"url"`
			Headers map[string]string `json:"headers"`
			Body    string            `json:"body"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)

		var reqBodyBytes []byte
		if body.Body != "" {
			if decoded, err := base64.StdEncoding.DecodeString(body.Body); err == nil {
				reqBodyBytes = decoded
			} else {
				reqBodyBytes = []byte(body.Body)
			}
		}

		select {
		case pr.ActionChan <- mitm.BreakpointAction{
			Type:    mitm.ActionResume,
			Method:  body.Method,
			URL:     body.URL,
			Headers: body.Headers,
			Body:    reqBodyBytes,
		}:
			writeJSON(w, map[string]any{"success": true})
		default:
			writeAPIError(w, http.StatusConflict, "Breakpoint already processed")
		}

	case "respond":
		var body struct {
			StatusCode int               `json:"statusCode"`
			Headers    map[string]string `json:"headers"`
			Body       string            `json:"body"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)

		var respBodyBytes []byte
		if body.Body != "" {
			if decoded, err := base64.StdEncoding.DecodeString(body.Body); err == nil {
				respBodyBytes = decoded
			} else {
				respBodyBytes = []byte(body.Body)
			}
		}

		select {
		case pr.ActionChan <- mitm.BreakpointAction{
			Type:        mitm.ActionRespond,
			StatusCode:  body.StatusCode,
			RespHeaders: body.Headers,
			RespBody:    respBodyBytes,
		}:
			writeJSON(w, map[string]any{"success": true})
		default:
			writeAPIError(w, http.StatusConflict, "Breakpoint already processed")
		}

	case "abort":
		select {
		case pr.ActionChan <- mitm.BreakpointAction{Type: mitm.ActionAbort}:
			writeJSON(w, map[string]any{"success": true})
		default:
			writeAPIError(w, http.StatusConflict, "Breakpoint already processed")
		}

	default:
		writeAPIError(w, http.StatusBadRequest, "Unknown action")
	}
}
