package docker

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync/atomic"
)

// TunnelSocks is a SOCKS5 listener that forwards via HTTP CONNECT to the main proxy.
type TunnelSocks struct {
	proxyPort int
	port      int
	ln        net.Listener
	running   atomic.Bool
}

func NewTunnelSocks(proxyPort int) *TunnelSocks {
	return &TunnelSocks{proxyPort: proxyPort}
}

func (t *TunnelSocks) Port() int { return t.port }

func (t *TunnelSocks) Start() error {
	if t.running.Load() {
		return nil
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	t.ln = ln
	t.port = ln.Addr().(*net.TCPAddr).Port
	t.running.Store(true)
	go t.acceptLoop()
	return nil
}

func (t *TunnelSocks) Stop() error {
	t.running.Store(false)
	if t.ln != nil {
		err := t.ln.Close()
		t.ln = nil
		t.port = 0
		return err
	}
	return nil
}

func (t *TunnelSocks) acceptLoop() {
	for t.running.Load() {
		conn, err := t.ln.Accept()
		if err != nil {
			if t.running.Load() {
				continue
			}
			return
		}
		go t.handleConn(conn)
	}
}

func (t *TunnelSocks) handleConn(client net.Conn) {
	defer client.Close()
	target, err := socksHandshake(client)
	if err != nil {
		return
	}
	upstream, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", t.proxyPort))
	if err != nil {
		return
	}
	defer upstream.Close()
	req := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", target, target)
	if _, err := upstream.Write([]byte(req)); err != nil {
		return
	}
	buf := make([]byte, 1024)
	n, err := upstream.Read(buf)
	if err != nil || n == 0 {
		return
	}
	go io.Copy(upstream, client)
	io.Copy(client, upstream)
}

func socksHandshake(conn net.Conn) (string, error) {
	buf := make([]byte, 512)
	if _, err := io.ReadAtLeast(conn, buf, 2); err != nil {
		return "", err
	}
	nmethods := int(buf[1])
	if _, err := io.ReadFull(conn, buf[:nmethods]); err != nil {
		return "", err
	}
	if _, err := conn.Write([]byte{0x05, 0x00}); err != nil {
		return "", err
	}
	if _, err := io.ReadFull(conn, buf[:4]); err != nil {
		return "", err
	}
	if buf[1] != 0x01 {
		return "", fmt.Errorf("unsupported command")
	}
	var host string
	switch buf[3] {
	case 0x01:
		if _, err := io.ReadFull(conn, buf[:4]); err != nil {
			return "", err
		}
		host = net.IP(buf[:4]).String()
	case 0x03:
		if _, err := io.ReadFull(conn, buf[:1]); err != nil {
			return "", err
		}
		dlen := int(buf[0])
		if _, err := io.ReadFull(conn, buf[:dlen]); err != nil {
			return "", err
		}
		host = string(buf[:dlen])
	case 0x04:
		if _, err := io.ReadFull(conn, buf[:16]); err != nil {
			return "", err
		}
		host = net.IP(buf[:16]).String()
	default:
		return "", fmt.Errorf("bad atyp")
	}
	if _, err := io.ReadFull(conn, buf[:2]); err != nil {
		return "", err
	}
	port := binary.BigEndian.Uint16(buf[:2])
	reply := []byte{0x05, 0x00, 0x00, 0x01, 127, 0, 0, 1, 0, 0}
	if _, err := conn.Write(reply); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", host, port), nil
}
