package mcp

import "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"

func Status() map[string]any {
	return map[string]any{
		"bridgeAvailable": true,
		"mcpCommand":      "htk-mcp",
		"version":         config.ServerVersion(),
	}
}

func Tools() []map[string]any {
	return BuiltInTools()
}
