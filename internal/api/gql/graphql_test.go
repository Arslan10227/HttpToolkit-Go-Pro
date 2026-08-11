package gql_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/api/gql"
)

// mockProvider implements gql.Provider with stub data.
type mockProvider struct{}

func (m *mockProvider) Version() string                   { return "1.0.0-go-test" }
func (m *mockProvider) CertificatePath() string           { return "/tmp/cert.pem" }
func (m *mockProvider) CertificateContent() string        { return "-----BEGIN CERTIFICATE-----" }
func (m *mockProvider) CertificateFingerprint() string    { return "sha256/abc123" }
func (m *mockProvider) NetworkInterfaces() map[string]any { return map[string]any{} }
func (m *mockProvider) SystemProxy() map[string]any       { return nil }
func (m *mockProvider) DNSServers(_ int) []string         { return []string{"127.0.0.1"} }
func (m *mockProvider) RuleParameterKeys() []string       { return []string{"key1"} }
func (m *mockProvider) GetInterceptors(_ int) []map[string]any {
	return []map[string]any{{"id": "test", "version": "1", "isActivable": true}}
}
func (m *mockProvider) GetInterceptor(id string) (map[string]any, error) {
	return map[string]any{"id": id, "version": "1", "isActivable": true}, nil
}
func (m *mockProvider) IsInterceptorActive(_ string, _ int) bool { return false }
func (m *mockProvider) InterceptorMetadata(_ string, _ string) (any, error) {
	return map[string]any{"info": "stub"}, nil
}
func (m *mockProvider) ActivateInterceptor(_ string, _ int, _ map[string]any) (any, error) {
	return map[string]any{"success": true}, nil
}
func (m *mockProvider) DeactivateInterceptor(_ string, _ int, _ map[string]any) (bool, error) {
	return true, nil
}
func (m *mockProvider) TriggerUpdate()   {}
func (m *mockProvider) TriggerShutdown() {}

func TestGraphQLVersionQuery(t *testing.T) {
	h := gql.Handler(&mockProvider{})
	body, _ := json.Marshal(map[string]string{"query": `{ version }`})
	req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rr.Code, rr.Body.String())
	}
	var result map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	data, ok := result["data"].(map[string]any)
	if !ok {
		t.Fatalf("no data field: %v", result)
	}
	if data["version"] != "1.0.0-go-test" {
		t.Errorf("unexpected version: %v", data["version"])
	}
}

func TestGraphQLInterceptorsQuery(t *testing.T) {
	h := gql.Handler(&mockProvider{})
	body, _ := json.Marshal(map[string]string{
		"query": `{ interceptors { id version isActivable } }`,
	})
	req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rr.Code, rr.Body.String())
	}
	var result map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&result)
	data := result["data"].(map[string]any)
	interceptors, ok := data["interceptors"].([]any)
	if !ok || len(interceptors) == 0 {
		t.Fatalf("expected interceptors list, got: %v", data)
	}
}

func TestGraphQLConfigQuery(t *testing.T) {
	h := gql.Handler(&mockProvider{})
	body, _ := json.Marshal(map[string]string{
		"query": `{ config { certificatePath certificateContent certificateFingerprint } }`,
	})
	req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rr.Code, rr.Body.String())
	}
	var result map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&result)
	data := result["data"].(map[string]any)
	cfg, ok := data["config"].(map[string]any)
	if !ok {
		t.Fatalf("expected config object, got: %v", data)
	}
	if cfg["certificateFingerprint"] != "sha256/abc123" {
		t.Errorf("unexpected fingerprint: %v", cfg["certificateFingerprint"])
	}
}

func TestGraphQLMethodNotAllowed(t *testing.T) {
	h := gql.Handler(&mockProvider{})
	req := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 got %d", rr.Code)
	}
}
