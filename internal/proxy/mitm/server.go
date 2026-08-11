package mitm

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	certmgr "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/logger"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/events"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/rules"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/rtc"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/session"
	"golang.org/x/net/http2"
)

type rtcHook interface {
	IsRunning() bool
	NotifyPeerConnected(sessionID string) string
	NotifySDPOffer(httpSessionID, sdp string) string
}

type Server struct {
	certs        *certmgr.Manager
	rules        *rules.Engine
	events       *events.Bus
	sessions     *session.Manager
	rtc          rtcHook
	bpMgr        *BreakpointManager
	port         int
	ln           net.Listener
	socks        *SocksServer
	running      atomic.Bool
	http2Enabled atomic.Bool
	socksEnabled atomic.Bool
}

// bufConn wraps a net.Conn so that any bytes already buffered by a
// bufio.Reader are consumed first. This is essential for proxy connections:
// http.ReadRequest may read-ahead into the TLS ClientHello that follows the
// CONNECT request, and those buffered bytes must be passed to the inner TLS
// handshake rather than re-read from the underlying connection.
type bufConn struct {
	r *bufio.Reader
	net.Conn
}

func (bc *bufConn) Read(p []byte) (int, error) {
	return bc.r.Read(p)
}

func NewServer(certs *certmgr.Manager, eng *rules.Engine, bus *events.Bus) *Server {
	s := &Server{certs: certs, rules: eng, events: bus}
	s.socks = NewSocksServer(s)
	return s
}

func (s *Server) SetBreakpointManager(mgr *BreakpointManager) {
	s.bpMgr = mgr
}

func (s *Server) SetSession(sess *session.Manager) {
	s.sessions = sess
}

func (s *Server) SetRTC(mgr rtcHook) {
	s.rtc = mgr
}

func (s *Server) SetSocksEnabled(on bool) { s.socksEnabled.Store(on) }

func (s *Server) SocksPort() int { return s.socks.Port() }

func (s *Server) startSocksIfEnabled() {
	if s.socksEnabled.Load() {
		_ = s.socks.Start()
	}
}

func (s *Server) Port() int { return s.port }

func (s *Server) SetHTTP2Enabled(on bool) {
	s.http2Enabled.Store(on)
}

func (s *Server) Start(port int) error {
	if s.running.Load() {
		return nil
	}
	logger.Debug("Starting MITM proxy server", map[string]any{"port": port})
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Debug("Port occupied, falling back to dynamic port", map[string]any{"port": port, "error": err.Error()})
		ln, err = net.Listen("tcp", ":0")
		if err != nil {
			logger.Error(err, map[string]any{"msg": "failed to bind fallback port"})
			return err
		}
	}
	tlsCfg := &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			name := info.ServerName
			if name == "" {
				name = "127.0.0.1"
			}
			logger.Debug("Leaf cert requested on connection", map[string]any{"serverName": name})
			cert, err := s.certs.LeafCertChain(name)
			if err != nil {
				logger.Error(err, map[string]any{"msg": "failed to generate leaf cert", "serverName": name})
				return nil, err
			}
			return &cert, nil
		},
		MinVersion: tls.VersionTLS12,
	}
	s.ln = &muxListener{Listener: ln, tlsConfig: tlsCfg}
	s.port = ln.Addr().(*net.TCPAddr).Port
	s.running.Store(true)
	go s.acceptLoop()
	s.startSocksIfEnabled()
	logger.Info("MITM proxy server successfully listening", map[string]any{"boundPort": s.port})
	return nil
}

func (s *Server) StartPortRange(start, end int) error {
	logger.Debug("Starting proxy on port range", map[string]any{"start": start, "end": end})
	if start <= 0 {
		start = 8000
	}
	if end <= start {
		end = start + 100
	}
	var lastErr error
	for port := start; port <= end; port++ {
		if err := s.Start(port); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("no port available in range")
}

func (s *Server) Stop() error {
	logger.Info("Stopping MITM proxy server", map[string]any{"port": s.port})
	s.running.Store(false)
	_ = s.socks.Stop()
	if s.ln != nil {
		err := s.ln.Close()
		s.ln = nil
		return err
	}
	return nil
}

func (s *Server) IsRunning() bool { return s.running.Load() }

func (s *Server) acceptLoop() {
	for s.running.Load() {
		conn, err := s.ln.Accept()
		if err != nil {
			if s.running.Load() {
				logger.Error(err, map[string]any{"msg": "error in accept loop"})
				continue
			}
			return
		}
		logger.Debug("Accepted TCP connection", map[string]any{"remoteAddr": conn.RemoteAddr().String()})
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(5 * time.Minute))
	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		logger.Debug("Failed to read HTTP request from connection", map[string]any{"remoteAddr": conn.RemoteAddr().String(), "error": err.Error()})
		return
	}
	logger.Debug("Handling incoming HTTP request", map[string]any{"method": req.Method, "host": req.Host, "url": req.URL.String()})
	if req.Method == http.MethodConnect {
		s.handleConnect(conn, reader, req)
		return
	}
	s.handleHTTP(conn, req)
}

func (s *Server) handleConnect(client net.Conn, reader *bufio.Reader, req *http.Request) {
	host := req.Host
	hostOnly, port, err := net.SplitHostPort(host)
	if err != nil {
		// No port (or an IPv6 literal without a port). Default to 443.
		hostOnly = strings.Trim(req.Host, "[]")
		port = "443"
	}
	host = net.JoinHostPort(hostOnly, port)

	// Wrap the connection so the inner TLS/raw bytes that http.ReadRequest
	// may have buffered are not lost.
	bc := &bufConn{r: reader, Conn: client}

	logger.Debug("Dispatching CONNECT request", map[string]any{"destination": host})

	// If it is a TLS port (like 443 or 8443), decrypt and intercept via MITM.
	if isTLSDefaultPort(host) {
		logger.Debug("MITM Interception active for target TLS destination", map[string]any{"host": host})
		// 1. Respond 200 Connection Established to client
		_, err := bc.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
		if err != nil {
			logger.Error(err, map[string]any{"msg": "failed to write connection established response"})
			return
		}

		// 2. Generate leaf certificate (including the CA in the chain for SPKI pinning)
		cert, err := s.certs.LeafCertChain(hostOnly)
		if err != nil {
			logger.Error(err, map[string]any{"msg": "failed to create leaf certificate", "host": hostOnly})
			return
		}

		nextProtos := []string{"http/1.1"}
		if s.http2Enabled.Load() {
			nextProtos = []string{"h2", "http/1.1"}
		}
		tlsCfg := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
			NextProtos:   nextProtos,
		}

		// 3. Perform TLS Server Handshake
		tlsClient := tls.Server(bc, tlsCfg)
		if err := tlsClient.Handshake(); err != nil {
			certInfo := certmgr.DescribeCertDER(cert.Certificate[0])
			logger.Info("TLS Client Handshake failed", map[string]any{
				"host":       hostOnly,
				"error":      err.Error(),
				"cert":       certInfo,
				"serverName": tlsCfg.ServerName,
			})
			s.events.PublishHTTP("tls-client-error", map[string]any{
				"id":       uuid.NewString(),
				"hostname": hostOnly,
				"error":    structuredError(err.Error(), "TLS_HANDSHAKE"),
				"cert":     certInfo,
			})
			return
		}
		defer tlsClient.Close()

		// 4. Handle HTTP2 or HTTP1.1 loop
		negotiated := tlsClient.ConnectionState().NegotiatedProtocol
		logger.Debug("TLS Client Handshake completed successfully", map[string]any{"host": hostOnly, "protocol": negotiated})
		if negotiated == "h2" && s.http2Enabled.Load() {
			logger.Debug("Upgrading connection loop to multiplexed HTTP/2", map[string]any{"host": hostOnly})
			s.serveHTTP2(tlsClient, hostOnly)
			return
		}

		logger.Debug("Serving decrypted HTTP/1.1 requests", map[string]any{"host": hostOnly})
		tlsReader := bufio.NewReader(tlsClient)
		for {
			innerReq, err := http.ReadRequest(tlsReader)
			if err != nil {
				logger.Debug("Decrypted HTTP/1.1 client connection closed", map[string]any{"host": hostOnly})
				return
			}
			innerReq.URL.Scheme = "https"
			innerReq.URL.Host = hostOnly
			innerReq.RequestURI = ""

			// Generate a NEW unique exchange ID for this request!
			exchangeID := uuid.NewString()
			logger.Debug("MITM Intercepted HTTPS Request", map[string]any{"exchangeId": exchangeID, "method": innerReq.Method, "url": innerReq.URL.String()})
			s.serveMITMRequest(tlsClient, innerReq, exchangeID)
		}

	} else {
		// Passthrough path for non-TLS/non-default ports
		id := uuid.NewString()
		logger.Debug("Non-TLS port detected, falling back to raw TCP passthrough", map[string]any{"id": id, "destination": host})
		s.emitRawPassthrough(id, host, true)
		defer s.emitRawPassthrough(id, host, false)

		target, err := net.DialTimeout("tcp", host, 30*time.Second)
		if err != nil {
			logger.Error(err, map[string]any{"msg": "failed to connect to passthrough destination", "destination": host})
			s.events.PublishHTTP("client-error", map[string]any{
				"id":    id,
				"error": structuredError(err.Error(), "ECONNREFUSED"),
			})
			return
		}
		defer target.Close()

		_, err = bc.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
		if err != nil {
			return
		}

		done := make(chan struct{}, 2)
		go func() {
			_, _ = io.Copy(target, bc)
			done <- struct{}{}
		}()
		go func() {
			_, _ = io.Copy(bc, target)
			done <- struct{}{}
		}()
		<-done
		logger.Debug("Raw TCP tunnel connection terminated", map[string]any{"id": id, "destination": host})
	}
}

func (s *Server) serveHTTP2(conn net.Conn, host string) {
	h2srv := &http2.Server{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Scheme = "https"
		r.URL.Host = host
		exchangeID := uuid.NewString()
		logger.Debug("MITM Intercepted Multiplexed HTTP/2 Request", map[string]any{"exchangeId": exchangeID, "method": r.Method, "url": r.URL.String()})
		s.serveMITMRequest(w, r, exchangeID)
	})
	h2srv.ServeConn(conn, &http2.ServeConnOpts{Handler: handler})
}

func (s *Server) handleHTTP(client net.Conn, req *http.Request) {
	id := uuid.NewString()
	if req.URL != nil && req.URL.Scheme == "" {
		// Non-CONNECT requests inside an HTTP proxy are for http:// targets
		// (https:// uses CONNECT). Do not assume https just because the proxy
		// connection itself is TLS.
		req.URL.Scheme = "http"
	}
	if req.URL != nil && req.URL.Host == "" {
		// For non-proxy-style requests (GET /path rather than GET http://host/path),
		// populate the URL host from the Host header so downstream code can
		// resolve the target correctly.
		req.URL.Host = req.Host
	}
	logger.Debug("Direct non-CONNECT HTTP request intercepted", map[string]any{"exchangeId": id, "method": req.Method, "url": req.URL.String()})
	s.serveMITMRequest(client, req, id)
}

func (s *Server) serveMITMRequest(w io.Writer, req *http.Request, id string) {
	start := time.Now()
	if isAmiusingHost(req.URL.Host) || isAmiusingHost(req.Host) {
		s.serveAmiusingFallback(w, req, id)
		return
	}
	urlStr := req.URL.String()
	if req.URL.Scheme == "" {
		urlStr = "http://" + req.Host + req.URL.Path
	}
	headers := flattenHeaders(req.Header)

	if isWebSocketUpgrade(req) {
		logger.Debug("WebSocket upgrade request detected", map[string]any{"exchangeId": id, "url": urlStr})
		s.handleWebSocketUpgrade(w, req, id, urlStr, headers)
		return
	}

	logger.Debug("Processing MITM request exchange", map[string]any{"exchangeId": id, "method": req.Method, "url": urlStr})
	s.events.PublishHTTP("request-initiated", mergeMaps(map[string]any{
		"id": id, "method": req.Method, "url": urlStr, "headers": headers,
	}, eventTiming(start)))

	bodyBytes, _ := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	s.events.PublishHTTP("request", mergeMaps(map[string]any{
		"id": id, "method": req.Method, "url": urlStr, "headers": headers,
		"body": map[string]string{"buffer": base64.StdEncoding.EncodeToString(bodyBytes)},
	}, eventTiming(start)))

	s.maybeNotifySDP(id, bodyBytes, req.Header.Get("Content-Type"))

	ruleIdx := s.rules.MatchFirst(req.Method, urlStr, headers, bodyBytes)
	action := s.rules.ActionForRule(ruleIdx)
	if ruleIdx >= 0 {
		phase := "request"
		if action.IsCallback {
			phase = "callback"
		}
		if action.Kind == "breakpoint" {
			phase = "breakpoint"
		}

		if s.bpMgr != nil && (action.Kind == "breakpoint" || action.IsCallback) {
			_, ch := s.bpMgr.Pause(id, action.RuleID, req, bodyBytes)
			s.events.PublishHTTP("rule-event", map[string]any{
				"id":            id,
				"requestId":     id,
				"matchedRuleId": action.RuleID,
				"phase":         phase,
				"paused":        true,
				"request": map[string]any{
					"method":  req.Method,
					"url":     urlStr,
					"headers": headers,
					"body":    base64.StdEncoding.EncodeToString(bodyBytes),
				},
			})

			select {
			case bpAction := <-ch:
				switch bpAction.Type {
				case ActionResume:
					if bpAction.Method != "" {
						req.Method = bpAction.Method
					}
					if bpAction.URL != "" {
						urlStr = bpAction.URL
						newURL, err := url.Parse(urlStr)
						if err == nil {
							req.URL = newURL
						}
					}
					if bpAction.Headers != nil {
						req.Header = make(http.Header)
						for k, v := range bpAction.Headers {
							req.Header.Set(k, v)
						}
						headers = flattenHeaders(req.Header)
					}
					if bpAction.Body != nil {
						bodyBytes = bpAction.Body
						req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
					}
					action.Kind = "passthrough"

				case ActionRespond:
					statusCode := bpAction.StatusCode
					if statusCode == 0 {
						statusCode = 200
					}
					action = rules.StepAction{
						Kind:       "fixed-response",
						StatusCode: statusCode,
						Body:       string(bpAction.RespBody),
						Headers:    bpAction.RespHeaders,
					}

				case ActionAbort:
					action.Kind = "abort"
				}

			case <-time.After(60 * time.Second):
				// Auto-timeout: fall back to normal passthrough
				action.Kind = "passthrough"
			}
			s.bpMgr.Remove(id)
		} else {
			s.events.PublishHTTP("rule-event", map[string]any{
				"id": id, "matchedRuleId": action.RuleID, "phase": phase,
			})
			if action.Webhook.URL != "" {
				go rules.FireWebhook(action.Webhook, map[string]any{
					"id": id, "method": req.Method, "url": urlStr,
				})
			}
		}
	}

	switch action.Kind {
	case "stream":
		writeStreamResponse(w, id, s.events, action, start)
		return
	case "fixed-response":
		writeFixedResponse(w, id, s.events, action, start)
		return
	case "redirect":
		action.Headers = map[string]string{"Location": action.Body}
		action.Body = ""
		writeFixedResponse(w, id, s.events, action, start)
		return
	case "file":
		if data, err := os.ReadFile(action.FilePath); err == nil {
			action.Body = string(data)
			writeFixedResponse(w, id, s.events, action, start)
		} else {
			s.events.PublishHTTP("abort", map[string]any{"id": id, "error": structuredError(err.Error(), "FILE_ERROR")})
		}
		return
	case "close-connection":
		s.events.PublishHTTP("abort", map[string]any{"id": id})
		return
	case "reset-connection":
		s.events.PublishHTTP("abort", map[string]any{"id": id})
		if rc, ok := w.(net.Conn); ok {
			resetTCPConn(rc)
		}
		return
	case "abort":
		s.events.PublishHTTP("abort", map[string]any{"id": id, "error": structuredError("request aborted by rule", "ABORTED")})
		if rc, ok := w.(net.Conn); ok {
			resetTCPConn(rc)
		}
		return
	case "timeout":
		time.Sleep(time.Duration(action.DelayMs) * time.Millisecond)
		s.events.PublishHTTP("abort", map[string]any{"id": id, "error": structuredError("timeout", "TIMEOUT")})
		return
	}
	if action.DelayMs > 0 {
		time.Sleep(time.Duration(action.DelayMs) * time.Millisecond)
	}

	targetURL, _ := url.Parse(urlStr)
	if action.Passthrough.TransformRequest != nil {
		var err error
		bodyBytes, err = s.applyRequestTransforms(req, targetURL, bodyBytes, action.Passthrough.TransformRequest)
		if err != nil {
			s.events.PublishHTTP("abort", map[string]any{"id": id, "error": structuredError(err.Error(), "TRANSFORM_ERROR")})
			return
		}
	}

	transport := buildPassthroughTransportWithSession(targetURL, action.Passthrough, s.sessions)
	proxyReq, _ := http.NewRequest(req.Method, targetURL.String(), bytes.NewReader(bodyBytes))
	proxyReq.Header = req.Header
	logger.Debug("MITM RoundTrip", map[string]any{"target": targetURL.String()})
	resp, err := transport.RoundTrip(proxyReq)
	if err != nil {
		logger.Error(err, map[string]any{"msg": "MITM RoundTrip error", "target": targetURL.String()})
		s.events.PublishHTTP("abort", map[string]any{"id": id, "error": structuredError(err.Error(), "UPSTREAM_ERROR")})
		return
	}
	logger.Debug("MITM RoundTrip response", map[string]any{"target": targetURL.String(), "status": resp.StatusCode})
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	respStatusCode := resp.StatusCode
	if action.Passthrough.TransformResponse != nil {
		var transformErr error
		respBody, respStatusCode, transformErr = s.applyResponseTransforms(resp, respBody, action.Passthrough.TransformResponse)
		if transformErr != nil {
			s.events.PublishHTTP("abort", map[string]any{"id": id, "error": structuredError(transformErr.Error(), "TRANSFORM_ERROR")})
			return
		}
		resp.StatusCode = respStatusCode
		resp.Status = fmt.Sprintf("%d %s", respStatusCode, httpStatusMessage(respStatusCode))
	}

	resp.Body = io.NopCloser(bytes.NewReader(respBody))
	resp.ContentLength = int64(len(respBody))
	if resp.Header.Get("Content-Length") != "" {
		resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(respBody)))
	}

	respHeaders := flattenHeaders(resp.Header)
	s.events.PublishHTTP("response", mergeMaps(map[string]any{
		"id": id, "statusCode": respStatusCode, "statusMessage": httpStatusMessage(respStatusCode),
		"headers": respHeaders,
		"body":    map[string]string{"buffer": base64.StdEncoding.EncodeToString(respBody)},
	}, eventTimingWithResponse(start, time.Now())))

	s.maybeNotifySDP(id, respBody, resp.Header.Get("Content-Type"))

	if rw, ok := w.(http.ResponseWriter); ok {
		for k, vals := range resp.Header {
			for _, v := range vals {
				rw.Header().Add(k, v)
			}
		}
		rw.WriteHeader(respStatusCode)
		_, _ = rw.Write(respBody)
		return
	}
	_ = resp.Write(w)
}

func flattenHeaders(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for k, vals := range h {
		out[strings.ToLower(k)] = strings.Join(vals, ", ")
	}
	return out
}

func writeFixedResponse(w io.Writer, id string, bus *events.Bus, action rules.StepAction, start time.Time) {
	// If the writer is an http.ResponseWriter (HTTP/2 or tested handler),
	// use its methods instead of writing raw HTTP/1.1 bytes, which would
	// corrupt the HTTP/2 framing.
	if rw, ok := w.(http.ResponseWriter); ok {
		statusCode := action.StatusCode
		if statusCode == 0 {
			statusCode = 200
		}
		for k, v := range action.Headers {
			rw.Header().Set(k, v)
		}
		rw.WriteHeader(statusCode)
		_, _ = rw.Write([]byte(action.Body))
		bus.PublishHTTP("response", mergeMaps(map[string]any{
			"id": id, "statusCode": statusCode, "statusMessage": httpStatusMessage(statusCode),
			"headers": action.Headers,
			"body":    map[string]string{"buffer": base64.StdEncoding.EncodeToString([]byte(action.Body))},
		}, eventTimingWithResponse(start, time.Now())))
		return
	}

	statusText := httpStatusMessage(action.StatusCode)
	if statusText == "" {
		statusText = "OK"
	}
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", action.StatusCode, statusText)
	resp := statusLine
	for k, v := range action.Headers {
		resp += k + ": " + v + "\r\n"
	}
	if !hasHeader(action.Headers, "content-length") {
		resp += fmt.Sprintf("Content-Length: %d\r\n", len(action.Body))
	}
	resp += "\r\n" + action.Body
	_, _ = w.Write([]byte(resp))
	bus.PublishHTTP("response", mergeMaps(map[string]any{
		"id": id, "statusCode": action.StatusCode, "statusMessage": httpStatusMessage(action.StatusCode),
		"headers": action.Headers,
		"body":    map[string]string{"buffer": base64.StdEncoding.EncodeToString([]byte(action.Body))},
	}, eventTimingWithResponse(start, time.Now())))
}

func writeStreamResponse(w io.Writer, id string, bus *events.Bus, action rules.StepAction, start time.Time) {
	data, err := os.ReadFile(action.StreamFile)
	if err != nil {
		bus.PublishHTTP("abort", map[string]any{"id": id, "error": structuredError(err.Error(), "FILE_ERROR")})
		return
	}
	status := 200
	resp := fmt.Sprintf("HTTP/1.1 %d OK\r\nTransfer-Encoding: chunked\r\nContent-Type: application/octet-stream\r\n\r\n", status)
	_, _ = w.Write([]byte(resp))
	chunk := fmt.Sprintf("%x\r\n%s\r\n", len(data), data)
	_, _ = w.Write([]byte(chunk))
	_, _ = w.Write([]byte("0\r\n\r\n"))
	bus.PublishHTTP("response", mergeMaps(map[string]any{
		"id": id, "statusCode": status, "statusMessage": "OK",
		"headers": map[string]string{"transfer-encoding": "chunked"},
		"body":    map[string]string{"buffer": base64.StdEncoding.EncodeToString(data)},
	}, eventTimingWithResponse(start, time.Now())))
}

func hasHeader(headers map[string]string, name string) bool {
	for k := range headers {
		if strings.EqualFold(k, name) {
			return true
		}
	}
	return false
}

func (s *Server) maybeNotifySDP(httpExchangeID string, body []byte, contentType string) {
	if s.rtc == nil || !s.rtc.IsRunning() || !rtc.LooksLikeSDP(body, contentType) {
		return
	}
	peerID := s.rtc.NotifySDPOffer(httpExchangeID, string(body))
	s.events.PublishRTC("session-metadata", map[string]any{
		"id": httpExchangeID, "peerId": peerID, "timestamp": time.Now().UnixMilli(), "type": "offer",
	})
}

func mergeMaps(a, b map[string]any) map[string]any {
	out := make(map[string]any, len(a)+len(b))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}

func isWebSocketUpgrade(req *http.Request) bool {
	return strings.EqualFold(req.Header.Get("Upgrade"), "websocket") &&
		strings.Contains(strings.ToLower(req.Header.Get("Connection")), "upgrade")
}

func (s *Server) applyRequestTransforms(req *http.Request, targetURL *url.URL, body []byte, tr *rules.TransformRequestOptions) ([]byte, error) {
	if tr.ReplaceHost != nil && tr.ReplaceHost.TargetHost != "" {
		targetURL.Host = tr.ReplaceHost.TargetHost
		if tr.ReplaceHost.UpdateHostHeader {
			req.Header.Set("Host", tr.ReplaceHost.TargetHost)
			req.Host = tr.ReplaceHost.TargetHost
		}
	}
	if tr.SetProtocol != "" {
		targetURL.Scheme = tr.SetProtocol
	}

	for _, p := range tr.MatchReplaceHost {
		re, err := regexp.Compile(p.Match)
		if err == nil {
			targetURL.Host = re.ReplaceAllString(targetURL.Host, p.Replace)
			if tr.ReplaceHost != nil && tr.ReplaceHost.UpdateHostHeader {
				req.Header.Set("Host", targetURL.Host)
				req.Host = targetURL.Host
			}
		}
	}
	for _, p := range tr.MatchReplacePath {
		re, err := regexp.Compile(p.Match)
		if err == nil {
			targetURL.Path = re.ReplaceAllString(targetURL.Path, p.Replace)
		}
	}
	for _, p := range tr.MatchReplaceQuery {
		re, err := regexp.Compile(p.Match)
		if err == nil {
			query := targetURL.RawQuery
			query = re.ReplaceAllString(query, p.Replace)
			targetURL.RawQuery = query
		}
	}

	if tr.ReplaceHeaders != nil {
		newHeader := http.Header{}
		for k, v := range tr.ReplaceHeaders {
			newHeader.Set(k, v)
		}
		req.Header = newHeader
	}
	for k, v := range tr.UpdateHeaders {
		if v == "" {
			req.Header.Del(k)
		} else {
			req.Header.Set(k, v)
		}
	}

	var err error
	body, err = applyBodyTransforms(body, tr.ReplaceBody, tr.ReplaceBodyFile, tr.UpdateJsonBody, tr.PatchJsonBody, tr.MatchReplaceBody)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Server) applyResponseTransforms(resp *http.Response, body []byte, tr *rules.TransformResponseOptions) ([]byte, int, error) {
	statusCode := resp.StatusCode
	if tr.ReplaceStatus > 0 {
		statusCode = tr.ReplaceStatus
	}

	if tr.ReplaceHeaders != nil {
		newHeader := http.Header{}
		for k, v := range tr.ReplaceHeaders {
			newHeader.Set(k, v)
		}
		resp.Header = newHeader
	}
	for k, v := range tr.UpdateHeaders {
		if v == "" {
			resp.Header.Del(k)
		} else {
			resp.Header.Set(k, v)
		}
	}

	var err error
	body, err = applyBodyTransforms(body, tr.ReplaceBody, tr.ReplaceBodyFile, tr.UpdateJsonBody, tr.PatchJsonBody, tr.MatchReplaceBody)
	if err != nil {
		return nil, statusCode, err
	}

	return body, statusCode, nil
}

func applyBodyTransforms(body []byte, replaceBody *string, replaceBodyFile *string, updateJsonBody map[string]any, patchJsonBody []any, matchReplaceBody []rules.MatchReplacePair) ([]byte, error) {
	if replaceBody != nil {
		body = []byte(*replaceBody)
	}

	if replaceBodyFile != nil {
		data, err := os.ReadFile(*replaceBodyFile)
		if err != nil {
			return nil, fmt.Errorf("read replace file error: %w", err)
		}
		body = data
	}

	if updateJsonBody != nil {
		var target map[string]any
		if len(body) > 0 {
			_ = json.Unmarshal(body, &target)
		}
		if target == nil {
			target = make(map[string]any)
		}
		mergeJSON(target, updateJsonBody)
		newBody, err := json.Marshal(target)
		if err != nil {
			return nil, fmt.Errorf("marshal merge json error: %w", err)
		}
		body = newBody
	}

	if len(patchJsonBody) > 0 {
		var target any
		if len(body) > 0 {
			_ = json.Unmarshal(body, &target)
		}
		res, err := applyJSONPatch(target, patchJsonBody)
		if err != nil {
			return nil, fmt.Errorf("apply json patch error: %w", err)
		}
		newBody, err := json.Marshal(res)
		if err != nil {
			return nil, fmt.Errorf("marshal patch json error: %w", err)
		}
		body = newBody
	}

	for _, p := range matchReplaceBody {
		re, err := regexp.Compile(p.Match)
		if err == nil {
			bodyStr := string(body)
			bodyStr = re.ReplaceAllString(bodyStr, p.Replace)
			body = []byte(bodyStr)
		}
	}

	return body, nil
}

func mergeJSON(target map[string]any, patch map[string]any) {
	for k, v := range patch {
		if v == nil {
			delete(target, k)
		} else if mPatch, ok := v.(map[string]any); ok {
			if mTarget, ok := target[k].(map[string]any); ok {
				mergeJSON(mTarget, mPatch)
			} else {
				target[k] = mPatch
			}
		} else {
			target[k] = v
		}
	}
}

func applyJSONPatch(target any, patch []any) (any, error) {
	var err error
	for _, opVal := range patch {
		opMap, ok := opVal.(map[string]any)
		if !ok {
			continue
		}
		op, _ := opMap["op"].(string)
		pathStr, _ := opMap["path"].(string)
		value := opMap["value"]

		if pathStr == "" || pathStr == "/" {
			if op == "add" || op == "replace" {
				target = value
			} else if op == "remove" {
				target = nil
			}
			continue
		}

		parts := strings.Split(strings.TrimPrefix(pathStr, "/"), "/")
		target, err = setPathValue(target, parts, op, value)
		if err != nil {
			return nil, err
		}
	}
	return target, nil
}

func setPathValue(obj any, path []string, op string, value any) (any, error) {
	if len(path) == 0 {
		if op == "replace" || op == "add" {
			return value, nil
		}
		return nil, fmt.Errorf("invalid op at root: %s", op)
	}

	current := path[0]
	if m, ok := obj.(map[string]any); ok {
		if len(path) == 1 {
			if op == "add" || op == "replace" {
				m[current] = value
			} else if op == "remove" {
				delete(m, current)
			}
			return m, nil
		}
		sub, ok := m[current]
		if !ok {
			if op == "add" {
				sub = make(map[string]any)
				m[current] = sub
			} else {
				return nil, fmt.Errorf("path not found: %s", current)
			}
		}
		res, err := setPathValue(sub, path[1:], op, value)
		if err != nil {
			return nil, err
		}
		m[current] = res
		return m, nil
	}

	if arr, ok := obj.([]any); ok {
		idx := -1
		if current == "-" {
			idx = len(arr)
		} else {
			_, _ = fmt.Sscanf(current, "%d", &idx)
		}
		if idx < 0 || idx > len(arr) {
			return nil, fmt.Errorf("invalid array index: %s", current)
		}
		if len(path) == 1 {
			if op == "add" {
				if idx == len(arr) {
					arr = append(arr, value)
				} else {
					arr = append(arr[:idx+1], arr[idx:]...)
					arr[idx] = value
				}
			} else if op == "replace" {
				if idx < len(arr) {
					arr[idx] = value
				}
			} else if op == "remove" {
				if idx < len(arr) {
					arr = append(arr[:idx], arr[idx+1:]...)
				}
			}
			return arr, nil
		}
		if idx < len(arr) {
			res, err := setPathValue(arr[idx], path[1:], op, value)
			if err != nil {
				return nil, err
			}
			arr[idx] = res
			return arr, nil
		}
	}

	return obj, nil
}

type bufferedConn struct {
	net.Conn
	r io.Reader
}

func (bc *bufferedConn) Read(b []byte) (int, error) {
	return bc.r.Read(b)
}

type muxListener struct {
	net.Listener
	tlsConfig *tls.Config
}

func (l *muxListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	// Read 1 byte to check if it's TLS ClientHello (0x16)
	firstByte := make([]byte, 1)
	_ = c.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, err = io.ReadFull(c, firstByte)
	_ = c.SetReadDeadline(time.Time{}) // reset
	if err != nil {
		_ = c.Close()
		return nil, err
	}

	var r io.Reader = io.MultiReader(bytes.NewReader(firstByte), c)
	bc := &bufferedConn{Conn: c, r: r}

	if firstByte[0] == 0x16 {
		tlsConn := tls.Server(bc, l.tlsConfig)
		return tlsConn, nil
	}

	return bc, nil
}
