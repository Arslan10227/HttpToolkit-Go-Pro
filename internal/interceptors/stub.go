package interceptors

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
)

type activator interface {
	IsActivable() (bool, error)
	IsActive(proxyPort int) (bool, error)
	Activate(proxyPort int, options map[string]any) (map[string]any, error)
	Deactivate(proxyPort int, options map[string]any) error
	Metadata(kind string) (any, error)
}

type subMetadataProvider interface {
	SubMetadata(subID string) (any, error)
}

type stubInterceptor struct {
	id     string
	cfg    *config.Config
	spki   string
	mu     sync.Mutex
	active map[int]bool
	impl   activator
}

func newStub(cfg *config.Config, id, spki string, certs *cert.Manager) Interceptor {
	s := &stubInterceptor{id: id, cfg: cfg, spki: spki, active: make(map[int]bool)}
	switch id {
	case "system-proxy":
		s.impl = &systemProxy{base: s}
	case "fresh-chrome", "fresh-chromium", "fresh-edge", "fresh-brave", "fresh-opera",
		"fresh-chrome-beta", "fresh-chrome-dev", "fresh-chrome-canary",
		"fresh-chromium-dev", "fresh-edge-beta", "fresh-edge-dev", "fresh-edge-canary":
		s.impl = &chromiumFresh{base: s, browser: browserForID(id), spki: s.spki, certs: certs}
	case "existing-chrome", "existing-chromium", "existing-arc":
		s.impl = &chromiumExisting{base: s, browser: browserForID(id), spki: s.spki, certs: certs}
	case "fresh-firefox", "fresh-firefox-dev", "fresh-firefox-nightly":
		s.impl = &firefoxFresh{base: s, certs: certs}
	case "fresh-terminal", "existing-terminal":
		s.impl = &terminalInterceptor{base: s, fresh: id == "fresh-terminal"}
	case "attach-jvm":
		s.impl = &jvmInterceptor{base: s}
	case "android-adb":
		s.impl = &androidAdb{base: s}
	case "android-frida":
		s.impl = &androidFrida{base: s}
	case "ios-frida":
		s.impl = &iosFrida{base: s}
	case "electron":
		s.impl = &electronInterceptor{base: s}
	case "docker-attach":
		s.impl = &dockerAttach{base: s}
	case "fresh-safari":
		s.impl = &safariFresh{base: s}
	default:
		s.impl = s
	}
	return s
}

func (s *stubInterceptor) ID() string      { return s.id }
func (s *stubInterceptor) Version() string { return "1.0.0" }

func (s *stubInterceptor) setActive(proxyPort int, active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if active {
		s.active[proxyPort] = true
	} else {
		delete(s.active, proxyPort)
	}
}

func (s *stubInterceptor) IsActivable() (bool, error) {
	if s.impl != s {
		return s.impl.IsActivable()
	}
	return false, nil
}

func (s *stubInterceptor) IsActive(proxyPort int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.active[proxyPort], nil
}

func (s *stubInterceptor) Activate(proxyPort int, options map[string]any) (map[string]any, error) {
	if s.impl != s {
		meta, err := s.impl.Activate(proxyPort, options)
		if err != nil {
			return nil, err
		}
		s.mu.Lock()
		s.active[proxyPort] = true
		s.mu.Unlock()
		return meta, nil
	}
	return nil, fmt.Errorf("interceptor not available on this platform")
}

func (s *stubInterceptor) Deactivate(proxyPort int, options map[string]any) error {
	if s.impl != s {
		if err := s.impl.Deactivate(proxyPort, options); err != nil {
			return err
		}
	}
	s.mu.Lock()
	delete(s.active, proxyPort)
	s.mu.Unlock()
	return nil
}

func (s *stubInterceptor) Metadata(kind string) (any, error) {
	if s.impl != s {
		if meta, err := s.impl.Metadata(kind); err == nil {
			return meta, nil
		}
	}
	return map[string]any{"id": s.id, "type": kind}, nil
}

func (s *stubInterceptor) SubMetadata(subID string) (any, error) {
	if s.impl != s {
		if sp, ok := s.impl.(subMetadataProvider); ok {
			return sp.SubMetadata(subID)
		}
	}
	return s.Metadata("detailed")
}

func assetPath(cfg *config.Config, parts ...string) string {
	return filepath.Join(append([]string{cfg.AssetsDir, "overrides"}, parts...)...)
}

func javaAgentPath(cfg *config.Config) string {
	// Allow an explicit override for development, but never default to a local path.
	if overrideJar := os.Getenv("HTK_JAVA_AGENT_JAR"); overrideJar != "" {
		if _, err := os.Stat(overrideJar); err == nil {
			return overrideJar
		}
	}
	p := assetPath(cfg, "java-agent.jar")
	if _, err := os.Stat(p); err == nil {
		return p
	}
	return filepath.Join(cfg.AssetsDir, "java-agent.jar")
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func isWindows() bool { return runtime.GOOS == "windows" }
