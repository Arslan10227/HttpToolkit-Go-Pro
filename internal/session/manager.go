package session

import (
	"fmt"
	"strings"
	"sync"
)

// Manager tracks active proxy session state exposed via GET /config.
type Manager struct {
	mu                sync.RWMutex
	proxyPort         int
	dnsServers        []string
	ruleParameterKeys []string
	http2Enabled      bool
	webrtcEnabled     bool
	socksPort         int
	dockerTunnelPorts map[int]int // proxyPort -> local socks tunnel port
	dockerHostCheck   func(string) bool
}

func NewManager() *Manager {
	return &Manager{
		dnsServers:        []string{},
		ruleParameterKeys: []string{},
		dockerTunnelPorts: make(map[int]int),
	}
}

func (m *Manager) SetActive(port int, dns []string, ruleKeys []string, http2, webrtc bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.proxyPort = port
	if dns != nil {
		m.dnsServers = append([]string(nil), dns...)
	} else {
		m.dnsServers = []string{}
	}
	if ruleKeys != nil {
		m.ruleParameterKeys = append([]string(nil), ruleKeys...)
	} else {
		m.ruleParameterKeys = []string{}
	}
	m.http2Enabled = http2
	m.webrtcEnabled = webrtc
}

func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.proxyPort = 0
	m.dnsServers = []string{}
	m.ruleParameterKeys = []string{}
	m.http2Enabled = false
	m.webrtcEnabled = false
	m.socksPort = 0
	m.dockerTunnelPorts = make(map[int]int)
	m.dockerHostCheck = nil
}

func (m *Manager) SetSocksPort(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.socksPort = port
}

func (m *Manager) SocksPort() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.socksPort
}

func (m *Manager) SetDockerTunnelPort(proxyPort, tunnelPort int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dockerTunnelPorts[proxyPort] = tunnelPort
}

func (m *Manager) DockerTunnelPort(proxyPort int) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.dockerTunnelPorts[proxyPort]
}

func (m *Manager) SetDockerHostCheck(fn func(string) bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dockerHostCheck = fn
}

func (m *Manager) ResolveRuleParamProxy(paramKey, hostname string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if strings.HasPrefix(paramKey, "docker-tunnel-proxy-") {
		if m.dockerHostCheck != nil && !m.dockerHostCheck(hostname) {
			return ""
		}
		if port, ok := m.dockerTunnelPorts[m.proxyPort]; ok && port > 0 {
			return fmt.Sprintf("socks5://127.0.0.1:%d", port)
		}
	}
	return ""
}

func (m *Manager) ProxyPort() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.proxyPort
}

func (m *Manager) DNSServers() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]string(nil), m.dnsServers...)
}

func (m *Manager) RuleParameterKeys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]string(nil), m.ruleParameterKeys...)
}

func (m *Manager) HTTP2Enabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.http2Enabled
}

func (m *Manager) WebRTCEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.webrtcEnabled
}

func (m *Manager) AddRuleParameterKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, k := range m.ruleParameterKeys {
		if k == key {
			return
		}
	}
	m.ruleParameterKeys = append(m.ruleParameterKeys, key)
}

func (m *Manager) RemoveRuleParameterKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := m.ruleParameterKeys[:0]
	for _, k := range m.ruleParameterKeys {
		if k != key {
			out = append(out, k)
		}
	}
	m.ruleParameterKeys = out
}
