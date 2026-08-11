package admin

import (
	"net/http"

	"github.com/coder/websocket"
)

func (s *Server) handleEventsWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{"*"},
	})
	if err != nil {
		Error(err, map[string]any{"msg": "WebSocket Accept failed", "remoteAddr": r.RemoteAddr, "origin": r.Header.Get("Origin")})
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	Info("WebSocket event connection established", map[string]any{"remoteAddr": r.RemoteAddr})

	ch := s.events.Subscribe()
	defer s.events.Unsubscribe(ch)

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			Info("WebSocket event connection context done", map[string]any{"remoteAddr": r.RemoteAddr})
			return
		case msg, ok := <-ch:
			if !ok {
				Info("WebSocket events channel closed", map[string]any{"remoteAddr": r.RemoteAddr})
				return
			}
			if err := conn.Write(ctx, websocket.MessageText, msg); err != nil {
				Error(err, map[string]any{"msg": "WebSocket Write failed", "remoteAddr": r.RemoteAddr})
				return
			}
		}
	}
}
