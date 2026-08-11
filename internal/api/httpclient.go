package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"
)

type sendOptions struct {
	IgnoreHostHttpsErrors []string
	ClientCertHostMap     map[string]clientCert
	ProxyConfig           any
	LookupServers         []string
	TrustedCAs            []string
}

type clientCert struct {
	PFX        []byte
	Passphrase string
}

func parseSendOptions(raw map[string]any) sendOptions {
	opts := sendOptions{
		IgnoreHostHttpsErrors: stringSlice(raw["ignoreHostHttpsErrors"]),
	}
	if m, ok := raw["clientCertificateHostMap"].(map[string]any); ok {
		opts.ClientCertHostMap = parseClientCertMap(m)
	}
	opts.ProxyConfig = raw["proxyConfig"]
	if lo, ok := raw["lookupOptions"].(map[string]any); ok {
		opts.LookupServers = stringSlice(lo["servers"])
	}
	if cas, ok := raw["trustAdditionalCAs"].([]any); ok {
		for _, c := range cas {
			if cm, ok := c.(map[string]any); ok {
				if cert, ok := cm["cert"].(string); ok {
					opts.TrustedCAs = append(opts.TrustedCAs, cert)
				}
			}
		}
	}
	if cas, ok := raw["additionalTrustedCAs"].([]any); ok {
		for _, c := range cas {
			if cm, ok := c.(map[string]any); ok {
				if cert, ok := cm["cert"].(string); ok {
					opts.TrustedCAs = append(opts.TrustedCAs, cert)
				}
			}
		}
	}
	return opts
}

func parseClientCertMap(m map[string]any) map[string]clientCert {
	out := make(map[string]clientCert, len(m))
	for host, v := range m {
		cm, ok := v.(map[string]any)
		if !ok {
			continue
		}
		pfxB64, _ := cm["pfx"].(string)
		pfx, _ := base64.StdEncoding.DecodeString(pfxB64)
		pass, _ := cm["passphrase"].(string)
		out[host] = clientCert{PFX: pfx, Passphrase: pass}
	}
	return out
}

func stringSlice(v any) []string {
	switch s := v.(type) {
	case []string:
		return s
	case []any:
		out := make([]string, 0, len(s))
		for _, item := range s {
			if str, ok := item.(string); ok {
				out = append(out, str)
			}
		}
		return out
	default:
		return nil
	}
}

func buildSendTransport(targetURL *url.URL, opts sendOptions) (*http.Transport, error) {
	tlsCfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if len(opts.TrustedCAs) > 0 {
		pool, err := x509.SystemCertPool()
		if err != nil || pool == nil {
			pool = x509.NewCertPool()
		}
		for _, pem := range opts.TrustedCAs {
			pool.AppendCertsFromPEM([]byte(pem))
		}
		tlsCfg.RootCAs = pool
	}
	host := targetURL.Hostname()
	for _, h := range opts.IgnoreHostHttpsErrors {
		if h == "*" || strings.EqualFold(h, host) {
			tlsCfg.InsecureSkipVerify = true
			break
		}
	}
	if cert := pickClientCert(host, opts.ClientCertHostMap); len(cert.PFX) > 0 {
		if pair, err := clientCertFromPFX(cert.PFX, cert.Passphrase); err == nil && pair != nil {
			tlsCfg.Certificates = []tls.Certificate{*pair}
		}
	}

	transport := &http.Transport{
		TLSClientConfig:       tlsCfg,
		ExpectContinueTimeout: time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
	}

	if len(opts.LookupServers) > 0 {
		transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 30 * time.Second}
			host, port, err := net.SplitHostPort(address)
			if err != nil {
				return d.DialContext(ctx, network, address)
			}
			ips, err := lookupViaServers(host, opts.LookupServers)
			if err != nil || len(ips) == 0 {
				return d.DialContext(ctx, network, address)
			}
			return d.DialContext(ctx, network, net.JoinHostPort(ips[0], port))
		}
	}

	if proxyURL := resolveProxyURL(opts.ProxyConfig, targetURL); proxyURL != nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	return transport, nil
}

func pickClientCert(host string, m map[string]clientCert) clientCert {
	if c, ok := m[host]; ok {
		return c
	}
	if c, ok := m["*"]; ok {
		return c
	}
	return clientCert{}
}

func resolveProxyURL(cfg any, target *url.URL) *url.URL {
	switch v := cfg.(type) {
	case nil:
		return nil
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
	list := stringSlice(raw)
	if len(list) == 0 {
		return false
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

func lookupViaServers(host string, servers []string) ([]string, error) {
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

func clientCertFromPFX(pfx []byte, passphrase string) (*tls.Certificate, error) {
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

func normalizeSendHeaders(raw any) (map[string][]string, error) {
	switch h := raw.(type) {
	case map[string]any:
		out := make(map[string][]string)
		for k, v := range h {
			out[k] = stringSlice(v)
		}
		return out, nil
	case map[string][]string:
		return h, nil
	case []any:
		out := make(map[string][]string)
		for _, pair := range h {
			tuple, ok := pair.([]any)
			if !ok || len(tuple) < 2 {
				continue
			}
			k, _ := tuple[0].(string)
			v, _ := tuple[1].(string)
			if k != "" {
				out[k] = append(out[k], v)
			}
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported headers shape")
	}
}
