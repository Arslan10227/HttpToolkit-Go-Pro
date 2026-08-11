package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/mitm"
)

func TestRESTBreakpointResume(t *testing.T) {
	bpMgr := mitm.NewBreakpointManager()
	cfg := &config.Config{AuthToken: "test-token"}
	s := New(cfg, nil, "", nil, nil, nil, nil, bpMgr)

	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	_, ch := bpMgr.Pause("req-rest-1", "rule-1", req, []byte("original"))
	defer bpMgr.Remove("req-rest-1")

	// Call REST resume endpoint
	body := `{"method":"PUT","url":"http://example.com/new","body":"bW9kaWZpZWQ="}` // "modified" in base64
	w := httptest.NewRecorder()
	restReq, _ := http.NewRequest("POST", "/session/breakpoint/req-rest-1/resume", bytes.NewReader([]byte(body)))
	restReq.Header.Set("Authorization", "Bearer test-token")

	// Trigger via ServeHTTP
	s.mux.ServeHTTP(w, restReq)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	select {
	case action := <-ch:
		if action.Type != mitm.ActionResume {
			t.Fatalf("Expected ActionResume, got %v", action.Type)
		}
		if action.Method != "PUT" {
			t.Fatalf("Expected Method PUT, got %s", action.Method)
		}
		if action.URL != "http://example.com/new" {
			t.Fatalf("Expected URL, got %s", action.URL)
		}
		if string(action.Body) != "modified" {
			t.Fatalf("Expected Body 'modified', got %s", string(action.Body))
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for ActionResume")
	}
}

func TestRESTBreakpointRespondAndAbort(t *testing.T) {
	bpMgr := mitm.NewBreakpointManager()
	cfg := &config.Config{AuthToken: "test-token"}
	s := New(cfg, nil, "", nil, nil, nil, nil, bpMgr)

	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	_, ch := bpMgr.Pause("req-rest-2", "rule-2", req, []byte("original"))
	defer bpMgr.Remove("req-rest-2")

	// Test Respond directly
	body := `{"statusCode":200,"headers":{"Content-Type":"application/json"},"body":"eyJvayI6dHJ1ZX0="}` // '{"ok":true}' in base64
	w := httptest.NewRecorder()
	restReq, _ := http.NewRequest("POST", "/session/breakpoint/req-rest-2/respond", bytes.NewReader([]byte(body)))
	restReq.Header.Set("Authorization", "Bearer test-token")

	s.mux.ServeHTTP(w, restReq)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	select {
	case action := <-ch:
		if action.Type != mitm.ActionRespond {
			t.Fatalf("Expected ActionRespond, got %v", action.Type)
		}
		if action.StatusCode != 200 {
			t.Fatalf("Expected StatusCode 200, got %d", action.StatusCode)
		}
		if string(action.RespBody) != `{"ok":true}` {
			t.Fatalf("Expected RespBody '{\"ok\":true}', got %s", string(action.RespBody))
		}
		if action.RespHeaders["Content-Type"] != "application/json" {
			t.Fatalf("Expected Content-Type header")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for ActionRespond")
	}
}
