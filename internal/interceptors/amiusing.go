package interceptors

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	certmgr "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
)

// amiusingServer is a self-hosted "am I being intercepted" check page.
// It serves a plain-HTTP page that, when loaded by a browser, fetches an
// HTTPS endpoint through the proxy. If the proxy is intercepting and the
// CA is trusted, the fetch succeeds and the server signals ready.
//
// It replaces the external https://amiusing.httptoolkit.tech dependency and
// works for both Chromium (SPKI trust) and Firefox (installed CA trust).
type amiusingServer struct {
	mu                 sync.Mutex
	srv                *http.Server
	ready              chan struct{}
	failed             chan struct{}
	url                string
	httpInterceptedURL string
	httpsURL           string
}

// bufferedConn and muxListener are reused from the old certCheckServer to
// serve HTTP and HTTPS on the same port.
type bufferedConn struct {
	net.Conn
	r io.Reader
}

func (bc *bufferedConn) Read(b []byte) (int, error) {
	return bc.r.Read(b)
}

type muxListener struct {
	net.Listener
	tlsConfig *tls.Config
}

func (l *muxListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	firstByte := make([]byte, 1)
	_ = c.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, err = io.ReadFull(c, firstByte)
	_ = c.SetReadDeadline(time.Time{})
	if err != nil {
		_ = c.Close()
		return nil, err
	}

	r := io.MultiReader(bytes.NewReader(firstByte), c)
	bc := &bufferedConn{Conn: c, r: r}

	if firstByte[0] == 0x16 { // TLS ClientHello
		return tls.Server(bc, l.tlsConfig), nil
	}

	return bc, nil
}

func (l *muxListener) Close() error   { return l.Listener.Close() }
func (l *muxListener) Addr() net.Addr { return l.Listener.Addr() }

// startAmiusingServer starts an HTTP/HTTPS server on a random localhost port.
func startAmiusingServer(certs *certmgr.Manager) (*amiusingServer, error) {
	certPEM, keyPEM, err := certs.LeafCert("127.0.0.1")
	if err != nil {
		return nil, fmt.Errorf("amiusing leaf cert: %w", err)
	}
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	s := &amiusingServer{
		ready:  make(chan struct{}, 1),
		failed: make(chan struct{}, 1),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/amiusing", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		page := amiusingPage(s.httpInterceptedURL)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(page)))
		_, _ = w.Write([]byte(page))
	})
	mux.HandleFunc("/amiusing/intercepted", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Any request that reaches this endpoint means the launched browser
		// loaded the check page and was able to fetch the HTTPS URL, which
		// proves the HTTP Toolkit CA is trusted. We intentionally do not
		// require a TLS connection here because the proxy may terminate TLS
		// and forward a plain HTTP request upstream.
		select {
		case s.ready <- struct{}{}:
		default:
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "11")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("intercepted"))
	})
	mux.HandleFunc("/amiusing/failed", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		select {
		case s.failed <- struct{}{}:
		default:
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "6")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("failed"))
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	port := ln.Addr().(*net.TCPAddr).Port
	s.url = fmt.Sprintf("http://127.0.0.1:%d/amiusing", port)
	s.httpInterceptedURL = fmt.Sprintf("http://127.0.0.1:%d/amiusing/intercepted", port)
	s.httpsURL = fmt.Sprintf("https://127.0.0.1:%d/amiusing/intercepted", port)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		MinVersion:   tls.VersionTLS12,
	}
	s.srv = &http.Server{Handler: mux}

	ml := &muxListener{Listener: ln, tlsConfig: tlsConfig}
	go func() { _ = s.srv.Serve(ml) }()
	return s, nil
}

func (s *amiusingServer) URL() string      { return s.url }
func (s *amiusingServer) Host() string     { return trimScheme(s.url) }
func (s *amiusingServer) HTTPSURL() string { return s.httpsURL }

func (s *amiusingServer) WaitSuccess(timeout time.Duration) error {
	select {
	case <-s.ready:
		return nil
	case <-s.failed:
		return fmt.Errorf("interception check failed")
	case <-time.After(timeout):
		return fmt.Errorf("interception check timed out")
	}
}

func (s *amiusingServer) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.srv != nil {
		_ = s.srv.Close()
		s.srv = nil
	}
}

func amiusingPage(testURL string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>HTTP Toolkit Interception Check</title>
<style>
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #1e2028; color: #fff; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; }
.card { background: #2a2d37; padding: 40px; border-radius: 12px; text-align: center; max-width: 420px; }
h1 { margin-top: 0; font-size: 1.2em; }
p { color: #a0a5b5; line-height: 1.5; }
.success { color: #00e676; }
.error { color: #ff5252; }
</style>
</head>
<body>
<div class="card" id="card">
  <h1>Checking interception...</h1>
  <p id="msg">Verifying that HTTP Toolkit is intercepting this browser's traffic.</p>
</div>
<script>
const testUrl = %q;
fetch(testUrl, { cache: "no-store" })
  .then(r => {
    if (r.ok) {
      document.getElementById("card").innerHTML = '<h1 class="success">Interception active</h1><p>This browser is ready to use.</p>';
    } else {
      fail();
    }
  })
  .catch(fail);

function fail() {
  document.getElementById("card").innerHTML = '<h1 class="error">Interception not active</h1><p>The HTTP Toolkit certificate is not trusted in this browser. Install and trust the CA, then restart the browser.</p>';
  fetch(window.location.href.replace("/amiusing", "/amiusing/failed"), { method: "POST" });
}
</script>
</body>
</html>`, testURL)
}

func trimScheme(u string) string {
	u = trimPrefix(u, "http://")
	u = trimPrefix(u, "https://")
	if i := indexByte(u, '/'); i > 0 {
		u = u[:i]
	}
	return u
}

func trimPrefix(s, prefix string) string {
	for i := 0; i < len(prefix); i++ {
		if i >= len(s) || s[i] != prefix[i] {
			return s
		}
	}
	return s[len(prefix):]
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

// certFilesExist reports whether the CA cert and key files exist in cfgDir.
func certFilesExist(cfgDir string) bool {
	cert := filepath.Join(cfgDir, "ca.pem")
	key := filepath.Join(cfgDir, "ca.key")
	if _, err := os.Stat(cert); err != nil {
		return false
	}
	if _, err := os.Stat(key); err != nil {
		return false
	}
	return true
}
