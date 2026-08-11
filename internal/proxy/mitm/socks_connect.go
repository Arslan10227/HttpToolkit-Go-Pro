package mitm

import (
	"bufio"
	"net"
	"net/http"
)

func (s *Server) handleSocksConnect(client net.Conn, host string) {
	req, err := http.NewRequest(http.MethodConnect, "https://"+host, nil)
	if err != nil {
		return
	}
	req.Host = host
	s.handleConnect(client, bufio.NewReader(client), req)
}
