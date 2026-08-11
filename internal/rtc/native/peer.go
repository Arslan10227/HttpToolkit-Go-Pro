// Package native implements an in-process, Go-native WebRTC mocking handler
// using github.com/pion/webrtc/v4.
//
// The bundled Chromium webextension (assets/overrides/webextension) talks to
// the MockRTC admin session using a small GraphQL protocol
// (createOffer/createExternalOffer/answerOffer/answerExternalOffer/
// completeOffer/getSeenMessages — see MockRTCRemotePeer / RemoteSessionApi in
// background.js). This package implements a Go-native peer that answers that
// exact wire protocol (see graphql.go) using real pion PeerConnections.
package native

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pion/webrtc/v4"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/events"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/rtc"
)

// SessionDescription mirrors mockrtc's GraphQL SessionDescription{type, sdp}.
type SessionDescription struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

// seenMessage records one data-channel message for later retrieval via
// getSeenMessages (used by page JS to assert on traffic in tests).
type seenMessage struct {
	channel string
	value   []byte
	isText  bool
}

// mockSession wraps one pion PeerConnection created on behalf of the mock
// peer (either answering a browser-initiated offer, or offering to a
// browser that then completes it).
type mockSession struct {
	id       string
	peerID   string
	pc       *webrtc.PeerConnection
	mu       sync.Mutex
	messages []seenMessage
}

// Handler manages all mock WebRTC sessions for one proxy/admin session. It
// implements the same event surface (contracts/mockrtc-events.json) the
// rtc.Manager already publishes, so no UI changes are needed.
type Handler struct {
	bus   *events.Bus
	rules *RuleEngine

	mu       sync.Mutex
	sessions map[string]*mockSession // sessionID -> session
}

// NewHandler creates a native WebRTC mock handler publishing events onto bus.
func NewHandler(bus *events.Bus) *Handler {
	return &Handler{
		bus:      bus,
		rules:    NewRuleEngine(),
		sessions: make(map[string]*mockSession),
	}
}

// SetRules updates the data-channel matcher/step rules (from PUT /rules/rtc).
func (h *Handler) SetRules(rules []RTCRule) {
	h.rules.SetRules(rules)
}

func peerConnectionConfig() webrtc.Configuration {
	// No STUN/TURN servers by default: mock peers are answered/offered
	// in-process, so ICE only ever needs to negotiate local/host candidates
	// with the real remote browser peer.
	return webrtc.Configuration{}
}

func (h *Handler) newSession(peerID string) (*mockSession, error) {
	pc, err := webrtc.NewPeerConnection(peerConnectionConfig())
	if err != nil {
		return nil, fmt.Errorf("create peer connection: %w", err)
	}
	sess := &mockSession{id: uuid.NewString(), peerID: peerID, pc: pc}

	h.mu.Lock()
	h.sessions[sess.id] = sess
	h.mu.Unlock()

	h.wireEvents(sess)
	return sess, nil
}

// wireEvents attaches data-channel/connection-state handlers that publish
// mockrtc-compatible events and apply rule-driven behaviour.
func (h *Handler) wireEvents(sess *mockSession) {
	pc := sess.pc

	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		switch state {
		case webrtc.PeerConnectionStateConnected:
			h.bus.PublishRTC("peer-connected", rtc.NormalizeEvent("peer-connected",
				map[string]any{"peerId": sess.peerID, "id": sess.id}, sess.id))
		case webrtc.PeerConnectionStateClosed, webrtc.PeerConnectionStateFailed, webrtc.PeerConnectionStateDisconnected:
			h.bus.PublishRTC("peer-disconnected", rtc.NormalizeEvent("peer-disconnected",
				map[string]any{"peerId": sess.peerID, "id": sess.id}, sess.id))
		}
	})

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		h.handleDataChannel(sess, dc)
	})

	pc.OnTrack(func(track *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		trackID := uuid.NewString()
		h.bus.PublishRTC("media-track-opened", rtc.NormalizeEvent("media-track-opened",
			map[string]any{"peerId": sess.peerID, "id": trackID, "kind": track.Kind().String()}, sess.id))
	})
}

func (h *Handler) handleDataChannel(sess *mockSession, dc *webrtc.DataChannel) {
	channelID := uuid.NewString()
	dc.OnOpen(func() {
		h.bus.PublishRTC("data-channel-opened", rtc.NormalizeEvent("data-channel-opened",
			map[string]any{"peerId": sess.peerID, "id": channelID, "label": dc.Label()}, sess.id))

		// Rule-driven "send" actions fire as soon as the channel opens.
		for _, action := range h.rules.MatchOpen(dc.Label()) {
			applyDataChannelAction(dc, action)
		}
	})

	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		sess.mu.Lock()
		sess.messages = append(sess.messages, seenMessage{channel: dc.Label(), value: msg.Data, isText: msg.IsString})
		sess.mu.Unlock()

		h.bus.PublishRTC("data-channel-message-received", rtc.NormalizeEvent("data-channel-message-received",
			map[string]any{
				"peerId":  sess.peerID,
				"id":      channelID,
				"channel": dc.Label(),
				"content": msg.Data,
			}, sess.id))

		for _, action := range h.rules.MatchMessage(dc.Label(), msg.Data) {
			applyDataChannelAction(dc, action)
			if action.Kind == ActionEcho {
				h.bus.PublishRTC("data-channel-message-sent", rtc.NormalizeEvent("data-channel-message-sent",
					map[string]any{"peerId": sess.peerID, "id": channelID, "channel": dc.Label(), "content": msg.Data}, sess.id))
			}
		}
	})

	dc.OnClose(func() {
		h.bus.PublishRTC("data-channel-closed", rtc.NormalizeEvent("data-channel-closed",
			map[string]any{"peerId": sess.peerID, "id": channelID, "label": dc.Label()}, sess.id))
	})
}

func applyDataChannelAction(dc *webrtc.DataChannel, action DataChannelAction) {
	switch action.Kind {
	case ActionEcho:
		// Handled by caller (re-sends the same message that was received).
		if action.LastMessage != nil {
			_ = dc.Send(action.LastMessage)
		}
	case ActionSend:
		if action.Text != "" {
			_ = dc.SendText(action.Text)
		} else if len(action.Binary) > 0 {
			_ = dc.Send(action.Binary)
		}
	case ActionClose:
		_ = dc.Close()
	}
}

// AnswerOffer implements the peer-level `answerOffer(peerId, offer, options)`
// GraphQL mutation: the browser is initiating a connection to the mock peer,
// so we take its offer and produce a matching answer.
func (h *Handler) AnswerOffer(peerID string, offer SessionDescription) (sessionID string, answer SessionDescription, err error) {
	sess, err := h.newSession(peerID)
	if err != nil {
		return "", SessionDescription{}, err
	}
	if err := sess.pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.NewSDPType(offer.Type), SDP: offer.SDP,
	}); err != nil {
		return "", SessionDescription{}, fmt.Errorf("set remote description: %w", err)
	}
	ans, err := sess.pc.CreateAnswer(nil)
	if err != nil {
		return "", SessionDescription{}, fmt.Errorf("create answer: %w", err)
	}
	gatherComplete := webrtc.GatheringCompletePromise(sess.pc)
	if err := sess.pc.SetLocalDescription(ans); err != nil {
		return "", SessionDescription{}, fmt.Errorf("set local description: %w", err)
	}
	waitForGather(gatherComplete)

	local := sess.pc.LocalDescription()
	return sess.id, SessionDescription{Type: local.Type.String(), SDP: local.SDP}, nil
}

// CreateOffer implements `createOffer(peerId, options)`: the mock peer
// initiates a connection towards the browser (e.g. for session.createOffer
// followed later by session.completeOffer(answer)).
func (h *Handler) CreateOffer(peerID string) (sessionID string, offer SessionDescription, err error) {
	sess, err := h.newSession(peerID)
	if err != nil {
		return "", SessionDescription{}, err
	}
	// A data channel is required for an SDP offer to include the
	// application m-line that data-channel-only mocks need.
	if _, err := sess.pc.CreateDataChannel("mockrtc", nil); err != nil {
		return "", SessionDescription{}, fmt.Errorf("create data channel: %w", err)
	}
	off, err := sess.pc.CreateOffer(nil)
	if err != nil {
		return "", SessionDescription{}, fmt.Errorf("create offer: %w", err)
	}
	gatherComplete := webrtc.GatheringCompletePromise(sess.pc)
	if err := sess.pc.SetLocalDescription(off); err != nil {
		return "", SessionDescription{}, fmt.Errorf("set local description: %w", err)
	}
	waitForGather(gatherComplete)

	local := sess.pc.LocalDescription()
	return sess.id, SessionDescription{Type: local.Type.String(), SDP: local.SDP}, nil
}

// CompleteOffer implements `completeOffer(peerId, sessionId, answer)`: the
// browser has answered an offer the mock peer generated via CreateOffer.
func (h *Handler) CompleteOffer(sessionID string, answer SessionDescription) error {
	sess, err := h.session(sessionID)
	if err != nil {
		return err
	}
	return sess.pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.NewSDPType(answer.Type), SDP: answer.SDP,
	})
}

// GetSeenMessages implements `getSeenMessages(peerId, channelName)`.
// Returns entries shaped like mockrtc: plain strings for text messages,
// {"type":"buffer","value":"<base64>"} for binary ones.
func (h *Handler) GetSeenMessages(peerID string, channelName string) []any {
	h.mu.Lock()
	sessions := make([]*mockSession, 0, len(h.sessions))
	for _, s := range h.sessions {
		if s.peerID == peerID {
			sessions = append(sessions, s)
		}
	}
	h.mu.Unlock()

	var out []any
	for _, sess := range sessions {
		sess.mu.Lock()
		for _, m := range sess.messages {
			if channelName != "" && m.channel != channelName {
				continue
			}
			if m.isText {
				out = append(out, string(m.value))
			} else {
				out = append(out, map[string]any{"type": "buffer", "value": m.value})
			}
		}
		sess.mu.Unlock()
	}
	return out
}

func (h *Handler) session(id string) (*mockSession, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	sess, ok := h.sessions[id]
	if !ok {
		return nil, fmt.Errorf("unknown mockrtc session %q", id)
	}
	return sess, nil
}

// Close tears down every open peer connection (called on session/stop).
func (h *Handler) Close() {
	h.mu.Lock()
	sessions := make([]*mockSession, 0, len(h.sessions))
	for _, s := range h.sessions {
		sessions = append(sessions, s)
	}
	h.sessions = make(map[string]*mockSession)
	h.mu.Unlock()

	for _, s := range sessions {
		_ = s.pc.Close()
	}
}

// waitForGather blocks (briefly) for ICE gathering to complete so the
// returned SDP contains resolved candidates, bounded so a slow/broken
// network stack can never hang a GraphQL request indefinitely.
func waitForGather(done <-chan struct{}) {
	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}
}
