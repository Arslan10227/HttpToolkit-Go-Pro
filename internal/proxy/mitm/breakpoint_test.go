package mitm

import (
	"net/http"
	"testing"
	"time"
)

func TestBreakpointManager(t *testing.T) {
	bm := NewBreakpointManager()

	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	body := []byte("original-body")

	// 1. Pause request
	pr, ch := bm.Pause("req-1", "rule-abc", req, body)
	if pr.ID != "req-1" || pr.RuleID != "rule-abc" {
		t.Fatalf("Expected PausedRequest with id 'req-1' and ruleID 'rule-abc', got %+v", pr)
	}

	// Verify Get
	got := bm.Get("req-1")
	if got != pr {
		t.Fatalf("Expected Get to return correct PausedRequest")
	}

	// 2. Test Resume Action
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch <- BreakpointAction{
			Type:   ActionResume,
			Method: "POST",
			Body:   []byte("modified-body"),
		}
	}()

	select {
	case action := <-pr.ActionChan:
		if action.Type != ActionResume {
			t.Fatalf("Expected ActionResume, got %v", action.Type)
		}
		if action.Method != "POST" {
			t.Fatalf("Expected modified Method 'POST', got %s", action.Method)
		}
		if string(action.Body) != "modified-body" {
			t.Fatalf("Expected modified Body 'modified-body', got %s", string(action.Body))
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for ActionResume")
	}

	// 3. Verify Remove
	bm.Remove("req-1")
	if bm.Get("req-1") != nil {
		t.Fatal("Expected Get to return nil after Remove")
	}
}

func TestBreakpointRespondAndAbort(t *testing.T) {
	bm := NewBreakpointManager()

	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	body := []byte("original-body")

	// Test Respond Action
	_, chRespond := bm.Pause("req-2", "rule-xyz", req, body)
	go func() {
		chRespond <- BreakpointAction{
			Type:        ActionRespond,
			StatusCode:  201,
			RespBody:    []byte("created"),
			RespHeaders: map[string]string{"Content-Type": "text/plain"},
		}
	}()

	select {
	case action := <-chRespond:
		if action.Type != ActionRespond {
			t.Fatalf("Expected ActionRespond, got %v", action.Type)
		}
		if action.StatusCode != 201 {
			t.Fatalf("Expected StatusCode 201, got %d", action.StatusCode)
		}
		if string(action.RespBody) != "created" {
			t.Fatalf("Expected RespBody 'created', got %s", string(action.RespBody))
		}
		if action.RespHeaders["Content-Type"] != "text/plain" {
			t.Fatalf("Expected Content-Type 'text/plain'")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for ActionRespond")
	}

	// Test Abort Action
	_, chAbort := bm.Pause("req-3", "rule-123", req, body)
	go func() {
		chAbort <- BreakpointAction{Type: ActionAbort}
	}()

	select {
	case action := <-chAbort:
		if action.Type != ActionAbort {
			t.Fatalf("Expected ActionAbort, got %v", action.Type)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for ActionAbort")
	}
}
