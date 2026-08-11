package api

import (
	"net/http"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/mcp"
)

func (s *Server) handleMCPStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, mcp.Status())
}

func (s *Server) handleMCPTools(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{"tools": mcp.Tools()})
}
