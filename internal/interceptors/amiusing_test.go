package interceptors

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	certmgr "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/events"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/mitm"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/rules"
)

func TestAmiusingServer(t *testing.T) {
	cfg := &config.Config{ConfigDir: t.TempDir()}
	certs, err := certmgr.NewManager(cfg)
	if err != nil {
		t.Fatalf("cert manager: %v", err)
	}

	// Start a real MITM proxy so we can verify the amiusing check end-to-end.
	bus := events.NewBus()
	eng := rules.NewEngine()
	proxy := mitm.NewServer(certs, eng, bus)
	if err := proxy.Start(0); err != nil {
		t.Fatalf("proxy start: %v", err)
	}
	defer func() { _ = proxy.Stop() }()

	srv, err := startAmiusingServer(certs)
	if err != nil {
		t.Fatalf("startAmiusingServer: %v", err)
	}
	defer srv.Stop()

	// The HTTP page loads through the proxy as well (Chrome loads it as the
	// initial tab and the proxy passes it through to the amiusing server).
	proxyURL, _ := url.Parse(fmt.Sprintf("https://127.0.0.1:%d", proxy.Port()))
	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM([]byte(certs.CertPEM())); !ok {
		t.Fatal("failed to load CA into pool")
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{RootCAs: caPool, MinVersion: tls.VersionTLS12},
		},
		Timeout: 10 * time.Second,
	}

	// 1. HTTP page should load.
	pageResp, err := httpClient.Get(srv.URL())
	if err != nil {
		t.Fatalf("HTTP amiusing page: %v", err)
	}
	body, _ := io.ReadAll(pageResp.Body)
	_ = pageResp.Body.Close()
	if pageResp.StatusCode != http.StatusOK {
		t.Fatalf("HTTP page status: %d, body: %s", pageResp.StatusCode, body)
	}
	if !strings.Contains(string(body), "Checking interception") {
		t.Fatalf("unexpected page body: %s", body)
	}

	// 2. HTTPS test endpoint should succeed, proving interception works.
	testResp, err := httpClient.Get(srv.HTTPSURL())
	if err != nil {
		t.Fatalf("HTTPS intercepted endpoint: %v", err)
	}
	testBody, _ := io.ReadAll(testResp.Body)
	_ = testResp.Body.Close()
	if testResp.StatusCode != http.StatusOK {
		t.Fatalf("HTTPS test status: %d, body: %s", testResp.StatusCode, testBody)
	}
	if string(testBody) != "intercepted" {
		t.Fatalf("unexpected test body: %q", testBody)
	}

	// 3. The success should be observable via WaitSuccess.
	if err := srv.WaitSuccess(5 * time.Second); err != nil {
		t.Fatalf("WaitSuccess: %v", err)
	}
}
