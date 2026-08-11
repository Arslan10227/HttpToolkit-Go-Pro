package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"software.sslmate.com/src/go-pkcs12"
)

type Manager struct {
	dir         string
	certPath    string
	keyPath     string
	certPEM     string
	keyPEM      []byte
	fingerprint string

	leafCacheMu sync.RWMutex
	leafCache   map[string]tls.Certificate
}

func NewManager(cfg *config.Config) (*Manager, error) {
	m := &Manager{
		dir:       cfg.ConfigDir,
		certPath:  filepath.Join(cfg.ConfigDir, "ca.pem"),
		keyPath:   filepath.Join(cfg.ConfigDir, "ca.key"),
		leafCache: make(map[string]tls.Certificate),
	}
	if err := os.MkdirAll(cfg.ConfigDir, 0o700); err != nil {
		return nil, err
	}
	if err := m.loadOrGenerate(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Manager) CertPath() string    { return m.certPath }
func (m *Manager) KeyPath() string     { return m.keyPath }
func (m *Manager) CertPEM() string     { return m.certPEM }
func (m *Manager) KeyPEM() []byte      { return m.keyPEM }
func (m *Manager) Fingerprint() string { return m.fingerprint }
func (m *Manager) CertFiles() map[string]string {
	return map[string]string{
		"pem": m.certPath,
		"crt": filepath.Join(m.dir, "ca.crt"),
		"cer": filepath.Join(m.dir, "ca.cer"),
		"p12": filepath.Join(m.dir, "ca.p12"),
	}
}

func (m *Manager) loadOrGenerate() error {
	certData, certErr := os.ReadFile(m.certPath)
	keyData, keyErr := os.ReadFile(m.keyPath)
	if certErr == nil && keyErr == nil {
		if err := m.validateExisting(certData); err != nil {
			return m.generate()
		}
		m.certPEM = string(certData)
		m.keyPEM = keyData
		m.fingerprint = sha256Fingerprint(certData)
		return m.exportFormats()
	}
	return m.generate()
}

func (m *Manager) validateExisting(certData []byte) error {
	block, _ := pem.Decode(certData)
	if block == nil {
		return fmt.Errorf("invalid pem")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}
	if cert.Subject.CommonName != config.CACommonName {
		return fmt.Errorf("subject mismatch")
	}
	if time.Until(cert.NotAfter) < 48*time.Hour {
		return fmt.Errorf("expiring soon")
	}
	return nil
}

func (m *Manager) generate() error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName:   config.CACommonName,
			Organization: []string{config.CACommonName},
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	// Compute a stable Subject Key Identifier so leaf certs can reference the CA
	// via AuthorityKeyId and Chrome can build a complete chain.
	template.SubjectKeyId = subjectKeyID(&key.PublicKey)
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	if err := os.WriteFile(m.certPath, certPEM, 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(m.keyPath, keyPEM, 0o600); err != nil {
		return err
	}
	m.certPEM = string(certPEM)
	m.keyPEM = keyPEM
	m.fingerprint = sha256Fingerprint(certPEM)
	return m.exportFormats()
}

func (m *Manager) exportFormats() error {
	certData := []byte(m.certPEM)
	_ = os.WriteFile(filepath.Join(m.dir, "ca.crt"), certData, 0o644)
	_ = os.WriteFile(filepath.Join(m.dir, "ca.cer"), certData, 0o644)
	if p12, err := m.ExportP12("httptoolkit"); err == nil {
		_ = os.WriteFile(filepath.Join(m.dir, "ca.p12"), p12, 0o644)
	}
	return nil
}

// ExportP12 returns a PKCS#12 bundle of the CA cert and key.
func (m *Manager) ExportP12(password string) ([]byte, error) {
	block, _ := pem.Decode([]byte(m.certPEM))
	if block == nil {
		return nil, fmt.Errorf("invalid ca pem")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	keyBlock, _ := pem.Decode(m.keyPEM)
	if keyBlock == nil {
		return nil, fmt.Errorf("invalid ca key")
	}
	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return pkcs12.Modern.Encode(key, cert, nil, password)
}

// LeafCertChain returns a tls.Certificate for the given host, signed by the
// CA and with the CA appended to the chain. Results are cached per host to
// avoid regenerating a leaf cert (and a new RSA key pair) on every TLS
// handshake — the leaf cert is valid for 365 days so caching is safe.
func (m *Manager) LeafCertChain(host string) (tls.Certificate, error) {
	// Fast path: check the cache.
	m.leafCacheMu.RLock()
	if cached, ok := m.leafCache[host]; ok {
		m.leafCacheMu.RUnlock()
		return cached, nil
	}
	m.leafCacheMu.RUnlock()

	certPEM, keyPEM, err := m.LeafCert(host)
	if err != nil {
		return tls.Certificate{}, err
	}
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return tls.Certificate{}, err
	}
	if block, _ := pem.Decode([]byte(m.certPEM)); block != nil {
		// Append the CA certificate so clients that rely on SPKI pinning or
		// an otherwise unknown root can build a complete chain.
		cert.Certificate = append(cert.Certificate, block.Bytes)
	}

	// Store in cache.
	m.leafCacheMu.Lock()
	m.leafCache[host] = cert
	m.leafCacheMu.Unlock()

	return cert, nil
}

func (m *Manager) LeafCert(host string) (certPEM, keyPEM []byte, err error) {
	block, _ := pem.Decode([]byte(m.certPEM))
	if block == nil {
		return nil, nil, fmt.Errorf("ca cert invalid")
	}
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, err
	}
	caKeyBlock, _ := pem.Decode(m.keyPEM)
	if caKeyBlock == nil {
		return nil, nil, fmt.Errorf("ca key invalid")
	}
	caKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	leafKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName: host,
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		AuthorityKeyId:        caCert.SubjectKeyId,
		SubjectKeyId:          subjectKeyID(&leafKey.PublicKey),
		BasicConstraintsValid: true,
	}
	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = []net.IP{ip}
	} else {
		template.DNSNames = []string{host}
	}
	der, err := x509.CreateCertificate(rand.Reader, template, caCert, &leafKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(leafKey)})
	return certPEM, keyPEM, nil
}

func subjectKeyID(pub *rsa.PublicKey) []byte {
	b, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil
	}
	sum := sha1.Sum(b)
	return sum[:]
}

func sha256Fingerprint(certPEM []byte) string {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return ""
	}
	sum := sha256.Sum256(block.Bytes)
	return hex.EncodeToString(sum[:])
}

// DescribeCertDER returns a human-readable summary of a certificate from its
// DER bytes. It is used for TLS handshake diagnostics.
func DescribeCertDER(der []byte) map[string]any {
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	return map[string]any{
		"subject":        cert.Subject.String(),
		"issuer":         cert.Issuer.String(),
		"notBefore":      cert.NotBefore.Format(time.RFC3339),
		"notAfter":       cert.NotAfter.Format(time.RFC3339),
		"dnsNames":       cert.DNSNames,
		"ipAddresses":    cert.IPAddresses,
		"isCA":           cert.IsCA,
		"keyUsage":       cert.KeyUsage.String(),
		"extKeyUsage":    cert.ExtKeyUsage,
		"subjectKeyId":   hex.EncodeToString(cert.SubjectKeyId),
		"authorityKeyId": hex.EncodeToString(cert.AuthorityKeyId),
		"serialNumber":   cert.SerialNumber.String(),
		"signatureAlg":   cert.SignatureAlgorithm.String(),
		"publicKeyAlg":   cert.PublicKeyAlgorithm.String(),
		"fingerprint":    hex.EncodeToString(sha256.New().Sum(nil)),
	}
}
