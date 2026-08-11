package api

import (
	"net"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
)

// ──────────────────────────────────────────────────────────────────────────
// gql.Provider implementation on *Server
// ──────────────────────────────────────────────────────────────────────────

func (s *Server) Version() string {
	return config.ServerVersion()
}

func (s *Server) CertificatePath() string {
	return s.certs.CertPath()
}

func (s *Server) CertificateContent() string {
	return s.certs.CertPEM()
}

func (s *Server) CertificateFingerprint() string {
	return s.spki
}

func (s *Server) NetworkInterfaces() map[string]any {
	ifaces, _ := net.Interfaces()
	out := make(map[string]any, len(ifaces))
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		entries := make([]any, 0, len(addrs))
		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil {
				continue
			}
			family := "IPv6"
			if v4 := ip.To4(); v4 != nil {
				ip = v4
				family = "IPv4"
			}
			entries = append(entries, map[string]any{
				"address":  ip.String(),
				"family":   family,
				"internal": ip.IsLoopback() || ip.IsLinkLocalUnicast(),
			})
		}
		if len(entries) > 0 {
			out[iface.Name] = entries
		}
	}
	return out
}

func (s *Server) SystemProxy() map[string]any {
	return systemProxy()
}

func (s *Server) DNSServers(proxyPort int) []string {
	if s.sessions == nil {
		return []string{}
	}
	return s.sessions.DNSServers()
}

func (s *Server) RuleParameterKeys() []string {
	if s.sessions == nil {
		return []string{}
	}
	return s.sessions.RuleParameterKeys()
}

func (s *Server) GetInterceptors(proxyPort int) []map[string]any {
	return s.interceptors.List(proxyPort)
}

func (s *Server) GetInterceptor(id string) (map[string]any, error) {
	list := s.interceptors.List(0)
	for _, item := range list {
		if item["id"] == id {
			return item, nil
		}
	}
	return nil, nil
}

func (s *Server) IsInterceptorActive(id string, proxyPort int) bool {
	list := s.interceptors.List(proxyPort)
	for _, item := range list {
		if item["id"] == id {
			if active, ok := item["isActive"].(bool); ok {
				return active
			}
		}
	}
	return false
}

func (s *Server) InterceptorMetadata(id string, metaType string) (any, error) {
	return s.interceptors.Metadata(id)
}

func (s *Server) ActivateInterceptor(id string, proxyPort int, options map[string]any) (any, error) {
	result, err := s.interceptors.Activate(id, proxyPort, options)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) DeactivateInterceptor(id string, proxyPort int, options map[string]any) (bool, error) {
	result, err := s.interceptors.Deactivate(id, proxyPort, options)
	if err != nil {
		return false, err
	}
	if success, ok := result["success"].(bool); ok {
		return success, nil
	}
	return true, nil
}

func (s *Server) TriggerUpdate() {
	// Update signal: no-op in embedded mode; extend as needed.
}

func (s *Server) TriggerShutdown() {
	s.shutdownOnce.Do(func() {
		if s.onShutdown != nil {
			go s.onShutdown()
		}
	})
}
