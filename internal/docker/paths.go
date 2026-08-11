package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// ProxyListenPath returns the Docker API proxy socket/pipe path for a proxy port.
func ProxyListenPath(proxyPort int) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf(`\\.\pipe\httptoolkit-%d-docker`, proxyPort)
	}
	return filepath.Join(os.TempDir(), fmt.Sprintf("httptoolkit-%d-docker.sock", proxyPort))
}

// ProxyHostEnv returns the DOCKER_HOST value clients should use.
func ProxyHostEnv(proxyPort int) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("npipe:////./pipe/httptoolkit-%d-docker", proxyPort)
	}
	return "unix://" + ProxyListenPath(proxyPort)
}

func interceptionSuffix(proxyPort int) string {
	return fmt.Sprintf("+httptoolkit:%d", proxyPort)
}
