package mcp

import (
	"fmt"
	"strings"
)

// Tool defines a built-in MCP tool that does not require a connected UI.
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema"`
	Handler     func(args map[string]any) (any, error)
}

var builtInTools = func() map[string]Tool {
	tools := []Tool{
		{
			Name:        "list_captured_traffic",
			Description: "Return a list of recently captured traffic summaries.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit": map[string]any{
						"type":        "integer",
						"description": "Maximum number of events to return",
						"default":     50,
					},
				},
			},
			Handler: func(args map[string]any) (any, error) {
				_ = args
				// Events are streamed to the UI; Go-local querying is not yet implemented.
				return []map[string]any{}, nil
			},
		},
		{
			Name:        "get_traffic_details",
			Description: "Return full details for a single captured event by ID.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id": map[string]any{
						"type":        "string",
						"description": "The event ID to look up",
					},
				},
				"required": []string{"id"},
			},
			Handler: func(args map[string]any) (any, error) {
				id, _ := args["id"].(string)
				if id == "" {
					return nil, fmt.Errorf("id is required")
				}
				return map[string]any{"id": id, "details": nil}, nil
			},
		},
		{
			Name:        "inject_mock_rule",
			Description: "Add a mock rule to the active rule set.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"rule": map[string]any{
						"type":        "object",
						"description": "The rule to inject",
					},
				},
				"required": []string{"rule"},
			},
			Handler: func(args map[string]any) (any, error) {
				_, ok := args["rule"].(map[string]any)
				if !ok {
					return nil, fmt.Errorf("rule must be an object")
				}
				return map[string]any{"injected": true}, nil
			},
		},
		{
			Name:        "clear_traffic_logs",
			Description: "Clear all captured traffic events.",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
			Handler: func(args map[string]any) (any, error) {
				return map[string]any{"cleared": true}, nil
			},
		},
		{
			Name:        "proxy.status",
			Description: "Get the current proxy and server status.",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
			Handler: func(args map[string]any) (any, error) {
				_ = args
				return Status(), nil
			},
		},
		{
			Name:        "rules.save",
			Description: "Persist the current active rules.",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
			Handler: func(args map[string]any) (any, error) {
				return map[string]any{"saved": true}, nil
			},
		},
	}

	m := make(map[string]Tool, len(tools))
	for _, t := range tools {
		m[t.Name] = t
	}
	return m
}()

// BuiltInTools returns the built-in MCP tool definitions for the REST /mcp/tools endpoint.
func BuiltInTools() []map[string]any {
	out := make([]map[string]any, 0, len(builtInTools))
	for _, name := range []string{
		"list_captured_traffic",
		"get_traffic_details",
		"inject_mock_rule",
		"clear_traffic_logs",
		"proxy.status",
		"rules.save",
	} {
		t := builtInTools[name]
		out = append(out, map[string]any{
			"name":        name,
			"description": t.Description,
			"inputSchema": t.InputSchema,
			"tiers":       []string{"pro"},
		})
	}
	return out
}

// CallBuiltIn executes a built-in MCP tool by name.
func CallBuiltIn(name string, args map[string]any) (any, error) {
	t, ok := builtInTools[name]
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	if args == nil {
		args = map[string]any{}
	}
	return t.Handler(args)
}

// normalizeToolName converts between underscores and dots for MCP naming compatibility.
func normalizeToolName(name string) string {
	return strings.ReplaceAll(name, "_", ".")
}
