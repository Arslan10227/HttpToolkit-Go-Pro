package admin

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/rtc"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/rtc/native"
)

func (s *Server) handleMockRTCSession(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/session/") || r.URL.Path == "/session/start" || r.URL.Path == "/session/stop" {
		http.NotFound(w, r)
		return
	}
	sessionID := strings.TrimPrefix(r.URL.Path, "/session/")
	sessionID = strings.Trim(sessionID, "/")
	if idx := strings.Index(sessionID, "/"); idx >= 0 {
		sessionID = sessionID[:idx]
	}
	if sessionID == "" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet, http.MethodOptions:
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		writeJSON(w, map[string]any{
			"sessionId": sessionID,
			"running":   s.webrtcOn && s.rtc.IsRunning(),
		})
	case http.MethodPost:
		if s.nativeRTC != nil && s.webrtcOn {
			// The bundled webextension's MockRTC client speaks GraphQL
			// against this endpoint (createOffer/answerOffer/etc — see
			// assets/overrides/webextension/build/background.js). Detect
			// that shape and dispatch to the native pion-backed handler;
			// fall back to the legacy raw-event POST shape otherwise (e.g.
			// admin tooling that still posts { event, data } directly).
			body, err := io.ReadAll(r.Body)
			if err != nil {
				writeErr(w, http.StatusBadRequest, err.Error())
				return
			}
			var probe struct {
				Query string `json:"query"`
			}
			if err := json.Unmarshal(body, &probe); err == nil && probe.Query != "" {
				r.Body = io.NopCloser(bytes.NewReader(body))
				native.GraphQLHandler(s.nativeRTC).ServeHTTP(w, r)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(body))
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		event := rtc.EventFromPayload(body)
		if event != "" && rtc.ValidEvent(event) {
			s.rtc.PublishForSession(event, body, sessionID)
		}
		writeJSON(w, map[string]any{"ok": true})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) beginProxySession() string {
	sessionID := uuid.NewString()
	rtcPort := s.cfg.AdminPort
	s.mu.Lock()
	s.sessionID = sessionID
	s.mu.Unlock()
	if s.webext != nil {
		_ = s.webext.UpdateConfig(s.proxy.Port(), sessionID, s.webrtcOn, rtcPort)
	}
	return sessionID
}

func (s *Server) endProxySession() {
	if s.webext == nil {
		return
	}
	port := s.proxy.Port()
	if port > 0 {
		s.webext.ClearConfig(port)
	}
}
