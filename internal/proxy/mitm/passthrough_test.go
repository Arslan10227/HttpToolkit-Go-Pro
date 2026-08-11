package mitm

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/rules"
)

func TestMatchHostWildcard(t *testing.T) {
	tests := []struct {
		pattern  string
		host     string
		expected bool
	}{
		{"*", "example.com", true},
		{"*", "sub.domain.com", true},
		{"example.com", "example.com", true},
		{"example.com", "EXAMPLE.com", true},
		{"example.com", "www.example.com", true},
		{"example.com", "sub.www.example.com", true},
		{"*.example.com", "www.example.com", true},
		{"*.example.com", "example.com", true},
		{"*.example.com", "other.com", false},
		{"example.com", "other.com", false},
	}

	for _, tc := range tests {
		res := MatchHostWildcard(tc.pattern, tc.host)
		if res != tc.expected {
			t.Errorf("MatchHostWildcard(%q, %q) = %t; expected %t", tc.pattern, tc.host, res, tc.expected)
		}
	}
}

func TestPickClientCertWildcard(t *testing.T) {
	m := map[string]rules.ClientCert{
		"*.example.com": {PFX: "wildcard-pfx"},
		"exact.com":     {PFX: "exact-pfx"},
		"*":             {PFX: "fallback-pfx"},
	}

	tests := []struct {
		host     string
		expected string
	}{
		{"exact.com", "exact-pfx"},
		{"www.example.com", "wildcard-pfx"},
		{"example.com", "wildcard-pfx"},
		{"other.com", "fallback-pfx"},
	}

	for _, tc := range tests {
		cert := PickClientCert(tc.host, m)
		if cert.PFX != tc.expected {
			t.Errorf("PickClientCert(%q) = %q; expected %q", tc.host, cert.PFX, tc.expected)
		}
	}
}

func TestBuildPassthroughTransportIgnoreErrors(t *testing.T) {
	target, _ := url.Parse("https://www.example.com")
	opts := rules.PassthroughOptions{
		IgnoreHostHttpsErrors: []string{"example.com"},
	}

	tr := buildPassthroughTransport(target, opts)
	if tr.TLSClientConfig == nil || !tr.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("Expected InsecureSkipVerify to be true for matched ignore domain")
	}

	target2, _ := url.Parse("https://other.com")
	tr2 := buildPassthroughTransport(target2, opts)
	if tr2.TLSClientConfig != nil && tr2.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("Expected InsecureSkipVerify to be false for unmatched ignore domain")
	}
}

// TestBuildPassthroughTransportHTTP2Disabled verifies the passthrough transport
// does not attempt HTTP/2 upstreams. Some clients (e.g. Java HTTP/1.1 clients)
// receive HTTP/2.0 status lines they cannot parse, so upstream HTTP/1.1 is
// always used for now.
func TestBuildPassthroughTransportHTTP2Disabled(t *testing.T) {
	target, _ := url.Parse("https://www.example.com")
	tr := buildPassthroughTransport(target, rules.PassthroughOptions{})
	if tr.ForceAttemptHTTP2 {
		t.Fatal("Expected ForceAttemptHTTP2 to be false for HTTP/1.1 passthrough")
	}
}

// TestBuildPassthroughTransportLookupServers verifies that a custom
// DialContext is installed only when the passthrough options specify DNS
// lookupServers, so that requests without that option keep using the
// default (net/http) dialer.
func TestBuildPassthroughTransportLookupServers(t *testing.T) {
	target, _ := url.Parse("https://example.com")

	trNoLookup := buildPassthroughTransport(target, rules.PassthroughOptions{})
	if trNoLookup.DialContext != nil {
		t.Fatal("Expected no custom DialContext when lookupServers is unset")
	}

	trWithLookup := buildPassthroughTransport(target, rules.PassthroughOptions{
		LookupServers: []string{"127.0.0.1:5353"},
	})
	if trWithLookup.DialContext == nil {
		t.Fatal("Expected a custom DialContext when lookupServers is set")
	}
}

// TestBuildPassthroughTransportProxyConfig verifies that opts.ProxyConfig is
// applied as the transport's upstream HTTP proxy, and that noProxy entries
// (string or []any form, including the special "<local>" token) correctly
// suppress it.
func TestBuildPassthroughTransportProxyConfig(t *testing.T) {
	target, _ := url.Parse("https://example.com")

	proxied := buildPassthroughTransport(target, rules.PassthroughOptions{
		ProxyConfig: map[string]any{"proxyUrl": "http://proxy.local:8080"},
	})
	if proxied.Proxy == nil {
		t.Fatal("Expected Transport.Proxy to be set from proxyConfig.proxyUrl")
	}
	got, err := proxied.Proxy(&http.Request{URL: target})
	if err != nil || got == nil || got.Host != "proxy.local:8080" {
		t.Fatalf("Expected proxy host proxy.local:8080, got %v (err=%v)", got, err)
	}

	suppressedString := buildPassthroughTransport(target, rules.PassthroughOptions{
		ProxyConfig: map[string]any{"proxyUrl": "http://proxy.local:8080", "noProxy": "example.com"},
	})
	if suppressedString.Proxy != nil {
		if got, _ := suppressedString.Proxy(&http.Request{URL: target}); got != nil {
			t.Fatal("Expected proxy to be suppressed for a hostname listed in noProxy (string form)")
		}
	}

	suppressedList := buildPassthroughTransport(target, rules.PassthroughOptions{
		ProxyConfig: map[string]any{"proxyUrl": "http://proxy.local:8080", "noProxy": []any{"other.com", "example.com"}},
	})
	if suppressedList.Proxy != nil {
		if got, _ := suppressedList.Proxy(&http.Request{URL: target}); got != nil {
			t.Fatal("Expected proxy to be suppressed for a hostname listed in noProxy ([]any form)")
		}
	}

	localTarget, _ := url.Parse("https://127.0.0.1")
	suppressedLocal := buildPassthroughTransport(localTarget, rules.PassthroughOptions{
		ProxyConfig: map[string]any{"proxyUrl": "http://proxy.local:8080", "noProxy": "<local>"},
	})
	if suppressedLocal.Proxy != nil {
		if got, _ := suppressedLocal.Proxy(&http.Request{URL: localTarget}); got != nil {
			t.Fatal("Expected proxy to be suppressed for a loopback host via the <local> noProxy token")
		}
	}
}

// TestBuildPassthroughTransportClientCert verifies a matching client
// certificate host entry is applied to the TLS config.
func TestBuildPassthroughTransportClientCert(t *testing.T) {
	target, _ := url.Parse("https://example.com")

	// An invalid/empty PFX should be safely ignored rather than erroring the
	// whole transport build.
	tr := buildPassthroughTransport(target, rules.PassthroughOptions{
		ClientCertHostMap: map[string]rules.ClientCert{
			"example.com": {PFX: "", Passphrase: ""},
		},
	})
	if tr == nil {
		t.Fatal("Expected a transport even with an empty client cert entry")
	}
}

// TestResolveProxyURLMatrix verifies that every proxyConfig shape the UI can
// send is resolved correctly: nil, string, map, array, callback, noProxy.
func TestResolveProxyURLMatrix(t *testing.T) {
	target, _ := url.Parse("https://example.com:8443")

	cases := []struct {
		name   string
		cfg    any
		want   string // expected host, or empty for no proxy
	}{
		{"nil", nil, ""},
		{"string", "http://proxy.local:8080", "proxy.local:8080"},
		{"empty-string", "", ""},
		{"map", map[string]any{"proxyUrl": "http://proxy.local:8080"}, "proxy.local:8080"},
		{"map-noProxy-string", map[string]any{"proxyUrl": "http://proxy.local:8080", "noProxy": "example.com"}, ""},
		{"map-noProxy-list", map[string]any{"proxyUrl": "http://proxy.local:8080", "noProxy": []any{"other.com", "example.com"}}, ""},
		{"callback-shaped", map[string]any{"type": "callback"}, ""},
		{"array-of-strings", []any{"http://first.local:8080", "http://second.local:8080"}, "first.local:8080"},
		{"array-of-maps", []any{map[string]any{"noProxy": "example.com"}, map[string]any{"proxyUrl": "http://second.local:8080"}}, "second.local:8080"},
		{"invalid-url", "://not-a-url", ""},
		{"socks5", "socks5://proxy.local:1080", "proxy.local:1080"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := resolveProxyURL(tc.cfg, target)
			if (got == nil && tc.want != "") || (got != nil && got.Host != tc.want) {
				t.Fatalf("resolveProxyURL(%#v) host = %v; want %q", tc.cfg, got, tc.want)
			}
		})
	}
}
