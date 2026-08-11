package rules

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type WebhookAction struct {
	URL     string
	Method  string
	Headers map[string]string
}

func parseWebhookStep(obj map[string]any) WebhookAction {
	wh := WebhookAction{Method: http.MethodPost}
	if u, _ := obj["url"].(string); u != "" {
		wh.URL = u
	}
	if u, _ := obj["webhookUrl"].(string); u != "" {
		wh.URL = u
	}
	if m, _ := obj["method"].(string); m != "" {
		wh.Method = m
	}
	if h, ok := obj["headers"].(map[string]any); ok {
		wh.Headers = make(map[string]string, len(h))
		for k, v := range h {
			wh.Headers[k] = toString(v)
		}
	}
	return wh
}

func extractWebhookFromSteps(steps []json.RawMessage) WebhookAction {
	for _, s := range steps {
		var obj map[string]any
		if err := json.Unmarshal(s, &obj); err != nil {
			continue
		}
		if typ, _ := obj["type"].(string); typ == "webhook" {
			return parseWebhookStep(obj)
		}
	}
	return WebhookAction{}
}

// FireWebhook sends a best-effort notification for a matched rule.
func FireWebhook(wh WebhookAction, payload map[string]any) {
	if wh.URL == "" {
		return
	}
	data, _ := json.Marshal(payload)
	req, err := http.NewRequest(wh.Method, wh.URL, bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range wh.Headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err == nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
}
