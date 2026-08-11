//go:build windows

package docker

import (
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

func listenProxy(path string) (net.Listener, error) {
	return winio.ListenPipe(path, &winio.PipeConfig{
		SecurityDescriptor: "D:P(A;;GA;;;WD)",
		MessageMode:        false,
		InputBufferSize:    65536,
		OutputBufferSize:   65536,
	})
}

var _ = time.Second
