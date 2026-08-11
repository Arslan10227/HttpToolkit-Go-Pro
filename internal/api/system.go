package api

import "runtime"

func systemProxy() map[string]any {
	if runtime.GOOS != "windows" {
		return nil
	}
	return map[string]any{
		"proxyUrl": "system",
		"noProxy":  "<local>",
	}
}
