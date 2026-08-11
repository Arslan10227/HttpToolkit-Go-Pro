package native

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pion/webrtc/v4"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/events"
)

// newTestBrowserPeer creates a plain pion PeerConnection standing in for a
// real browser tab, so we can validate the native mock Handler's GraphQL
// wire protocol end-to-end without needing an actual browser.
func newTestBrowserPeer(t *testing.T) *webrtc.PeerConnection {
	t.Helper()
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatalf("browser peer connection: %v", err)
	}
	t.Cleanup(func() { _ = pc.Close() })
	return pc
}

func waitConnected(t *testing.T, pc *webrtc.PeerConnection) {
	t.Helper()
	done := make(chan struct{})
	pc.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		if s == webrtc.PeerConnectionStateConnected {
			select {
			case <-done:
			default:
				close(done)
			}
		}
	})
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatalf("peer connection did not reach connected state in time (state=%s)", pc.ConnectionState())
	}
}

// TestAnswerOffer_BrowserInitiated simulates the common case: a page calls
// `new RTCPeerConnection()` and creates an offer, the extension forwards it
// to answerOffer(peerId, offer), and the native handler answers using a real
// pion PeerConnection — the two sides should be able to fully connect and
// exchange a data-channel message.
func TestAnswerOffer_BrowserInitiated(t *testing.T) {
	bus := events.NewBus()
	sub := bus.Subscribe()
	defer bus.Unsubscribe(sub)

	h := NewHandler(bus)
	defer h.Close()

	browser := newTestBrowserPeer(t)
	dc, err := browser.CreateDataChannel("test-channel", nil)
	if err != nil {
		t.Fatalf("create data channel: %v", err)
	}

	received := make(chan string, 1)
	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		received <- string(msg.Data)
	})
	dcOpen := make(chan struct{})
	dc.OnOpen(func() { close(dcOpen) })

	offer, err := browser.CreateOffer(nil)
	if err != nil {
		t.Fatalf("create offer: %v", err)
	}
	gatherComplete := webrtc.GatheringCompletePromise(browser)
	if err := browser.SetLocalDescription(offer); err != nil {
		t.Fatalf("set local description: %v", err)
	}
	<-gatherComplete

	local := browser.LocalDescription()
	sessionID, answer, err := h.AnswerOffer("matching-peer", SessionDescription{
		Type: local.Type.String(), SDP: local.SDP,
	})
	if err != nil {
		t.Fatalf("AnswerOffer: %v", err)
	}
	if sessionID == "" {
		t.Fatal("expected non-empty session id")
	}
	if answer.SDP == "" || answer.Type != "answer" {
		t.Fatalf("unexpected answer: %+v", answer)
	}

	if err := browser.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer, SDP: answer.SDP,
	}); err != nil {
		t.Fatalf("browser set remote description: %v", err)
	}

	waitConnected(t, browser)

	// Rule: echo anything received on "test-channel" back to the sender.
	h.SetRules([]RTCRule{
		{
			Matchers: []json.RawMessage{json.RawMessage(`{"type":"has-data-channel","label":"test-channel"}`)},
			Steps:    []json.RawMessage{json.RawMessage(`{"type":"echo"}`)},
		},
	})

	select {
	case <-dcOpen:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for data channel to open")
	}

	if err := dc.SendText("ping"); err != nil {
		t.Fatalf("send: %v", err)
	}

	select {
	case msg := <-received:
		if msg != "ping" {
			t.Fatalf("expected echoed 'ping', got %q", msg)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for echoed message")
	}

	// The mock peer's session should have recorded the incoming message.
	deadline := time.Now().Add(2 * time.Second)
	var seen []any
	for time.Now().Before(deadline) {
		seen = h.GetSeenMessages("matching-peer", "test-channel")
		if len(seen) > 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if len(seen) == 0 {
		t.Fatal("expected at least one seen message")
	}
	if seen[0] != "ping" {
		t.Fatalf("expected seen message 'ping', got %v", seen[0])
	}
}

// TestCreateOffer_MockInitiated covers the reverse flow: the mock peer
// generates an offer (session.createOffer), the browser answers, and the
// mock peer completes with session.completeOffer(answer).
func TestCreateOffer_MockInitiated(t *testing.T) {
	bus := events.NewBus()
	h := NewHandler(bus)
	defer h.Close()

	browser := newTestBrowserPeer(t)
	browser.OnDataChannel(func(dc *webrtc.DataChannel) {})

	sessionID, offer, err := h.CreateOffer("matching-peer")
	if err != nil {
		t.Fatalf("CreateOffer: %v", err)
	}
	if offer.SDP == "" || offer.Type != "offer" {
		t.Fatalf("unexpected offer: %+v", offer)
	}

	if err := browser.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer, SDP: offer.SDP,
	}); err != nil {
		t.Fatalf("browser set remote description: %v", err)
	}
	answer, err := browser.CreateAnswer(nil)
	if err != nil {
		t.Fatalf("browser create answer: %v", err)
	}
	gatherComplete := webrtc.GatheringCompletePromise(browser)
	if err := browser.SetLocalDescription(answer); err != nil {
		t.Fatalf("browser set local description: %v", err)
	}
	<-gatherComplete

	local := browser.LocalDescription()
	if err := h.CompleteOffer(sessionID, SessionDescription{Type: local.Type.String(), SDP: local.SDP}); err != nil {
		t.Fatalf("CompleteOffer: %v", err)
	}

	waitConnected(t, browser)
}
