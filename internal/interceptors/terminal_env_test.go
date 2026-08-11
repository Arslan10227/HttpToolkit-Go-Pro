package interceptors

import (
	"strings"
	"testing"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
)

func TestBuildTerminalEnv(t *testing.T) {
	cfg := &config.Config{
		ConfigDir: t.TempDir(),
		AssetsDir: t.TempDir(),
	}
	env := buildTerminalEnv(cfg, 8000)
	text := strings.Join(env, "\n")
	for _, key := range []string{"HTTP_PROXY=", "NODE_OPTIONS=", "JAVA_TOOL_OPTIONS=", "HTTP_TOOLKIT_ACTIVE=true"} {
		if !strings.Contains(text, key) {
			t.Fatalf("missing %s in env", key)
		}
	}
}
