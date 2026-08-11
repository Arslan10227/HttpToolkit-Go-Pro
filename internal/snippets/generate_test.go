package snippets

import (
	"strings"
	"testing"
)

func TestGenerateShellCurl(t *testing.T) {
	req := HarRequest{
		Method: "GET",
		URL:    "https://api.example.com/users?page=1",
		Headers: []HarHeader{
			{Name: "Accept", Value: "application/json"},
		},
	}
	out, err := Generate(req, "shell", "curl")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.ToLower(out), "curl") {
		t.Fatalf("expected curl snippet, got %q", out)
	}
	if !strings.Contains(out, "api.example.com") {
		t.Fatalf("expected url in snippet, got %q", out)
	}
}

func TestGenerateGoNative(t *testing.T) {
	req := HarRequest{
		Method: "POST",
		URL:    "https://example.com/api",
		Headers: []HarHeader{
			{Name: "Content-Type", Value: "application/json"},
		},
		PostData: &HarPostData{MimeType: "application/json", Text: `{"ok":true}`},
	}
	out, err := Generate(req, "go", "native")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "http.NewRequest") {
		t.Fatalf("expected Go http client, got %q", out)
	}
}

func TestGenerateHTTPRaw(t *testing.T) {
	req := HarRequest{
		Method: "GET",
		URL:    "https://example.com/hello",
	}
	out, err := Generate(req, "http", "http1.1")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "GET /hello HTTP/1.1") {
		t.Fatalf("expected raw HTTP request, got %q", out)
	}
}

func TestGeneratePythonRequests(t *testing.T) {
	req := HarRequest{Method: "GET", URL: "https://example.com"}
	out, err := Generate(req, "python", "requests")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "requests.request") {
		t.Fatalf("expected requests snippet, got %q", out)
	}
}
