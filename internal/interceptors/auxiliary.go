package interceptors

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

var auxServers sync.Map // key: port -> *http.Server

func stopAuxiliaryServers() {
	auxServers.Range(func(key, value any) bool {
		if srv, ok := value.(*http.Server); ok {
			_ = srv.Close()
		}
		auxServers.Delete(key)
		return true
	})
}

func startMessageServer() int {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("HttpToolkit interception active\n"))
	})
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0
	}
	port := ln.Addr().(*net.TCPAddr).Port
	srv := &http.Server{Handler: mux}
	auxServers.Store(fmt.Sprintf("message-%d", port), srv)
	go func() { _ = srv.Serve(ln) }()
	return port
}

type waitableSetupServer struct {
	port  int
	ready chan struct{}
	srv   *http.Server
}

func startWaitableSetupServer() *waitableSetupServer {
	ready := make(chan struct{}, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		select {
		case ready <- struct{}{}:
		default:
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<!DOCTYPE html><html><body><p>HTTP Toolkit is preparing a Firefox profile, please wait...</p></body></html>"))
	})
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return &waitableSetupServer{ready: ready}
	}
	port := ln.Addr().(*net.TCPAddr).Port
	srv := &http.Server{Handler: mux}
	auxServers.Store(fmt.Sprintf("setup-%d", port), srv)
	go func() { _ = srv.Serve(ln) }()
	return &waitableSetupServer{port: port, ready: ready, srv: srv}
}

func (s *waitableSetupServer) url() string {
	return fmt.Sprintf("http://127.0.0.1:%d/", s.port)
}

func (s *waitableSetupServer) wait(timeout time.Duration) error {
	if s.port == 0 {
		return fmt.Errorf("setup server failed to bind")
	}
	select {
	case <-s.ready:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timed out waiting for firefox")
	}
}

func (s *waitableSetupServer) stop() {
	if s.srv != nil {
		_ = s.srv.Close()
		auxServers.Delete(fmt.Sprintf("setup-%d", s.port))
	}
}
