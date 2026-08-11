package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/logger"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/settings"
)

func (s *Server) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	if s.settings == nil {
		writeJSON(w, settings.Default())
		return
	}
	writeJSON(w, s.settings.Get())
}

func (s *Server) handlePostSettings(w http.ResponseWriter, r *http.Request) {
	if s.settings == nil {
		writeAPIError(w, http.StatusServiceUnavailable, "settings not available")
		return
	}
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var st settings.Settings
	if err := json.Unmarshal(body, &st); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.settings.Save(st); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, s.settings.Get())
}

// backupExport is the shape exported by POST /backup/export.
type backupExport struct {
	Version       string         `json:"version"`
	ExportedAt    string         `json:"exportedAt"`
	Settings      map[string]any `json:"settings"`
	Notifications map[string]any `json:"notifications"`
	Cloud         map[string]any `json:"cloud"`
}

func (s *Server) handleBackupExport(w http.ResponseWriter, r *http.Request) {
	b := backupExport{
		Version:       "1.0.0-go",
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		Settings:      map[string]any{},
		Notifications: map[string]any{},
		Cloud:         map[string]any{},
	}
	if s.settings != nil {
		st := s.settings.Get()
		b.Settings["confirmBeforeClose"] = st.ConfirmBeforeClose
	}
	writeJSON(w, b)
}

func (s *Server) handleBackupImport(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var b backupExport
	if err := json.Unmarshal(body, &b); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	if s.settings != nil {
		if stMap, ok := b.Settings["confirmBeforeClose"].(bool); ok {
			st := s.settings.Get()
			st.ConfirmBeforeClose = stMap
			if err := s.settings.Save(st); err != nil {
				writeAPIError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}
	writeJSON(w, map[string]any{"ok": true, "version": b.Version, "importedAt": time.Now().UTC().Format(time.RFC3339)})
}

// notificationRegister is the payload for POST /notifications/register.
type notificationRegister struct {
	Token       string   `json:"token"`
	Enabled     bool     `json:"enabled"`
	EventTypes  []string `json:"eventTypes"`
	PushService string   `json:"pushService"`
}

func (s *Server) handleNotificationRegister(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var n notificationRegister
	if err := json.Unmarshal(body, &n); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	logger.Info("Notification registration received", map[string]any{
		"enabled":     n.Enabled,
		"service":     n.PushService,
		"eventTypes":  n.EventTypes,
		"tokenPrefix": n.Token[:minInt(8, len(n.Token))],
	})
	writeJSON(w, map[string]any{
		"ok":         true,
		"registered": n.Enabled,
		"eventTypes": n.EventTypes,
	})
}

func (s *Server) handleNotificationStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]any{
		"registered":  false,
		"enabled":     true,
		"eventTypes":  []string{},
		"pushService": "",
	})
}

func (s *Server) handleCloudStatus(w http.ResponseWriter, r *http.Request) {
	d, _ := s.loadDesktopSession()
	connected := d.AccessToken != "" && !isTokenExpired(d.AccessToken)
	resp := map[string]any{
		"connected":  connected,
		"lastSyncAt": d.LastSyncAt,
		"pending":    0,
	}
	if exp, ok := tokenExpiry(d.AccessToken); ok {
		resp["expiresAt"] = time.Unix(exp, 0).UTC().Format(time.RFC3339)
	}
	writeJSON(w, resp)
}

func (s *Server) handleCloudSync(w http.ResponseWriter, r *http.Request) {
	d, _ := s.loadDesktopSession()
	if d.AccessToken == "" {
		writeAPIError(w, http.StatusBadRequest, "no cloud session")
		return
	}
	d.LastSyncAt = time.Now().UTC().Format(time.RFC3339)
	if err := s.writeDesktopSession(d); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, map[string]any{
		"ok":       true,
		"syncedAt": d.LastSyncAt,
		"pending":  0,
	})
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// tokenExpiry parses a JWT access token and returns its exp claim if present.
// It tolerates opaque/non-JWT tokens by returning (0, false).
func tokenExpiry(token string) (int64, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, false
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, false
	}
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return 0, false
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return 0, false
	}
	return int64(exp), true
}

func isTokenExpired(token string) bool {
	exp, ok := tokenExpiry(token)
	if !ok {
		return false
	}
	return exp < time.Now().Unix()
}
