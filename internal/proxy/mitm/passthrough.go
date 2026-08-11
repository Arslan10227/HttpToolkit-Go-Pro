package mitm

import (
	"io"
	"net"
	"time"
)

func tunnelWithEvents(client, target net.Conn, bus eventPublisher, id, kind string) {
	eventName := kind + "-opened"
	closeName := kind + "-closed"
	bus.PublishHTTP(eventName, map[string]any{
		"id": id, "timestamp": time.Now().UnixMilli(),
	})
	done := make(chan struct{}, 2)
	go func() {
		_, _ = io.Copy(target, client)
		done <- struct{}{}
	}()
	go func() {
		_, _ = io.Copy(client, target)
		done <- struct{}{}
	}()
	<-done
	bus.PublishHTTP(closeName, map[string]any{
		"id": id, "timestamp": time.Now().UnixMilli(),
	})
}

type eventPublisher interface {
	PublishHTTP(event string, data map[string]any)
}

func resetConn(conn net.Conn) {
	if conn == nil {
		return
	}
	if tc, ok := conn.(*net.TCPConn); ok {
		_ = tc.SetLinger(0)
	}
	_ = conn.Close()
}
