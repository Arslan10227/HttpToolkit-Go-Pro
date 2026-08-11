package cert

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/system"
)

func IsSystemCertInstalled(certPEM string) (bool, error) {
	switch runtime.GOOS {
	case "windows":
		// Check both User store and System store. We install to User store by default.
		out, err := system.Command("certutil", "-user", "-store", "Root").CombinedOutput()
		if err != nil || !strings.Contains(strings.ToLower(string(out)), "httptoolkit pro") {
			// Fallback/check machine store as well
			machineOut, err2 := system.Command("certutil", "-store", "Root").CombinedOutput()
			if err2 == nil {
				out = machineOut
			}
		}
		return strings.Contains(strings.ToLower(string(out)), "httptoolkit pro"), nil
	case "darwin":
		out, err := system.Command("security", "find-certificate", "-c", "HttpToolkit Pro").CombinedOutput()
		return err == nil && len(out) > 0, nil
	default:
		return false, nil
	}
}

func InstallSystemCert(certPath string) (installed bool, message string, err error) {
	switch runtime.GOOS {
	case "windows":
		cmd := system.Command("certutil", "-addstore", "-user", "Root", certPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return false, string(out), err
		}
		return true, string(out), nil
	case "darwin":
		cmd := system.Command("security", "add-trusted-cert", "-r", "trustRoot", "-p", "ssl", certPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return false, string(out), err
		}
		return true, string(out), nil
	default:
		// Debian/Ubuntu/Fedora: copy to local trust store
		local := filepath.Join(os.Getenv("HOME"), ".local/share/httptoolkit-ca.crt")
		if err := copyFile(certPath, local); err != nil {
			return false, err.Error(), err
		}
		if _, err := system.Command("update-ca-certificates", "--fresh").CombinedOutput(); err == nil {
			return true, "Installed via update-ca-certificates", nil
		}
		// Fallback: certutil (Firefox/system NSS)
		if certutil, err := exec.LookPath("certutil"); err == nil {
			cmd := system.Command(certutil, "-A", "-n", "HttpToolkit Pro", "-t", "C,,", "-i", certPath, "-d", "sql:"+filepath.Join(os.Getenv("HOME"), ".pki", "nssdb"))
			out, err := cmd.CombinedOutput()
			if err == nil {
				return true, string(out), nil
			}
		}
		return false, "Use manual install on this platform", fmt.Errorf("unsupported platform")
	}
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

type JavaVersionInfo struct {
	JavaPath    string `json:"javaPath"`
	KeytoolPath string `json:"keytoolPath"`
	Version     string `json:"version"`
}

func DetectJavaVersions() ([]JavaVersionInfo, error) {
	java := system.FindJava()
	if java == "java" {
		if _, err := exec.LookPath("java"); err != nil {
			return nil, nil
		}
	}
	out, err := system.Command(java, "-version").CombinedOutput()
	if err != nil {
		return nil, nil
	}
	version := strings.TrimSpace(string(out))
	keytool := strings.Replace(java, "java", "keytool", 1)
	if runtime.GOOS == "windows" {
		keytool = strings.Replace(java, "java.exe", "keytool.exe", 1)
	}
	return []JavaVersionInfo{{
		JavaPath: java, KeytoolPath: keytool, Version: version,
	}}, nil
}

func IsJavaCertInstalled(javaPath string) (bool, error) {
	keytool := keytoolFor(javaPath)
	out, err := system.Command(keytool, "-list", "-cacerts", "-storepass", "changeit", "-alias", "httptoolkit-pro").CombinedOutput()
	if err != nil {
		return false, nil
	}
	return strings.Contains(strings.ToLower(string(out)), "httptoolkit-pro"), nil
}

func InstallJavaCert(javaPath, certPath string) (map[string]any, error) {
	keytool := keytoolFor(javaPath)
	targets := javaKeystoreTargets(javaPath)
	tried := make([]string, 0, len(targets))
	var lastErr error
	for _, target := range targets {
		tried = append(tried, target.label)
		var cmd *exec.Cmd
		if target.cacerts {
			cmd = system.Command(keytool, "-importcert", "-noprompt", "-trustcacerts",
				"-cacerts", "-storepass", target.storepass, "-alias", "httptoolkit-pro", "-file", certPath)
		} else {
			cmd = system.Command(keytool, "-importcert", "-noprompt", "-trustcacerts",
				"-keystore", target.keystore, "-storepass", target.storepass,
				"-alias", "httptoolkit-pro", "-file", certPath)
		}
		out, err := cmd.CombinedOutput()
		if err == nil {
			javaHome := filepath.Dir(filepath.Dir(javaPath))
			return map[string]any{
				"success": true, "javaHome": javaHome,
				"installedTo": target.label, "fallbackTargetsTried": tried,
			}, nil
		}
		lastErr = fmt.Errorf("%s: %s", err, out)
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return map[string]any{"success": false, "fallbackTargetsTried": tried}, fmt.Errorf("no keystore succeeded")
}

type keystoreTarget struct {
	label     string
	keystore  string
	storepass string
	cacerts   bool
}

func javaKeystoreTargets(javaPath string) []keystoreTarget {
	home := filepath.Dir(filepath.Dir(javaPath))
	targets := []keystoreTarget{
		{label: "cacerts", cacerts: true, storepass: "changeit"},
	}
	if runtime.GOOS == "windows" {
		targets = append(targets, keystoreTarget{
			label: "jssecacerts", keystore: filepath.Join(home, "lib", "security", "jssecacerts"), storepass: "changeit",
		})
	}
	targets = append(targets, keystoreTarget{
		label: "lib/security/cacerts", keystore: filepath.Join(home, "lib", "security", "cacerts"), storepass: "changeit",
	})
	return targets
}

func keytoolFor(javaPath string) string {
	dir := filepath.Dir(javaPath)
	name := "keytool"
	if runtime.GOOS == "windows" {
		name = "keytool.exe"
	}
	p := filepath.Join(dir, name)
	if _, err := os.Stat(p); err == nil {
		return p
	}
	return filepath.Join(dir, name)
}
