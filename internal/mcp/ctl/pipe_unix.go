//go:build !windows

package ctl

import (
	"fmt"
	"net"
)

func winListenPipe(path string) (net.Listener, error) {
	return nil, fmt.Errorf("not windows")
}
