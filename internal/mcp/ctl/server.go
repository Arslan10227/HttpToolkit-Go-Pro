package ctl

import (
	"encoding/json"
	"net"
	"net/http"
	"runtime"
	"strings"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/uibridge"
)

type Server struct {
	bridge *uibridge.Bridge
}

func New(bridge *uibridge.Bridge) *Server {
	return &Server{bridge: bridge}
}

func (s *Server) Listen() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]any{"ready": s.bridge.IsReady()})
	})
	mux.HandleFunc("/api/operations", func(w http.ResponseWriter, _ *http.Request) {
		if !s.bridge.IsReady() {
			w.WriteHeader(http.StatusServiceUnavailable)
			writeJSON(w, map[string]any{"error": "not_ready"})
			return
		}
		writeJSON(w, s.bridge.Operations())
	})
	mux.HandleFunc("/api/execute", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name   string         `json:"name"`
			Args   map[string]any `json:"args"`
			Source string         `json:"source"`
			JWT    string         `json:"jwt"`
			Tier   string         `json:"tier"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if !operationAllowed(s.bridge, body.Name, body.Tier) {
			w.WriteHeader(http.StatusForbidden)
			writeJSON(w, map[string]any{"error": "tier_not_allowed"})
			return
		}
		result, err := s.bridge.Execute(body.Name, body.Args)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			writeJSON(w, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, result)
	})

	if runtime.GOOS == "windows" {
		ln, err := winListenPipe(`\\.\pipe\httptoolkit-ctl`)
		if err != nil {
			return err
		}
		return http.Serve(ln, mux)
	}
	ln, err := net.Listen("unix", socketPath())
	if err != nil {
		return err
	}
	return http.Serve(ln, mux)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func operationAllowed(bridge *uibridge.Bridge, name, tier string) bool {
	if tier == "" || tier == "pro" {
		return true
	}
	for _, op := range bridge.Operations() {
		if op.Name != name {
			continue
		}
		if len(op.Tiers) == 0 {
			return true
		}
		for _, t := range op.Tiers {
			if strings.EqualFold(t, tier) {
				return true
			}
		}
		return false
	}
	return true
}

func socketPath() string {
	return "/tmp/httptoolkit-ctl.sock"
}
