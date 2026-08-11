package webextension

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Manager installs the Chromium webextension and writes per-session MockRTC config.
type Manager struct {
	mu          sync.Mutex
	assetsDir   string
	installPath string
	configDir   string
}

func NewManager(assetsDir string) *Manager {
	return &Manager{assetsDir: assetsDir}
}

func (m *Manager) ExtensionPath() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.installPath != "" {
		return m.installPath
	}
	return filepath.Join(m.assetsDir, "overrides", "webextension")
}

func (m *Manager) EnsureInstalled() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.installPath != "" {
		return nil
	}
	src := filepath.Join(m.assetsDir, "overrides", "webextension")
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("webextension assets: %w", err)
	}
	dst := filepath.Join(os.TempDir(), "httptoolkit-webextension")
	if err := copyDir(src, dst); err != nil {
		return err
	}
	cfg := filepath.Join(dst, "config")
	if err := os.MkdirAll(cfg, 0o755); err != nil {
		return err
	}
	m.installPath = dst
	m.configDir = cfg
	return nil
}

type configFile struct {
	MockRTC any `json:"mockRtc"`
}

// UpdateConfig writes extension config for a proxy session.
func (m *Manager) UpdateConfig(proxyPort int, sessionID string, webrtcEnabled bool, adminPort int) error {
	if err := m.EnsureInstalled(); err != nil {
		return err
	}
	var cfg configFile
	if webrtcEnabled {
		cfg.MockRTC = map[string]string{
			"peerId":       "matching-peer",
			"adminBaseUrl": fmt.Sprintf("http://internal.httptoolkit.localhost:%d/session/%s", adminPort, sessionID),
		}
	} else {
		cfg.MockRTC = false
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(m.configPath(proxyPort), data, 0o644)
}

func (m *Manager) ClearConfig(proxyPort int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.configDir == "" {
		return
	}
	_ = os.Remove(m.configPathLocked(proxyPort))
}

func (m *Manager) Cleanup() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.installPath == "" {
		return nil
	}
	err := os.RemoveAll(m.installPath)
	m.installPath = ""
	m.configDir = ""
	return err
}

func (m *Manager) configPath(proxyPort int) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.configPathLocked(proxyPort)
}

func (m *Manager) configPathLocked(proxyPort int) string {
	key := fmt.Sprintf("127_0_0_1.%d", proxyPort)
	return filepath.Join(m.configDir, key)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
