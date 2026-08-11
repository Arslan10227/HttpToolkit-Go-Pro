package docker

import "testing"

func TestTransformComposeLabels(t *testing.T) {
	labels := map[string]string{
		"com.docker.compose.config-hash": "abc123",
		"com.docker.compose.service":     "web",
	}
	created := TransformComposeCreationLabels(8000, labels)
	if created["com.docker.compose.config-hash"] != "abc123+httptoolkit:8000" {
		t.Fatalf("unexpected create hash: %q", created["com.docker.compose.config-hash"])
	}
	response := TransformComposeResponseLabels(8000, created)
	if response["com.docker.compose.config-hash"] != "abc123" {
		t.Fatalf("unexpected response hash: %q", response["com.docker.compose.config-hash"])
	}
	unintercepted := TransformComposeResponseLabels(8000, labels)
	if unintercepted["com.docker.compose.config-hash"] != "abc123+unintercepted" {
		t.Fatalf("unexpected unintercepted hash: %q", unintercepted["com.docker.compose.config-hash"])
	}
}

func TestTransformContainerCreateConfig(t *testing.T) {
	raw := []byte(`{"Image":"alpine","Env":["FOO=bar"],"Labels":{},"HostConfig":{"Binds":[]}}`)
	out, err := TransformContainerCreateConfig(raw, InterceptionSettings{ProxyPort: 8000})
	if err != nil {
		t.Fatal(err)
	}
	if !containsAll(string(out), "HTTP_PROXY=", "tech.httptoolkit.docker.proxy", "httptoolkit-injection-8000") {
		t.Fatalf("missing injection fields: %s", out)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !contains(s, p) {
			return false
		}
	}
	return true
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func TestProxyHostEnv(t *testing.T) {
	if ProxyHostEnv(8000) == "" {
		t.Fatal("expected proxy host env")
	}
}
