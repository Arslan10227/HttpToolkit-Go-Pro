package mitm

import (
	"net"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/logger"
)

func (s *Server) emitRawPassthrough(id, host string, open bool) {
	event := "raw-passthrough-closed"
	if open {
		event = "raw-passthrough-opened"
	}
	s.events.PublishHTTP(event, map[string]any{
		"id": id, "destination": host, "timestamp": time.Now().UnixMilli(),
	})
}

func (s *Server) emitTLSPassthrough(id, host string, open bool) {
	event := "tls-passthrough-closed"
	if open {
		event = "tls-passthrough-opened"
	}
	s.events.PublishHTTP(event, map[string]any{
		"id": id, "destination": host, "timestamp": time.Now().UnixMilli(),
	})
}

func resetTCPConn(conn net.Conn) {
	if tc, ok := conn.(*net.TCPConn); ok {
		if err := tc.SetLinger(0); err != nil {
			logger.Error(err, map[string]any{"msg": "failed to reset TCP linger"})
		}
	}
}

func isTLSDefaultPort(host string) bool {
	_, port, err := net.SplitHostPort(host)
	if err != nil {
		// No port specified - assume standard TLS.
		return true
	}
	return port == "443" || port == "8443"
}
