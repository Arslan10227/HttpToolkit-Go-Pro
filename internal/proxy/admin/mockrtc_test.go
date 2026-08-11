package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/events"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/rtc"
)

func TestMockRTCEventRelay(t *testing.T) {
	bus := events.NewBus()
	ch := bus.Subscribe()
	defer bus.Unsubscribe(ch)

	rtcMgr := rtc.NewManager(bus)
	_ = rtcMgr.Start(8000)

	s := &Server{
		cfg:      &config.Config{AdminPort: 45456},
		rtc:      rtcMgr,
		events:   bus,
		webrtcOn: true,
	}

	req := httptest.NewRequest(http.MethodPost, "/session/test-session",
		strings.NewReader(`{"type":"data-channel-opened","label":"chat"}`))
	rec := httptest.NewRecorder()
	s.handleMockRTCSession(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	msg := <-ch
	var envelope struct {
		Plugin string         `json:"plugin"`
		Event  string         `json:"event"`
		Data   map[string]any `json:"data"`
	}
	if err := json.Unmarshal(msg, &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.Plugin != "webrtc" || envelope.Event != "data-channel-opened" {
		t.Fatalf("unexpected envelope: %#v", envelope)
	}
	if envelope.Data["sessionId"] != "test-session" {
		t.Fatalf("expected sessionId, got %#v", envelope.Data)
	}
	if envelope.Data["id"] == "" {
		t.Fatal("expected id on rtc event")
	}
}

func TestMockRTCGetSession(t *testing.T) {
	rtcMgr := rtc.NewManager(events.NewBus())
	_ = rtcMgr.Start(8000)
	s := &Server{
		cfg:      &config.Config{AdminPort: 45456},
		rtc:      rtcMgr,
		webrtcOn: true,
	}
	req := httptest.NewRequest(http.MethodGet, "/session/abc-123", nil)
	rec := httptest.NewRecorder()
	s.handleMockRTCSession(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["sessionId"] != "abc-123" || body["running"] != true {
		t.Fatalf("unexpected body: %#v", body)
	}
}
