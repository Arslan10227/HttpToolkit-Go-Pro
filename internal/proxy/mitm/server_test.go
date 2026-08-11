package mitm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	certmgr "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/events"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/rules"
)

func TestMITMInterception(t *testing.T) {
	// 1. Setup target test HTTP server
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test-Response", "Yes")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello from target!"))
	}))
	defer targetServer.Close()

	// 2. Setup cert manager
	cfg := &config.Config{ConfigDir: t.TempDir()}
	certs, err := certmgr.NewManager(cfg)
	if err != nil {
		t.Fatalf("failed to create cert manager: %v", err)
	}

	// 3. Setup bus, engine, and MITM proxy server
	bus := events.NewBus()
	eng := rules.NewEngine()
	s := NewServer(certs, eng, bus)

	// Subscribe to bus events
	ch := bus.Subscribe()
	defer bus.Unsubscribe(ch)

	// Start proxy on a dynamic port
	err = s.Start(0)
	if err != nil {
		t.Fatalf("failed to start proxy server: %v", err)
	}
	defer func() { _ = s.Stop() }()

	proxyURL, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", s.Port()))
	if err != nil {
		t.Fatalf("failed to parse proxy URL: %v", err)
	}

	// 4. Configure HTTP client to use our proxy
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 5 * time.Second,
	}

	// 5. Send an HTTP request via the proxy to the target server
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, targetServer.URL+"/test-path", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("X-Test-Request", "Incoming")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed via proxy: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// 6. Gather events from the bus
	var gotRequestInitiated, gotRequest, gotResponse bool
	deadline := time.After(2 * time.Second)

loop:
	for {
		select {
		case msg := <-ch:
			msgStr := string(msg)
			if strings.Contains(msgStr, "request-initiated") && strings.Contains(msgStr, "/test-path") {
				gotRequestInitiated = true
			}
			if strings.Contains(msgStr, "\"event\":\"request\"") && strings.Contains(msgStr, "/test-path") {
				gotRequest = true
			}
			if strings.Contains(msgStr, "\"event\":\"response\"") && strings.Contains(msgStr, "Yes") {
				gotResponse = true
			}
			if gotRequestInitiated && gotRequest && gotResponse {
				break loop
			}
		case <-deadline:
			break loop
		}
	}

	if !gotRequestInitiated {
		t.Error("missing request-initiated event")
	}
	if !gotRequest {
		t.Error("missing request event")
	}
	if !gotResponse {
		t.Error("missing response event")
	}
}

// TestMITMFixedResponseRule verifies the rule engine actually intercepts a
// real HTTP request end-to-end: a "fixed-response" (`simple`) step for a
// matching rule should be returned by the live proxy instead of forwarding
// to the real upstream target, and a "response" event should still be
// published on the bus.
func TestMITMFixedResponseRule(t *testing.T) {
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("real upstream response"))
	}))
	defer targetServer.Close()

	cfg := &config.Config{ConfigDir: t.TempDir()}
	certs, err := certmgr.NewManager(cfg)
	if err != nil {
		t.Fatalf("failed to create cert manager: %v", err)
	}

	bus := events.NewBus()
	eng := rules.NewEngine()

	eng.SetHTTPRules([]rules.RequestRuleData{
		{
			ID: "mock-rule-1",
			Matchers: mustRawMessages(t,
				`{"type":"method","method":"GET"}`,
				fmt.Sprintf(`{"type":"simple-path","path":"%s"}`, "/mocked-path"),
			),
			Steps: mustRawMessages(t,
				`{"type":"simple","status":200,"data":"mocked response body","headers":{"x-mock":"yes"}}`,
			),
		},
	})

	s := NewServer(certs, eng, bus)
	ch := bus.Subscribe()
	defer bus.Unsubscribe(ch)

	if err := s.Start(0); err != nil {
		t.Fatalf("failed to start proxy server: %v", err)
	}
	defer func() { _ = s.Stop() }()

	proxyURL, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", s.Port()))
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
		Timeout:   5 * time.Second,
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, targetServer.URL+"/mocked-path", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed via proxy: %v", err)
	}
	defer resp.Body.Close()

	body := make([]byte, 256)
	n, _ := resp.Body.Read(body)
	got := string(body[:n])

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if resp.Header.Get("x-mock") != "yes" {
		t.Errorf("expected mocked x-mock header, got headers: %v", resp.Header)
	}
	if !strings.Contains(got, "mocked response body") {
		t.Errorf("expected mocked body, got: %q (real upstream would say 'real upstream response')", got)
	}
}

func mustRawMessages(t *testing.T, jsons ...string) []json.RawMessage {
	t.Helper()
	out := make([]json.RawMessage, len(jsons))
	for i, j := range jsons {
		out[i] = json.RawMessage(j)
	}
	return out
}
