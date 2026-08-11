package cert

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
)

func TestNewManagerGeneratesCACert(t *testing.T) {
	m := newTestManager(t)

	block, _ := pem.Decode([]byte(m.CertPEM()))
	if block == nil {
		t.Fatal("cert PEM invalid")
	}
	ca, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse cert: %v", err)
	}

	if !ca.IsCA {
		t.Error("CA cert is not marked as CA")
	}
	if ca.SignatureAlgorithm != x509.SHA256WithRSA {
		t.Errorf("CA signature algorithm = %v, want SHA256WithRSA", ca.SignatureAlgorithm)
	}
	if len(ca.SubjectKeyId) == 0 {
		t.Error("CA cert missing SubjectKeyId")
	}
	now := time.Now()
	if now.Before(ca.NotBefore) || now.After(ca.NotAfter) {
		t.Error("CA cert not valid for current time")
	}
}

func TestLeafCertChain(t *testing.T) {
	m := newTestManager(t)

	cert, err := m.LeafCertChain("example.com")
	if err != nil {
		t.Fatalf("LeafCertChain: %v", err)
	}
	if len(cert.Certificate) != 2 {
		t.Fatalf("expected leaf + CA chain, got %d certs", len(cert.Certificate))
	}

	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatalf("parse leaf: %v", err)
	}
	ca, err := x509.ParseCertificate(cert.Certificate[1])
	if err != nil {
		t.Fatalf("parse ca: %v", err)
	}

	if leaf.Subject.CommonName != "example.com" {
		t.Errorf("leaf CN = %q, want example.com", leaf.Subject.CommonName)
	}
	if len(leaf.DNSNames) != 1 || leaf.DNSNames[0] != "example.com" {
		t.Errorf("leaf DNSNames = %v, want [example.com]", leaf.DNSNames)
	}
	if len(leaf.SubjectKeyId) == 0 {
		t.Error("leaf missing SubjectKeyId")
	}
	if string(leaf.AuthorityKeyId) != string(ca.SubjectKeyId) {
		t.Error("leaf AuthorityKeyId does not match CA SubjectKeyId")
	}
	if len(leaf.ExtKeyUsage) != 1 || leaf.ExtKeyUsage[0] != x509.ExtKeyUsageServerAuth {
		t.Errorf("leaf ExtKeyUsage = %v", leaf.ExtKeyUsage)
	}
	if leaf.SignatureAlgorithm != x509.SHA256WithRSA {
		t.Errorf("leaf signature algorithm = %v, want SHA256WithRSA", leaf.SignatureAlgorithm)
	}
}

func TestLeafCertChainForIP(t *testing.T) {
	m := newTestManager(t)

	cert, err := m.LeafCertChain("127.0.0.1")
	if err != nil {
		t.Fatalf("LeafCertChain: %v", err)
	}
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatalf("parse leaf: %v", err)
	}
	if len(leaf.IPAddresses) != 1 || leaf.IPAddresses[0].String() != "127.0.0.1" {
		t.Errorf("leaf IPAddresses = %v", leaf.IPAddresses)
	}
}

func TestSPKIFingerprintMatchesRaw(t *testing.T) {
	m := newTestManager(t)
	spki, err := m.SPKIFingerprint()
	if err != nil {
		t.Fatalf("SPKIFingerprint: %v", err)
	}

	block, _ := pem.Decode([]byte(m.CertPEM()))
	ca, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse cert: %v", err)
	}
	rawSum := sha256.Sum256(ca.RawSubjectPublicKeyInfo)
	raw := base64.StdEncoding.EncodeToString(rawSum[:])
	if raw != spki {
		t.Errorf("SPKI mismatch: %q vs %q", raw, spki)
	}
}

// TestLeafCertChainCaching verifies that LeafCertChain returns the same
// tls.Certificate for the same host on repeated calls (i.e., the cache
// works and we don't regenerate a new RSA key pair on every call).
func TestLeafCertChainCaching(t *testing.T) {
	m := newTestManager(t)

	cert1, err := m.LeafCertChain("example.com")
	if err != nil {
		t.Fatalf("LeafCertChain #1: %v", err)
	}
	cert2, err := m.LeafCertChain("example.com")
	if err != nil {
		t.Fatalf("LeafCertChain #2: %v", err)
	}

	// The cached cert should be identical (same leaf DER bytes).
	if len(cert1.Certificate) == 0 || len(cert2.Certificate) == 0 {
		t.Fatal("certificate chain is empty")
	}
	if len(cert1.Certificate[0]) != len(cert2.Certificate[0]) {
		t.Fatalf("leaf cert DER length differs: %d vs %d", len(cert1.Certificate[0]), len(cert2.Certificate[0]))
	}
	for i := range cert1.Certificate[0] {
		if cert1.Certificate[0][i] != cert2.Certificate[0][i] {
			t.Fatal("cached leaf cert DER bytes differ — cache not working")
		}
	}

	// A different host should get a different cert.
	cert3, err := m.LeafCertChain("other.com")
	if err != nil {
		t.Fatalf("LeafCertChain other.com: %v", err)
	}
	if len(cert3.Certificate[0]) == len(cert1.Certificate[0]) {
		// Check if they're actually the same (unlikely but possible with small key sizes)
		same := true
		for i := range cert3.Certificate[0] {
			if cert3.Certificate[0][i] != cert1.Certificate[0][i] {
				same = false
				break
			}
		}
		if same {
			t.Fatal("different host returned same leaf cert — cache key is wrong")
		}
	}
}

func newTestManager(t *testing.T) *Manager {
	t.Helper()
	cfg := &config.Config{ConfigDir: t.TempDir()}
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return m
}
