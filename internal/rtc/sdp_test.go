package rtc

import (
	"encoding/json"
	"testing"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/events"
)

func TestNotifySDPOffer(t *testing.T) {
	bus := events.NewBus()
	ch := bus.Subscribe()
	defer bus.Unsubscribe(ch)

	m := NewManager(bus)
	_ = m.Start(8000)
	sdp := "v=0\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nm=audio 9 UDP/TLS/RTP/SAVPF 111\r\n"
	peerID := m.NotifySDPOffer("req-1", sdp)
	if peerID == "" {
		t.Fatal("expected peer id")
	}

	var names []string
	for i := 0; i < 3; i++ {
		msg := <-ch
		var envelope struct {
			Event string `json:"event"`
		}
		_ = json.Unmarshal(msg, &envelope)
		names = append(names, envelope.Event)
	}
	if len(names) < 3 {
		t.Fatalf("expected at least 3 rtc events, got %v", names)
	}
}
