package rtc

import "testing"

func TestLooksLikeSDP(t *testing.T) {
	sdp := []byte("v=0\r\no=- 0 0 IN IP4 127.0.0.1\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\n")
	if !LooksLikeSDP(sdp, "") {
		t.Fatal("expected SDP body to match")
	}
	if LooksLikeSDP([]byte("hello world"), "") {
		t.Fatal("expected non-SDP body to be rejected")
	}
	if !LooksLikeSDP([]byte("not sdp"), "application/sdp") {
		t.Fatal("expected content-type application/sdp to match")
	}
}

func TestNormalizeEvent(t *testing.T) {
	out := NormalizeEvent("peer-connected", map[string]any{
		"type": "peer-connected",
	}, "sess-1")
	if out["id"] == "" || out["peerId"] == "" {
		t.Fatalf("expected peer id fields, got %#v", out)
	}
	if out["sessionId"] != "sess-1" {
		t.Fatalf("expected sessionId, got %#v", out)
	}
}

func TestEventFromPayload(t *testing.T) {
	if got := EventFromPayload(map[string]any{"event": "data-channel-opened"}); got != "data-channel-opened" {
		t.Fatalf("unexpected event: %s", got)
	}
	if got := EventFromPayload(map[string]any{"type": "peer-disconnected"}); got != "peer-disconnected" {
		t.Fatalf("unexpected event: %s", got)
	}
}
