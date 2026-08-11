package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/snippets"
)

func (s *Server) handleClientSend(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Request struct {
			Method  string `json:"method"`
			URL     string `json:"url"`
			Headers any    `json:"headers"`
			RawBody string `json:"rawBody"`
		} `json:"request"`
		Options map[string]any `json:"options"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/x-ndjson")
	flusher, _ := w.(http.Flusher)
	enc := func(v any) {
		_ = json.NewEncoder(w).Encode(v)
		if flusher != nil {
			flusher.Flush()
		}
	}

	ts := time.Now().UnixMilli()
	enc(map[string]any{
		"type": "request-start", "timestamp": ts, "startTime": ts,
		"request": body.Request,
	})

	headerMap, err := normalizeSendHeaders(body.Request.Headers)
	if err != nil {
		enc(map[string]any{"type": "error", "timestamp": time.Now().UnixMilli(), "error": map[string]string{"message": err.Error()}})
		return
	}

	reqBody, _ := base64.StdEncoding.DecodeString(body.Request.RawBody)
	req, err := http.NewRequest(body.Request.Method, body.Request.URL, bytes.NewReader(reqBody))
	if err != nil {
		enc(map[string]any{"type": "error", "timestamp": time.Now().UnixMilli(), "error": map[string]string{"message": err.Error()}})
		return
	}
	for k, vals := range headerMap {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	opts := parseSendOptions(body.Options)
	transport, err := buildSendTransport(req.URL, opts)
	if err != nil {
		enc(map[string]any{"type": "error", "timestamp": time.Now().UnixMilli(), "error": map[string]string{"message": err.Error()}})
		return
	}
	client := &http.Client{Transport: transport, Timeout: 120 * time.Second}

	resp, err := client.Do(req.WithContext(r.Context()))
	if err != nil {
		enc(map[string]any{
			"type": "error", "timestamp": time.Now().UnixMilli(),
			"error": map[string]any{"message": err.Error(), "code": "REQUEST_FAILED"},
		})
		return
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	rawHeaders := make([]string, 0, len(resp.Header)*2)
	for k, vals := range resp.Header {
		for _, v := range vals {
			rawHeaders = append(rawHeaders, k, v)
		}
	}
	statusMsg := resp.Status
	if parts := splitStatus(resp.Status); len(parts) > 1 {
		statusMsg = parts[1]
	}

	enc(map[string]any{
		"type": "response-head", "timestamp": time.Now().UnixMilli(),
		"statusCode": resp.StatusCode, "statusMessage": statusMsg,
		"headers": rawHeaders,
	})
	enc(map[string]any{
		"type": "response-body-part", "timestamp": time.Now().UnixMilli(),
		"rawBody": base64.StdEncoding.EncodeToString(respBody),
	})
	enc(map[string]any{"type": "response-end", "timestamp": time.Now().UnixMilli()})
}

func splitStatus(status string) []string {
	parts := make([]string, 0, 2)
	for i, c := range status {
		if c == ' ' && i > 0 {
			parts = append(parts, status[:i], status[i+1:])
			break
		}
	}
	if len(parts) == 0 {
		parts = append(parts, status)
	}
	return parts
}

func (s *Server) handleSnippets(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Request json.RawMessage `json:"request"`
		Target  string          `json:"target"`
		Client  string          `json:"client"`
		Format  string          `json:"format"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	harReq, err := snippets.ParseHarRequest(body.Request)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	target, client := body.Target, body.Client
	if target == "" && body.Format != "" {
		if parts := strings.SplitN(body.Format, "~~", 2); len(parts) == 2 {
			target, client = parts[0], parts[1]
		} else {
			target, client = "shell", body.Format
		}
	}
	if target == "" || client == "" {
		writeAPIError(w, http.StatusBadRequest, "Missing required snippet parameters")
		return
	}
	snippet, err := snippets.Generate(harReq, target, client)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, map[string]any{"snippet": snippet})
}
