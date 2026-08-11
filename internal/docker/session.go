package docker

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/session"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/system"
)

// Session manages Docker interception rule parameters for a proxy port.
type Session struct {
	mu         sync.Mutex
	port       int
	running    bool
	sess       *session.Manager
	tunnel     *TunnelSocks
	network    *NetworkMonitor
	apiProxy   *APIProxy
	assetsDir  string
	tunnelPort int
	certPath   string
}

func NewSession(sess *session.Manager, assetsDir string) *Session {
	return &Session{sess: sess, assetsDir: assetsDir}
}

func (d *Session) RuleParamKey(port int) string {
	return fmt.Sprintf("docker-tunnel-proxy-%d", port)
}

func (d *Session) Start(port int, certPath string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.running {
		return nil
	}
	if !commandExists("docker") {
		return nil
	}
	key := d.RuleParamKey(port)
	d.sess.AddRuleParameterKey(key)
	d.port = port
	d.certPath = certPath
	d.running = true

	// Local SOCKS5 tunnel for docker-routed hostnames (matches Node docker-tunnel-proxy-{port}).
	d.tunnel = NewTunnelSocks(port)
	if err := d.tunnel.Start(); err == nil {
		d.tunnelPort = d.tunnel.Port()
	}
	d.network = NewNetworkMonitor()
	d.network.Start()
	d.sess.SetDockerHostCheck(d.network.HasAlias)

	_ = d.ensureInjectionVolume(certPath)
	d.apiProxy = NewAPIProxy(port, certPath, d.assetsDir)
	_ = d.apiProxy.Start()
	return nil
}

func (d *Session) TunnelPort() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.tunnelPort
}

func (d *Session) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.running {
		return
	}
	if d.tunnel != nil {
		_ = d.tunnel.Stop()
		d.tunnel = nil
		d.tunnelPort = 0
	}
	if d.apiProxy != nil {
		_ = d.apiProxy.Stop()
		d.apiProxy = nil
	}
	_ = DeleteInterceptedContainers(d.port)
	if d.network != nil {
		d.network.Stop()
		d.network = nil
	}
	d.sess.SetDockerHostCheck(nil)
	d.sess.RemoveRuleParameterKey(d.RuleParamKey(d.port))
	d.running = false
	d.port = 0
}

func (d *Session) ensureInjectionVolume(certPath string) error {
	if certPath == "" {
		return nil
	}
	vol := fmt.Sprintf("httptoolkit-injection-%d", d.port)
	_ = system.Command("docker", "volume", "create", vol).Run()
	if d.certPath == "" {
		return nil
	}
	run := system.Command("docker", "run", "--rm", "-v", vol+":/inject",
		"-v", certPath+":/ca.pem:ro", "alpine", "sh", "-c",
		"cp /ca.pem /inject/httptoolkit-ca.pem && chmod 644 /inject/httptoolkit-ca.pem")
	_, _ = run.CombinedOutput()
	return nil
}

func (d *Session) InjectCA(containerID, certPath string) error {
	if containerID == "" || certPath == "" {
		return fmt.Errorf("containerId and cert required")
	}
	cmd := system.Command("docker", "cp", certPath, containerID+":/usr/local/share/ca-certificates/httptoolkit.crt")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("docker cp: %s: %w", out, err)
	}
	update := system.Command("docker", "exec", containerID, "update-ca-certificates")
	if out, err := update.CombinedOutput(); err != nil {
		return fmt.Errorf("update-ca-certificates: %s: %w", out, err)
	}
	return nil
}

func IsAvailable() bool {
	if _, err := exec.LookPath("docker"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	out, err := system.CommandContext(ctx, "docker", "info", "--format", "{{.OSType}}").Output()
	if err != nil {
		return false
	}
	ostype := strings.TrimSpace(string(out))
	return strings.EqualFold(ostype, "linux")
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func ListContainers() ([]map[string]any, error) {
	cmd := system.Command("docker", "ps", "--format", "{{.ID}}\t{{.Names}}\t{{.Image}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var list []map[string]any
	for _, line := range splitLines(string(out)) {
		parts := splitTab(line)
		if len(parts) < 2 {
			continue
		}
		list = append(list, map[string]any{
			"id": parts[0], "name": parts[1], "image": parts[len(parts)-1],
		})
	}
	return list, nil
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if line != "" {
				out = append(out, line)
			}
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

func splitTab(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\t' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}
