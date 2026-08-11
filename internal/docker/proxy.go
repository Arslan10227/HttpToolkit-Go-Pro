package docker

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

var (
	createContainerRe  = regexp.MustCompile(`^/v[\d.]+/containers/create`)
	startContainerRe   = regexp.MustCompile(`^/v[\d.]+/containers/([^/]+)/start`)
	buildImageRe       = regexp.MustCompile(`^/v[\d.]+/build`)
	containerInspectRe = regexp.MustCompile(`^/v[\d.]+/containers/[^/]+/json`)
	containerListRe    = regexp.MustCompile(`^/v[\d.]+/containers/json`)
)

// APIProxy forwards Docker API traffic with interception hooks.
type APIProxy struct {
	mu        sync.Mutex
	proxyPort int
	certPath  string
	assetsDir string
	ln        net.Listener
	server    *http.Server
	client    *http.Client
}

func NewAPIProxy(proxyPort int, certPath, assetsDir string) *APIProxy {
	return &APIProxy{
		proxyPort: proxyPort,
		certPath:  certPath,
		assetsDir: assetsDir,
		client:    dockerHTTPClient(),
	}
}

func (p *APIProxy) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.ln != nil {
		return nil
	}
	if !IsAvailable() {
		return nil
	}
	ln, err := listenProxy(ProxyListenPath(p.proxyPort))
	if err != nil {
		return err
	}
	p.ln = ln
	srv := &http.Server{Handler: http.HandlerFunc(p.serveHTTP)}
	if runtime.GOOS == "windows" {
		srv.ReadHeaderTimeout = 0
	}
	p.server = srv
	go func() { _ = srv.Serve(ln) }()
	return nil
}

func (p *APIProxy) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.server != nil {
		_ = p.server.Close()
		p.server = nil
	}
	if p.ln != nil {
		err := p.ln.Close()
		p.ln = nil
		return err
	}
	return nil
}

func (p *APIProxy) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if !IsAvailable() {
		http.Error(w, "HTTP Toolkit could not connect to Docker", http.StatusGatewayTimeout)
		return
	}
	if r.Header.Get("Upgrade") != "" {
		p.handleUpgrade(w, r)
		return
	}

	path := r.URL.Path
	query := r.URL.RawQuery
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = r.Body.Close()

	settings := InterceptionSettings{ProxyPort: p.proxyPort, CertPath: p.certPath}
	extraBuildCommands := 0

	if createContainerRe.MatchString(path) {
		body, err = TransformContainerCreateConfig(body, settings)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if m := startContainerRe.FindStringSubmatch(path); len(m) == 2 {
		if blocked, msg := p.blockNonInterceptedStart(m[1]); blocked {
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
	}

	if buildImageRe.MatchString(path) {
		u, _ := url.Parse("?" + query)
		if remote := u.Query().Get("remote"); remote != "" {
			if remote == "client-session" {
				http.Error(w, "HTTP Toolkit does not yet support BuildKit-powered builds", http.StatusBadRequest)
			} else {
				http.Error(w, "HTTP Toolkit does not support intercepting remote build sources", http.StatusBadRequest)
			}
			return
		}
		dockerfileName := u.Query().Get("dockerfile")
		injected, added, err := InjectBuildContext(body, dockerfileName, BuildSettings{
			ProxyPort: p.proxyPort,
			CertPath:  p.certPath,
			AssetsDir: p.assetsDir,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		body = injected
		extraBuildCommands = added
	}

	target := dockerAPIURL(path)
	if query != "" {
		target += "?" + query
	}
	upReq, err := http.NewRequestWithContext(r.Context(), r.Method, target, bytes.NewReader(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	upReq.Header = r.Header.Clone()
	upReq.Header.Del("Content-Length")
	upReq.ContentLength = int64(len(body))

	resp, err := p.client.Do(upReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	isInspect := containerInspectRe.MatchString(path)
	isComposeList := containerListRe.MatchString(path) && strings.Contains(query, "com.docker.compose")
	isBuild := buildImageRe.MatchString(path)
	if isInspect || isComposeList {
		if isInspect {
			respBody, _ = remapInspectResponse(p.proxyPort, respBody)
		} else {
			respBody, _ = remapListResponse(p.proxyPort, respBody)
		}
		resp.Header.Del("Content-Length")
	}
	if isBuild && resp.StatusCode == http.StatusOK {
		respBody = TransformBuildOutput(respBody, extraBuildCommands)
		resp.Header.Del("Content-Length")
	}

	for k, vals := range resp.Header {
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = w.Write(respBody)
}

func (p *APIProxy) blockNonInterceptedStart(containerID string) (bool, string) {
	inspect, err := inspectContainer(containerID)
	if err != nil || inspect == nil {
		return false, ""
	}
	if IsInterceptedContainer(inspect.Config.Labels, p.proxyPort) {
		return false, ""
	}
	return true, "HTTP Toolkit cannot intercept startup of preexisting non-intercepted containers. " +
		"The container must be recreated here first - try `docker run <image>` instead."
}

func (p *APIProxy) handleUpgrade(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "hijack unsupported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hj.Hijack()
	if err != nil {
		return
	}
	defer clientConn.Close()

	path := r.URL.Path
	query := r.URL.RawQuery
	target := dockerAPIURL(path)
	if query != "" {
		target += "?" + query
	}

	upReq, err := http.NewRequestWithContext(r.Context(), r.Method, target, nil)
	if err != nil {
		return
	}
	upReq.Header = r.Header.Clone()

	upConn, err := dialDockerDaemon(r.Context(), "", "")
	if err != nil {
		return
	}
	defer upConn.Close()

	if err := upReq.Write(upConn); err != nil {
		return
	}

	upResp, err := http.ReadResponse(bufio.NewReader(upConn), upReq)
	if err != nil {
		return
	}
	defer upResp.Body.Close()

	statusLine := fmt.Sprintf("HTTP/1.1 %s\r\n", upResp.Status)
	_, _ = clientConn.Write([]byte(statusLine))
	for k, vals := range upResp.Header {
		for _, v := range vals {
			_, _ = clientConn.Write([]byte(k + ": " + v + "\r\n"))
		}
	}
	_, _ = clientConn.Write([]byte("\r\n"))

	if upResp.StatusCode == http.StatusSwitchingProtocols {
		go func() { _, _ = io.Copy(upConn, clientConn) }()
		_, _ = io.Copy(clientConn, upConn)
		return
	}
	_, _ = io.Copy(clientConn, upResp.Body)
}
