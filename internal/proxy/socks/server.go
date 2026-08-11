package socks

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"
)

// Handler receives a CONNECT target after SOCKS5 negotiation.
type Handler func(client net.Conn, targetHost string, targetPort uint16)

// Server is a minimal SOCKS5 server (CONNECT only, no auth).
type Server struct {
	handler Handler
	port    int
	ln      net.Listener
	running atomic.Bool
}

func New(handler Handler) *Server {
	return &Server{handler: handler}
}

func (s *Server) Port() int { return s.port }

func (s *Server) Start() error {
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

func (s *Server) Stop() error {
	s.running.Store(false)
	if s.ln != nil {
		err := s.ln.Close()
		s.ln = nil
		return err
	}
	return nil
}

func (s *Server) acceptLoop() {
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

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(2 * time.Minute))

	if err := socksHandshake(conn); err != nil {
		return
	}
	host, port, err := readConnectRequest(conn)
	if err != nil {
		return
	}
	_, _ = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if s.handler != nil {
		s.handler(conn, host, port)
	}
}

func socksHandshake(conn net.Conn) error {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return err
	}
	if buf[0] != 0x05 {
		return fmt.Errorf("unsupported socks version")
	}
	methods := make([]byte, buf[1])
	if _, err := io.ReadFull(conn, methods); err != nil {
		return err
	}
	_, err := conn.Write([]byte{0x05, 0x00})
	return err
}

func readConnectRequest(conn net.Conn) (string, uint16, error) {
	hdr := make([]byte, 4)
	if _, err := io.ReadFull(conn, hdr); err != nil {
		return "", 0, err
	}
	if hdr[0] != 0x05 || hdr[1] != 0x01 {
		return "", 0, fmt.Errorf("unsupported command")
	}
	var host string
	switch hdr[3] {
	case 0x01:
		ip := make([]byte, 4)
		if _, err := io.ReadFull(conn, ip); err != nil {
			return "", 0, err
		}
		host = net.IP(ip).String()
	case 0x03:
		lb := make([]byte, 1)
		if _, err := io.ReadFull(conn, lb); err != nil {
			return "", 0, err
		}
		name := make([]byte, lb[0])
		if _, err := io.ReadFull(conn, name); err != nil {
			return "", 0, err
		}
		host = string(name)
	case 0x04:
		ip := make([]byte, 16)
		if _, err := io.ReadFull(conn, ip); err != nil {
			return "", 0, err
		}
		host = net.IP(ip).String()
	default:
		return "", 0, fmt.Errorf("unsupported address type")
	}
	portBuf := make([]byte, 2)
	if _, err := io.ReadFull(conn, portBuf); err != nil {
		return "", 0, err
	}
	port := binary.BigEndian.Uint16(portBuf)
	return host, port, nil
}
