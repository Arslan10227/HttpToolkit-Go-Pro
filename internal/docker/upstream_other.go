//go:build !windows

package docker

import (
	"context"
	"fmt"
	"net"
)

func dialNpipe(ctx context.Context, path string) (net.Conn, error) {
	return nil, fmt.Errorf("named pipes not supported on this platform")
}
