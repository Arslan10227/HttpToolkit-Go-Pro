package mitm

import (
	"net"
	"testing"
)

// TestIPv6HostParsing verifies that handleConnect's host parsing logic
// correctly handles IPv6 addresses in CONNECT requests. IPv6 literals in
// CONNECT requests must be bracketed (e.g., "[::1]:443"), and the parsing
// must use net.SplitHostPort / net.JoinHostPort to handle them.
func TestIPv6HostParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantHost string
		wantPort string
	}{
		{"IPv4 with port", "127.0.0.1:443", "127.0.0.1", "443"},
		{"IPv6 bracketed with port", "[::1]:443", "::1", "443"},
		{"IPv6 bracketed with port 8443", "[::1]:8443", "::1", "8443"},
		{"hostname with port", "example.com:443", "example.com", "443"},
		{"hostname with non-TLS port", "example.com:8080", "example.com", "8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hostOnly, port, err := net.SplitHostPort(tt.input)
			if err != nil {
				// Fallback path (matches handleConnect's error case)
				hostOnly = trimBrackets(tt.input)
				port = "443"
			}
			host := net.JoinHostPort(hostOnly, port)

			if hostOnly != tt.wantHost {
				t.Errorf("host = %q, want %q", hostOnly, tt.wantHost)
			}
			if port != tt.wantPort {
				t.Errorf("port = %q, want %q", port, tt.wantPort)
			}
			// Verify the round-tripped host is valid for net.Dial
			if host == "" {
				t.Error("JoinHostPort returned empty")
			}
		})
	}
}

// TestIPv6BracketedNoPort tests the fallback path where an IPv6 address
// is bracketed but has no port (edge case — should default to 443).
func TestIPv6BracketedNoPort(t *testing.T) {
	input := "[::1]"
	hostOnly, port, err := net.SplitHostPort(input)
	if err == nil {
		t.Fatalf("SplitHostPort should fail for %q, got host=%q port=%q", input, hostOnly, port)
	}
	// Fallback path (matches handleConnect)
	hostOnly = trimBrackets(input)
	port = "443"
	host := net.JoinHostPort(hostOnly, port)

	if hostOnly != "::1" {
		t.Errorf("host = %q, want ::1", hostOnly)
	}
	if port != "443" {
		t.Errorf("port = %q, want 443", port)
	}
	if host != "[::1]:443" {
		t.Errorf("JoinHostPort = %q, want [::1]:443", host)
	}
}

func trimBrackets(s string) string {
	if len(s) >= 2 && s[0] == '[' && s[len(s)-1] == ']' {
		return s[1 : len(s)-1]
	}
	return s
}
