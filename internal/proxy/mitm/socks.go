package mitm

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync/atomic"
)

const (
	socksVersion5     = 0x05
	socksCmdConnect   = 0x01
	socksAddrIPv4     = 0x01
	socksAddrDomain   = 0x03
	socksAddrIPv6     = 0x04
	socksReplySuccess = 0x00
	socksReplyFailure = 0x01
)

type SocksServer struct {
	mitm    *Server
	port    int
	ln      net.Listener
	running atomic.Bool
}

func NewSocksServer(m *Server) *SocksServer {
	return &SocksServer{mitm: m}
}

func (s *SocksServer) Port() int { return s.port }

func (s *SocksServer) Start() error {
	if s.running.Load() {
		return nil
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	s.ln = ln
	s.port = ln.Addr().(*net.TCPAddr).Port
	s.running.Store(true)
	go s.acceptLoop()
	return nil
}

func (s *SocksServer) Stop() error {
	s.running.Store(false)
	if s.ln != nil {
		err := s.ln.Close()
		s.ln = nil
		s.port = 0
		return err
	}
	return nil
}

func (s *SocksServer) acceptLoop() {
	for s.running.Load() {
		conn, err := s.ln.Accept()
		if err != nil {
			if s.running.Load() {
				continue
			}
			return
		}
		go s.handleConn(conn)
	}
}

func (s *SocksServer) handleConn(conn net.Conn) {
	defer conn.Close()
	host, err := s.handshake(conn)
	if err != nil {
		return
	}
	s.mitm.handleSocksConnect(conn, host)
}
func (s *SocksServer) handshake(conn net.Conn) (string, error) {
	r := bufio.NewReader(conn)
	buf := make([]byte, 258)
	if _, err := io.ReadFull(r, buf[:2]); err != nil {
		_, _ = conn.Write([]byte{socksVersion5, 0xFF})
		return "", fmt.Errorf("socks version: %w", err)
	}
	if buf[0] != socksVersion5 {
		_, _ = conn.Write([]byte{socksVersion5, 0xFF})
		return "", fmt.Errorf("socks version")
	}
	nmethods := int(buf[1])
	if _, err := io.ReadFull(r, buf[:nmethods]); err != nil {
		return "", err
	}
	if _, err := conn.Write([]byte{socksVersion5, 0x00}); err != nil {
		return "", err
	}

	if _, err := io.ReadFull(r, buf[:4]); err != nil {
		return "", err
	}
	if buf[0] != socksVersion5 || buf[1] != socksCmdConnect {
		return "", fmt.Errorf("unsupported socks command")
	}

	var host string
	switch buf[3] {
	case socksAddrIPv4:
		if _, err := io.ReadFull(r, buf[:4]); err != nil {
			return "", err
		}
		host = net.IP(buf[:4]).String()
	case socksAddrDomain:
		if _, err := io.ReadFull(r, buf[:1]); err != nil {
			return "", err
		}
		dlen := int(buf[0])
		if _, err := io.ReadFull(r, buf[:dlen]); err != nil {
			return "", err
		}
		host = string(buf[:dlen])
	case socksAddrIPv6:
		if _, err := io.ReadFull(r, buf[:16]); err != nil {
			return "", err
		}
		host = net.IP(buf[:16]).String()
	default:
		return "", fmt.Errorf("unsupported address type")
	}
	if _, err := io.ReadFull(r, buf[:2]); err != nil {
		return "", err
	}
	port := binary.BigEndian.Uint16(buf[:2])
	target := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	reply := []byte{socksVersion5, socksReplySuccess, 0x00, socksAddrIPv4, 127, 0, 0, 1, 0, 0}
	if _, err := conn.Write(reply); err != nil {
		return "", err
	}
	return target, nil
}
