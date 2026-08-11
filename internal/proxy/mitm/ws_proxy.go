package mitm

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/rules"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/session"
)

func buildPassthroughTransportWithSession(target *url.URL, opts rules.PassthroughOptions, sess *session.Manager) *http.Transport {
	if sess == nil {
		return buildPassthroughTransport(target, opts)
	}
	optsCopy := opts
	optsCopy.ProxyConfig = resolvePassthroughProxyConfig(opts.ProxyConfig, target, sess)
	return buildPassthroughTransport(target, optsCopy)
}

func resolvePassthroughProxyConfig(cfg any, target *url.URL, sess *session.Manager) any {
	switch v := cfg.(type) {
	case nil:
		return nil
	case map[string]any:
		if key, ok := ruleParamRef(v); ok {
			if u := sess.ResolveRuleParamProxy(key, target.Hostname()); u != "" {
				return map[string]any{"proxyUrl": u}
			}
			return nil
		}
		return v
	case []any:
		out := make([]any, 0, len(v))
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				if key, ok := ruleParamRef(m); ok {
					if u := sess.ResolveRuleParamProxy(key, target.Hostname()); u != "" {
						out = append(out, map[string]any{"proxyUrl": u})
					}
					continue
				}
			}
			out = append(out, item)
		}
		if len(out) == 0 {
			return nil
		}
		return out
	default:
		return cfg
	}
}

func ruleParamRef(m map[string]any) (string, bool) {
	for _, k := range []string{"ruleParameterKey", "ruleParameter", "mockttpParamRef"} {
		if v, ok := m[k].(string); ok && v != "" {
			return v, true
		}
	}
	for k, v := range m {
		if strings.Contains(k, "ParamRef") || strings.Contains(k, "paramRef") {
			if s, ok := v.(string); ok {
				return s, true
			}
		}
	}
	return "", false
}

func (s *Server) handleWebSocketUpgrade(w io.Writer, req *http.Request, id, urlStr string, headers map[string]string) {
	start := time.Now()
	wsRuleIdx := s.rules.MatchWS(urlStr, headers)
	action := s.rules.WSAction(wsRuleIdx)

	s.events.PublishHTTP("websocket-request", mergeMaps(map[string]any{
		"id": id, "url": urlStr, "headers": headers,
	}, eventTiming(start)))

	targetURL, err := url.Parse(urlStr)
	if err != nil {
		s.events.PublishHTTP("client-error", map[string]any{"id": id, "error": structuredError(err.Error(), "WS_URL")})
		return
	}
	targetURL = normalizeWebSocketTargetURL(targetURL)

	opts := action.Passthrough
	if s.sessions != nil {
		optsCopy := opts
		optsCopy.ProxyConfig = resolvePassthroughProxyConfig(opts.ProxyConfig, targetURL, s.sessions)
		opts = optsCopy
	}

	clientConn := connFromWriter(w)
	if clientConn == nil {
		s.events.PublishHTTP("client-error", map[string]any{"id": id, "error": structuredError("cannot hijack websocket connection", "WS_HIJACK")})
		return
	}

	if action.Kind == "ws-echo" {
		if err := writeWebSocketAccepted(clientConn, req); err != nil {
			s.events.PublishHTTP("client-error", map[string]any{"id": id, "error": structuredError(err.Error(), "WS_ACCEPT")})
			return
		}
		s.events.PublishHTTP("websocket-accepted", mergeMaps(map[string]any{
			"id": id, "url": urlStr, "statusCode": 101,
		}, eventTiming(start)))
		s.relayWebSocketEcho(clientConn, id)
		return
	}

	upstream, statusCode, respBytes, err := connectWebSocketUpstream(targetURL, req, opts)
	if err != nil {
		s.events.PublishHTTP("client-error", map[string]any{"id": id, "error": structuredError(err.Error(), "WS_UPGRADE")})
		return
	}

	if statusCode != http.StatusSwitchingProtocols {
		defer upstream.Close()
		_, _ = clientConn.Write(respBytes)
		s.events.PublishHTTP("websocket-accepted", mergeMaps(map[string]any{
			"id": id, "url": urlStr, "statusCode": statusCode,
		}, eventTiming(start)))
		return
	}

	if err := writeWebSocketAccepted(clientConn, req); err != nil {
		upstream.Close()
		s.events.PublishHTTP("client-error", map[string]any{"id": id, "error": structuredError(err.Error(), "WS_ACCEPT")})
		return
	}

	s.events.PublishHTTP("websocket-accepted", mergeMaps(map[string]any{
		"id": id, "url": urlStr, "statusCode": 101,
	}, eventTiming(start)))

	s.runWebSocketRelay(clientConn, upstream, id)
}

func normalizeWebSocketTargetURL(target *url.URL) *url.URL {
	u := *target
	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	}
	if u.Host == "" {
		u.Host = target.Host
	}
	return &u
}

func connFromWriter(w io.Writer) net.Conn {
	if c, ok := w.(net.Conn); ok {
		return c
	}
	if hj, ok := w.(http.Hijacker); ok {
		conn, _, err := hj.Hijack()
		if err != nil {
			return nil
		}
		return conn
	}
	return nil
}

func writeWebSocketAccepted(client net.Conn, req *http.Request) error {
	key := req.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		return fmt.Errorf("missing Sec-WebSocket-Key")
	}
	accept := computeWebSocketAccept(key)
	resp := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + accept + "\r\n\r\n"
	_, err := client.Write([]byte(resp))
	return err
}

func connectWebSocketUpstream(target *url.URL, req *http.Request, opts rules.PassthroughOptions) (net.Conn, int, []byte, error) {
	conn, err := dialPassthroughConn(target, opts)
	if err != nil {
		return nil, 0, nil, err
	}

	upReq := req.Clone(req.Context())
	upReq.URL = target
	upReq.RequestURI = target.RequestURI()
	upReq.Header.Del("Proxy-Connection")
	if err := upReq.Write(conn); err != nil {
		conn.Close()
		return nil, 0, nil, err
	}

	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, upReq)
	if err != nil {
		conn.Close()
		return nil, 0, nil, err
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
		resp.ContentLength = int64(len(respBody))
		if resp.Header.Get("Content-Length") != "" {
			resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(respBody)))
		}
		var buf bytes.Buffer
		_ = resp.Write(&buf)
		_ = resp.Body.Close()
		return conn, resp.StatusCode, buf.Bytes(), nil
	}
	_ = resp.Body.Close()
	return conn, resp.StatusCode, nil, nil
}

func dialPassthroughConn(target *url.URL, opts rules.PassthroughOptions) (net.Conn, error) {
	host := target.Host
	if !strings.Contains(host, ":") {
		if target.Scheme == "https" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	dialHost := host
	if len(opts.LookupServers) > 0 {
		h, port, err := net.SplitHostPort(host)
		if err == nil {
			ips, err := LookupViaServers(h, opts.LookupServers)
			if err == nil && len(ips) > 0 {
				dialHost = net.JoinHostPort(ips[0], port)
			}
		}
	}

	var conn net.Conn
	var err error

	if proxyURL := resolveProxyURL(opts.ProxyConfig, target); proxyURL != nil {
		switch strings.ToLower(proxyURL.Scheme) {
		case "socks5", "socks5h":
			conn, err = dialViaSocks5(proxyURL.Host, dialHost)
		case "http", "https":
			conn, err = dialViaHTTPConnect(proxyURL.Host, dialHost, target.Scheme == "https")
		default:
			d := net.Dialer{Timeout: 30 * time.Second}
			conn, err = d.Dial("tcp", dialHost)
		}
	} else {
		d := net.Dialer{Timeout: 30 * time.Second}
		conn, err = d.Dial("tcp", dialHost)
	}

	if err != nil {
		return nil, err
	}

	if target.Scheme == "https" {
		tlsCfg := &tls.Config{MinVersion: tls.VersionTLS12, ServerName: target.Hostname()}
		for _, h := range opts.IgnoreHostHttpsErrors {
			if MatchHostWildcard(h, target.Hostname()) {
				tlsCfg.InsecureSkipVerify = true
				break
			}
		}
		if cert := PickClientCert(target.Hostname(), opts.ClientCertHostMap); cert.PFX != "" {
			pfx, _ := base64.StdEncoding.DecodeString(cert.PFX)
			if pair, err := ClientCertFromPFX(pfx, cert.Passphrase); err == nil && pair != nil {
				tlsCfg.Certificates = []tls.Certificate{*pair}
			}
		}
		tlsConn := tls.Client(conn, tlsCfg)
		if err := tlsConn.Handshake(); err != nil {
			conn.Close()
			return nil, err
		}
		return tlsConn, nil
	}

	return conn, nil
}

func dialViaHTTPConnect(proxyHost, targetHost string, tlsUpstream bool) (net.Conn, error) {
	if !strings.Contains(proxyHost, ":") {
		proxyHost += ":80"
	}
	conn, err := net.DialTimeout("tcp", proxyHost, 30*time.Second)
	if err != nil {
		return nil, err
	}
	req := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", targetHost, targetHost)
	if _, err := conn.Write([]byte(req)); err != nil {
		conn.Close()
		return nil, err
	}
	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, &http.Request{Method: http.MethodConnect})
	if err != nil {
		conn.Close()
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		conn.Close()
		return nil, fmt.Errorf("CONNECT failed: %s", resp.Status)
	}
	if tlsUpstream {
		hostOnly, _, _ := net.SplitHostPort(targetHost)
		return tls.Client(conn, &tls.Config{ServerName: hostOnly, InsecureSkipVerify: true}), nil
	}
	return conn, nil
}

func dialViaSocks5(proxyHost, targetHost string) (net.Conn, error) {
	if !strings.Contains(proxyHost, ":") {
		proxyHost += ":1080"
	}
	conn, err := net.DialTimeout("tcp", proxyHost, 30*time.Second)
	if err != nil {
		return nil, err
	}
	if _, err := conn.Write([]byte{0x05, 0x01, 0x00}); err != nil {
		conn.Close()
		return nil, err
	}
	resp := make([]byte, 2)
	if _, err := io.ReadFull(conn, resp); err != nil || resp[0] != 0x05 || resp[1] != 0x00 {
		conn.Close()
		return nil, fmt.Errorf("socks5 auth failed")
	}

	host, portStr, err := net.SplitHostPort(targetHost)
	if err != nil {
		conn.Close()
		return nil, err
	}
	var port uint16
	fmt.Sscanf(portStr, "%d", &port)

	req := []byte{0x05, 0x01, 0x00, 0x03, byte(len(host))}
	req = append(req, []byte(host)...)
	req = append(req, byte(port>>8), byte(port))
	if _, err := conn.Write(req); err != nil {
		conn.Close()
		return nil, err
	}
	reply := make([]byte, 4)
	if _, err := io.ReadFull(conn, reply); err != nil || reply[1] != 0x00 {
		conn.Close()
		return nil, fmt.Errorf("socks5 connect failed")
	}
	switch reply[3] {
	case 0x01:
		discard := make([]byte, 4+2)
		_, _ = io.ReadFull(conn, discard)
	case 0x03:
		lenBuf := make([]byte, 1)
		_, _ = io.ReadFull(conn, lenBuf)
		discard := make([]byte, int(lenBuf[0])+2)
		_, _ = io.ReadFull(conn, discard)
	case 0x04:
		discard := make([]byte, 16+2)
		_, _ = io.ReadFull(conn, discard)
	}
	return conn, nil
}

func (s *Server) runWebSocketRelay(client, upstream net.Conn, id string) {
	defer client.Close()
	defer upstream.Close()
	defer s.events.PublishHTTP("websocket-close", map[string]any{"id": id, "timestamp": time.Now().UnixMilli()})

	done := make(chan struct{}, 2)
	go func() {
		s.copyWSWithEvents(client, upstream, id, "websocket-message-received")
		done <- struct{}{}
	}()
	s.copyWSWithEvents(upstream, client, id, "websocket-message-sent")
	<-done
}

func (s *Server) relayWebSocketEcho(client net.Conn, id string) {
	defer client.Close()
	defer s.events.PublishHTTP("websocket-close", map[string]any{"id": id, "timestamp": time.Now().UnixMilli()})

	buf := make([]byte, 32*1024)
	for {
		n, err := client.Read(buf)
		if n > 0 {
			payload, binary := decodeWSFramePayload(buf[:n])
			s.events.PublishHTTP("websocket-message-received", map[string]any{
				"id": id, "streamId": id, "content": base64.StdEncoding.EncodeToString(payload),
				"isBinary": binary, "timestamp": time.Now().UnixMilli(),
				"timingEvents": map[string]any{"startTime": time.Now().UnixMilli()},
			})
			s.events.PublishHTTP("websocket-message-sent", map[string]any{
				"id": id, "streamId": id, "content": base64.StdEncoding.EncodeToString(payload),
				"isBinary": binary, "timestamp": time.Now().UnixMilli(),
				"timingEvents": map[string]any{"startTime": time.Now().UnixMilli()},
			})
			echoFrame := encodeWSFrame(payload, binary)
			_, _ = client.Write(echoFrame)
		}
		if err != nil {
			return
		}
	}
}

func (s *Server) copyWSWithEvents(src, dst net.Conn, id, event string) {
	buf := make([]byte, 32*1024)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			payload, binary := decodeWSFramePayload(buf[:n])
			s.events.PublishHTTP(event, map[string]any{
				"id": id, "streamId": id, "content": base64.StdEncoding.EncodeToString(payload),
				"isBinary": binary, "timestamp": time.Now().UnixMilli(),
				"timingEvents": map[string]any{"startTime": time.Now().UnixMilli()},
			})
			_, _ = dst.Write(buf[:n])
		}
		if err != nil {
			return
		}
	}
}

func decodeWSFramePayload(raw []byte) ([]byte, bool) {
	if len(raw) < 2 {
		return raw, false
	}
	opcode := raw[0] & 0x0F
	binary := opcode == 0x02
	if opcode != 0x01 && opcode != 0x02 {
		return raw, binary
	}
	masked := raw[1]&0x80 != 0
	payloadLen := int(raw[1] & 0x7f)
	idx := 2
	if payloadLen == 126 {
		if len(raw) < 4 {
			return raw, binary
		}
		payloadLen = int(raw[2])<<8 | int(raw[3])
		idx = 4
	} else if payloadLen == 127 {
		if len(raw) < 10 {
			return raw, binary
		}
		payloadLen = 0
		for i := 0; i < 8; i++ {
			payloadLen = payloadLen<<8 | int(raw[2+i])
		}
		idx = 10
	}
	if masked {
		if len(raw) < idx+4+payloadLen {
			return raw, binary
		}
		mask := raw[idx : idx+4]
		idx += 4
		payload := make([]byte, payloadLen)
		copy(payload, raw[idx:idx+payloadLen])
		for i := range payload {
			payload[i] ^= mask[i%4]
		}
		return payload, binary
	}
	if len(raw) < idx+payloadLen {
		return raw, binary
	}
	return raw[idx : idx+payloadLen], binary
}

func encodeWSFrame(payload []byte, binary bool) []byte {
	opcode := byte(0x01)
	if binary {
		opcode = 0x02
	}
	frame := []byte{0x80 | opcode}
	// Server-to-client frames MUST NOT be masked per RFC 6455 §5.3.
	if len(payload) < 126 {
		frame = append(frame, byte(len(payload)))
	} else if len(payload) < 65536 {
		frame = append(frame, byte(126), byte(len(payload)>>8), byte(len(payload)))
	} else {
		frame = append(frame, byte(127), 0, 0, 0, 0, 0, 0, 0, 0)
	}
	return append(frame, payload...)
}

func computeWebSocketAccept(key string) string {
	const magic = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	h := sha1.Sum([]byte(key + magic))
	return base64.StdEncoding.EncodeToString(h[:])
}
