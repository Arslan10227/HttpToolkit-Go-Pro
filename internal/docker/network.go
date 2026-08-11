package docker

import (
	"bufio"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/system"
)

// NetworkMonitor tracks Docker container network aliases for DNS/tunnel routing.
type NetworkMonitor struct {
	mu      sync.RWMutex
	aliases map[string]struct{}
	stop    chan struct{}
	once    sync.Once
}

func NewNetworkMonitor() *NetworkMonitor {
	return &NetworkMonitor{
		aliases: make(map[string]struct{}),
		stop:    make(chan struct{}),
	}
}

func (m *NetworkMonitor) Start() {
	m.refresh()
	go m.poll()
	go m.watchEvents()
}

func (m *NetworkMonitor) Stop() {
	m.once.Do(func() { close(m.stop) })
}

func (m *NetworkMonitor) HasAlias(host string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.aliases[host]
	return ok
}

func (m *NetworkMonitor) addAlias(name string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}
	m.mu.Lock()
	m.aliases[name] = struct{}{}
	m.mu.Unlock()
}

func (m *NetworkMonitor) removeAlias(name string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}
	m.mu.Lock()
	delete(m.aliases, name)
	m.mu.Unlock()
}

func (m *NetworkMonitor) poll() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-m.stop:
			return
		case <-ticker.C:
			m.refresh()
		}
	}
}

func (m *NetworkMonitor) watchEvents() {
	cmd := system.Command("docker", "events", "--filter", "type=container", "--format", "{{json .}}")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	if err := cmd.Start(); err != nil {
		return
	}
	go func() {
		<-m.stop
		_ = cmd.Process.Kill()
	}()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		select {
		case <-m.stop:
			return
		default:
		}
		var evt struct {
			Action string `json:"Action"`
			Actor  struct {
				Attributes map[string]string `json:"Attributes"`
			} `json:"Actor"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &evt); err != nil {
			continue
		}
		name := evt.Actor.Attributes["name"]
		switch evt.Action {
		case "start", "rename":
			m.addAlias(name)
		case "destroy", "die", "remove":
			m.removeAlias(name)
		}
	}
}

func (m *NetworkMonitor) refresh() {
	out, err := system.Command("docker", "ps", "--format", "{{.Names}}").Output()
	if err != nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for k := range m.aliases {
		delete(m.aliases, k)
	}
	for _, name := range strings.Split(string(out), "\n") {
		name = strings.TrimSpace(name)
		if name != "" {
			m.aliases[name] = struct{}{}
		}
	}
}
