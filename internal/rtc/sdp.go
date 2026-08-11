package rtc

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// NotifySDPOffer emits MockRTC-style events from an intercepted SDP offer body.
func (m *Manager) NotifySDPOffer(httpSessionID, sdp string) string {
	peerID := m.NotifyPeerConnected(httpSessionID)
	now := time.Now().UnixMilli()
	for _, line := range strings.Split(sdp, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "m=application"):
			m.Publish("data-channel-opened", map[string]any{
				"id": uuid.NewString(), "peerId": peerID, "sessionId": httpSessionID,
				"timestamp": now, "label": "datachannel",
			})
		case strings.HasPrefix(line, "m=audio"):
			m.Publish("media-track-opened", map[string]any{
				"id": uuid.NewString(), "peerId": peerID, "sessionId": httpSessionID,
				"timestamp": now, "kind": "audio",
			})
		case strings.HasPrefix(line, "m=video"):
			m.Publish("media-track-opened", map[string]any{
				"id": uuid.NewString(), "peerId": peerID, "sessionId": httpSessionID,
				"timestamp": now, "kind": "video",
			})
		}
	}
	return peerID
}
