package rules

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFireWebhook(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	FireWebhook(WebhookAction{URL: srv.URL, Method: http.MethodPost}, map[string]any{"ok": true})
	if !called {
		t.Fatal("expected webhook to be called")
	}
}

func TestExtractWebhookFromSteps(t *testing.T) {
	steps := []json.RawMessage{[]byte(`{"type":"webhook","url":"https://example.com/hook"}`)}
	wh := extractWebhookFromSteps(steps)
	if wh.URL != "https://example.com/hook" {
		t.Fatalf("unexpected webhook url: %q", wh.URL)
	}
}
