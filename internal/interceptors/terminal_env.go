package interceptors

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	dockersvc "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/docker"
)

func buildTerminalEnv(cfg *config.Config, proxyPort int) []string {
	overridePath := filepath.Join(cfg.AssetsDir, "overrides")
	certPath := filepath.Join(cfg.ConfigDir, "ca.pem")
	proxyURL := fmt.Sprintf("http://127.0.0.1:%d", proxyPort)
	sep := ":"
	if runtime.GOOS == "windows" {
		sep = ";"
	}

	binPath := filepath.Join(overridePath, "path")
	rubyPath := filepath.Join(overridePath, "gems")
	pythonPath := filepath.Join(overridePath, "pythonpath")
	phpPath := filepath.Join(overridePath, "php")
	nodeScript := filepath.ToSlash(filepath.Join(overridePath, "js", "prepend-node.js"))
	javaAgent := javaAgentPath(cfg)

	env := map[string]string{
		"HTTP_PROXY":                           proxyURL,
		"HTTPS_PROXY":                          proxyURL,
		"http_proxy":                           proxyURL,
		"https_proxy":                          proxyURL,
		"WS_PROXY":                             proxyURL,
		"WSS_PROXY":                            proxyURL,
		"GLOBAL_AGENT_HTTP_PROXY":              proxyURL,
		"CGI_HTTP_PROXY":                       proxyURL,
		"npm_config_proxy":                     proxyURL,
		"npm_config_https_proxy":               proxyURL,
		"npm_config_scripts_prepend_node_path": "false",
		"SSL_CERT_FILE":                        certPath,
		"NODE_EXTRA_CA_CERTS":                  certPath,
		"DENO_CERT":                            certPath,
		"PERL_LWP_SSL_CA_FILE":                 certPath,
		"GIT_SSL_CAINFO":                       certPath,
		"CARGO_HTTP_CAINFO":                    certPath,
		"CURL_CA_BUNDLE":                       certPath,
		"AWS_CA_BUNDLE":                        certPath,
		"HTTP_TOOLKIT_ACTIVE":                  "true",
		"HTTP_TOOLKIT_OVERRIDE_PATH":           overridePath,
		"NODE_OPTIONS":                         "--require " + quoteIfNeeded(nodeScript),
		// Keep the option minimal: the agent reads proxy/cert from HTTPS_PROXY and
		// SSL_CERT_FILE, so paths with spaces do not break the JVM option parser.
		"JAVA_TOOL_OPTIONS":                    fmt.Sprintf(`-javaagent:"%s"`, javaAgent),
		"PHP_INI_SCAN_DIR":                     phpPath,
		"DOCKER_BUILDKIT":                      "0",
	}

	if dockersvc.IsAvailable() {
		env["DOCKER_HOST"] = dockersvc.ProxyHostEnv(proxyPort)
	}

	prependPath(env, "PATH", binPath, sep)
	prependPath(env, "RUBYLIB", rubyPath, sep)
	prependPath(env, "PYTHONPATH", pythonPath, sep)

	if existingPHP := os.Getenv("PHP_INI_SCAN_DIR"); existingPHP != "" {
		env["PHP_INI_SCAN_DIR"] = existingPHP + sep + phpPath
	}
	if existingJava := os.Getenv("JAVA_TOOL_OPTIONS"); existingJava != "" {
		env["JAVA_TOOL_OPTIONS"] = existingJava + " " + env["JAVA_TOOL_OPTIONS"]
	}

	out := make([]string, 0, len(env)+16)
	seen := make(map[string]struct{}, len(env))
	for k, v := range env {
		out = append(out, k+"="+v)
		seen[k] = struct{}{}
	}
	for _, kv := range os.Environ() {
		key := kv
		if idx := strings.IndexByte(kv, '='); idx >= 0 {
			key = kv[:idx]
		}
		if _, ok := seen[key]; ok {
			continue
		}
		if skipTerminalInherit(key) {
			continue
		}
		out = append(out, kv)
	}
	return out
}

func prependPath(env map[string]string, key, prefix, sep string) {
	if existing := os.Getenv(key); existing != "" {
		env[key] = prefix + sep + existing
	} else {
		env[key] = prefix
	}
}

func quoteIfNeeded(s string) string {
	if strings.Contains(s, " ") {
		return `"` + s + `"`
	}
	return s
}

func skipTerminalInherit(key string) bool {
	switch strings.ToUpper(key) {
	case "NODE_SKIP_PLATFORM_CHECK", "HTTPTOOLKIT_SERVER_BINPATH", "HTK_DESKTOP_EXE",
		"HTK_DESKTOP_RESOURCES", "NO_AT_BRIDGE", "ORIGINAL_XDG_CURRENT_DESKTOP",
		"GDK_BACKEND", "CHROME_DESKTOP":
		return true
	default:
		return false
	}
}
