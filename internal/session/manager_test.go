package session

import "testing"

func TestResolveRuleParamProxyDockerTunnel(t *testing.T) {
	m := NewManager()
	m.SetActive(8000, nil, []string{"docker-tunnel-proxy-8000"}, false, false)
	m.SetDockerTunnelPort(8000, 19000)
	m.SetDockerHostCheck(func(host string) bool { return host == "mycontainer" })

	if got := m.ResolveRuleParamProxy("docker-tunnel-proxy-8000", "mycontainer"); got != "socks5://127.0.0.1:19000" {
		t.Fatalf("got %q", got)
	}
	if got := m.ResolveRuleParamProxy("docker-tunnel-proxy-8000", "example.com"); got != "" {
		t.Fatalf("expected empty for non-docker host, got %q", got)
	}
}
