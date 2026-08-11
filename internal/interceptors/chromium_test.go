package interceptors

import (
	"strings"
	"testing"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
)

func TestChromiumArgs(t *testing.T) {
	cfg := &config.Config{
		ConfigDir: t.TempDir(),
		AssetsDir: t.TempDir(),
	}
	base := &stubInterceptor{id: "fresh-chrome", cfg: cfg, spki: "test-spki"}

	amiusing := &amiusingServer{url: "http://127.0.0.1:54321/amiusing"}
	args := chromiumArgs(base, "test-spki", 8000, "C:\\tmp\\profile", amiusing, "chrome")

	want := []string{
		"--user-data-dir=C:\\tmp\\profile",
		"--proxy-server=https://127.0.0.1:8000",
		"--proxy-bypass-list=",
		"--test-type",
		"--ignore-certificate-errors-spki-list=test-spki",
		"--disable-quic",
		"--disable-features=",
		"--disable-background-networking",
		"--no-default-browser-check",
		"--no-first-run",
		"--disable-popup-blocking",
		"--component-updater=url-source=http://disabled-chromium-update.localhost:0",
		"--check-for-update-interval=31536000",
		"http://127.0.0.1:54321/amiusing",
	}

	for _, w := range want {
		found := false
		for _, a := range args {
			if strings.HasPrefix(a, w) || a == w {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing expected arg prefix %q in %v", w, args)
		}
	}

	// Brave should not get the component-updater flags.
	braveArgs := chromiumArgs(base, "test-spki", 8000, "", nil, "brave")
	for _, a := range braveArgs {
		if strings.Contains(a, "component-updater") || strings.Contains(a, "check-for-update-interval") {
			t.Errorf("Brave should not get component-updater flags, got %q", a)
		}
	}

	// Verify the bypass list contains the internal patterns but not the amiusing
	// server (we want the check page to go through the proxy).
	var bypassList string
	for _, a := range args {
		if strings.HasPrefix(a, "--proxy-bypass-list=") {
			bypassList = a
			break
		}
	}
	if bypassList == "" {
		t.Fatalf("missing --proxy-bypass-list")
	}
	for _, pattern := range []string{
		"<-loopback>",
		"internal.httptoolkit.localhost",
		"tauri.localhost",
		"localhost",
		"127.0.0.1",
		"vercel.app",
	} {
		if !strings.Contains(bypassList, pattern) {
			t.Errorf("proxy bypass list missing %q: %s", pattern, bypassList)
		}
	}
	if strings.Contains(bypassList, "127.0.0.1:54321") {
		t.Errorf("amiusing server should not be in proxy bypass list: %s", bypassList)
	}
}
