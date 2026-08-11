package docker

import (
	"fmt"
	"sort"
	"strings"
)

const BuildLabel = "tech.httptoolkit.docker.build-proxy"

const contextInjectPath = "/.http-toolkit-injections"
const contextOverridesPath = contextInjectPath + "/overrides"
const contextCAPath = contextInjectPath + "/ca.pem"

type BuildSettings struct {
	ProxyPort   int
	CertPath    string
	CertContent []byte
	AssetsDir   string
}

// InjectIntoDockerfile inserts interception steps after each FROM instruction.
func InjectIntoDockerfile(contents string, proxyPort int, env map[string]string) (string, int) {
	lines := strings.Split(contents, "\n")
	block := dockerfileInjectionBlock(proxyPort, env)
	var out []string
	added := 0
	for _, line := range lines {
		out = append(out, line)
		if isDockerFromLine(line) {
			out = append(out, block...)
			added += len(block)
		}
	}
	return strings.Join(out, "\n"), added
}

func isDockerFromLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return false
	}
	upper := strings.ToUpper(trimmed)
	return strings.HasPrefix(upper, "FROM ")
}

func dockerfileInjectionBlock(proxyPort int, env map[string]string) []string {
	block := []string{
		fmt.Sprintf("LABEL %s=started-%d", BuildLabel, proxyPort),
		fmt.Sprintf("COPY %s %s", contextInjectPath, injectedPath),
	}
	if envLine := formatDockerEnvLine(env); envLine != "" {
		block = append(block, envLine)
	}
	block = append(block, fmt.Sprintf("LABEL %s=%d", BuildLabel, proxyPort))
	return block
}

func formatDockerEnvLine(env map[string]string) string {
	if len(env) == 0 {
		return ""
	}
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%q", k, env[k]))
	}
	return "ENV " + strings.Join(parts, " ")
}

func buildInjectionEnv(proxyPort int) map[string]string {
	host := DockerHostAddress("")
	proxyURL := fmt.Sprintf("http://%s:%d", host, proxyPort)
	return injectEnvVars(proxyURL, injectedCAPath)
}
