package webextension

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigPath(t *testing.T) {
	m := NewManager(t.TempDir())
	m.configDir = filepath.Join(t.TempDir(), "config")
	p := m.configPathLocked(8000)
	if filepath.Base(p) != "127_0_0_1.8000" {
		t.Fatalf("unexpected config path: %s", p)
	}
}

func TestUpdateConfigDisabled(t *testing.T) {
	root := t.TempDir()
	ext := filepath.Join(root, "overrides", "webextension")
	if err := os.MkdirAll(ext, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ext, "manifest.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	m := NewManager(root)
	if err := m.UpdateConfig(8000, "sess-1", false, 45456); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(m.configPath(8000))
	if err != nil {
		t.Fatal(err)
	}
	if !contains(string(data), `"mockRtc":false`) {
		t.Fatalf("expected disabled config, got %s", data)
	}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
