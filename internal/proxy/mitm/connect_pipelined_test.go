package mitm

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	certmgr "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/events"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/rules"
)

// TestPipelinedConnectAndTLS verifies that the MITM proxy correctly handles
// the case where a client sends a CONNECT request and the inner TLS
// ClientHello in a single TCP write (pipelining). The bufio.Reader in
// handleHTTP may read-ahead past the CONNECT request line into the TLS bytes;
// those buffered bytes must be passed to the inner TLS handshake via bufConn,
// not lost or re-read from the underlying connection.
//
// This is a regression test for the bug where the inner TLS handshake would
// fail because the first byte(s) of the ClientHello were consumed by the
// bufio.Reader and never forwarded to tls.Server.
func TestPipelinedConnectAndTLS(t *testing.T) {
	// 1. Setup a target HTTPS server (the real upstream we intercept).
	targetServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test-Response", "Yes")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello from HTTPS target!"))
	}))
	defer targetServer.Close()

	targetURL, _ := url.Parse(targetServer.URL)
	targetHost, targetPort, _ := net.SplitHostPort(targetURL.Host)

	// 2. Setup cert manager and MITM proxy.
	cfg := &config.Config{ConfigDir: t.TempDir()}
	certs, err := certmgr.NewManager(cfg)
	if err != nil {
		t.Fatalf("failed to create cert manager: %v", err)
	}

	bus := events.NewBus()
	eng := rules.NewEngine()
	s := NewServer(certs, eng, bus)

	if err := s.Start(0); err != nil {
		t.Fatalf("failed to start proxy: %v", err)
	}
	defer func() { _ = s.Stop() }()

	proxyAddr := fmt.Sprintf("127.0.0.1:%d", s.Port())

	// 3. Open a raw TCP connection to the proxy.
	rawConn, err := net.DialTimeout("tcp", proxyAddr, 5*time.Second)
	if err != nil {
		t.Fatalf("failed to dial proxy: %v", err)
	}
	defer rawConn.Close()

	// 4. Write the CONNECT request directly to the raw connection.
	connectReq := fmt.Sprintf("CONNECT %s:%s HTTP/1.1\r\nHost: %s:%s\r\n\r\n",
		targetHost, targetPort, targetHost, targetPort)
	if _, err := rawConn.Write([]byte(connectReq)); err != nil {
		t.Fatalf("failed to write CONNECT: %v", err)
	}

	// 5. Immediately start the TLS handshake in a goroutine WITHOUT waiting
	// for the 200 response. This simulates pipelining: the client sends the
	// TLS ClientHello right after the CONNECT request, and the proxy's
	// bufio.Reader may read-ahead and buffer the ClientHello bytes.
	// Use a shared buffered reader with a mutex so we can read the 200
	// response in the main goroutine and then let the TLS handshake read
	// the ServerHello from the same reader.
	sharedReader := bufio.NewReader(rawConn)
	var readMu sync.Mutex

	// We use InsecureSkipVerify because the leaf cert is generated for the
	// target host (127.0.0.1) but the httptest TLS server's cert is different.
	// The pipelining test is about verifying the TLS handshake completes,
	// not about cert verification (that's tested in the cert package).
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}

	tlsConn := tls.Client(&syncReaderConn{r: sharedReader, mu: &readMu, Conn: rawConn}, tlsConfig)
	defer tlsConn.Close()

	handshakeErr := make(chan error, 1)
	go func() {
		handshakeErr <- tlsConn.Handshake()
	}()

	// 6. Read the 200 Connection Established response. This must happen
	// before the TLS handshake can complete (the proxy sends 200 first,
	// then the ServerHello). The shared mutex ensures the TLS handshake
	// goroutine doesn't read the 200 response bytes.
	readMu.Lock()
	resp, err := http.ReadResponse(sharedReader, &http.Request{Method: http.MethodConnect})
	readMu.Unlock()
	if err != nil {
		t.Fatalf("failed to read CONNECT response: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 Connection Established, got %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	// 7. Wait for the TLS handshake to complete. If the proxy's bufConn
	// doesn't properly pass the buffered ClientHello bytes to tls.Server,
	// this will fail with a handshake error.
	select {
	case err := <-handshakeErr:
		if err != nil {
			t.Fatalf("TLS handshake failed (pipelined ClientHello was likely lost): %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("TLS handshake timed out")
	}

	// 8. Send an HTTP request over the established TLS connection and verify
	// the response comes back through the MITM proxy.
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet,
		fmt.Sprintf("https://%s:%s/pipelined-test", targetHost, targetPort), nil)
	if err := req.Write(tlsConn); err != nil {
		t.Fatalf("failed to write HTTP request over TLS: %v", err)
	}

	httpResp, err := http.ReadResponse(bufio.NewReader(tlsConn), req)
	if err != nil {
		t.Fatalf("failed to read HTTP response over MITM TLS: %v", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", httpResp.StatusCode)
	}
	body, _ := io.ReadAll(httpResp.Body)
	if !strings.Contains(string(body), "Hello from HTTPS target") {
		t.Errorf("expected response body from target, got: %q", body)
	}
	if httpResp.Header.Get("X-Test-Response") != "Yes" {
		t.Errorf("expected X-Test-Response header, got headers: %v", httpResp.Header)
	}
}

// syncReaderConn wraps a bufio.Reader + net.Conn with a mutex so that reads
// from the shared reader are safe across goroutines. This allows us to read
// the CONNECT response and then hand the same buffered reader to tls.Client.
type syncReaderConn struct {
	r  *bufio.Reader
	mu *sync.Mutex
	net.Conn
}

func (s *syncReaderConn) Read(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.r.Read(p)
}
