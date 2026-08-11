package interceptors

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	dockersvc "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/docker"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/system"
)

var execCommand = system.Command

// probeTimeout bounds short-lived diagnostic/probe commands (frida-ps, adb devices,
// jps, etc.) so that a misbehaving or hung external tool can never block interceptor
// Metadata()/List() calls indefinitely (these are invoked from the REST API on every
// UI refresh of the interceptors list).
const probeTimeout = 5 * time.Second

// runWithTimeout starts cmd and kills it if it hasn't exited within timeout,
// returning a timeout error in that case.
func runWithTimeout(cmd *exec.Cmd, timeout time.Duration) error {
	if err := cmd.Start(); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		<-done
		return fmt.Errorf("%s: timed out after %s", cmd.Path, timeout)
	}
}

// execRunTimeout runs name with args, discarding output, bounded by timeout.
func execRunTimeout(timeout time.Duration, name string, args ...string) error {
	return runWithTimeout(execCommand(name, args...), timeout)
}

// execOutputTimeout runs name with args and captures stdout, bounded by timeout.
func execOutputTimeout(timeout time.Duration, name string, args ...string) ([]byte, error) {
	cmd := execCommand(name, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err := runWithTimeout(cmd, timeout)
	return buf.Bytes(), err
}

// execCombinedOutputTimeout runs name with args and captures stdout+stderr, bounded by timeout.
func execCombinedOutputTimeout(timeout time.Duration, name string, args ...string) ([]byte, error) {
	cmd := execCommand(name, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := runWithTimeout(cmd, timeout)
	return buf.Bytes(), err
}

type jvmInterceptor struct {
	base        *stubInterceptor
	mu          sync.Mutex
	intercepted map[string]int // maps PID -> proxyPort
}

func (j *jvmInterceptor) IsActivable() (bool, error) {
	if _, err := os.Stat(javaAgentPath(j.base.cfg)); err != nil {
		return false, nil
	}
	// Attach requires a JDK (com.sun.tools.attach / jdk.attach module).
	return system.FindJDK() != "", nil
}

func (j *jvmInterceptor) IsActive(p int) (bool, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.intercepted == nil {
		return false, nil
	}
	for _, port := range j.intercepted {
		if port == p {
			return true, nil
		}
	}
	return false, nil
}

func (j *jvmInterceptor) Activate(proxyPort int, options map[string]any) (map[string]any, error) {
	pidVal := optionFloat(options, "pid")
	if pidVal == 0 {
		pidVal = optionFloat(options, "targetPid")
	}
	pidStr, _ := options["targetPid"].(string)
	if pidStr == "" {
		pidStr, _ = options["pid"].(string)
	}
	if pidStr == "" && pidVal != 0 {
		pidStr = fmt.Sprintf("%.0f", pidVal)
	}
	if pidStr == "" {
		return nil, fmt.Errorf("pid required for JVM attach")
	}

	agent := javaAgentPath(j.base.cfg)
	certPath := filepath.Join(j.base.cfg.ConfigDir, "ca.pem")

	javaExe := system.FindJDK()
	if javaExe == "" {
		return nil, fmt.Errorf("JVM attach requires a JDK (not a JRE). Install a JDK or set JAVA_HOME to a JDK.")
	}

	out, err := execCombinedOutputTimeout(30*time.Second, javaExe, "-jar", agent, pidStr, "127.0.0.1", fmt.Sprint(proxyPort), certPath)

	attachOutput := string(out)
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = -1
		}
	}

	attachSucceeded := err == nil || (exitCode == 3 && (strings.Contains(strings.ToLower(attachOutput), "failed to load agent library: 0") ||
		strings.Contains(strings.ToLower(attachOutput), "agentloadexception: 0")))

	if !attachSucceeded {
		return nil, friendlyAttachError(exitCode, attachOutput)
	}

	j.mu.Lock()
	if j.intercepted == nil {
		j.intercepted = make(map[string]int)
	}
	j.intercepted[pidStr] = proxyPort
	j.mu.Unlock()

	return map[string]any{"attached": true, "pid": pidStr}, nil
}

func (j *jvmInterceptor) Deactivate(proxyPort int, options map[string]any) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.intercepted == nil {
		return nil
	}

	pidStr, _ := options["targetPid"].(string)
	if pidStr == "" {
		pidStr, _ = options["pid"].(string)
	}
	if pidStr != "" {
		delete(j.intercepted, pidStr)
		return nil
	}

	for pid, port := range j.intercepted {
		if port == proxyPort {
			delete(j.intercepted, pid)
		}
	}
	return nil
}

func (j *jvmInterceptor) Metadata(kind string) (any, error) {
	if kind == "detailed" {
		processes := listJvmProcesses(javaAgentPath(j.base.cfg))
		jvmTargets := make(map[string]any)
		j.mu.Lock()
		defer j.mu.Unlock()
		for _, proc := range processes {
			pid := proc["pid"]
			name := proc["name"]
			target := map[string]any{
				"pid":  pid,
				"name": name,
			}
			if port, ok := j.intercepted[pid]; ok {
				target["interceptedByProxy"] = port
			}
			jvmTargets[pid] = target
		}
		return map[string]any{
			"agentPath":  javaAgentPath(j.base.cfg),
			"jvmTargets": jvmTargets,
		}, nil
	}
	return map[string]any{"agentPath": javaAgentPath(j.base.cfg)}, nil
}

func (j *jvmInterceptor) SubMetadata(subID string) (any, error) {
	if subID == "targets" || subID == "processes" {
		return map[string]any{"targets": listJvmProcesses(javaAgentPath(j.base.cfg))}, nil
	}
	return j.Metadata("detailed")
}

func listJvmProcesses(agentPath string) []map[string]string {
	var targets []map[string]string
	pidsSeen := make(map[string]bool)

	javaExe := system.FindJava()
	// Stage 1: Try running `java -jar java-agent.jar list-targets`
	out, err := execCombinedOutputTimeout(probeTimeout, javaExe, "-jar", agentPath, "list-targets")
	if err == nil {
		for _, line := range strings.Split(string(out), "\n") {
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
			pidsSeen[pid] = true
		}
	}

	// Stage 2: Try running `jps -l` (since `jps` is on path in standard JDK environments)
	jpsOut, jpsErr := execCombinedOutputTimeout(probeTimeout, "jps", "-l")
	if jpsErr == nil {
		for _, line := range strings.Split(string(jpsOut), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) < 1 {
				continue
			}
			pid := parts[0]
			if pidsSeen[pid] {
				continue
			}
			name := "Java Process"
			if len(parts) >= 2 {
				name = parts[1]
			}
			if strings.Contains(name, "sun.tools.jps.Jps") || strings.Contains(name, "java-agent.jar") {
				continue
			}
			targets = append(targets, map[string]string{
				"pid":  pid,
				"name": name,
			})
			pidsSeen[pid] = true
		}
	}

	// Stage 3: On Windows, fallback to querying powershell for java/javaw processes
	if isWindows() {
		psOut, psErr := execCombinedOutputTimeout(probeTimeout, "powershell", "-NoProfile", "-Command", "Get-Process -Name java, javaw -ErrorAction SilentlyContinue | ForEach-Object { \"$($_.Id):$($_.Path)\" }")
		if psErr == nil {
			for _, line := range strings.Split(string(psOut), "\n") {
				line = strings.TrimSpace(line)
				if line == "" || !strings.Contains(line, ":") {
					continue
				}
				parts := strings.SplitN(line, ":", 2)
				pid := parts[0]
				if pidsSeen[pid] {
					continue
				}
				pathVal := strings.TrimSpace(parts[1])
				name := "Java Process"
				if pathVal != "" {
					name = filepath.Base(pathVal)
				}
				targets = append(targets, map[string]string{
					"pid":  pid,
					"name": name,
				})
				pidsSeen[pid] = true
			}
		}
	} else {
		// Linux / macOS fallback using ps
		psOut, psErr := execCombinedOutputTimeout(probeTimeout, "sh", "-c", "ps -eo pid,comm | grep -E 'java|javaw'")
		if psErr == nil {
			for _, line := range strings.Split(string(psOut), "\n") {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				parts := strings.Fields(line)
				if len(parts) < 2 {
					continue
				}
				pid := parts[0]
				if pidsSeen[pid] {
					continue
				}
				name := strings.Join(parts[1:], " ")
				targets = append(targets, map[string]string{
					"pid":  pid,
					"name": name,
				})
				pidsSeen[pid] = true
			}
		}
	}

	return targets
}

func friendlyAttachError(exitCode int, output string) error {
	lower := strings.ToLower(output)
	switch {
	case strings.Contains(lower, "jvm.dll not loaded"):
		return fmt.Errorf("JVM attach failed: 32/64-bit mismatch or the target is not a Java process")
	case strings.Contains(lower, "could not start") && strings.Contains(lower, "jdk"):
		return fmt.Errorf("JVM attach failed: the Java executable is not a JDK. Install a JDK or set JAVA_HOME.")
	case strings.Contains(lower, "can't scan for attachable jvms"):
		return fmt.Errorf("JVM attach failed: the JDK's attach module is not available. Use a JDK 8+.")
	case exitCode == 3:
		return fmt.Errorf("JVM attach failed (exit code 3): %s", strings.TrimSpace(output))
	default:
		return fmt.Errorf("JVM attach failed with exit code %d: %s", exitCode, strings.TrimSpace(output))
	}
}

var (
	adbTunnelsMutex sync.Mutex
	adbTunnels      = make(map[string]chan struct{}) // maps deviceId:port -> stopChannel
)

func getCertificateSubjectHash(certPEM string) (string, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}
	h := md5.Sum(cert.RawSubject)
	val := uint32(h[0]) | uint32(h[1])<<8 | uint32(h[2])<<16 | uint32(h[3])<<24
	return fmt.Sprintf("%08x", val), nil
}

func readAdbFile(options map[string]any, remotePath string) (string, error) {
	cmd := adbCommand(options, "shell", "cat", remotePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func isMatchingCert(certPEM string, expectedFingerprint string) bool {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return false
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false
	}
	h := sha1.Sum(cert.Raw)
	fingerprint := fmt.Sprintf("%x", h)
	return strings.ToLower(fingerprint) == strings.ToLower(expectedFingerprint)
}

func hasCertInstalled(options map[string]any, certHash string, expectedFingerprint string) bool {
	if expectedFingerprint == "" {
		return false
	}

	systemCertPath := fmt.Sprintf("/system/etc/security/cacerts/%s.0", certHash)
	systemCertPEM, err := readAdbFile(options, systemCertPath)
	if err != nil || !isMatchingCert(systemCertPEM, expectedFingerprint) {
		return false
	}

	cmdLs := adbCommand(options, "shell", "ls", "/apex/com.android.conscrypt")
	lsOut, err := cmdLs.CombinedOutput()
	if err == nil && strings.Contains(string(lsOut), "cacerts") {
		apexCertPath := fmt.Sprintf("/apex/com.android.conscrypt/cacerts/%s.0", certHash)
		apexCertPEM, err := readAdbFile(options, apexCertPath)
		if err != nil || !isMatchingCert(apexCertPEM, expectedFingerprint) {
			return false
		}
	}

	return true
}

func ensureAndroidAPK(cfg *config.Config) (string, error) {
	apkPath := filepath.Join(cfg.ConfigDir, "httptoolkit-latest.apk")
	if _, err := os.Stat(apkPath); err == nil {
		return apkPath, nil
	}

	url := "https://github.com/httptoolkit/httptoolkit-android/releases/latest/download/httptoolkit.apk"
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Android APK: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download Android APK: HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(apkPath)
	if err != nil {
		return "", fmt.Errorf("failed to create APK cache file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write APK to disk: %w", err)
	}

	return apkPath, nil
}

func getRootCommand(deviceId string, options map[string]any) ([]string, error) {
	cmd := adbCommand(options, "shell", "su", "-c", "id")
	out, err := cmd.CombinedOutput()
	if err == nil && strings.Contains(string(out), "uid=0(root)") {
		return []string{"su", "-c"}, nil
	}

	cmd = adbCommand(options, "shell", "su", "root", "id")
	out, err = cmd.CombinedOutput()
	if err == nil && strings.Contains(string(out), "uid=0(root)") {
		return []string{"su", "root"}, nil
	}

	cmd = adbCommand(options, "root")
	_, _ = cmd.CombinedOutput()

	for i := 0; i < 10; i++ {
		time.Sleep(200 * time.Millisecond)
		cmd = adbCommand(options, "shell", "id")
		out, err = cmd.CombinedOutput()
		if err == nil && strings.Contains(string(out), "uid=0(root)") {
			return []string{}, nil
		}
	}

	return nil, fmt.Errorf("root not available")
}

func buildCertInjectionScript(certPath, subjectHash string) string {
	return fmt.Sprintf(`
set -e
echo "Injecting certificate:"
mkdir -p /data/local/tmp/htk-ca-copy
chmod 700 /data/local/tmp/htk-ca-copy
rm -rf /data/local/tmp/htk-ca-copy/*

if [ -d "/apex/com.android.conscrypt/cacerts" ]; then
    cp /apex/com.android.conscrypt/cacerts/* /data/local/tmp/htk-ca-copy/
else
    cp /system/etc/security/cacerts/* /data/local/tmp/htk-ca-copy/
fi

mount -t tmpfs tmpfs /system/etc/security/cacerts
mv /data/local/tmp/htk-ca-copy/* /system/etc/security/cacerts/
mv %s /system/etc/security/cacerts/%s.0

chown root:root /system/etc/security/cacerts/*
chmod 644 /system/etc/security/cacerts/*
chcon u:object_r:system_file:s0 /system/etc/security/cacerts/
chcon u:object_r:system_file:s0 /system/etc/security/cacerts/*

echo "System cacerts setup completed"

if [ -d "/apex/com.android.conscrypt/cacerts" ]; then
    echo "Injecting certificates into APEX cacerts"
    mount --bind /system/etc/security/cacerts /apex/com.android.conscrypt/cacerts

    ZYGOTE_PID=$(pidof zygote || true)
    ZYGOTE64_PID=$(pidof zygote64 || true)
    Z_PIDS="$ZYGOTE_PID $ZYGOTE64_PID"

    for Z_PID in $Z_PIDS; do
        if [ -n "$Z_PID" ]; then
            nsenter --mount=/proc/$Z_PID/ns/mnt -- mount --bind /system/etc/security/cacerts /apex/com.android.conscrypt/cacerts || true
        fi
    done

    echo "Zygote APEX certificates remounted"

    APP_PIDS=$(ps -o PID -P || ps -A -o PID | grep -v PID || true)
    for PID in $APP_PIDS; do
        if [ -d "/proc/$PID/ns" ]; then
            nsenter --mount=/proc/$PID/ns/mnt -- mount --bind /system/etc/security/cacerts /apex/com.android.conscrypt/cacerts &
        fi
    done
    wait
    echo "APEX certificates remounted for apps"
fi

rm -rf /data/local/tmp/htk-ca-copy
echo "System cert successfully injected"
`, certPath, subjectHash)
}

func buildChromeFlagsScript(spkiFingerprint string) string {
	chromeFlagsLocations := []string{
		"/data/local/chrome-command-line",
		"/data/local/tmp/chrome-command-line",
		"/data/local/android-webview-command-line",
		"/data/local/tmp/android-webview-command-line",
		"/data/local/webview-command-line",
		"/data/local/tmp/webview-command-line",
		"/data/local/content-shell-command-line",
		"/data/local/tmp/content-shell-command-line",
	}

	var sb strings.Builder
	sb.WriteString("set -e\n")
	for _, path := range chromeFlagsLocations {
		sb.WriteString(fmt.Sprintf("echo \"chrome --ignore-certificate-errors-spki-list=%s\" > \"%s\"\n", spkiFingerprint, path))
		sb.WriteString(fmt.Sprintf("chmod 744 \"%s\"\n", path))
		sb.WriteString(fmt.Sprintf("chcon \"u:object_r:shell_data_file:s0\" \"%s\"\n", path))
	}
	sb.WriteString("echo \"Chrome flags script completed\"\n")
	return sb.String()
}

func startPersistentReverseTunnel(options map[string]any, localPort, remotePort int, deviceId string) {
	key := fmt.Sprintf("%s:%d", deviceId, localPort)
	adbTunnelsMutex.Lock()
	if _, exists := adbTunnels[key]; exists {
		adbTunnelsMutex.Unlock()
		return
	}
	stopChan := make(chan struct{})
	adbTunnels[key] = stopChan
	adbTunnelsMutex.Unlock()

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		cmd := adbCommand(options, "reverse", fmt.Sprintf("tcp:%d", localPort), fmt.Sprintf("tcp:%d", remotePort))
		_ = cmd.Run()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				cmd := adbCommand(options, "reverse", fmt.Sprintf("tcp:%d", localPort), fmt.Sprintf("tcp:%d", remotePort))
				_ = cmd.Run()
			}
		}
	}()
}

func stopPersistentReverseTunnel(options map[string]any, localPort int, deviceId string) {
	key := fmt.Sprintf("%s:%d", deviceId, localPort)
	adbTunnelsMutex.Lock()
	if stopChan, exists := adbTunnels[key]; exists {
		close(stopChan)
		delete(adbTunnels, key)
	}
	adbTunnelsMutex.Unlock()

	cmd := adbCommand(options, "reverse", "--remove", fmt.Sprintf("tcp:%d", localPort))
	_ = cmd.Run()
}

func getReachableIPv4s() []string {
	var ips []string
	ips = append(ips, "10.0.2.2", "10.0.3.2")
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					ips = append(ips, ip4.String())
				}
			}
		}
	}
	return ips
}

// androidCertFingerprint returns the SHA-256 fingerprint of the DER certificate,
// which is what the Android app expects. Falls back to the SPKI fingerprint.
func androidCertFingerprint(cfg *config.Config, fallback string) string {
	certPath := filepath.Join(cfg.ConfigDir, "ca.pem")
	data, err := os.ReadFile(certPath)
	if err == nil {
		block, _ := pem.Decode(data)
		if block != nil {
			h := sha256.Sum256(block.Bytes)
			parts := make([]string, 32)
			for i, b := range h {
				parts[i] = fmt.Sprintf("%02X", b)
			}
			return strings.Join(parts, ":")
		}
	}
	return fallback
}

// buildAndroidQRData builds the base64url-encoded JSON payload that the Android
// companion app scans from the UI's QR code.
func buildAndroidQRData(proxyPort int, certFingerprint string) string {
	params := map[string]any{
		"addresses":       getReachableIPv4s(),
		"port":            proxyPort,
		"localTunnelPort": proxyPort,
		"enableSocks":     false,
		"certFingerprint": certFingerprint,
	}
	b, err := json.Marshal(params)
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func isAppInstalled(options map[string]any, packageName string) bool {
	cmd := adbCommand(options, "shell", "pm", "list", "packages", packageName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "package:"+packageName)
}

func quoteSlice(s []string) []string {
	res := make([]string, len(s))
	for i, v := range s {
		res[i] = fmt.Sprintf(`"%s"`, v)
	}
	return res
}

type androidAdb struct{ base *stubInterceptor }

func (a *androidAdb) SubMetadata(subID string) (any, error) {
	if subID == "devices" || subID == "" {
		out, err := execOutputTimeout(probeTimeout, "adb", "devices", "-l")
		if err != nil {
			return nil, err
		}
		deviceIds, devices := parseAdbDevicesMetadata(string(out))
		return map[string]any{
			"deviceIds": deviceIds,
			"devices":   devices,
		}, nil
	}
	if subID == "qr" {
		a.base.mu.Lock()
		var proxyPort int
		for p, active := range a.base.active {
			if active {
				proxyPort = p
				break
			}
		}
		a.base.mu.Unlock()
		if proxyPort == 0 {
			return nil, fmt.Errorf("android-adb is not active")
		}
		return buildAndroidQRData(proxyPort, androidCertFingerprint(a.base.cfg, a.base.spki)), nil
	}
	return a.Metadata("detailed")
}

func parseAdbDevicesMetadata(raw string) ([]string, map[string]map[string]string) {
	deviceIds := make([]string, 0)
	devices := make(map[string]map[string]string)

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "List of devices") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		id := parts[0]
		state := parts[1]
		if state == "offline" || state == "unauthorized" {
			continue
		}

		deviceIds = append(deviceIds, id)

		deviceMap := map[string]string{
			"id":    id,
			"state": state,
		}

		for _, part := range parts[2:] {
			kv := strings.SplitN(part, ":", 2)
			if len(kv) == 2 {
				key := kv[0]
				val := kv[1]
				if key == "model" {
					deviceMap["ro.product.model"] = strings.ReplaceAll(val, "_", " ")
				}
				deviceMap[key] = val
			}
		}

		if _, exists := deviceMap["ro.product.model"]; !exists {
			deviceMap["ro.product.model"] = id
		}

		devices[id] = deviceMap
	}

	return deviceIds, devices
}

func parseAdbDevices(raw string) []map[string]string {
	var devices []map[string]string
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "List of devices") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		entry := map[string]string{"id": parts[0], "state": parts[1]}
		if len(parts) > 2 {
			entry["details"] = strings.Join(parts[2:], " ")
		}
		devices = append(devices, entry)
	}
	return devices
}

func (a *androidAdb) IsActivable() (bool, error)   { return commandExists("adb"), nil }
func (a *androidAdb) IsActive(p int) (bool, error) { return a.base.IsActive(p) }
func adbArgs(options map[string]any, args ...string) []string {
	out := []string{"adb"}
	if id, ok := options["deviceId"].(string); ok && id != "" {
		out = append(out, "-s", id)
	} else if id, ok := options["device"].(string); ok && id != "" {
		out = append(out, "-s", id)
	}
	return append(out, args...)
}

func adbCommand(options map[string]any, args ...string) *exec.Cmd {
	cmdArgs := adbArgs(options, args...)
	return execCommand(cmdArgs[0], cmdArgs[1:]...)
}

func (a *androidAdb) Activate(proxyPort int, options map[string]any) (map[string]any, error) {
	deviceId, _ := options["deviceId"].(string)
	if deviceId == "" {
		deviceId, _ = options["device"].(string)
	}

	if deviceId == "" {
		out, err := execCombinedOutputTimeout(probeTimeout, "adb", "devices")
		if err == nil {
			devices := parseAdbDevices(string(out))
			if len(devices) > 0 {
				deviceId = devices[0]["id"]
				options["deviceId"] = deviceId
			}
		}
	}

	if deviceId == "" {
		return nil, fmt.Errorf("no Android devices found via ADB")
	}

	certPEM, _ := options["certificateContent"].(string)
	if certPEM == "" {
		certPath := filepath.Join(a.base.cfg.ConfigDir, "ca.pem")
		if data, err := os.ReadFile(certPath); err == nil {
			certPEM = string(data)
		}
	}

	if certPEM != "" {
		subjectHash, err := getCertificateSubjectHash(certPEM)
		if err == nil {
			var expectedFingerprint string
			block, _ := pem.Decode([]byte(certPEM))
			if block != nil {
				if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
					h := sha1.Sum(cert.Raw)
					expectedFingerprint = fmt.Sprintf("%x", h)
				}
			}

			if !hasCertInstalled(options, subjectHash, expectedFingerprint) {
				if rootPrefix, err := getRootCommand(deviceId, options); err == nil {
					tmpCertPath := fmt.Sprintf("/data/local/tmp/%s.0", subjectHash)
					hostTmp, err := os.CreateTemp("", "htk-ca-*.pem")
					if err == nil {
						defer os.Remove(hostTmp.Name())
						_, _ = hostTmp.WriteString(certPEM)
						hostTmp.Close()

						cmdPush := adbCommand(options, "push", hostTmp.Name(), tmpCertPath)
						if err := cmdPush.Run(); err == nil {
							script := buildCertInjectionScript(tmpCertPath, subjectHash)
							hostScript, err := os.CreateTemp("", "htk-inject-*.sh")
							if err == nil {
								defer os.Remove(hostScript.Name())
								_, _ = hostScript.WriteString(script)
								hostScript.Close()

								cmdPushScript := adbCommand(options, "push", hostScript.Name(), "/data/local/tmp/htk-inject.sh")
								if err := cmdPushScript.Run(); err == nil {
									rootArgs := append([]string{"shell"}, append(rootPrefix, "sh", "/data/local/tmp/htk-inject.sh")...)
									cmdRunScript := adbCommand(options, rootArgs...)
									_, _ = cmdRunScript.CombinedOutput()

									cmdCleanup := adbCommand(options, "shell", "rm", "-f", "/data/local/tmp/htk-inject.sh")
									_ = cmdCleanup.Run()
								}
							}
						}
					}

					if spki := a.base.spki; spki != "" {
						script := buildChromeFlagsScript(spki)
						hostScript, err := os.CreateTemp("", "htk-chrome-flags-*.sh")
						if err == nil {
							defer os.Remove(hostScript.Name())
							_, _ = hostScript.WriteString(script)
							hostScript.Close()

							cmdPushFlags := adbCommand(options, "push", hostScript.Name(), "/data/local/tmp/htk-set-chrome-flags.sh")
							if err := cmdPushFlags.Run(); err == nil {
								rootArgs := append([]string{"shell"}, append(rootPrefix, "sh", "/data/local/tmp/htk-set-chrome-flags.sh")...)
								cmdRunFlags := adbCommand(options, rootArgs...)
								_, _ = cmdRunFlags.CombinedOutput()

								cmdCleanupFlags := adbCommand(options, "shell", "rm", "-f", "/data/local/tmp/htk-set-chrome-flags.sh")
								_ = cmdCleanupFlags.Run()
							}
						}
						rootStopArgs := append([]string{"shell"}, append(rootPrefix, "am", "force-stop", "com.android.chrome")...)
						cmdStopChrome := adbCommand(options, rootStopArgs...)
						_, _ = cmdStopChrome.CombinedOutput()
					}
				}
			}
		}
	}

	appID := "tech.httptoolkit.android.v1"
	if !isAppInstalled(options, appID) {
		apkPath, err := ensureAndroidAPK(a.base.cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch HTTP Toolkit Android APK: %w", err)
		}
		cmd := adbCommand(options, "install", apkPath)
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to install HTTP Toolkit Android VPN app: %w", err)
		}
	}

	_ = adbCommand(options, "shell", "input", "keyevent", "KEYCODE_WAKEUP").Run()
	time.Sleep(50 * time.Millisecond)
	_ = adbCommand(options, "shell", "am", "start", "--activity-single-top", "tech.httptoolkit.android.v1/tech.httptoolkit.android.MainActivity").Run()

	startPersistentReverseTunnel(options, proxyPort, proxyPort, deviceId)

	enableSocks := false
	if s, ok := options["enableSocks"].(bool); ok {
		enableSocks = s
	}

	// Compute SHA-256 fingerprint of the DER certificate in colon-separated uppercase hex format,
	// which is what the Android HttpToolkit app expects (not SPKI base64).
	certFingerprint := a.base.spki // fallback to SPKI if we can't compute the real fingerprint
	if certPEM != "" {
		block, _ := pem.Decode([]byte(certPEM))
		if block != nil {
			h := sha256.Sum256(block.Bytes)
			parts := make([]string, 32)
			for i, b := range h {
				parts[i] = fmt.Sprintf("%02X", b)
			}
			certFingerprint = strings.Join(parts, ":")
		}
	}

	// Filter addresses to only include LAN IPs reachable from a real device.
	// Exclude emulator-only addresses (10.0.2.2, 10.0.3.2) for physical devices.
	allAddresses := getReachableIPv4s()

	setupParams := fmt.Sprintf(`{"addresses":[%s],"port":%d,"localTunnelPort":%d,"enableSocks":%t,"certFingerprint":"%s"}`,
		strings.Join(quoteSlice(allAddresses), ","),
		proxyPort,
		proxyPort,
		enableSocks,
		certFingerprint,
	)

	intentData := base64.RawURLEncoding.EncodeToString([]byte(setupParams))
	intentURL := fmt.Sprintf("https://android.httptoolkit.tech/connect/?data=%s", intentData)
	_ = adbCommand(options, "shell", "am", "start", "--activity-single-top",
		"-a", "tech.httptoolkit.android.ACTIVATE",
		"-d", intentURL,
	).Run()

	return map[string]any{
		"ok":       true,
		"deviceId": deviceId,
	}, nil
}

func (a *androidAdb) Deactivate(proxyPort int, options map[string]any) error {
	deviceId, _ := options["deviceId"].(string)
	if deviceId == "" {
		deviceId, _ = options["device"].(string)
	}

	stopPersistentReverseTunnel(options, proxyPort, deviceId)

	_ = adbCommand(options, "shell", "am", "start", "--activity-single-top", "tech.httptoolkit.android.v1/tech.httptoolkit.android.MainActivity").Run()
	time.Sleep(100 * time.Millisecond)

	_ = adbCommand(options, "shell", "am", "start", "--activity-single-top",
		"-a", "tech.httptoolkit.android.DEACTIVATE",
	).Run()

	return nil
}

// getDeviceCertificateStatus checks if the CA certificate is installed in the
// Android system trust store for the given device. Mirrors the Node backend's
// AndroidAdbInterceptor.getDeviceCertificateStatus.
func (a *androidAdb) getDeviceCertificateStatus(deviceId, certPEM string) (map[string]any, error) {
	options := map[string]any{"deviceId": deviceId}

	rootPrefix, err := getRootCommand(deviceId, options)
	if err != nil {
		return map[string]any{
			"deviceId":      deviceId,
			"installed":     false,
			"rootAvailable": false,
			"message":       "Root access is required to install the CA in the Android system trust store.",
		}, nil
	}
	_ = rootPrefix

	subjectHash, err := getCertificateSubjectHash(certPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to compute certificate subject hash: %w", err)
	}

	var expectedFingerprint string
	if block, _ := pem.Decode([]byte(certPEM)); block != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			h := sha1.Sum(cert.Raw)
			expectedFingerprint = fmt.Sprintf("%x", h)
		}
	}

	installed := hasCertInstalled(options, subjectHash, expectedFingerprint)

	message := "HttpToolkit CA is not installed. Install it before intercepting HTTPS."
	if installed {
		message = "HttpToolkit CA is installed in the system trust store."
	}

	return map[string]any{
		"deviceId":      deviceId,
		"installed":     installed,
		"rootAvailable": true,
		"message":       message,
	}, nil
}

// installDeviceCertificate installs the CA certificate into the Android system
// trust store for the given device (requires root). Mirrors the Node backend's
// AndroidAdbInterceptor.installDeviceCertificate.
func (a *androidAdb) installDeviceCertificate(deviceId, certPEM, certPath string) (map[string]any, error) {
	options := map[string]any{"deviceId": deviceId}

	rootPrefix, err := getRootCommand(deviceId, options)
	if err != nil {
		return map[string]any{
			"deviceId":      deviceId,
			"installed":     false,
			"rootAvailable": false,
			"message":       "Root access is required to install the CA in the Android system trust store.",
		}, nil
	}

	subjectHash, err := getCertificateSubjectHash(certPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to compute certificate subject hash: %w", err)
	}

	tmpCertPath := fmt.Sprintf("/data/local/tmp/%s.0", subjectHash)
	hostTmp, err := os.CreateTemp("", "htk-ca-*.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp cert file: %w", err)
	}
	defer os.Remove(hostTmp.Name())
	_, _ = hostTmp.WriteString(certPEM)
	hostTmp.Close()

	cmdPush := adbCommand(options, "push", hostTmp.Name(), tmpCertPath)
	if err := cmdPush.Run(); err != nil {
		return nil, fmt.Errorf("failed to push cert to device: %w", err)
	}

	script := buildCertInjectionScript(tmpCertPath, subjectHash)
	hostScript, err := os.CreateTemp("", "htk-inject-*.sh")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp script: %w", err)
	}
	defer os.Remove(hostScript.Name())
	_, _ = hostScript.WriteString(script)
	hostScript.Close()

	cmdPushScript := adbCommand(options, "push", hostScript.Name(), "/data/local/tmp/htk-inject.sh")
	if err := cmdPushScript.Run(); err != nil {
		return nil, fmt.Errorf("failed to push injection script: %w", err)
	}

	rootArgs := append([]string{"shell"}, append(rootPrefix, "sh", "/data/local/tmp/htk-inject.sh")...)
	cmdRunScript := adbCommand(options, rootArgs...)
	_, _ = cmdRunScript.CombinedOutput()

	cmdCleanup := adbCommand(options, "shell", "rm", "-f", "/data/local/tmp/htk-inject.sh")
	_ = cmdCleanup.Run()

	// Verify installation
	status, _ := a.getDeviceCertificateStatus(deviceId, certPEM)
	return status, nil
}

func (a *androidAdb) Metadata(string) (any, error) {
	out, err := execOutputTimeout(probeTimeout, "adb", "devices", "-l")
	if err != nil {
		return map[string]any{
			"type":      "adb",
			"deviceIds": []string{},
			"devices":   map[string]any{},
		}, nil
	}
	deviceIds, devices := parseAdbDevicesMetadata(string(out))
	return map[string]any{
		"type":      "adb",
		"deviceIds": deviceIds,
		"devices":   devices,
	}, nil
}

type androidFrida struct {
	base        *stubInterceptor
	mu          sync.Mutex
	cmds        map[int]*exec.Cmd
	hostIDs     map[int]string
	scriptPaths map[int]string
	stdins      map[int]io.WriteCloser
}

func (a *androidFrida) init() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.cmds == nil {
		a.cmds = make(map[int]*exec.Cmd)
	}
	if a.hostIDs == nil {
		a.hostIDs = make(map[int]string)
	}
	if a.scriptPaths == nil {
		a.scriptPaths = make(map[int]string)
	}
	if a.stdins == nil {
		a.stdins = make(map[int]io.WriteCloser)
	}
}

func parseAdbDevicesForFrida(base *stubInterceptor) []map[string]string {
	out, err := execOutputTimeout(probeTimeout, "adb", "devices", "-l")
	if err != nil {
		return nil
	}
	rawDevices := parseAdbDevices(string(out))
	var devices []map[string]string

	for _, dev := range rawDevices {
		id := dev["id"]
		state := dev["state"]
		if state != "device" {
			continue
		}

		name := dev["details"]
		if name == "" {
			name = "Android Device"
		} else {
			fields := strings.Fields(name)
			for _, f := range fields {
				if strings.HasPrefix(f, "model:") {
					name = strings.ReplaceAll(strings.TrimPrefix(f, "model:"), "_", " ")
					break
				}
			}
		}

		fridaState := "setup-required"
		if err := execRunTimeout(probeTimeout, "frida-ps", "-D", id); err == nil {
			fridaState = "available"
		} else {
			if err := execRunTimeout(probeTimeout, "adb", "-s", id, "shell", "ls", "/data/local/tmp/frida-server"); err == nil {
				fridaState = "launch-required"
			}
		}

		devices = append(devices, map[string]string{
			"id":    id,
			"name":  name,
			"state": fridaState,
		})
	}
	return devices
}

func parseFridaPsTabular(raw string) []map[string]any {
	var targets []map[string]any
	lines := strings.Split(raw, "\n")
	if len(lines) < 2 {
		return targets
	}
	header := lines[0]
	nameIdx := strings.Index(header, "Name")
	idIdx := strings.Index(header, "Identifier")

	if nameIdx == -1 || idIdx == -1 {
		for _, line := range lines[1:] {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "---") {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				id := fields[len(fields)-1]
				name := fields[0]
				if len(fields) > 2 {
					name = strings.Join(fields[1:len(fields)-1], " ")
				}
				targets = append(targets, map[string]any{
					"id":   id,
					"name": name,
				})
			}
		}
		return targets
	}

	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "---") {
			continue
		}
		var name, id string
		if len(line) > idIdx {
			id = strings.TrimSpace(line[idIdx:])
			name = strings.TrimSpace(line[nameIdx:idIdx])
		} else if len(line) > nameIdx {
			name = strings.TrimSpace(line[nameIdx:])
		}
		if id != "" {
			targets = append(targets, map[string]any{
				"id":   id,
				"name": name,
			})
		}
	}
	return targets
}

func (a *androidFrida) SubMetadata(hostID string) (any, error) {
	if hostID == "" {
		return map[string]any{"targets": []any{}, "hostId": hostID}, nil
	}

	out, err := execOutputTimeout(probeTimeout, "frida-ps", "-D", hostID, "-a", "-i", "--json")
	if err == nil {
		var apps []map[string]any
		if err := json.Unmarshal(out, &apps); err == nil {
			targets := make([]map[string]any, 0, len(apps))
			for _, app := range apps {
				id, _ := app["identifier"].(string)
				name, _ := app["name"].(string)
				if id == "" {
					id, _ = app["id"].(string)
				}
				if id != "" {
					targets = append(targets, map[string]any{
						"id":   id,
						"name": name,
					})
				}
			}
			return map[string]any{"targets": targets, "hostId": hostID}, nil
		}
	}

	outRaw, err := execOutputTimeout(probeTimeout, "frida-ps", "-D", hostID, "-a", "-i")
	if err != nil {
		return map[string]any{"targets": []any{}, "hostId": hostID}, nil
	}

	targets := parseFridaPsTabular(string(outRaw))
	return map[string]any{"targets": targets, "hostId": hostID}, nil
}

func (a *androidFrida) IsActivable() (bool, error) {
	return commandExists("adb") && commandExists("frida"), nil
}
func (a *androidFrida) IsActive(p int) (bool, error) { return a.base.IsActive(p) }

func (a *androidFrida) Activate(proxyPort int, options map[string]any) (map[string]any, error) {
	a.init()
	action, _ := options["action"].(string)
	hostID, _ := options["hostId"].(string)

	if action == "setup" {
		return a.setupFridaServer(hostID)
	}

	if action == "launch" {
		if hostID == "" {
			return nil, fmt.Errorf("hostId required to launch frida-server")
		}
		launchCmd := execCommand("adb", "-s", hostID, "shell", "su", "-c", "/data/local/tmp/frida-server &")
		_ = launchCmd.Start()

		for i := 0; i < 10; i++ {
			time.Sleep(500 * time.Millisecond)
			if err := execRunTimeout(probeTimeout, "frida-ps", "-D", hostID); err == nil {
				return map[string]any{"success": true}, nil
			}
		}
		return map[string]any{"success": true}, nil
	}

	target, _ := options["targetId"].(string)
	if target == "" {
		target, _ = options["target"].(string)
	}
	if target == "" {
		target = "Gadget"
	}

	certPEM, _ := options["certificateContent"].(string)
	if certPEM == "" {
		certPath := filepath.Join(a.base.cfg.ConfigDir, "ca.pem")
		if data, err := os.ReadFile(certPath); err == nil {
			certPEM = string(data)
		}
	}

	enableSocks := false
	if s, ok := options["enableSocks"].(bool); ok {
		enableSocks = s
	}

	ignoredPorts := knownAppProblematicPorts[target]

	proxyHost := "127.0.0.1"
	if hostID != "" {
		a.mu.Lock()
		a.hostIDs[proxyPort] = hostID
		a.mu.Unlock()
		startPersistentReverseTunnel(options, proxyPort, proxyPort, hostID)
	} else {
		proxyHost = getPrimaryLANIP()
	}

	scriptContent, err := buildAndroidFridaScript(a.base.cfg, certPEM, proxyHost, proxyPort, ignoredPorts, enableSocks)
	if err != nil {
		if hostID != "" {
			stopPersistentReverseTunnel(options, proxyPort, hostID)
			a.mu.Lock()
			delete(a.hostIDs, proxyPort)
			a.mu.Unlock()
		}
		return nil, fmt.Errorf("failed to build Frida script: %w", err)
	}

	tmpScript, err := os.CreateTemp("", "htk-android-frida-*.js")
	if err != nil {
		if hostID != "" {
			stopPersistentReverseTunnel(options, proxyPort, hostID)
			a.mu.Lock()
			delete(a.hostIDs, proxyPort)
			a.mu.Unlock()
		}
		return nil, fmt.Errorf("failed to create temp script file: %w", err)
	}
	defer tmpScript.Close()

	if _, err := tmpScript.WriteString(scriptContent); err != nil {
		os.Remove(tmpScript.Name())
		if hostID != "" {
			stopPersistentReverseTunnel(options, proxyPort, hostID)
			a.mu.Lock()
			delete(a.hostIDs, proxyPort)
			a.mu.Unlock()
		}
		return nil, fmt.Errorf("failed to write temp script: %w", err)
	}

	a.mu.Lock()
	a.scriptPaths[proxyPort] = tmpScript.Name()
	a.mu.Unlock()

	var cmd *exec.Cmd
	var stderrBuf bytes.Buffer
	if hostID != "" {
		// -D: connect to device by ADB serial; -f: spawn the app; --no-pause: auto-resume after script loads
		cmd = execCommand("frida", "-D", hostID, "-f", target, "--no-pause", "-l", tmpScript.Name())
	} else {
		cmd = execCommand("frida", "-U", "-f", target, "--no-pause", "-l", tmpScript.Name())
	}

	pr, pw := io.Pipe()
	cmd.Stdin = pr
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		pr.Close()
		pw.Close()
		os.Remove(tmpScript.Name())
		if hostID != "" {
			stopPersistentReverseTunnel(options, proxyPort, hostID)
		}
		a.mu.Lock()
		delete(a.scriptPaths, proxyPort)
		delete(a.hostIDs, proxyPort)
		a.mu.Unlock()
		return nil, fmt.Errorf("failed to start frida: %w", err)
	}

	// Wait briefly to detect immediate exits (e.g. package not found, frida-server not running)
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case exitErr := <-done:
		// Frida exited almost immediately — this is an error
		pr.Close()
		pw.Close()
		os.Remove(tmpScript.Name())
		if hostID != "" {
			stopPersistentReverseTunnel(options, proxyPort, hostID)
		}
		a.mu.Lock()
		delete(a.scriptPaths, proxyPort)
		delete(a.hostIDs, proxyPort)
		a.mu.Unlock()
		errMsg := strings.TrimSpace(stderrBuf.String())
		if errMsg == "" {
			errMsg = fmt.Sprintf("frida exited immediately (target app '%s' may not be installed or frida-server may not be running on device)", target)
		}
		if exitErr != nil {
			return nil, fmt.Errorf("failed to intercept app %s: %s", target, errMsg)
		}
		return nil, fmt.Errorf("frida exited unexpectedly for %s: %s", target, errMsg)
	case <-time.After(3 * time.Second):
		// Frida is still running after 3 seconds — assume success
	}

	a.mu.Lock()
	a.cmds[proxyPort] = cmd
	a.stdins[proxyPort] = pw
	a.mu.Unlock()

	return map[string]any{"started": true, "target": target}, nil
}

// setupFridaServer pushes a host-side frida-server binary to the device if it can
// be found in PATH. If the binary is already on the device it returns success.
func (a *androidFrida) setupFridaServer(hostID string) (map[string]any, error) {
	if hostID == "" {
		return nil, fmt.Errorf("hostId required to set up frida-server")
	}

	// If frida-server is already present, just mark it as ready to launch.
	lsCmd := adbCommand(map[string]any{"deviceId": hostID}, "shell", "ls", "/data/local/tmp/frida-server")
	if err := lsCmd.Run(); err == nil {
		return map[string]any{"success": true, "message": "frida-server already present"}, nil
	}

	serverPath, err := exec.LookPath("frida-server")
	if err != nil {
		return nil, fmt.Errorf("frida-server not found in PATH; download it from https://frida.re and add it to PATH")
	}

	pushCmd := adbCommand(map[string]any{"deviceId": hostID}, "push", serverPath, "/data/local/tmp/frida-server")
	if out, err := pushCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to push frida-server: %w\n%s", err, string(out))
	}

	chmodCmd := adbCommand(map[string]any{"deviceId": hostID}, "shell", "su", "-c", "chmod 755 /data/local/tmp/frida-server")
	_ = chmodCmd.Run()

	return map[string]any{"success": true, "message": "frida-server pushed"}, nil
}

func (a *androidFrida) Deactivate(proxyPort int, options map[string]any) error {
	a.init()
	a.mu.Lock()
	hostID := a.hostIDs[proxyPort]
	delete(a.hostIDs, proxyPort)

	scriptPath := a.scriptPaths[proxyPort]
	delete(a.scriptPaths, proxyPort)

	pw, okPw := a.stdins[proxyPort]
	if okPw && pw != nil {
		_ = pw.Close()
	}
	delete(a.stdins, proxyPort)

	cmd, ok := a.cmds[proxyPort]
	if ok && cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
		// NOTE: Do NOT call cmd.Process.Wait() here — the goroutine in Activate already
		// calls cmd.Wait() to detect early exits. Double-Wait causes a panic.
	}
	delete(a.cmds, proxyPort)
	a.mu.Unlock()

	if scriptPath != "" {
		os.Remove(scriptPath)
	}

	if hostID != "" {
		stopPersistentReverseTunnel(options, proxyPort, hostID)
	}

	return nil
}

func (a *androidFrida) Metadata(string) (any, error) {
	devices := parseAdbDevicesForFrida(a.base)
	hosts := make(map[string]any)
	for _, dev := range devices {
		id := dev["id"]
		hosts[id] = map[string]any{
			"id":    id,
			"name":  dev["name"],
			"state": dev["state"],
		}
	}
	return map[string]any{
		"type":  "frida-android",
		"hosts": hosts,
	}, nil
}

type iosFrida struct {
	base        *stubInterceptor
	mu          sync.Mutex
	cmds        map[int]*exec.Cmd
	scriptPaths map[int]string
	stdins      map[int]io.WriteCloser
}

func (i *iosFrida) init() {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.cmds == nil {
		i.cmds = make(map[int]*exec.Cmd)
	}
	if i.scriptPaths == nil {
		i.scriptPaths = make(map[int]string)
	}
	if i.stdins == nil {
		i.stdins = make(map[int]io.WriteCloser)
	}
}

func (i *iosFrida) SubMetadata(hostID string) (any, error) {
	if hostID == "" {
		return map[string]any{"targets": []any{}, "hostId": hostID}, nil
	}

	out, err := execOutputTimeout(probeTimeout, "frida-ps", "-U", "-a", "-i", "--json")
	if err == nil {
		var apps []map[string]any
		if err := json.Unmarshal(out, &apps); err == nil {
			targets := make([]map[string]any, 0, len(apps))
			for _, app := range apps {
				id, _ := app["identifier"].(string)
				name, _ := app["name"].(string)
				if id == "" {
					id, _ = app["id"].(string)
				}
				if id != "" {
					targets = append(targets, map[string]any{
						"id":   id,
						"name": name,
					})
				}
			}
			return map[string]any{"targets": targets, "hostId": hostID}, nil
		}
	}

	outRaw, err := execOutputTimeout(probeTimeout, "frida-ps", "-U", "-a", "-i")
	if err != nil {
		return map[string]any{"targets": []any{}, "hostId": hostID}, nil
	}

	targets := parseFridaPsTabular(string(outRaw))
	return map[string]any{"targets": targets, "hostId": hostID}, nil
}

func (i *iosFrida) IsActivable() (bool, error)   { return commandExists("frida"), nil }
func (i *iosFrida) IsActive(p int) (bool, error) { return i.base.IsActive(p) }

func (i *iosFrida) Activate(proxyPort int, options map[string]any) (map[string]any, error) {
	i.init()
	action, _ := options["action"].(string)
	if action == "setup" || action == "launch" {
		return map[string]any{"success": true}, nil
	}

	target, _ := options["targetId"].(string)
	if target == "" {
		target, _ = options["target"].(string)
	}
	if target == "" {
		return nil, fmt.Errorf("targetId required")
	}

	certPEM, _ := options["certificateContent"].(string)
	if certPEM == "" {
		certPath := filepath.Join(i.base.cfg.ConfigDir, "ca.pem")
		if data, err := os.ReadFile(certPath); err == nil {
			certPEM = string(data)
		}
	}

	enableSocks := false
	if s, ok := options["enableSocks"].(bool); ok {
		enableSocks = s
	}

	ignoredPorts := knownAppProblematicPorts[target]

	proxyHost, _ := options["proxyHost"].(string)
	if proxyHost == "" {
		proxyHost, _ = options["host"].(string)
	}
	if proxyHost == "" {
		proxyHost = getPrimaryLANIP()
	}

	scriptContent, err := buildIosFridaScript(i.base.cfg, certPEM, proxyHost, proxyPort, ignoredPorts, enableSocks)
	if err != nil {
		return nil, fmt.Errorf("failed to build Frida script: %w", err)
	}

	tmpScript, err := os.CreateTemp("", "htk-ios-frida-*.js")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp script file: %w", err)
	}
	defer tmpScript.Close()

	if _, err := tmpScript.WriteString(scriptContent); err != nil {
		os.Remove(tmpScript.Name())
		return nil, fmt.Errorf("failed to write temp script: %w", err)
	}

	i.mu.Lock()
	i.scriptPaths[proxyPort] = tmpScript.Name()
	i.mu.Unlock()

	var stderrBufIOS bytes.Buffer
	cmd := execCommand("frida", "-U", "-f", target, "--no-pause", "-l", tmpScript.Name())
	pr, pw := io.Pipe()
	cmd.Stdin = pr
	cmd.Stderr = &stderrBufIOS

	if err := cmd.Start(); err != nil {
		pr.Close()
		pw.Close()
		os.Remove(tmpScript.Name())
		i.mu.Lock()
		delete(i.scriptPaths, proxyPort)
		i.mu.Unlock()
		return nil, fmt.Errorf("failed to start frida: %w", err)
	}

	// Wait briefly to detect immediate exits (e.g. package not found, frida-server not running)
	doneIOS := make(chan error, 1)
	go func() { doneIOS <- cmd.Wait() }()

	select {
	case exitErr := <-doneIOS:
		pr.Close()
		pw.Close()
		os.Remove(tmpScript.Name())
		i.mu.Lock()
		delete(i.scriptPaths, proxyPort)
		i.mu.Unlock()
		errMsg := strings.TrimSpace(stderrBufIOS.String())
		if errMsg == "" {
			errMsg = fmt.Sprintf("frida exited immediately (target app '%s' may not be installed or iOS device may not have frida-server running)", target)
		}
		if exitErr != nil {
			return nil, fmt.Errorf("failed to intercept iOS app %s: %s", target, errMsg)
		}
		return nil, fmt.Errorf("frida exited unexpectedly for iOS %s: %s", target, errMsg)
	case <-time.After(3 * time.Second):
		// Frida is still running — assume success
	}

	i.mu.Lock()
	i.cmds[proxyPort] = cmd
	i.stdins[proxyPort] = pw
	i.mu.Unlock()

	return map[string]any{"started": true, "target": target}, nil
}

func (i *iosFrida) Deactivate(proxyPort int, _ map[string]any) error {
	i.init()
	i.mu.Lock()
	scriptPath := i.scriptPaths[proxyPort]
	delete(i.scriptPaths, proxyPort)

	pw, okPw := i.stdins[proxyPort]
	if okPw && pw != nil {
		_ = pw.Close()
	}
	delete(i.stdins, proxyPort)

	cmd, ok := i.cmds[proxyPort]
	if ok && cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
		// NOTE: Do NOT call cmd.Process.Wait() here — the goroutine in Activate already
		// calls cmd.Wait() to detect early exits. Double-Wait causes a panic.
	}
	delete(i.cmds, proxyPort)
	i.mu.Unlock()

	if scriptPath != "" {
		os.Remove(scriptPath)
	}

	return nil
}

func (i *iosFrida) Metadata(string) (any, error) {
	hosts := make(map[string]any)
	if err := execRunTimeout(probeTimeout, "frida-ps", "-U"); err == nil {
		hosts["usb"] = map[string]any{
			"id":    "usb",
			"name":  "iOS Device",
			"state": "available",
		}
	} else {
		hosts["usb"] = map[string]any{
			"id":    "usb",
			"name":  "iOS Device",
			"state": "unavailable",
		}
	}
	return map[string]any{
		"type":  "frida-ios",
		"hosts": hosts,
	}, nil
}

type electronInterceptor struct{ base *stubInterceptor }

func (e *electronInterceptor) IsActivable() (bool, error)   { return true, nil }
func (e *electronInterceptor) IsActive(p int) (bool, error) { return e.base.IsActive(p) }
func (e *electronInterceptor) Activate(proxyPort int, options map[string]any) (map[string]any, error) {
	appPath, _ := options["path"].(string)
	if appPath == "" {
		return nil, fmt.Errorf("electron app path required")
	}
	preload := assetPath(e.base.cfg, "js", "prepend-electron.js")
	cmd := execCommand(appPath)

	env := buildTerminalEnv(e.base.cfg, proxyPort)
	env = append(env, "HTTP_TOOLKIT_SPKI="+e.base.spki)
	env = append(env, "NODE_OPTIONS=--require="+preload)

	cmd.Env = append(os.Environ(), env...)
	return map[string]any{"pid": 0}, cmd.Start()
}
func (e *electronInterceptor) Deactivate(_ int, _ map[string]any) error { return nil }
func (e *electronInterceptor) Metadata(string) (any, error) {
	return map[string]any{"type": "electron"}, nil
}

type dockerAttach struct{ base *stubInterceptor }

func (d *dockerAttach) SubMetadata(subID string) (any, error) {
	list, err := dockersvc.ListDetailedContainers()
	if err != nil {
		return nil, err
	}
	return map[string]any{"containers": list, "subId": subID}, nil
}

func (d *dockerAttach) IsActivable() (bool, error) { return dockersvc.IsAvailable(), nil }

func (d *dockerAttach) IsActive(proxyPort int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	out, err := system.CommandContext(ctx, "docker", "ps", "-q", "--filter",
		fmt.Sprintf("label=%s=%d", dockersvc.ContainerLabel, proxyPort)).Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) != "", nil
}

func (d *dockerAttach) Activate(proxyPort int, options map[string]any) (map[string]any, error) {
	container, _ := options["containerId"].(string)
	if container == "" {
		return nil, fmt.Errorf("containerId required")
	}
	certPath := filepath.Join(d.base.cfg.ConfigDir, "ca.pem")
	if err := dockersvc.RestartAndInjectContainer(container, dockersvc.InterceptionSettings{
		ProxyPort: proxyPort,
		CertPath:  certPath,
	}); err != nil {
		return nil, err
	}
	return map[string]any{"ok": true, "containerId": container}, nil
}

func (d *dockerAttach) Deactivate(proxyPort int, _ map[string]any) error {
	return dockersvc.DeleteInterceptedContainers(proxyPort)
}

func (d *dockerAttach) Metadata(_ string) (any, error) {
	targets, err := dockersvc.ListDetailedContainers()
	if err != nil {
		return nil, err
	}
	return map[string]any{"targets": targets}, nil
}

func optionFloat(options map[string]any, key string) float64 {
	if v, ok := options[key].(float64); ok {
		return v
	}
	return 0
}

var knownAppProblematicPorts = map[string][]int{
	"com.spotify.music":  {4070},
	"com.spotify.client": {4070},
}

func getPrimaryLANIP() string {
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					ipStr := ip4.String()
					if !strings.HasPrefix(ipStr, "10.0.2.") && !strings.HasPrefix(ipStr, "10.0.3.") && !strings.HasPrefix(ipStr, "169.254.") {
						return ipStr
					}
				}
			}
		}
	}
	if err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					return ip4.String()
				}
			}
		}
	}
	return "127.0.0.1"
}

func buildFridaConfig(configTemplate string, caCertContent string, proxyHost string, proxyPort int, ignoredPorts []int, enableSocks bool) string {
	cleanCert := strings.ReplaceAll(strings.TrimSpace(caCertContent), "\r", "")
	ignoredPortsJSON := "[]"
	if len(ignoredPorts) > 0 {
		var portStrs []string
		for _, p := range ignoredPorts {
			portStrs = append(portStrs, strconv.Itoa(p))
		}
		ignoredPortsJSON = "[" + strings.Join(portStrs, ",") + "]"
	}

	res := configTemplate

	// Replace CERT_PEM = `...`
	reCert := regexp.MustCompile("(?s)const CERT_PEM = `[^`]*`;?")
	res = reCert.ReplaceAllString(res, "const CERT_PEM = `"+cleanCert+"`;")

	// Replace PROXY_HOST = '...'
	reHost := regexp.MustCompile(`const PROXY_HOST = '[^']*';?`)
	res = reHost.ReplaceAllString(res, "const PROXY_HOST = '"+proxyHost+"';")

	// Replace PROXY_PORT = ...
	rePort := regexp.MustCompile(`const PROXY_PORT = \d+;?`)
	res = rePort.ReplaceAllString(res, "const PROXY_PORT = "+strconv.Itoa(proxyPort)+";")

	// Replace PROXY_SUPPORTS_SOCKS5 = ...
	reSocks := regexp.MustCompile(`const PROXY_SUPPORTS_SOCKS5 = (true|false);?`)
	res = reSocks.ReplaceAllString(res, "const PROXY_SUPPORTS_SOCKS5 = "+strconv.FormatBool(enableSocks)+";")

	// Replace IGNORED_NON_HTTP_PORTS = ...
	reIgnored := regexp.MustCompile(`const IGNORED_NON_HTTP_PORTS = \[[^\]]*\];?`)
	res = reIgnored.ReplaceAllString(res, "const IGNORED_NON_HTTP_PORTS = "+ignoredPortsJSON+";")

	return res
}

func buildAndroidFridaScript(cfg *config.Config, caCertContent string, proxyHost string, proxyPort int, ignoredPorts []int, enableSocks bool) (string, error) {
	javaBridge, err := os.ReadFile(assetPath(cfg, "frida", "frida-java-bridge.js"))
	if err != nil {
		return "", err
	}

	configTemplate, err := os.ReadFile(assetPath(cfg, "frida", "config.js"))
	if err != nil {
		return "", err
	}
	templatedConfig := buildFridaConfig(string(configTemplate), caCertContent, proxyHost, proxyPort, ignoredPorts, enableSocks)

	nativeConnect, err := os.ReadFile(assetPath(cfg, "frida", "native-connect-hook.js"))
	if err != nil {
		return "", err
	}

	nativeTls, err := os.ReadFile(assetPath(cfg, "frida", "native-tls-hook.js"))
	if err != nil {
		return "", err
	}

	proxyOverride, err := os.ReadFile(assetPath(cfg, "frida", "android", "android-proxy-override.js"))
	if err != nil {
		return "", err
	}

	systemCertInject, err := os.ReadFile(assetPath(cfg, "frida", "android", "android-system-certificate-injection.js"))
	if err != nil {
		return "", err
	}

	certUnpinning, err := os.ReadFile(assetPath(cfg, "frida", "android", "android-certificate-unpinning.js"))
	if err != nil {
		return "", err
	}

	certUnpinningFallback, err := os.ReadFile(assetPath(cfg, "frida", "android", "android-certificate-unpinning-fallback.js"))
	if err != nil {
		return "", err
	}

	disableRoot, err := os.ReadFile(assetPath(cfg, "frida", "android", "android-disable-root-detection.js"))
	if err != nil {
		return "", err
	}

	scripts := []string{
		string(javaBridge),
		templatedConfig,
		string(nativeConnect),
		string(nativeTls),
		string(proxyOverride),
		string(systemCertInject),
		string(certUnpinning),
		string(certUnpinningFallback),
		string(disableRoot),
	}

	return strings.Join(scripts, "\n"), nil
}

func buildIosFridaScript(cfg *config.Config, caCertContent string, proxyHost string, proxyPort int, ignoredPorts []int, enableSocks bool) (string, error) {
	objcBridge, err := os.ReadFile(assetPath(cfg, "frida", "frida-objc-bridge.js"))
	if err != nil {
		return "", err
	}

	configTemplate, err := os.ReadFile(assetPath(cfg, "frida", "config.js"))
	if err != nil {
		return "", err
	}
	templatedConfig := buildFridaConfig(string(configTemplate), caCertContent, proxyHost, proxyPort, ignoredPorts, enableSocks)

	iosConnect, err := os.ReadFile(assetPath(cfg, "frida", "ios", "ios-connect-hook.js"))
	if err != nil {
		return "", err
	}

	iosDisable, err := os.ReadFile(assetPath(cfg, "frida", "ios", "ios-disable-detection.js"))
	if err != nil {
		return "", err
	}

	nativeTls, err := os.ReadFile(assetPath(cfg, "frida", "native-tls-hook.js"))
	if err != nil {
		return "", err
	}

	nativeConnect, err := os.ReadFile(assetPath(cfg, "frida", "native-connect-hook.js"))
	if err != nil {
		return "", err
	}

	scripts := []string{
		string(objcBridge),
		templatedConfig,
		string(iosConnect),
		string(iosDisable),
		string(nativeTls),
		string(nativeConnect),
	}

	return strings.Join(scripts, "\n"), nil
}
