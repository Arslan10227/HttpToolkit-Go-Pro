package interceptors

import (
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
)

func buildAll(cfg *config.Config, spki string, certs *cert.Manager) []Interceptor {
	ids := []string{
		"fresh-chrome", "existing-chrome", "fresh-chrome-beta", "fresh-chrome-dev", "fresh-chrome-canary",
		"fresh-chromium", "existing-chromium", "fresh-chromium-dev",
		"fresh-edge", "fresh-edge-beta", "fresh-edge-dev", "fresh-edge-canary",
		"fresh-brave", "fresh-opera", "existing-arc",
		"fresh-firefox", "fresh-firefox-dev", "fresh-firefox-nightly",
		"fresh-safari", "system-proxy",
		"fresh-terminal", "existing-terminal",
		"electron", "attach-jvm",
		"android-adb", "android-frida", "ios-frida",
		"docker-attach",
	}
	out := make([]Interceptor, 0, len(ids))
	for _, id := range ids {
		out = append(out, newStub(cfg, id, spki, certs))
	}
	return out
}
