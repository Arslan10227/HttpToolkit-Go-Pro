package mitm

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	certmgr "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/events"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/rules"
)

// TestWebSocketEchoE2E tests the WebSocket echo mode through the MITM proxy.
// It sets up a mock HTTP server that responds with 101 Switching Protocols,
// then connects through the proxy and verifies the echo round-trip.
func TestWebSocketEchoE2E(t *testing.T) {
	cfg := &config.Config{ConfigDir: t.TempDir()}
	certs, err := certmgr.NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	bus := events.NewBus()
	eng := rules.NewEngine()
	srv := NewServer(certs, eng, bus)

	if err := srv.Start(0); err != nil {
		t.Fatalf("Server.Start: %v", err)
	}
	defer srv.Stop()

	// Set a ws-echo rule for any WebSocket upgrade
	wsStep, _ := json.Marshal(map[string]any{"type": "ws-echo"})
	eng.SetWSRules([]rules.WebSocketRuleData{
		{ID: "echo", Steps: []json.RawMessage{wsStep}},
	})

	// Connect to the proxy and send a WebSocket upgrade request
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", srv.Port()))
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()

	key := "dGhlIHNhbXBsZSBub25jZQ=="
	upgradeReq := "GET / HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Key: " + key + "\r\n" +
		"Sec-WebSocket-Version: 13\r\n" +
		"\r\n"
	if _, err := conn.Write([]byte(upgradeReq)); err != nil {
		t.Fatalf("Write upgrade: %v", err)
	}

	// Read the 101 response
	br := bufio.NewReader(conn)
	resp, err := http.ReadResponse(br, nil)
	if err != nil {
		t.Fatalf("ReadResponse: %v", err)
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("expected 101, got %d", resp.StatusCode)
	}

	// Verify Sec-WebSocket-Accept
	expectedAccept := computeWebsocketAccept(key)
	if resp.Header.Get("Sec-WebSocket-Accept") != expectedAccept {
		t.Errorf("Sec-WebSocket-Accept = %q, want %q",
			resp.Header.Get("Sec-WebSocket-Accept"), expectedAccept)
	}

	// Send a text frame "hello"
	frame := buildWSTextFrame("hello", true)
	if _, err := conn.Write(frame); err != nil {
		t.Fatalf("Write frame: %v", err)
	}

	// Read the echoed frame
	echoFrame := make([]byte, 10)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := br.Read(echoFrame)
	if err != nil {
		t.Fatalf("Read echo: %v", err)
	}

	// Parse the echoed frame (opcode 0x01 = text, masked=false from server)
	if n < 2 {
		t.Fatalf("echo frame too short: %d bytes", n)
	}
	opcode := echoFrame[0] & 0x0F
	if opcode != 0x01 {
		t.Errorf("opcode = %d, want 1 (text)", opcode)
	}
	payloadLen := int(echoFrame[1] & 0x7F)
	if payloadLen != 5 {
		t.Errorf("payload length = %d, want 5", payloadLen)
	}
	if n < 2+payloadLen {
		t.Fatalf("echo frame truncated: got %d bytes, need %d", n, 2+payloadLen)
	}
	payload := string(echoFrame[2 : 2+payloadLen])
	if payload != "hello" {
		t.Errorf("echo payload = %q, want %q", payload, "hello")
	}
}

// TestWebSocketPassthroughE2E tests WebSocket relay to a real upstream server
// through the MITM proxy (no echo rule).
func TestWebSocketPassthroughE2E(t *testing.T) {
	// Set up a fake upstream WebSocket server
	upstreamLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	defer upstreamLn.Close()

	go func() {
		conn, err := upstreamLn.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		br := bufio.NewReader(conn)
		req, err := http.ReadRequest(br)
		if err != nil {
			return
		}
		_ = req
		// Respond with 101 Switching Protocols
		acceptKey := computeWebsocketAccept(req.Header.Get("Sec-WebSocket-Key"))
		resp := "HTTP/1.1 101 Switching Protocols\r\n" +
			"Upgrade: websocket\r\n" +
			"Connection: Upgrade\r\n" +
			"Sec-WebSocket-Accept: " + acceptKey + "\r\n" +
			"\r\n"
		conn.Write([]byte(resp))
		// Read a frame and echo it back
		buf := make([]byte, 256)
		n, _ := br.Read(buf)
		if n > 0 {
			// Echo back as an unmasked text frame
			frame := buildWSTextFrame(string(buf[2:2+int(buf[1]&0x7F)]), false)
			conn.Write(frame)
		}
	}()

	cfg := &config.Config{ConfigDir: t.TempDir()}
	certs, err := certmgr.NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	bus := events.NewBus()
	eng := rules.NewEngine()
	srv := NewServer(certs, eng, bus)

	if err := srv.Start(0); err != nil {
		t.Fatalf("Server.Start: %v", err)
	}
	defer srv.Stop()

	// No ws-echo rule — traffic should go to the upstream
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", srv.Port()))
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()

	key := "dGhlIHNhbXBsZSBub25jZQ=="
	upgradeReq := "GET / HTTP/1.1\r\n" +
		"Host: " + upstreamLn.Addr().String() + "\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Key: " + key + "\r\n" +
		"Sec-WebSocket-Version: 13\r\n" +
		"\r\n"
	if _, err := conn.Write([]byte(upgradeReq)); err != nil {
		t.Fatalf("Write upgrade: %v", err)
	}

	br := bufio.NewReader(conn)
	resp, err := http.ReadResponse(br, nil)
	if err != nil {
		t.Fatalf("ReadResponse: %v", err)
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("expected 101, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Send a text frame
	frame := buildWSTextFrame("ping", true)
	if _, err := conn.Write(frame); err != nil {
		t.Fatalf("Write frame: %v", err)
	}

	// Read the echoed response from upstream
	echoBuf := make([]byte, 10)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := br.Read(echoBuf)
	if err != nil {
		t.Fatalf("Read echo: %v", err)
	}
	if n < 2 {
		t.Fatalf("echo frame too short: %d bytes", n)
	}
	payloadLen := int(echoBuf[1] & 0x7F)
	if payloadLen != 4 {
		t.Errorf("payload length = %d, want 4", payloadLen)
	}
}

// TestWebSocketNonUpgradeRequest tests that a non-WebSocket request through
// the WebSocket handler path is handled gracefully.
func TestWebSocketNonUpgradeRequest(t *testing.T) {
	cfg := &config.Config{ConfigDir: t.TempDir()}
	certs, err := certmgr.NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	bus := events.NewBus()
	eng := rules.NewEngine()
	srv := NewServer(certs, eng, bus)

	if err := srv.Start(0); err != nil {
		t.Fatalf("Server.Start: %v", err)
	}
	defer srv.Stop()

	// Send a regular HTTP GET (not a WebSocket upgrade)
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", srv.Port()))
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()

	req := "GET / HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Connection: close\r\n" +
		"\r\n"
	if _, err := conn.Write([]byte(req)); err != nil {
		t.Fatalf("Write: %v", err)
	}

	// Read response — should get some HTTP response (not a crash)
	br := bufio.NewReader(conn)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	resp, err := http.ReadResponse(br, nil)
	if err != nil {
		// Timeout or connection close is acceptable — the proxy may
		// try to connect to example.com and fail, which is fine.
		return
	}
	resp.Body.Close()
}

// computeWebsocketAccept implements the RFC 6455 Sec-WebSocket-Accept computation.
func computeWebsocketAccept(key string) string {
	const magic = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	h := sha1.Sum([]byte(key + magic))
	return base64.StdEncoding.EncodeToString(h[:])
}

// buildWSTextFrame builds a WebSocket text frame with the given payload.
// If masked is true, the frame is client-to-server masked.
func buildWSTextFrame(payload string, masked bool) []byte {
	payloadBytes := []byte(payload)
	frame := []byte{0x81} // FIN + text opcode
	maskBit := byte(0)
	if masked {
		maskBit = 0x80
	}
	if len(payloadBytes) < 126 {
		frame = append(frame, maskBit|byte(len(payloadBytes)))
	} else {
		frame = append(frame, maskBit|126)
		frame = append(frame, byte(len(payloadBytes)>>8), byte(len(payloadBytes)))
	}
	if masked {
		mask := []byte{0x12, 0x34, 0x56, 0x78}
		frame = append(frame, mask...)
		for i, b := range payloadBytes {
			frame = append(frame, b^mask[i%4])
		}
	} else {
		frame = append(frame, payloadBytes...)
	}
	return frame
}

// TestIsWebSocketUpgrade tests the isWebSocketUpgrade detection function.
func TestIsWebSocketUpgrade(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    bool
	}{
		{"standard upgrade", map[string]string{"Upgrade": "websocket", "Connection": "Upgrade"}, true},
		{"case insensitive", map[string]string{"Upgrade": "WebSocket", "Connection": "upgrade"}, true},
		{"connection includes upgrade", map[string]string{"Upgrade": "websocket", "Connection": "keep-alive, Upgrade"}, true},
		{"no upgrade header", map[string]string{"Connection": "Upgrade"}, false},
		{"no connection header", map[string]string{"Upgrade": "websocket"}, false},
		{"wrong upgrade", map[string]string{"Upgrade": "h2c", "Connection": "Upgrade"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			if got := isWebSocketUpgrade(req); got != tt.want {
				t.Errorf("isWebSocketUpgrade = %v, want %v", got, tt.want)
			}
		})
	}
}

// Ensure io is used (for the non-upgrade test's bufio.Reader)
var _ = io.EOF
var _ = strings.Contains
