package mitm

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/rules"
)

func TestApplyRequestTransforms(t *testing.T) {
	s := &Server{}

	req, _ := http.NewRequest("POST", "http://original-host.com/original-path?original-query=value", bytes.NewReader([]byte(`{"hello":"world"}`)))
	req.Header.Set("x-original-header", "original-val")
	req.Header.Set("x-remove-header", "remove-val")

	targetURL, _ := url.Parse(req.URL.String())
	bodyBytes := []byte(`{"hello":"world"}`)

	tr := &rules.TransformRequestOptions{
		ReplaceHost: &rules.ReplaceHostOptions{
			TargetHost:       "new-host.com",
			UpdateHostHeader: true,
		},
		SetProtocol: "https",
		ReplaceHeaders: map[string]string{
			"x-new-header": "new-val",
		},
		UpdateHeaders: map[string]string{
			"x-remove-header": "",
			"x-update-header": "update-val",
		},
		UpdateJsonBody: map[string]any{
			"hello": "rewritten",
			"added": 123.0,
		},
		MatchReplacePath: []rules.MatchReplacePair{
			{Match: "original-path", Replace: "new-path"},
		},
		MatchReplaceQuery: []rules.MatchReplacePair{
			{Match: "original-query", Replace: "new-query"},
		},
	}

	transformedBody, err := s.applyRequestTransforms(req, targetURL, bodyBytes, tr)
	if err != nil {
		t.Fatal(err)
	}

	if targetURL.Host != "new-host.com" {
		t.Fatalf("Expected targetURL.Host to be new-host.com, got %s", targetURL.Host)
	}
	if targetURL.Scheme != "https" {
		t.Fatalf("Expected targetURL.Scheme to be https, got %s", targetURL.Scheme)
	}
	if targetURL.Path != "/new-path" {
		t.Fatalf("Expected targetURL.Path to be /new-path, got %s", targetURL.Path)
	}
	if targetURL.RawQuery != "new-query=value" {
		t.Fatalf("Expected targetURL.RawQuery to be new-query=value, got %s", targetURL.RawQuery)
	}

	if req.Header.Get("x-new-header") != "new-val" {
		t.Fatalf("Expected x-new-header to be new-val, got %s", req.Header.Get("x-new-header"))
	}
	if req.Header.Get("x-original-header") != "" {
		t.Fatal("Expected x-original-header to be replaced/removed")
	}
	if req.Header.Get("x-remove-header") != "" {
		t.Fatal("Expected x-remove-header to be removed")
	}
	if req.Header.Get("x-update-header") != "update-val" {
		t.Fatalf("Expected x-update-header to be update-val, got %s", req.Header.Get("x-update-header"))
	}

	expectedBody := `{"added":123,"hello":"rewritten"}`
	if string(transformedBody) != expectedBody {
		t.Fatalf("Expected body to be %s, got %s", expectedBody, string(transformedBody))
	}
}

func TestApplyResponseTransforms(t *testing.T) {
	s := &Server{}

	resp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{},
	}
	resp.Header.Set("x-original-res", "val")
	bodyBytes := []byte(`{"status":"ok"}`)

	tr := &rules.TransformResponseOptions{
		ReplaceStatus: 201,
		UpdateHeaders: map[string]string{
			"x-original-res": "",
			"x-added-res":    "added",
		},
		ReplaceBody: func() *string { s := "replaced-body"; return &s }(),
	}

	transformedBody, status, err := s.applyResponseTransforms(resp, bodyBytes, tr)
	if err != nil {
		t.Fatal(err)
	}

	if status != 201 {
		t.Fatalf("Expected status to be 201, got %d", status)
	}
	if resp.Header.Get("x-original-res") != "" {
		t.Fatal("Expected x-original-res to be deleted")
	}
	if resp.Header.Get("x-added-res") != "added" {
		t.Fatalf("Expected x-added-res to be added, got %s", resp.Header.Get("x-added-res"))
	}
	if string(transformedBody) != "replaced-body" {
		t.Fatalf("Expected body to be 'replaced-body', got %s", string(transformedBody))
	}
}
