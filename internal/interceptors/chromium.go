package interceptors

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	certmgr "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
)

type chromiumFresh struct {
	base       *stubInterceptor
	browser    browserSpec
	spki       string
	certs      *certmgr.Manager
	hideSrv    *amiusingServer
	cmd        *exec.Cmd
	profileDir string
}

type chromiumExisting struct {
	base    *stubInterceptor
	browser browserSpec
	spki    string
	certs   *certmgr.Manager
	cmd     *exec.Cmd
}

type browserSpec struct {
	name    string
	winPath []string
	darwin  []string
	linux   []string
}

func browserForID(id string) browserSpec {
	m := map[string]browserSpec{
		"fresh-chrome": {
			name:    "chrome",
			winPath: []string{`Google\Chrome\Application\chrome.exe`},
			darwin:  []string{"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"},
			linux:   []string{"/usr/bin/google-chrome", "/usr/bin/google-chrome-stable", "/usr/bin/chromium-browser"},
		},
		"existing-chrome": {
			name:    "chrome",
			winPath: []string{`Google\Chrome\Application\chrome.exe`},
			darwin:  []string{"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"},
			linux:   []string{"/usr/bin/google-chrome", "/usr/bin/google-chrome-stable"},
		},
		"fresh-chrome-beta": {
			name:    "chrome-beta",
			winPath: []string{`Google\Chrome Beta\Application\chrome.exe`},
			darwin:  []string{"/Applications/Google Chrome Beta.app/Contents/MacOS/Google Chrome Beta"},
		},
		"fresh-chrome-dev": {
			name:    "chrome-dev",
			winPath: []string{`Google\Chrome Dev\Application\chrome.exe`},
			darwin:  []string{"/Applications/Google Chrome Dev.app/Contents/MacOS/Google Chrome Dev"},
		},
		"fresh-chrome-canary": {
			name:    "chrome-canary",
			winPath: []string{`Google\Chrome SxS\Application\chrome.exe`},
			darwin:  []string{"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary"},
		},
		"fresh-chromium": {
			name:    "chromium",
			winPath: []string{`Chromium\Application\chrome.exe`},
			darwin:  []string{"/Applications/Chromium.app/Contents/MacOS/Chromium"},
			linux:   []string{"/usr/bin/chromium", "/usr/bin/chromium-browser", "/snap/bin/chromium"},
		},
		"existing-chromium": {
			name:    "chromium",
			winPath: []string{`Chromium\Application\chrome.exe`},
			darwin:  []string{"/Applications/Chromium.app/Contents/MacOS/Chromium"},
			linux:   []string{"/usr/bin/chromium", "/usr/bin/chromium-browser"},
		},
		"fresh-chromium-dev": {
			name:  "chromium-dev",
			linux: []string{"/usr/bin/chromium"},
		},
		"fresh-edge": {
			name:    "msedge",
			winPath: []string{`Microsoft\Edge\Application\msedge.exe`},
			darwin:  []string{"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge"},
		},
		"fresh-edge-beta": {
			name:    "msedge-beta",
			winPath: []string{`Microsoft\Edge Beta\Application\msedge.exe`},
		},
		"fresh-edge-dev": {
			name:    "msedge-dev",
			winPath: []string{`Microsoft\Edge Dev\Application\msedge.exe`},
		},
		"fresh-edge-canary": {
			name:    "msedge-canary",
			winPath: []string{`Microsoft\Edge SxS\Application\msedge.exe`},
		},
		"fresh-brave": {
			name:    "brave",
			winPath: []string{`BraveSoftware\Brave-Browser\Application\brave.exe`},
			darwin:  []string{"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser"},
			linux:   []string{"/usr/bin/brave-browser"},
		},
		"fresh-opera": {
			name:    "opera",
			winPath: []string{`Opera\launcher.exe`, `Opera\opera.exe`},
			darwin:  []string{"/Applications/Opera.app/Contents/MacOS/Opera"},
		},
		"existing-arc": {
			name: "arc",
			darwin: []string{
				"/Applications/Arc.app/Contents/MacOS/Arc",
				os.Getenv("HOME") + "/Applications/Arc.app/Contents/MacOS/Arc",
			},
		},
	}
	if s, ok := m[id]; ok {
		return s
	}
	return browserSpec{
		name:    id,
		winPath: []string{`Google\Chrome\Application\chrome.exe`},
		darwin:  []string{"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"},
	}
}

func (c *chromiumFresh) IsActivable() (bool, error) {
	_, err := findBrowser(c.browser)
	return err == nil, nil
}

func (c *chromiumFresh) IsActive(proxyPort int) (bool, error) {
	if c.cmd != nil && c.cmd.Process != nil && c.cmd.ProcessState != nil && c.cmd.ProcessState.Exited() {
		c.base.setActive(proxyPort, false)
		return false, nil
	}
	return c.base.IsActive(proxyPort)
}

func (c *chromiumFresh) Activate(proxyPort int, _ map[string]any) (map[string]any, error) {
	exe, err := findBrowser(c.browser)
	if err != nil {
		return nil, err
	}

	// Start a self-hosted "am I being intercepted" check server. This removes
	// the dependency on amiusing.httptoolkit.tech and verifies locally that the
	// browser is routing through the proxy and trusts the CA.
	c.hideSrv, err = startAmiusingServer(c.certs)
	if err != nil {
		return nil, fmt.Errorf("amiusing server: %w", err)
	}

	c.profileDir = freshProfileDir(c.base.id)
	_ = os.RemoveAll(c.profileDir)
	if err := os.MkdirAll(c.profileDir, 0o700); err != nil {
		c.hideSrv.Stop()
		c.hideSrv = nil
		return nil, err
	}

	args := chromiumArgs(c.base, c.spki, proxyPort, c.profileDir, c.hideSrv, c.browser.name)
	c.cmd = exec.Command(exe, args...)
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	if err := c.cmd.Start(); err != nil {
		c.hideSrv.Stop()
		c.hideSrv = nil
		c.cmd = nil
		return nil, err
	}
	go func() { _ = c.cmd.Wait() }()

	// Wait for the hide-warning page to be loaded by the browser. If Chrome
	// doesn't load it within the timeout, kill the process and fail activation.
	if err := c.hideSrv.WaitSuccess(25 * time.Second); err != nil {
		_ = c.cmd.Process.Kill()
		c.hideSrv.Stop()
		c.hideSrv = nil
		c.cmd = nil
		return nil, fmt.Errorf("Chrome hide-warning check failed: %w", err)
	}

	return map[string]any{"pid": c.cmd.Process.Pid}, nil
}

func (c *chromiumFresh) Deactivate(_ int, _ map[string]any) error {
	if c.cmd != nil && c.cmd.Process != nil {
		_ = c.cmd.Process.Kill()
		c.cmd = nil
	}
	if c.hideSrv != nil {
		c.hideSrv.Stop()
		c.hideSrv = nil
	}
	if c.profileDir != "" {
		// Opera's launcher exits immediately on Windows and may still hold the
		// profile, so skip cleanup there.
		if runtime.GOOS != "windows" || !strings.Contains(c.browser.name, "opera") {
			_ = os.RemoveAll(c.profileDir)
		}
		c.profileDir = ""
	}
	return nil
}

func (c *chromiumFresh) Metadata(string) (any, error) {
	return map[string]any{"browser": c.browser.name}, nil
}

func (c *chromiumExisting) IsActivable() (bool, error) {
	_, err := findBrowser(c.browser)
	return err == nil, nil
}
func (c *chromiumExisting) IsActive(p int) (bool, error) {
	if c.cmd != nil && c.cmd.Process != nil && c.cmd.ProcessState != nil && c.cmd.ProcessState.Exited() {
		c.base.setActive(p, false)
		return false, nil
	}
	return c.base.IsActive(p)
}
func (c *chromiumExisting) Activate(p int, _ map[string]any) (map[string]any, error) {
	exe, err := findBrowser(c.browser)
	if err != nil {
		return nil, err
	}
	// Existing Chrome also uses the self-hosted amiusing page to verify
	// interception is working before reporting activation success.
	hideSrv, err := startAmiusingServer(c.certs)
	if err != nil {
		return nil, fmt.Errorf("amiusing server: %w", err)
	}
	args := chromiumArgs(c.base, c.spki, p, "", hideSrv, c.browser.name)
	c.cmd = exec.Command(exe, args...)
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	if err := c.cmd.Start(); err != nil {
		hideSrv.Stop()
		return nil, err
	}
	go func() { _ = c.cmd.Wait() }()
	// Wait for the hide-warning page to be loaded, then stop the server.
	if err := hideSrv.WaitSuccess(25 * time.Second); err != nil {
		hideSrv.Stop()
		return nil, fmt.Errorf("Chrome hide-warning check failed: %w", err)
	}
	// Stop the server after a short delay (matching Node.js which delays 10s,
	// but we stop immediately since we don't need to keep it for Opera reloads).
	go func() {
		time.Sleep(10 * time.Second)
		hideSrv.Stop()
	}()
	return map[string]any{"pid": c.cmd.Process.Pid, "existingProfile": true}, nil
}
func (c *chromiumExisting) Deactivate(_ int, _ map[string]any) error {
	if c.cmd != nil && c.cmd.Process != nil {
		_ = c.cmd.Process.Kill()
		c.cmd = nil
	}
	return nil
}
func (c *chromiumExisting) Metadata(string) (any, error) {
	return map[string]any{"browser": c.browser.name, "existing": true}, nil
}

func chromiumArgs(base *stubInterceptor, spki string, proxyPort int, profileDir string, amiusing *amiusingServer, browserName string) []string {
	proxyURL := fmt.Sprintf("https://127.0.0.1:%d", proxyPort)

	bypass := []string{"<-loopback>",
		"internal.httptoolkit.localhost",
		"tauri.localhost",
		"localhost",
		"127.0.0.1",
		// INTERNAL_PROXY_BYPASS_SUFFIXES
		"vercel.app",
		"vercel-dns.com",
		"identitytoolkit.googleapis.com",
		"securetoken.googleapis.com",
		"www.googleapis.com",
		"www.gstatic.com",
		"gstatic.com",
		"googleapis.com",
		"google-analytics.com",
		"accounts.google.com",
	}

	args := []string{
		fmt.Sprintf("--proxy-server=%s", proxyURL),
		"--proxy-bypass-list=" + strings.Join(bypass, ","),
		"--test-type",
		"--ignore-certificate-errors-spki-list=" + strings.TrimSpace(spki),
		"--disable-quic",
		"--disable-features=" + strings.Join([]string{
			"CertificateTransparency",
			"ChromeWhatsNewUI",
			"SidePanelPinning",
			"OptimizationGuideModelDownloading",
			"OptimizationHintsFetching",
			"OptimizationTargetPrediction",
			"OptimizationHints",
		}, ","),
		"--disable-background-networking",
		"--no-default-browser-check",
		"--no-first-run",
		"--disable-popup-blocking",
	}

	if !strings.Contains(browserName, "brave") {
		args = append(args,
			"--component-updater=url-source=http://disabled-chromium-update.localhost:0",
			"--check-for-update-interval=31536000",
		)
	}

	if profileDir != "" {
		args = append([]string{"--user-data-dir=" + profileDir}, args...)
	}

	ext := chromiumExtensionPath(base.cfg)
	if _, err := os.Stat(ext); err == nil {
		args = append(args, "--load-extension="+ext)
	}

	if amiusing != nil {
		args = append(args, amiusing.URL())
	}

	return args
}

func findBrowser(b browserSpec) (string, error) {
	var candidates []string
	switch runtime.GOOS {
	case "windows":
		for _, rel := range b.winPath {
			candidates = append(candidates,
				filepath.Join(os.Getenv("ProgramFiles"), rel),
				filepath.Join(os.Getenv("ProgramFiles(x86)"), rel),
				filepath.Join(os.Getenv("LOCALAPPDATA"), rel),
			)
		}
	case "darwin":
		candidates = append(candidates, b.darwin...)
	default:
		candidates = append(candidates, b.linux...)
	}
	for _, p := range candidates {
		if p == "" {
			continue
		}
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("browser %s not found", b.name)
}

func freshProfileDir(id string) string {
	return filepath.Join(os.TempDir(), "httptoolkit-"+id)
}
