package dns

import (
	"fmt"
	"net"
	"sync"
	"time"

	mdns "github.com/miekg/dns"
)

// Server is a per-session UDP DNS forwarder (127.0.0.1:ephemeral).
type Server struct {
	mu       sync.Mutex
	ln       *net.UDPConn
	port     int
	upstream []string
}

func NewServer(upstream []string) *Server {
	if len(upstream) == 0 {
		upstream = []string{"8.8.8.8:53", "1.1.1.1:53"}
	}
	return &Server{upstream: upstream}
}

func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ln != nil {
		return nil
	}
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	ln, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	s.ln = ln
	s.port = ln.LocalAddr().(*net.UDPAddr).Port
	go s.serve()
	return nil
}

func (s *Server) Port() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.port
}

func (s *Server) Address() string {
	return fmt.Sprintf("127.0.0.1:%d", s.Port())
}

func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ln != nil {
		err := s.ln.Close()
		s.ln = nil
		s.port = 0
		return err
	}
	return nil
}

func (s *Server) serve() {
	buf := make([]byte, 4096)
	for {
		s.mu.Lock()
		ln := s.ln
		s.mu.Unlock()
		if ln == nil {
			return
		}
		n, client, err := ln.ReadFromUDP(buf)
		if err != nil {
			return
		}
		go s.forward(buf[:n], client, ln)
	}
}

func (s *Server) forward(query []byte, client *net.UDPAddr, ln *net.UDPConn) {
	for _, up := range s.upstream {
		raddr, err := net.ResolveUDPAddr("udp", up)
		if err != nil {
			continue
		}
		conn, err := net.DialUDP("udp", nil, raddr)
		if err != nil {
			continue
		}
		_, _ = conn.Write(query)
		resp := make([]byte, 4096)
		_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := conn.Read(resp)
		conn.Close()
		if err != nil || n == 0 {
			continue
		}
		_, _ = ln.WriteToUDP(resp[:n], client)
		return
	}
	msg := new(mdns.Msg)
	if err := msg.Unpack(query); err == nil {
		msg.Rcode = mdns.RcodeServerFailure
		out, _ := msg.Pack()
		_, _ = ln.WriteToUDP(out, client)
	}
}
