//go:build !windows

package ctlclient

import (
	"context"
	"net"
)

func pipeDialer() (func(ctx context.Context, network, addr string) (net.Conn, error), error) {
	return func(ctx context.Context, _, _ string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, "unix", "/tmp/httptoolkit-ctl.sock")
	}, nil
}
