package interceptors

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
)

func TestNSSPlatform(t *testing.T) {
	switch runtime.GOOS {
	case "windows":
		if nssPlatform() != "win32" {
			t.Fatalf("unexpected platform: %s", nssPlatform())
		}
	case "darwin":
		if nssPlatform() != "darwin" {
			t.Fatalf("unexpected platform: %s", nssPlatform())
		}
	default:
		if nssPlatform() != "linux" {
			t.Fatalf("unexpected platform: %s", nssPlatform())
		}
	}
}

func TestFirefoxProfilePath(t *testing.T) {
	cfg := &config.Config{ConfigDir: t.TempDir()}
	p := firefoxProfilePath(cfg, "")
	if filepath.Base(p) != "firefox-profile" {
		t.Fatalf("unexpected profile path: %s", p)
	}
}

func TestWriteFirefoxUserJS(t *testing.T) {
	dir := t.TempDir()
	if err := writeFirefoxUserJS(dir, 8000); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "user.js"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, "network.proxy.http_port") || !strings.Contains(text, "8000") {
		t.Fatalf("expected proxy port in user.js: %s", text)
	}
}

func TestBundledCertutilPath(t *testing.T) {
	assets := filepath.Join("..", "..", "assets")
	if _, err := os.Stat(nssAssetDir(assets)); err != nil {
		t.Skip("bundled nss assets not present")
	}
	bin := filepath.Join(nssAssetDir(assets), certutilBinName())
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("missing bundled certutil: %s", bin)
	}
}
