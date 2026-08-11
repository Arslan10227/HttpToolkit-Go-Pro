package rules

import (
	"encoding/json"
	"testing"
)

func TestMethodMatcher(t *testing.T) {
	e := NewEngine()
	e.SetHTTPRules([]RequestRuleData{{
		Matchers: mustMatchers(t, map[string]any{"type": "method", "method": "POST"}),
		Steps:    mustSteps(t, map[string]any{"type": "passthrough"}),
	}})
	if e.MatchFirst("GET", "http://x.com/a", nil, nil) != -1 {
		t.Fatal("GET should not match POST rule")
	}
	if e.MatchFirst("POST", "http://x.com/a", nil, nil) != 0 {
		t.Fatal("POST should match")
	}
}

func TestSimpleResponseStep(t *testing.T) {
	e := NewEngine()
	e.SetHTTPRules([]RequestRuleData{{
		Matchers: mustMatchers(t, map[string]any{"type": "wildcard"}),
		Steps: mustSteps(t, map[string]any{
			"type": "simple", "status": 418, "data": "teapot",
		}),
	}})
	a := e.ActionForRule(0)
	if a.Kind != "fixed-response" || a.StatusCode != 418 || a.Body != "teapot" {
		t.Fatalf("unexpected action: %+v", a)
	}
}

func TestRedirectAndAbortSteps(t *testing.T) {
	e := NewEngine()
	e.SetHTTPRules([]RequestRuleData{
		{
			Matchers: mustMatchers(t, map[string]any{"type": "wildcard"}),
			Steps: mustSteps(t, map[string]any{
				"type": "redirect", "status": 307, "location": "http://redirect-url.com",
			}),
		},
		{
			Matchers: mustMatchers(t, map[string]any{"type": "wildcard"}),
			Steps: mustSteps(t, map[string]any{
				"type": "abort",
			}),
		},
	})
	a1 := e.ActionForRule(0)
	if a1.Kind != "redirect" || a1.StatusCode != 307 || a1.Body != "http://redirect-url.com" {
		t.Fatalf("unexpected redirect action: %+v", a1)
	}
	a2 := e.ActionForRule(1)
	if a2.Kind != "abort" {
		t.Fatalf("unexpected abort action: %+v", a2)
	}
}

func TestTransformRequestPassthrough(t *testing.T) {
	e := NewEngine()
	e.SetHTTPRules([]RequestRuleData{{
		Matchers: mustMatchers(t, map[string]any{"type": "wildcard"}),
		Steps: mustSteps(t, map[string]any{
			"type": "passthrough",
			"transformRequest": map[string]any{
				"replaceHost": map[string]any{
					"targetHost":       "new-host.com",
					"updateHostHeader": true,
				},
				"setProtocol": "https",
			},
		}),
	}})
	a := e.ActionForRule(0)
	if a.Kind != "passthrough" {
		t.Fatalf("unexpected action: %+v", a)
	}
	tr := a.Passthrough.TransformRequest
	if tr == nil || tr.ReplaceHost == nil || tr.ReplaceHost.TargetHost != "new-host.com" || !tr.ReplaceHost.UpdateHostHeader || tr.SetProtocol != "https" {
		t.Fatalf("unexpected transform: %+v", tr)
	}
}

func TestFullRequestResponseTransform(t *testing.T) {
	e := NewEngine()
	e.SetHTTPRules([]RequestRuleData{{
		Matchers: mustMatchers(t, map[string]any{"type": "wildcard"}),
		Steps: mustSteps(t, map[string]any{
			"type": "passthrough",
			"transformRequest": map[string]any{
				"replaceHeaders": map[string]any{
					"x-req-test": "req-val",
				},
				"updateHeaders": map[string]any{
					"x-req-remove": nil,
				},
				"replaceBody": map[string]any{
					"type": "Buffer",
					"data": []any{97.0, 98.0, 99.0}, // "abc"
				},
				"matchReplaceBody": []any{
					[]any{map[string]any{"source": "match-req", "flags": "g"}, "replace-req"},
				},
			},
			"transformResponse": map[string]any{
				"replaceStatus": 204.0,
				"replaceHeaders": map[string]any{
					"x-res-test": "res-val",
				},
				"matchReplaceBody": []any{
					map[string]any{"match": "match-res", "replace": "replace-res"},
				},
			},
		}),
	}})
	a := e.ActionForRule(0)
	if a.Kind != "passthrough" {
		t.Fatalf("unexpected action: %+v", a)
	}
	tr := a.Passthrough.TransformRequest
	if tr == nil || tr.ReplaceHeaders == nil || tr.ReplaceHeaders["x-req-test"] != "req-val" {
		t.Fatalf("unexpected transform request: %+v", tr)
	}
	if tr.UpdateHeaders["x-req-remove"] != "" {
		t.Fatalf("expected header deletion update to be empty string")
	}
	if tr.ReplaceBody == nil || *tr.ReplaceBody != "abc" {
		t.Fatalf("expected replaceBody to be 'abc'")
	}
	if len(tr.MatchReplaceBody) != 1 || tr.MatchReplaceBody[0].Match != "match-req" || tr.MatchReplaceBody[0].Replace != "replace-req" {
		t.Fatalf("unexpected matchReplaceBody: %+v", tr.MatchReplaceBody)
	}

	tres := a.Passthrough.TransformResponse
	if tres == nil || tres.ReplaceStatus != 204 || tres.ReplaceHeaders["x-res-test"] != "res-val" {
		t.Fatalf("unexpected transform response: %+v", tres)
	}
	if len(tres.MatchReplaceBody) != 1 || tres.MatchReplaceBody[0].Match != "match-res" || tres.MatchReplaceBody[0].Replace != "replace-res" {
		t.Fatalf("unexpected response matchReplaceBody: %+v", tres.MatchReplaceBody)
	}
}

func TestPathMatcher(t *testing.T) {
	e := NewEngine()
	e.SetHTTPRules([]RequestRuleData{
		{
			Matchers: mustMatchers(t, map[string]any{"type": "path", "path": "http://android.httptoolkit.tech/config"}),
			Steps:    mustSteps(t, map[string]any{"type": "passthrough"}),
		},
		{
			Matchers: mustMatchers(t, map[string]any{"type": "path", "path": "amiusing.local/certificate"}),
			Steps:    mustSteps(t, map[string]any{"type": "passthrough"}),
		},
		{
			Matchers: mustMatchers(t, map[string]any{"type": "path", "path": "/only-path"}),
			Steps:    mustSteps(t, map[string]any{"type": "passthrough"}),
		},
	})

	// Test case 1: absolute URL
	if e.MatchFirst("GET", "http://android.httptoolkit.tech/config", nil, nil) != 0 {
		t.Fatal("Expected absolute URL match to succeed")
	}
	if e.MatchFirst("GET", "http://someotherhost.com/config", nil, nil) == 0 {
		t.Fatal("Expected absolute URL host mismatch to fail matching")
	}

	// Test case 2: host/path format
	if e.MatchFirst("GET", "https://amiusing.local/certificate", nil, nil) != 1 {
		t.Fatal("Expected host/path format match to succeed")
	}
	if e.MatchFirst("GET", "https://other.com/certificate", nil, nil) == 1 {
		t.Fatal("Expected host/path format host mismatch to fail matching")
	}

	// Test case 3: relative path
	if e.MatchFirst("GET", "https://anyhost.com/only-path", nil, nil) != 2 {
		t.Fatal("Expected relative path match to succeed on any host")
	}
}

func mustMatchers(t *testing.T, obj map[string]any) []json.RawMessage {
	t.Helper()
	b, err := json.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	}
	return []json.RawMessage{b}
}

func mustSteps(t *testing.T, obj map[string]any) []json.RawMessage {
	return mustMatchers(t, obj)
}
