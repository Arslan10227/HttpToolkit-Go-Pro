package cert

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

// SPKIFingerprintSHA256 returns the base64 SHA-256 SPKI fingerprint used by Chromium
// (--ignore-certificate-errors-spki-list), matching mockttp generateSPKIFingerprint.
func SPKIFingerprintSHA256(certPEM string) (string, error) {
	block, _ := pem.Decode([]byte(strings.TrimSpace(certPEM)))
	if block == nil {
		return "", fmt.Errorf("invalid certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}
	spkiDER, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(spkiDER)
	return strings.TrimSpace(base64.StdEncoding.EncodeToString(sum[:])), nil
}

func (m *Manager) SPKIFingerprint() (string, error) {
	return SPKIFingerprintSHA256(m.certPEM)
}
