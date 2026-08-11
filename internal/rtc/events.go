package rtc

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

var validEvents = map[string]struct{}{
	"peer-connected":                {},
	"peer-disconnected":             {},
	"external-peer-attached":        {},
	"data-channel-opened":           {},
	"data-channel-message-sent":     {},
	"data-channel-message-received": {},
	"data-channel-closed":           {},
	"media-track-opened":            {},
	"media-track-stats":             {},
	"media-track-closed":            {},
}

// ValidEvent reports whether name is a known MockRTC event.
func ValidEvent(name string) bool {
	_, ok := validEvents[name]
	return ok
}

// NormalizeEvent shapes extension/admin payloads for the UI event bus.
func NormalizeEvent(event string, payload map[string]any, sessionID string) map[string]any {
	out := make(map[string]any, len(payload)+4)
	for k, v := range payload {
		out[k] = v
	}
	if _, ok := out["type"]; ok && event == "" {
		if t, _ := out["type"].(string); t != "" {
			event = t
		}
	}
	delete(out, "type")

	peerID, _ := out["peerId"].(string)
	if peerID == "" {
		peerID, _ = out["id"].(string)
	}
	if peerID == "" {
		peerID = uuid.NewString()
	}
	out["peerId"] = peerID
	out["id"] = peerID

	if sessionID != "" {
		if sid, _ := out["sessionId"].(string); sid == "" {
			out["sessionId"] = sessionID
		}
	}
	if _, ok := out["timestamp"]; !ok {
		out["timestamp"] = time.Now().UnixMilli()
	}
	return out
}

// EventFromPayload resolves event name from admin POST bodies.
func EventFromPayload(payload map[string]any) string {
	if event, _ := payload["event"].(string); event != "" {
		return event
	}
	if event, _ := payload["type"].(string); event != "" {
		return event
	}
	return ""
}

func normalizeEventName(name string) string {
	return strings.TrimSpace(name)
}
