package rtc

import (
	"sync"

	"github.com/google/uuid"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/events"
)

// Manager provides MockRTC-compatible event emission for WebRTC sessions.
type Manager struct {
	bus     *events.Bus
	mu      sync.Mutex
	running bool
	port    int
}

func NewManager(bus *events.Bus) *Manager {
	return &Manager{bus: bus}
}

func (m *Manager) Start(proxyPort int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.running = true
	m.port = proxyPort
	return nil
}

func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.running = false
	return nil
}

func (m *Manager) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

// NotifyPeerConnected emits mockrtc peer-connected (called from extension hook or SDP intercept).
func (m *Manager) NotifyPeerConnected(sessionID string) string {
	id := uuid.NewString()
	m.bus.PublishRTC("peer-connected", NormalizeEvent("peer-connected", map[string]any{
		"id": id, "peerId": id, "sessionId": sessionID,
	}, sessionID))
	return id
}

func (m *Manager) NotifyPeerDisconnected(id string) {
	m.bus.PublishRTC("peer-disconnected", NormalizeEvent("peer-disconnected", map[string]any{
		"id": id, "peerId": id,
	}, ""))
}

func (m *Manager) Publish(event string, data map[string]any) {
	event = normalizeEventName(event)
	if !ValidEvent(event) {
		return
	}
	m.bus.PublishRTC(event, NormalizeEvent(event, data, ""))
}

// PublishForSession publishes a normalized event tied to an admin session id.
func (m *Manager) PublishForSession(event string, data map[string]any, sessionID string) {
	event = normalizeEventName(event)
	if !ValidEvent(event) {
		return
	}
	m.bus.PublishRTC(event, NormalizeEvent(event, data, sessionID))
}
