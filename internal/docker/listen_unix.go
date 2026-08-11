//go:build !windows

package docker

import (
	"net"
	"os"
)

func listenProxy(path string) (net.Listener, error) {
	_ = os.Remove(path)
	ln, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}
	_ = os.Chmod(path, 0o700)
	return ln, nil
}
