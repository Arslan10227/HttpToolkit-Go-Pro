package interceptors

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/logger"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/system"
)

const firefoxCANickname = "HTTP Toolkit"

type certutilCmd struct {
	Path string
	Env  []string
}

func nssPlatform() string {
	switch runtime.GOOS {
	case "windows":
		return "win32"
	case "darwin":
		return "darwin"
	default:
		return "linux"
	}
}

func nssAssetDir(assetsDir string) string {
	return filepath.Join(assetsDir, "nss", nssPlatform())
}

func certutilBinName() string {
	if runtime.GOOS == "windows" {
		return "certutil.exe"
	}
	return "certutil"
}

func mergeEnv(extra []string) []string {
	if len(extra) == 0 {
		return os.Environ()
	}
	return append(os.Environ(), extra...)
}

func certutilEnv(assetsDir string) []string {
	if runtime.GOOS != "linux" {
		return nil
	}
	libPath := nssAssetDir(assetsDir)
	if existing := os.Getenv("LD_LIBRARY_PATH"); existing != "" {
		return []string{"LD_LIBRARY_PATH=" + libPath + ":" + existing}
	}
	return []string{"LD_LIBRARY_PATH=" + libPath}
}

func testNSSCertutil(path string, extraEnv []string) bool {
	cmd := system.Command(path, "-h")
	cmd.Env = mergeEnv(extraEnv)
	out, err := cmd.CombinedOutput()
	if err == nil {
		return false
	}
	text := string(out)
	return strings.Contains(text, "NSS certificate databases") ||
		strings.Contains(text, "Utility to manipulate")
}

func resolveCertutil(assetsDir string) (certutilCmd, error) {
	if testNSSCertutil("certutil", nil) {
		return certutilCmd{Path: "certutil"}, nil
	}
	bundled := filepath.Join(nssAssetDir(assetsDir), certutilBinName())
	env := certutilEnv(assetsDir)
	if testNSSCertutil(bundled, env) {
		return certutilCmd{Path: bundled, Env: env}, nil
	}
	return certutilCmd{}, fmt.Errorf("no NSS certutil available")
}

func firefoxCertutilAvailable(assetsDir string) bool {
	_, err := resolveCertutil(assetsDir)
	return err == nil
}

func firefoxProfilePath(cfg *config.Config, firefoxPath string) string {
	if isFlatpak(firefoxPath) {
		return filepath.Join(os.Getenv("HOME"), ".var", "app", "org.mozilla.firefox", "data", "mozilla", "firefox", "httptoolkit")
	}
	if isSnap(firefoxPath) {
		return filepath.Join(snapConfigPath("firefox"), "profile")
	}
	return filepath.Join(cfg.ConfigDir, "firefox-profile")
}

func isFlatpak(bin string) bool {
	if runtime.GOOS != "linux" || bin == "" {
		return false
	}
	return filepath.Base(bin) == "flatpak"
}

func snapConfigPath(variant string) string {
	return filepath.Join(os.Getenv("HOME"), "snap", variant, "current", ".httptoolkit")
}

func isSnap(bin string) bool {
	if runtime.GOOS != "linux" || bin == "" {
		return false
	}
	if strings.HasPrefix(bin, "/snap/bin/") {
		return true
	}
	data, err := os.ReadFile(bin)
	if err != nil {
		return false
	}
	if len(data) > 100 {
		data = data[len(data)-100:]
	}
	return strings.Contains(string(data), "exec /snap/bin/")
}

func ensureFirefoxProfile(cfg *config.Config, firefoxPath, profilePath, certPath string) error {
	if err := os.MkdirAll(profilePath, 0o700); err != nil {
		return err
	}
	prefsPath := filepath.Join(profilePath, "prefs.js")
	if _, err := os.Stat(prefsPath); os.IsNotExist(err) {
		if err := bootstrapFirefoxProfile(firefoxPath, profilePath); err != nil {
			return err
		}
	}
	return injectFirefoxCA(cfg.AssetsDir, profilePath, certPath)
}

func bootstrapFirefoxProfile(firefoxPath, profilePath string) error {
	srv := startWaitableSetupServer()
	defer srv.stop()

	args := []string{"-profile", profilePath, "-no-remote", srv.url()}
	if isFlatpak(firefoxPath) {
		args = append([]string{"run", "org.mozilla.firefox"}, args...)
	}
	cmd := system.Command(firefoxPath, args...)
	if runtime.GOOS == "windows" {
		cmd.Env = append(os.Environ(), "MOZ_NO_REMOTE=1")
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start firefox for profile setup: %w", err)
	}

	waitErr := srv.wait(45 * time.Second)
	if err := cmd.Process.Kill(); err != nil {
		logger.Error(err, map[string]any{"msg": "failed to kill firefox profile setup process"})
	}
	if _, err := cmd.Process.Wait(); err != nil {
		logger.Debug("firefox profile setup wait returned error", map[string]any{"err": err.Error()})
	}

	if waitErr != nil {
		return fmt.Errorf("firefox profile setup: %w", waitErr)
	}
	return nil
}

func injectFirefoxCA(assetsDir, profilePath, certPath string) error {
	if certPath == "" {
		return fmt.Errorf("certificate path required")
	}
	if _, err := os.Stat(certPath); err != nil {
		return fmt.Errorf("certificate: %w", err)
	}
	cu, err := resolveCertutil(assetsDir)
	if err != nil {
		return err
	}
	dbPath := "sql:" + profilePath
	del := system.Command(cu.Path, "-D", "-d", dbPath, "-n", firefoxCANickname)
	del.Env = mergeEnv(cu.Env)
	if err := del.Run(); err != nil {
		logger.Debug("ignoring certutil delete error", map[string]any{"err": err.Error()})
	}

	add := system.Command(cu.Path, "-A", "-d", dbPath, "-t", "C,,", 
		"-i", certPath, "-n", firefoxCANickname)
	add.Env = mergeEnv(cu.Env)
	out, err := add.CombinedOutput()
	if err != nil {
		if _, statErr := os.Stat(filepath.Join(profilePath, "cert9.db")); statErr != nil {
			initDB := system.Command(cu.Path, "-N", "-d", dbPath, "--empty-password")
			initDB.Env = mergeEnv(cu.Env)
			if initOut, initErr := initDB.CombinedOutput(); initErr != nil {
				return fmt.Errorf("certutil -N: %w: %s", initErr, initOut)
			}
			add = system.Command(cu.Path, "-A", "-d", dbPath, "-t", "C,,",
				"-i", certPath, "-n", firefoxCANickname)
			add.Env = mergeEnv(cu.Env)
			out, err = add.CombinedOutput()
		}
		if err != nil {
			return fmt.Errorf("certutil -A: %w: %s", err, out)
		}
	}
	return nil
}

func writeFirefoxUserJS(profilePath string, proxyPort int) error {
	content := fmt.Sprintf(`// HTTP Toolkit interception preferences
user_pref("network.proxy.type", 1);
user_pref("network.proxy.http", "127.0.0.1");
user_pref("network.proxy.http_port", %d);
user_pref("network.proxy.ssl", "127.0.0.1");
user_pref("network.proxy.ssl_port", %d);
user_pref("network.proxy.allow_hijacking_localhost", true);
user_pref("network.proxy.no_proxies_on", "");
user_pref("network.captive-portal-service.enabled", false);
user_pref("app.update.auto", false);
user_pref("app.shield.optoutstudies.enabled", false);
user_pref("browser.startup.homepage_override.mstone", "ignore");
user_pref("browser.showQuitWarning", false);
user_pref("browser.tabs.warnOnClose", false);
user_pref("browser.tabs.warnOnCloseOtherTabs", false);
user_pref("browser.chrome.toolbar_tips", false);
user_pref("browser.uitour.enabled", false);
user_pref("browser.usedOnWindows10", true);
user_pref("browser.usedOnWindows10.introURL", "");
user_pref("browser.newtabpage.activity-stream.showSponsoredCheckboxes", false);
user_pref("browser.newtabpage.activity-stream.showSponsoredTopSites", false);
user_pref("datareporting.healthreport.uploadEnabled", false);
user_pref("datareporting.usage.uploadEnabled", false);
user_pref("datareporting.healthreport.service.firstRun", false);
user_pref("datareporting.policy.dataSubmissionEnabled", false);
user_pref("datareporting.policy.dataSubmissionPolicyAccepted", false);
user_pref("datareporting.policy.dataSubmissionPolicyBypassNotification", true);
user_pref("toolkit.telemetry.reportingpolicy.firstRun", false);
user_pref("browser.reader.detectedFirstArticle", false);
user_pref("trailhead.firstrun.didSeeAboutWelcome", true);
user_pref("privacy.history.custom", true);
user_pref("privacy.sanitize.sanitizeOnShutdown", true);
user_pref("privacy.clearOnShutdown.cache", true);
user_pref("privacy.clearOnShutdown.cookies", true);
user_pref("privacy.clearOnShutdown.downloads", true);
user_pref("privacy.clearOnShutdown.formdata", true);
user_pref("privacy.clearOnShutdown.history", true);
user_pref("privacy.clearOnShutdown.offlineApps", true);
user_pref("privacy.clearOnShutdown.siteSettings", true);
user_pref("privacy.clearOnShutdown.openWindows", false);
`, proxyPort, proxyPort)
	return os.WriteFile(filepath.Join(profilePath, "user.js"), []byte(content), 0o644)
}

func findFirefoxBinary() (string, error) {
	if runtime.GOOS == "windows" {
		candidates := []string{
			filepath.Join(os.Getenv("ProgramFiles"), "Mozilla Firefox", "firefox.exe"),
			filepath.Join(os.Getenv("ProgramFiles(x86)"), "Mozilla Firefox", "firefox.exe"),
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Mozilla Firefox", "firefox.exe"),
		}
		for _, p := range candidates {
			if p == "" {
				continue
			}
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}
		}
	}
	if p, err := exec.LookPath("firefox"); err == nil {
		return p, nil
	}
	if runtime.GOOS == "linux" {
		if fp, err := exec.LookPath("flatpak"); err == nil {
			// Check if the official Firefox flatpak is installed, with a short
			// timeout so a broken flatpak install cannot hang activation.
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := exec.CommandContext(ctx, fp, "info", "org.mozilla.firefox").Run(); err == nil {
				return fp, nil
			}
		}
	}
	if runtime.GOOS == "darwin" {
		p := "/Applications/Firefox.app/Contents/MacOS/firefox"
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("firefox not found")
}
