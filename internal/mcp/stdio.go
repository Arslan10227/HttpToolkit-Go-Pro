package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/mcp/ctlclient"
)

type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type jsonRPCNotification struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
}

type ctlOperation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema"`
	Annotations any    `json:"annotations"`
}

const pollInterval = 5 * time.Second

// RunStdio serves MCP JSON-RPC over stdin/stdout.
func RunStdio(log io.Writer) error {
	if log == nil {
		log = os.Stderr
	}
	client, err := ctlclient.New()
	if err != nil {
		return err
	}

	var (
		mu      sync.Mutex
		cached  []ctlOperation
		lastKey string
	)

	refresh := func() {
		ops, err := client.Operations()
		if err != nil {
			mu.Lock()
			cached = nil
			mu.Unlock()
			return
		}
		raw, _ := json.Marshal(ops)
		var parsed []ctlOperation
		_ = json.Unmarshal(raw, &parsed)
		mu.Lock()
		cached = parsed
		mu.Unlock()
	}
	refresh()

	send := func(v any) {
		data, _ := json.Marshal(v)
		fmt.Fprintf(os.Stdout, "%s\n", data)
	}
	result := func(id any, res any) {
		send(jsonRPCResponse{JSONRPC: "2.0", ID: id, Result: res})
	}
	rpcErr := func(id any, code int, msg string) {
		send(jsonRPCResponse{
			JSONRPC: "2.0", ID: id,
			Error: &struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{Code: code, Message: msg},
		})
	}

	toolsList := func() []map[string]any {
		mu.Lock()
		ops := append([]ctlOperation(nil), cached...)
		mu.Unlock()
		out := make([]map[string]any, 0, len(ops)+len(BuiltInTools()))
		for _, op := range ops {
			tool := map[string]any{
				"name":        strings.ReplaceAll(op.Name, ".", "_"),
				"description": op.Description,
				"inputSchema": map[string]any{
					"type":       "object",
					"properties": schemaProperties(op.InputSchema),
				},
			}
			if op.Annotations != nil {
				tool["annotations"] = op.Annotations
			}
			out = append(out, tool)
		}
		for _, t := range BuiltInTools() {
			out = append(out, t)
		}
		return out
	}

	callTool := func(name string, args map[string]any) map[string]any {
		if res, err := CallBuiltIn(name, args); err == nil {
			raw, _ := json.Marshal(res)
			return map[string]any{
				"content": []map[string]any{{"type": "text", "text": string(raw)}},
			}
		} else if !strings.Contains(err.Error(), "unknown tool") {
			return map[string]any{
				"content": []map[string]any{{"type": "text", "text": err.Error()}},
				"isError": true,
			}
		}

		mu.Lock()
		ops := append([]ctlOperation(nil), cached...)
		mu.Unlock()

		opName := name
		for _, op := range ops {
			if strings.ReplaceAll(op.Name, ".", "_") == name {
				opName = op.Name
				break
			}
		}
		if opName == name {
			opName = strings.ReplaceAll(name, "_", ".")
		}
		if args == nil {
			args = map[string]any{}
		}
		res, err := client.Execute(opName, args)
		if err != nil {
			return map[string]any{
				"content": []map[string]any{{"type": "text", "text": err.Error()}},
				"isError": true,
			}
		}
		raw, _ := json.Marshal(res)
		if m, ok := res.(map[string]any); ok {
			if success, ok := m["success"].(bool); ok && !success {
				if errObj, ok := m["error"].(map[string]any); ok {
					if msg, ok := errObj["message"].(string); ok {
						return map[string]any{
							"content": []map[string]any{{"type": "text", "text": msg}},
							"isError": true,
						}
					}
				}
			}
		}
		return map[string]any{
			"content": []map[string]any{{"type": "text", "text": string(raw)}},
		}
	}

	handle := func(msg jsonRPCRequest) {
		switch msg.Method {
		case "initialize":
			result(msg.ID, map[string]any{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]any{"tools": map[string]any{"listChanged": true}},
				"serverInfo":      map[string]any{"name": "httptoolkit", "version": config.ServerVersion()},
			})
		case "notifications/initialized":
		case "tools/list":
			result(msg.ID, map[string]any{"tools": toolsList()})
		case "tools/call":
			var params struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments"`
			}
			if err := json.Unmarshal(msg.Params, &params); err != nil || params.Name == "" {
				rpcErr(msg.ID, -32602, "Missing tool name")
				return
			}
			fmt.Fprintf(log, "[MCP] Tool called: %s\n", params.Name)
			result(msg.ID, callTool(params.Name, params.Arguments))
		default:
			if msg.ID != nil {
				rpcErr(msg.ID, -32601, "Method not found: "+msg.Method)
			}
		}
	}

	stopPoll := make(chan struct{})
	go func() {
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-stopPoll:
				return
			case <-ticker.C:
				refresh()
				mu.Lock()
				key := opsKey(cached)
				mu.Unlock()
				if key != lastKey {
					lastKey = key
					send(jsonRPCNotification{JSONRPC: "2.0", Method: "notifications/tools/list_changed"})
					fmt.Fprintln(log, "[MCP] Sent tools/list_changed")
				}
			}
		}
	}()

	fmt.Fprintln(log, "[MCP] MCP server started on stdio")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var msg jsonRPCRequest
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			rpcErr(nil, -32700, "Parse error")
			continue
		}
		handle(msg)
	}
	close(stopPoll)
	return scanner.Err()
}

func opsKey(ops []ctlOperation) string {
	names := make([]string, len(ops))
	for i, op := range ops {
		names[i] = op.Name
	}
	data, _ := json.Marshal(names)
	return string(data)
}

func schemaProperties(schema any) map[string]any {
	if schema == nil {
		return map[string]any{}
	}
	raw, _ := json.Marshal(schema)
	var parsed map[string]any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return map[string]any{}
	}
	if props, ok := parsed["properties"].(map[string]any); ok {
		return props
	}
	return map[string]any{}
}
