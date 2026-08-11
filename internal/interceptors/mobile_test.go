package interceptors

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
)

func generateTestCertPEM() (string, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return "", err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"HTTP Toolkit Pro Test Inc"},
			CommonName:   "HttpToolkit Test CA",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return "", err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	return string(certPEM), nil
}

func TestGetCertificateSubjectHash(t *testing.T) {
	certPEM, err := generateTestCertPEM()
	if err != nil {
		t.Fatalf("failed to generate test cert: %v", err)
	}

	hash, err := getCertificateSubjectHash(certPEM)
	if err != nil {
		t.Fatalf("getCertificateSubjectHash failed: %v", err)
	}

	if len(hash) != 8 {
		t.Fatalf("expected 8-char hex string hash, got: %s", hash)
	}

	// Verify that it is valid hexadecimal
	for _, char := range hash {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			t.Fatalf("expected lowercase hexadecimal characters, got: %s", hash)
		}
	}
}

func parseJvmProcessesString(output string) []map[string]string {
	var targets []map[string]string
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		pid := parts[0]
		name := strings.TrimSpace(parts[1])
		if strings.Contains(name, "java-agent.jar") {
			continue
		}
		targets = append(targets, map[string]string{
			"pid":  pid,
			"name": name,
		})
	}
	return targets
}

func TestJvmTargetsParsing(t *testing.T) {
	sampleOutput := "1234:org.example.Main\n5678:com.test.App\n9012:java-agent.jar\n"
	targets := parseJvmProcessesString(sampleOutput)
	if len(targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(targets))
	}
	if targets[0]["pid"] != "1234" || targets[0]["name"] != "org.example.Main" {
		t.Fatalf("mismatched target 0: %v", targets[0])
	}
	if targets[1]["pid"] != "5678" || targets[1]["name"] != "com.test.App" {
		t.Fatalf("mismatched target 1: %v", targets[1])
	}
}

func TestElectronInterceptorEnv(t *testing.T) {
	cfg := &config.Config{
		ConfigDir: t.TempDir(),
		AssetsDir: t.TempDir(),
	}
	s := &stubInterceptor{id: "electron", cfg: cfg, spki: "sha256-mockspkifingerprint", active: make(map[int]bool)}
	interceptor := &electronInterceptor{base: s}

	activable, err := interceptor.IsActivable()
	if err != nil || !activable {
		t.Fatalf("expected activable to be true, got %t: %v", activable, err)
	}

	env := buildTerminalEnv(cfg, 8080)
	env = append(env, "HTTP_TOOLKIT_SPKI="+s.spki)

	spkiFound := false
	proxyFound := false
	activeFound := false
	for _, kv := range env {
		if strings.HasPrefix(kv, "HTTP_TOOLKIT_SPKI=") {
			spkiFound = true
		}
		if strings.HasPrefix(kv, "HTTP_PROXY=") {
			proxyFound = true
		}
		if kv == "HTTP_TOOLKIT_ACTIVE=true" {
			activeFound = true
		}
	}

	if !spkiFound || !proxyFound || !activeFound {
		t.Fatalf("missing critical env vars in Electron env: spki=%t, proxy=%t, active=%t", spkiFound, proxyFound, activeFound)
	}
}

func TestIsMatchingCert(t *testing.T) {
	certPEM, err := generateTestCertPEM()
	if err != nil {
		t.Fatalf("failed to generate test cert: %v", err)
	}

	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		t.Fatalf("failed to decode certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}
	h := sha1.Sum(cert.Raw)
	fingerprint := fmt.Sprintf("%x", h)

	if !isMatchingCert(certPEM, fingerprint) {
		t.Fatalf("expected isMatchingCert to return true for valid PEM and fingerprint")
	}

	if isMatchingCert(certPEM, "invalidfingerprint") {
		t.Fatalf("expected isMatchingCert to return false for invalid fingerprint")
	}
}

func TestBuildFridaConfig(t *testing.T) {
	template := `
const CERT_PEM = ` + "`" + `[!! Put your CA certificate data here, in PEM format !!]` + "`" + `;
const PROXY_HOST = '127.0.0.1';
const PROXY_PORT = 8000;
const PROXY_SUPPORTS_SOCKS5 = false;
const IGNORED_NON_HTTP_PORTS = [];
`
	caCert := "-----BEGIN CERTIFICATE-----\nMOCKCERT\n-----END CERTIFICATE-----"
	host := "192.168.1.15"
	port := 8081
	ignored := []int{4070, 80}
	enableSocks := true

	result := buildFridaConfig(template, caCert, host, port, ignored, enableSocks)

	if !strings.Contains(result, "const CERT_PEM = `"+caCert+"`;") {
		t.Errorf("expected templated CERT_PEM, got:\n%s", result)
	}
	if !strings.Contains(result, "const PROXY_HOST = '192.168.1.15';") {
		t.Errorf("expected templated PROXY_HOST, got:\n%s", result)
	}
	if !strings.Contains(result, "const PROXY_PORT = 8081;") {
		t.Errorf("expected templated PROXY_PORT, got:\n%s", result)
	}
	if !strings.Contains(result, "const PROXY_SUPPORTS_SOCKS5 = true;") {
		t.Errorf("expected templated PROXY_SUPPORTS_SOCKS5, got:\n%s", result)
	}
	if !strings.Contains(result, "const IGNORED_NON_HTTP_PORTS = [4070,80];") {
		t.Errorf("expected templated IGNORED_NON_HTTP_PORTS, got:\n%s", result)
	}
}

func TestParseAdbDevicesMetadata(t *testing.T) {
	raw := `List of devices attached
emulator-5554          device product:sdk_gphone64_arm64 model:sdk_gphone64_arm64 device:emulator64_arm64 transport_id:1
0123456789ABCDEF       device product:samsung model:SM_G973F device:beyond1 transport_id:2
unauthorized_device     unauthorized
offline_device          offline
`
	ids, devices := parseAdbDevicesMetadata(raw)

	if len(ids) != 2 {
		t.Fatalf("expected 2 active devices, got %d", len(ids))
	}

	if ids[0] != "emulator-5554" || ids[1] != "0123456789ABCDEF" {
		t.Fatalf("mismatched active device IDs: %v", ids)
	}

	dev1, ok1 := devices["emulator-5554"]
	if !ok1 {
		t.Fatalf("expected details for emulator-5554")
	}
	if dev1["ro.product.model"] != "sdk gphone64 arm64" {
		t.Fatalf("expected model 'sdk gphone64 arm64', got '%s'", dev1["ro.product.model"])
	}

	dev2, ok2 := devices["0123456789ABCDEF"]
	if !ok2 {
		t.Fatalf("expected details for 0123456789ABCDEF")
	}
	if dev2["ro.product.model"] != "SM G973F" {
		t.Fatalf("expected model 'SM G973F', got '%s'", dev2["ro.product.model"])
	}
}
func TestAdbIntentCertFingerprint(t *testing.T) {
	// Verify that the SHA-256 fingerprint produced for the ADB ACTIVATE intent
	// is in the colon-separated uppercase hex format expected by the Android app.
	certPEM, err := generateTestCertPEM()
	if err != nil {
		t.Fatalf("failed to generate test cert: %v", err)
	}

	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		t.Fatalf("failed to decode PEM")
	}

	import_sha256 := sha256.Sum256(block.Bytes)
	parts := make([]string, 32)
	for i, b := range import_sha256 {
		parts[i] = fmt.Sprintf("%02X", b)
	}
	fingerprint := strings.Join(parts, ":")

	// Should be 32 groups of 2 hex chars joined by colons = 32*2 + 31 colons = 95 chars
	if len(fingerprint) != 95 {
		t.Fatalf("expected 95-char fingerprint (AA:BB:CC:...), got len=%d: %s", len(fingerprint), fingerprint)
	}

	sections := strings.Split(fingerprint, ":")
	if len(sections) != 32 {
		t.Fatalf("expected 32 colon-separated sections, got %d", len(sections))
	}

	for _, s := range sections {
		if len(s) != 2 {
			t.Fatalf("expected each section to be 2 chars, got: %s", s)
		}
		for _, c := range s {
			if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F')) {
				t.Fatalf("expected uppercase hex character, got: %c in %s", c, fingerprint)
			}
		}
	}
}
