package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/backup"
)

// Settings holds user-level preferences that are persisted on disk in the
// config directory, independent of cloud account settings.
type Settings struct {
	ConfirmBeforeClose bool `json:"confirmBeforeClose"`
}

// Default returns a Settings value with the standard defaults.
func Default() Settings {
	return Settings{
		ConfirmBeforeClose: false,
	}
}

// Manager loads and saves settings from a JSON file in the config directory.
type Manager struct {
	path   string
	mu     sync.RWMutex
	s      Settings
	backup backup.Backuper
}

// NewManager creates a settings manager for the given config directory.
func NewManager(configDir string) *Manager {
	return &Manager{
		path: filepath.Join(configDir, "settings.json"),
		s:    Default(),
	}
}

// SetBackup attaches a cloud backup provider. When set, Save will mirror the
// local settings to the cloud (best effort, async), and Load may restore from
// the cloud if the local file is missing.
func (m *Manager) SetBackup(b backup.Backuper) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.backup = b
}

// Load reads settings from disk, falling back to a cloud backup or defaults.
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			// No local file yet — try to pull a backup from Firebase.
			if m.backup != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				out := make(map[string]any)
				if err := m.backup.RestoreSettings(ctx, out); err == nil {
					if raw, err := json.Marshal(out); err == nil {
						var s Settings
						if err := json.Unmarshal(raw, &s); err == nil {
							m.s = s
						}
					}
				}
			}
			return nil
		}
		return fmt.Errorf("settings read: %w", err)
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("settings parse: %w", err)
	}
	m.s = s
	return nil
}

// Save persists the current settings to disk and, if a backup provider is
// attached, mirrors the data to the cloud asynchronously.
func (m *Manager) Save(s Settings) error {
	m.mu.Lock()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		m.mu.Unlock()
		return err
	}
	if err := os.MkdirAll(filepath.Dir(m.path), 0o700); err != nil {
		m.mu.Unlock()
		return err
	}
	if err := os.WriteFile(m.path, data, 0o600); err != nil {
		m.mu.Unlock()
		return err
	}
	m.s = s
	b := m.backup
	m.mu.Unlock()

	if b != nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			_ = b.SaveSettings(ctx, s)
		}()
	}
	return nil
}

// Get returns the current settings.
func (m *Manager) Get() Settings {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.s
}

// ConfirmBeforeClose reports whether the app should show a confirmation dialog
// before closing.
func (m *Manager) ConfirmBeforeClose() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.s.ConfirmBeforeClose
}
