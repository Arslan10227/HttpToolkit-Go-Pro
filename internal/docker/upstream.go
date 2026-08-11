package docker

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

func dockerHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext:       dialDockerDaemon,
			DisableKeepAlives: false,
		},
		Timeout: 0,
	}
}

func dialDockerDaemon(ctx context.Context, _, _ string) (net.Conn, error) {
	addr, err := resolveDockerDaemonAddress()
	if err != nil {
		return nil, err
	}
	d := net.Dialer{Timeout: 30 * time.Second}
	switch addr.kind {
	case "unix":
		return d.DialContext(ctx, "unix", addr.path)
	case "npipe":
		return dialNpipe(ctx, addr.path)
	case "tcp":
		return d.DialContext(ctx, "tcp", addr.host)
	default:
		return nil, fmt.Errorf("unsupported docker host")
	}
}

type daemonAddress struct {
	kind string
	path string
	host string
}

func resolveDockerDaemonAddress() (daemonAddress, error) {
	if host := strings.TrimSpace(os.Getenv("DOCKER_HOST")); host != "" {
		if strings.HasPrefix(host, "unix://") {
			return daemonAddress{kind: "unix", path: strings.TrimPrefix(host, "unix://")}, nil
		}
		if strings.HasPrefix(host, "npipe://") {
			path := strings.TrimPrefix(host, "npipe://")
			path = strings.TrimPrefix(path, "//")
			if !strings.HasPrefix(path, `\\.\pipe\`) {
				path = `\\.\pipe\` + path
			}
			return daemonAddress{kind: "npipe", path: path}, nil
		}
		if strings.HasPrefix(host, "tcp://") {
			return daemonAddress{kind: "tcp", host: strings.TrimPrefix(host, "tcp://")}, nil
		}
	}
	if runtime.GOOS == "windows" {
		return daemonAddress{kind: "npipe", path: `\\.\pipe\docker_engine`}, nil
	}
	return daemonAddress{kind: "unix", path: "/var/run/docker.sock"}, nil
}

func dockerAPIURL(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return "http://docker" + path
}
