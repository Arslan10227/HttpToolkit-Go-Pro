package mitm

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/rules"
	"golang.org/x/crypto/pkcs12"
)

func MatchHostWildcard(pattern, host string) bool {
	if pattern == "*" {
		return true
	}
	pattern = strings.ToLower(pattern)
	host = strings.ToLower(host)
	if pattern == host {
		return true
	}
	if strings.HasPrefix(pattern, "*.") {
		suffix := pattern[1:]
		return strings.HasSuffix(host, suffix) || host == pattern[2:]
	}
	return host == pattern || strings.HasSuffix(host, "."+pattern)
}

func buildPassthroughTransport(target *url.URL, opts rules.PassthroughOptions) *http.Transport {
	tlsCfg := &tls.Config{MinVersion: tls.VersionTLS12}
	host := target.Hostname()
	for _, h := range opts.IgnoreHostHttpsErrors {
		if MatchHostWildcard(h, host) {
			tlsCfg.InsecureSkipVerify = true
			break
		}
	}
	// Intercepted self-hosted services (e.g. the amiusing check server) run on
	// loopback with certificates signed by the private CA. Skip verification for
	// those to avoid chicken-and-egg problems where the proxy has to trust itself.
	if !tlsCfg.InsecureSkipVerify && isLoopbackHost(host) {
		tlsCfg.InsecureSkipVerify = true
	}
	if cert := PickClientCert(host, opts.ClientCertHostMap); cert.PFX != "" {
		pfx, _ := base64.StdEncoding.DecodeString(cert.PFX)
		if pair, err := ClientCertFromPFX(pfx, cert.Passphrase); err == nil && pair != nil {
			tlsCfg.Certificates = []tls.Certificate{*pair}
		}
	}

	transport := &http.Transport{
		TLSClientConfig:       tlsCfg,
		ForceAttemptHTTP2:     false, // Downgrade upstream to HTTP/1.1 so HTTP/1.1 clients don't see HTTP/2.0 status lines
		ExpectContinueTimeout: time.Second,
		ResponseHeaderTimeout: 120 * time.Second,
	}

	if len(opts.LookupServers) > 0 {
		servers := append([]string(nil), opts.LookupServers...)
		transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 30 * time.Second}
			h, port, err := net.SplitHostPort(address)
			if err != nil {
				return d.DialContext(ctx, network, address)
			}
			ips, err := LookupViaServers(h, servers)
			if err != nil || len(ips) == 0 {
				return d.DialContext(ctx, network, address)
			}
			return d.DialContext(ctx, network, net.JoinHostPort(ips[0], port))
		}
	}

	if proxyURL := resolveProxyURL(opts.ProxyConfig, target); proxyURL != nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	}
	return transport
}

func PickClientCert(host string, m map[string]rules.ClientCert) rules.ClientCert {
	host = strings.ToLower(host)
	if c, ok := m[host]; ok {
		return c
	}
	for pattern, c := range m {
		pattern = strings.ToLower(pattern)
		if pattern == "*" {
			continue
		}
		if MatchHostWildcard(pattern, host) {
			return c
		}
	}
	if c, ok := m["*"]; ok {
		return c
	}
	return rules.ClientCert{}
}

func ClientCertFromPFX(pfx []byte, passphrase string) (*tls.Certificate, error) {
	if len(pfx) == 0 {
		return nil, nil
	}
	key, cert, err := pkcs12.Decode(pfx, passphrase)
	if err != nil {
		return nil, err
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
	pair, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}
	return &pair, nil
}

func resolveProxyURL(cfg any, target *url.URL) *url.URL {
	switch v := cfg.(type) {
	case nil:
		return nil
	case string:
		if v == "" {
			return nil
		}
		pu, err := url.Parse(v)
		if err != nil {
			return nil
		}
		return pu
	case map[string]any:
		u, _ := v["proxyUrl"].(string)
		if u == "" {
			return nil
		}
		pu, err := url.Parse(u)
		if err != nil {
			return nil
		}
		if noProxyMatches(v["noProxy"], target.Hostname()) {
			return nil
		}
		return pu
	case []any:
		for _, item := range v {
			if u := resolveProxyURL(item, target); u != nil {
				return u
			}
		}
	}
	return nil
}

func noProxyMatches(raw any, host string) bool {
	var list []string
	switch s := raw.(type) {
	case string:
		list = strings.Split(s, ",")
	case []any:
		for _, item := range s {
			if str, ok := item.(string); ok {
				list = append(list, str)
			}
		}
	}
	for _, entry := range list {
		entry = strings.TrimSpace(entry)
		if entry == "" || entry == "<local>" {
			if host == "localhost" || strings.HasPrefix(host, "127.") {
				return true
			}
			continue
		}
		if strings.EqualFold(entry, host) || strings.HasSuffix(host, "."+entry) {
			return true
		}
	}
	return false
}

func isLoopbackHost(host string) bool {
	host = strings.ToLower(host)
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func LookupViaServers(host string, servers []string) ([]string, error) {
	if len(servers) == 0 {
		return net.LookupHost(host)
	}
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 5 * time.Second}
			return d.DialContext(ctx, "udp", servers[0])
		},
	}
	return r.LookupHost(context.Background(), host)
}

func structuredError(msg string, code string) map[string]any {
	if code == "" {
		code = "ERROR"
	}
	return map[string]any{"message": msg, "code": code}
}

func httpStatusMessage(code int) string {
	return http.StatusText(code)
}

func eventTiming(start time.Time) map[string]any {
	return map[string]any{
		"timestamp": time.Now().UnixMilli(),
		"timingEvents": map[string]any{
			"startTime": start.UnixMilli(),
		},
	}
}

func eventTimingWithResponse(start time.Time, end time.Time) map[string]any {
	return map[string]any{
		"timestamp": end.UnixMilli(),
		"timingEvents": map[string]any{
			"startTime":       start.UnixMilli(),
			"responseBodyEnd": end.UnixMilli(),
		},
	}
}
