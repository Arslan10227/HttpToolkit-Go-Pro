package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	dockersvc "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/docker"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/origins"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/dns"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/events"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/mitm"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/rules"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/rtc"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/rtc/native"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/session"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/webextension"
)

type Server struct {
	cfg        *config.Config
	proxy      *mitm.Server
	rules      *rules.Engine
	events     *events.Bus
	rtc        *rtc.Manager
	sessions   *session.Manager
	dns        *dns.Server
	docker     *dockersvc.Session
	webext     *webextension.Manager
	nativeRTC  *native.Handler
	certMgr    *cert.Manager
	mux        *http.ServeMux
	httpServer *http.Server
	mu         sync.Mutex
	sessionID  string
	webrtcOn   bool
	http2On    bool
}

func New(cfg *config.Config, proxy *mitm.Server, eng *rules.Engine, bus *events.Bus, rtcMgr *rtc.Manager, sessions *session.Manager, webext *webextension.Manager, certMgr *cert.Manager) *Server {
	s := &Server{
		cfg: cfg, proxy: proxy, rules: eng, events: bus, rtc: rtcMgr,
		sessions: sessions, dns: dns.NewServer(nil), docker: dockersvc.NewSession(sessions, cfg.AssetsDir),
		webext: webext, certMgr: certMgr,
		nativeRTC: native.NewHandler(bus),
	}
	s.mux = http.NewServeMux()
	s.routes()
	return s
}

func (s *Server) routes() {
	// Method-qualified ("METHOD /path") Go 1.22+ ServeMux patterns are used
	// for every route that only ever supports one HTTP method. Unlike
	// internal/api/server.go, this mux has no catch-all "/" pattern, so
	// ServeMux's built-in behavior applies cleanly: a request to a
	// registered path with the wrong method gets a real 405 Method Not
	// Allowed instead of every handler re-implementing that check by hand
	// (verified: PUT-only /rules/http correctly 405s a GET).
	s.mux.HandleFunc("GET /metadata", s.handleMetadata)
	s.mux.HandleFunc("GET /admin/health", s.handleHealth)
	s.mux.HandleFunc("POST /admin/stop", s.handleStop)
	s.mux.HandleFunc("GET /admin/certificate-status", s.handleCertificateStatus)
	s.mux.HandleFunc("POST /session/start", s.handleStart)
	s.mux.HandleFunc("POST /session/stop", s.handleStop)
	s.mux.HandleFunc("POST /start", s.handleStart)
	s.mux.HandleFunc("POST /stop", s.handleStop)
	// /session/{id} has variable depth and serves GET/POST/OPTIONS depending
	// on the request shape (metadata vs GraphQL vs legacy event POST), so it
	// keeps its own internal method dispatch.
	s.mux.HandleFunc("/session/", s.handleMockRTCSession)
	s.mux.HandleFunc("PUT /rules/http", s.handleHTTPRules)
	s.mux.HandleFunc("PUT /rules/ws", s.handleWSRules)
	s.mux.HandleFunc("PUT /rules/rtc", s.handleRTCRules)
	s.mux.HandleFunc("GET /events", s.handleEventsWS)
}

func (s *Server) ListenAndServe() error {
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", s.cfg.AdminPort),
		Handler: s.loggingMiddleware(s.corsMiddleware(s.mux)),
	}
	Info("Admin server listening", map[string]any{"addr": s.httpServer.Addr})
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	_ = s.stopSession()
	if s.httpServer != nil {
		return s.httpServer.Close()
	}
	return nil
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Info("Admin request", map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
			"origin": r.Header.Get("Origin"),
			"remote": r.RemoteAddr,
		})
		next.ServeHTTP(w, r)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if r.Header.Get("Access-Control-Request-Private-Network") != "" {
			w.Header().Set("Access-Control-Allow-Private-Network", "true")
		}
		allowed := origins.IsAllowed(origin, s.cfg.DevMode)
		Debug("CORS check", map[string]any{"origin": origin, "allowed": allowed, "devMode": s.cfg.DevMode})
		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	Info("Health check", map[string]any{"addr": s.cfg.AdminPort})
	writeJSON(w, map[string]any{"ok": true})
}

func (s *Server) handleMetadata(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{
		"running":       s.proxy.IsRunning(),
		"httpPort":      s.proxy.Port(),
		"webrtcEnabled": s.webrtcOn && s.rtc.IsRunning(),
	})
}

func (s *Server) handleStart(w http.ResponseWriter, r *http.Request) {
	var body struct {
		HTTP2     *bool `json:"http2"`
		Socks     *bool `json:"socks"`
		PortRange *struct {
			StartPort int `json:"startPort"`
			EndPort   int `json:"endPort"`
		} `json:"portRange"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	s.http2On = body.HTTP2 != nil && *body.HTTP2
	s.proxy.SetHTTP2Enabled(s.http2On)
	socksOn := body.Socks == nil || *body.Socks
	s.proxy.SetSocksEnabled(socksOn)

	var err error
	if body.PortRange != nil && body.PortRange.StartPort > 0 {
		err = s.proxy.StartPortRange(body.PortRange.StartPort, body.PortRange.EndPort)
	} else {
		err = s.proxy.Start(config.DefaultProxyPort)
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = s.dns.Start()
	dnsAddr := s.dns.Address()
	dnsServers := []string{dnsAddr}
	certPath := filepath.Join(s.cfg.ConfigDir, "ca.pem")
	_ = s.docker.Start(s.proxy.Port(), certPath)

	_ = s.rtc.Start(s.proxy.Port())
	s.webrtcOn = s.rtc.IsRunning()
	// Native backend needs no subprocess: the in-process pion-backed handler
	// (s.nativeRTC) is already ready to answer GraphQL requests on this same
	// admin server via /session/:id.
	s.beginProxySession()

	ruleKeys := s.sessions.RuleParameterKeys()
	s.sessions.SetActive(s.proxy.Port(), dnsServers, ruleKeys, s.http2On, s.webrtcOn)
	s.sessions.SetSocksPort(s.proxy.SocksPort())
	if tp := s.docker.TunnelPort(); tp > 0 {
		s.sessions.SetDockerTunnelPort(s.proxy.Port(), tp)
	}

	writeJSON(w, map[string]any{
		"httpProxyPort":    s.proxy.Port(),
		"socksProxyPort":   s.proxy.SocksPort(),
		"dockerTunnelPort": s.docker.TunnelPort(),
		"webrtcEnabled":    s.webrtcOn,
	})
}

func (s *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	Info("Stop request received", nil)
	_ = s.stopSession()
	writeJSON(w, map[string]any{"ok": true})
}

func (s *Server) stopSession() error {
	Info("Stopping session", nil)
	s.endProxySession()
	if s.nativeRTC != nil {
		s.nativeRTC.Close()
	}
	_ = s.proxy.Stop()
	_ = s.rtc.Stop()
	_ = s.dns.Stop()
	s.docker.Stop()
	s.sessions.Clear()
	s.webrtcOn = false
	s.mu.Lock()
	s.sessionID = ""
	s.mu.Unlock()
	return nil
}

func (s *Server) handleHTTPRules(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Rules []rules.RequestRuleData `json:"rules"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	s.rules.SetHTTPRules(body.Rules)
	writeJSON(w, map[string]any{"ok": true})
}

func (s *Server) handleWSRules(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Rules []rules.WebSocketRuleData `json:"rules"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	s.rules.SetWSRules(body.Rules)
	writeJSON(w, map[string]any{"ok": true})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	Error(fmt.Errorf("%s", msg), map[string]any{"code": code})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"message": msg}})
}

// handleCertificateStatus reports whether the CA certificate is installed.
func (s *Server) handleCertificateStatus(w http.ResponseWriter, r *http.Request) {
	installed := s.certMgr != nil
	writeJSON(w, map[string]any{"installed": installed, "justInstalled": false})
	Info("Certificate status requested", map[string]any{"installed": installed})
}
