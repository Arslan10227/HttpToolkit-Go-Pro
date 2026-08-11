package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/auth"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/backup"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
)

var desktopSessionMu sync.Mutex

type desktopSession struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	Picture      string `json:"picture"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	LastSyncAt   string `json:"lastSyncAt"`
}

func (s *Server) desktopSessionPath() string {
	return filepath.Join(s.cfg.ConfigDir, "desktop-session.json")
}

func (s *Server) loadDesktopSession() (desktopSession, error) {
	if s.cfg.ConfigDir == "" {
		return desktopSession{}, nil
	}
	desktopSessionMu.Lock()
	defer desktopSessionMu.Unlock()

	data, err := os.ReadFile(s.desktopSessionPath())
	if err != nil {
		if os.IsNotExist(err) {
			return desktopSession{}, nil
		}
		return desktopSession{}, err
	}
	var d desktopSession
	if err := json.Unmarshal(data, &d); err != nil {
		return desktopSession{}, err
	}
	return d, nil
}

func (s *Server) writeDesktopSession(d desktopSession) error {
	if s.cfg.ConfigDir == "" {
		return nil
	}
	desktopSessionMu.Lock()
	defer desktopSessionMu.Unlock()

	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return os.WriteFile(s.desktopSessionPath(), data, 0o600)
}

func (s *Server) clearDesktopSession() error {
	if s.cfg.ConfigDir == "" {
		return nil
	}
	desktopSessionMu.Lock()
	defer desktopSessionMu.Unlock()
	return os.Remove(s.desktopSessionPath())
}

func (s *Server) handleDesktopSessionPost(w http.ResponseWriter, r *http.Request) {
	var body struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.AccessToken == "" {
		writeAPIError(w, http.StatusBadRequest, "access_token required")
		return
	}

	// Verify the Google ID token and extract the user's profile.
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	g := auth.NewGoogle(s.cfg)
	user, err := g.VerifyIDToken(ctx, body.AccessToken)
	if err != nil {
		writeAPIError(w, http.StatusUnauthorized, "invalid token: "+err.Error())
		return
	}

	d, _ := s.loadDesktopSession()
	d.Email = user.Email
	d.Name = user.Name
	d.Picture = user.Picture
	d.AccessToken = body.AccessToken
	d.RefreshToken = body.RefreshToken
	d.LastSyncAt = time.Now().UTC().Format(time.RFC3339)
	if err := s.writeDesktopSession(d); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, map[string]any{"ok": true, "email": user.Email, "name": user.Name})
}

func (s *Server) handleDesktopSessionGet(w http.ResponseWriter, _ *http.Request) {
	d, _ := s.loadDesktopSession()
	if d.AccessToken == "" {
		writeAPIError(w, http.StatusNotFound, "No session")
		return
	}
	writeJSON(w, map[string]any{
		"email":        d.Email,
		"name":         d.Name,
		"picture":      d.Picture,
		"accessToken":  d.AccessToken,
		"refreshToken": d.RefreshToken,
	})
}

func (s *Server) handleDesktopSessionDelete(w http.ResponseWriter, _ *http.Request) {
	_ = s.clearDesktopSession()
	writeJSON(w, map[string]any{"ok": true})
}

func (s *Server) handleCertStatus(w http.ResponseWriter, _ *http.Request) {
	installed, _ := cert.IsSystemCertInstalled(s.certs.CertPEM())
	writeJSON(w, map[string]any{"installed": installed})
}

func (s *Server) handleCertExport(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "pem"
	}
	switch format {
	case "pem", "crt", "cer":
		data := []byte(s.certs.CertPEM())
		w.Header().Set("Content-Disposition", "attachment; filename=httptoolkit-ca."+format)
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(data)
	case "p12":
		data, err := s.certs.ExportP12("httptoolkit")
		if err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Disposition", "attachment; filename=httptoolkit-ca.p12")
		w.Header().Set("Content-Type", "application/x-pkcs12")
		_, _ = w.Write(data)
	default:
		writeAPIError(w, http.StatusBadRequest, "unsupported format: "+format)
	}
}

func (s *Server) handleCertInstall(w http.ResponseWriter, _ *http.Request) {
	ok, msg, err := cert.InstallSystemCert(s.certs.CertPath())
	writeJSON(w, map[string]any{
		"success": err == nil, "installed": ok, "message": msg,
	})
}

func (s *Server) handleJavaVersions(w http.ResponseWriter, _ *http.Request) {
	versions, _ := cert.DetectJavaVersions()
	writeJSON(w, map[string]any{"versions": versions})
}

func (s *Server) handleJavaCertStatus(w http.ResponseWriter, r *http.Request) {
	javaPath := r.URL.Query().Get("javaPath")
	if javaPath == "" {
		writeAPIError(w, http.StatusBadRequest, "javaPath required")
		return
	}
	installed, _ := cert.IsJavaCertInstalled(javaPath)
	writeJSON(w, map[string]any{"installed": installed})
}

func (s *Server) handleJavaCertInstall(w http.ResponseWriter, r *http.Request) {
	var body struct {
		JavaPath string `json:"javaPath"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.JavaPath == "" {
		writeAPIError(w, http.StatusBadRequest, "javaPath required")
		return
	}
	result, err := cert.InstallJavaCert(body.JavaPath, s.certs.CertPath())
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, result)
}

func (s *Server) handleInterceptors(w http.ResponseWriter, r *http.Request) {
	proxyPort := 0
	if p := r.URL.Query().Get("proxyPort"); p != "" {
		_, _ = fmt.Sscan(p, &proxyPort)
	}
	list := s.interceptors.List(proxyPort)
	writeJSON(w, map[string]any{"interceptors": list})
}

func (s *Server) handleOAuthStart(w http.ResponseWriter, r *http.Request) {
	desktop := r.URL.Query().Get("desktop") == "1"
	state := "web"
	if desktop {
		state = "desktop"
	}
	g := auth.NewGoogle(s.cfg)
	authURL, err := g.AuthURL(desktop, state)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, authURL, http.StatusFound)
}

// handleAuthCallback receives the OAuth 2.0 authorization code from Google,
// exchanges it for a token, and either deep-links back to the desktop app or
// redirects to the hosted web UI.
func (s *Server) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		writeAPIError(w, http.StatusBadRequest, "code required")
		return
	}
	state := r.URL.Query().Get("state")
	desktop := state == "desktop"

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	g := auth.NewGoogle(s.cfg)
	user, idToken, refreshToken, err := g.Exchange(ctx, code, desktop)
	if err != nil {
		writeAPIError(w, http.StatusUnauthorized, err.Error())
		return
	}

	d, _ := s.loadDesktopSession()
	d.Email = user.Email
	d.Name = user.Name
	d.Picture = user.Picture
	d.AccessToken = idToken
	d.RefreshToken = refreshToken
	d.LastSyncAt = time.Now().UTC().Format(time.RFC3339)
	_ = s.writeDesktopSession(d)

	if desktop {
		q := url.Values{}
		q.Set("access_token", idToken)
		q.Set("refresh_token", refreshToken)
		http.Redirect(w, r, "httptoolkitpro://auth/callback?"+q.Encode(), http.StatusFound)
		return
	}

	// Web flow: send the user back to the hosted callback with the tokens in
	// the hash, just like the old Supabase flow.
	q := url.Values{}
	q.Set("access_token", idToken)
	q.Set("refresh_token", refreshToken)
	http.Redirect(w, r, "https://httptoolkitpro.vercel.app/auth/callback#"+q.Encode(), http.StatusFound)
}

// handleAndroidAdbCertStatus checks if the HttpToolkit CA certificate is
// installed in the Android system trust store for a given ADB device.
func (s *Server) handleAndroidAdbCertStatus(w http.ResponseWriter, r *http.Request) {
	deviceId := r.URL.Query().Get("deviceId")
	if deviceId == "" {
		writeAPIError(w, http.StatusBadRequest, "deviceId query parameter is required")
		return
	}
	status, err := s.interceptors.AndroidAdbCertStatus(deviceId, s.certs.CertPEM())
	if err != nil {
		writeJSON(w, map[string]any{
			"status": map[string]any{
				"deviceId":      deviceId,
				"installed":     false,
				"rootAvailable": false,
				"message":       err.Error(),
			},
		})
		return
	}
	writeJSON(w, map[string]any{"status": status})
}

// handleAndroidAdbCertInstall installs the HttpToolkit CA certificate into
// the Android system trust store for a given ADB device (requires root).
func (s *Server) handleAndroidAdbCertInstall(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DeviceId string `json:"deviceId"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.DeviceId == "" {
		writeAPIError(w, http.StatusBadRequest, "deviceId is required")
		return
	}
	result, err := s.interceptors.AndroidAdbCertInstall(body.DeviceId, s.certs.CertPEM(), s.certs.CertPath())
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, map[string]any{"result": result})
}

func (s *Server) handleInterceptorSub(w http.ResponseWriter, r *http.Request) {
	// /interceptors/:id/metadata, /activate/:port, /deactivate/:port
	path := strings.TrimPrefix(r.URL.Path, "/interceptors/")
	parts := strings.Split(path, "/")
	if len(parts) < 1 {
		writeAPIError(w, http.StatusNotFound, "not found")
		return
	}
	id := parts[0]
	if len(parts) == 3 && parts[1] == "metadata" {
		meta, err := s.interceptors.SubMetadata(id, parts[2])
		if err != nil {
			writeAPIError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, map[string]any{"interceptorMetadata": meta})
		return
	}
	if len(parts) == 2 && parts[1] == "metadata" {
		meta, err := s.interceptors.Metadata(id)
		if err != nil {
			writeAPIError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, map[string]any{"interceptorMetadata": meta})
		return
	}
	if len(parts) >= 3 && parts[1] == "activate" {
		if r.Method != http.MethodPost {
			writeAPIError(w, http.StatusMethodNotAllowed, "activate requires POST")
			return
		}
		var proxyPort int
		_, _ = fmt.Sscan(parts[2], &proxyPort)
		if proxyPort <= 0 {
			writeAPIError(w, http.StatusBadRequest, "proxyPort is required")
			return
		}
		var options map[string]any
		_ = json.NewDecoder(r.Body).Decode(&options)
		result, err := s.interceptors.Activate(id, proxyPort, options)
		if err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, map[string]any{"result": result})
		return
	}
	if len(parts) >= 3 && parts[1] == "deactivate" {
		if r.Method != http.MethodPost {
			writeAPIError(w, http.StatusMethodNotAllowed, "deactivate requires POST")
			return
		}
		var proxyPort int
		_, _ = fmt.Sscan(parts[2], &proxyPort)
		if proxyPort <= 0 {
			writeAPIError(w, http.StatusBadRequest, "proxyPort is required")
			return
		}
		var options map[string]any
		_ = json.NewDecoder(r.Body).Decode(&options)
		result, err := s.interceptors.Deactivate(id, proxyPort, options)
		if err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, map[string]any{"result": result})
		return
	}
	writeAPIError(w, http.StatusNotFound, "not found")
}

// currentUserEmail loads the saved desktop session and returns the user's
// email. If the user is not signed in, it returns an empty string.
func (s *Server) currentUserEmail() string {
	d, _ := s.loadDesktopSession()
	return d.Email
}

func (s *Server) handleConfigBackupPost(w http.ResponseWriter, r *http.Request) {
	email := s.currentUserEmail()
	if email == "" {
		writeAPIError(w, http.StatusUnauthorized, "not signed in")
		return
	}

	var snapshot map[string]any
	if err := json.NewDecoder(r.Body).Decode(&snapshot); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	b := backup.NewUpstash(s.cfg)
	if b == nil {
		writeAPIError(w, http.StatusServiceUnavailable, "cloud backup not configured")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	if err := b.SaveUserConfig(ctx, email, snapshot); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, map[string]any{"ok": true, "savedAt": time.Now().UTC().Format(time.RFC3339)})
}

func (s *Server) handleConfigRestoreGet(w http.ResponseWriter, r *http.Request) {
	email := s.currentUserEmail()
	if email == "" {
		writeAPIError(w, http.StatusUnauthorized, "not signed in")
		return
	}

	b := backup.NewUpstash(s.cfg)
	if b == nil {
		writeAPIError(w, http.StatusServiceUnavailable, "cloud backup not configured")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	data, err := b.RestoreUserConfig(ctx, email)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if data == nil {
		writeAPIError(w, http.StatusNotFound, "no backup found")
		return
	}
	delete(data, "_savedAt")
	writeJSON(w, map[string]any{"ok": true, "payload": data})
}

func (s *Server) handleConfigSyncStatusGet(w http.ResponseWriter, r *http.Request) {
	email := s.currentUserEmail()
	if email == "" {
		writeJSON(w, map[string]any{"lastSyncAt": nil, "lastSyncError": nil})
		return
	}

	b := backup.NewUpstash(s.cfg)
	if b == nil {
		writeJSON(w, map[string]any{"lastSyncAt": nil, "lastSyncError": "cloud backup not configured"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	status, err := b.UserSyncStatus(ctx, email)
	if err != nil {
		writeJSON(w, map[string]any{"lastSyncAt": nil, "lastSyncError": err.Error()})
		return
	}
	writeJSON(w, status)
}
