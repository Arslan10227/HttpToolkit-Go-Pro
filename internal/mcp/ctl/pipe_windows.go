//go:build windows

package ctl

import (
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

func winListenPipe(path string) (net.Listener, error) {
	return winio.ListenPipe(path, &winio.PipeConfig{
		SecurityDescriptor: "D:P(A;;GA;;;WD)",
		MessageMode:        false,
		InputBufferSize:    65536,
		OutputBufferSize:   65536,
	})
}

// unused import guard
var _ = time.Second
