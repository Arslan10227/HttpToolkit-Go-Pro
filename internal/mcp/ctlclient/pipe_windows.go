//go:build windows

package ctlclient

import (
	"context"
	"net"

	"github.com/Microsoft/go-winio"
)

func pipeDialer() (func(ctx context.Context, network, addr string) (net.Conn, error), error) {
	return func(ctx context.Context, _, _ string) (net.Conn, error) {
		return winio.DialPipeContext(ctx, `\\.\pipe\httptoolkit-ctl`)
	}, nil
}
