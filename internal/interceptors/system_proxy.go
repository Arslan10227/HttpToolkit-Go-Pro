package interceptors

import (
	"fmt"
	"runtime"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/system"
)

type systemProxy struct{ base *stubInterceptor }

func (s *systemProxy) IsActivable() (bool, error) {
	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		return true, nil
	default:
		return false, nil
	}
}

func (s *systemProxy) IsActive(proxyPort int) (bool, error) {
	return s.base.IsActive(proxyPort)
}

func (s *systemProxy) Activate(proxyPort int, _ map[string]any) (map[string]any, error) {
	if err := setSystemProxy("127.0.0.1", proxyPort); err != nil {
		return nil, err
	}
	return map[string]any{"proxyPort": proxyPort}, nil
}

func (s *systemProxy) Deactivate(_ int, _ map[string]any) error {
	return clearSystemProxy()
}

func (s *systemProxy) Metadata(string) (any, error) {
	return map[string]any{"type": "system-proxy"}, nil
}

func setSystemProxy(host string, port int) error {
	switch runtime.GOOS {
	case "windows":
		return setWindowsProxy(host, port)
	case "darwin":
		if host == "" {
			_, err := system.Command("networksetup", "-setwebproxystate", "Wi-Fi", "off").CombinedOutput()
			_, err2 := system.Command("networksetup", "-setsecurewebproxystate", "Wi-Fi", "off").CombinedOutput()
			if err != nil {
				return err
			}
			return err2
		}
		proxy := fmt.Sprintf("%s:%d", host, port)
		if _, err := system.Command("networksetup", "-setwebproxy", "Wi-Fi", host, fmt.Sprint(port)).CombinedOutput(); err != nil {
			return err
		}
		_, err := system.Command("networksetup", "-setsecurewebproxy", "Wi-Fi", host, fmt.Sprint(port)).CombinedOutput()
		_ = proxy
		return err
	case "linux":
		if host == "" {
			_, err := system.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "none").CombinedOutput()
			return err
		}
		if _, err := system.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "manual").CombinedOutput(); err != nil {
			return err
		}
		_, err := system.Command("gsettings", "set", "org.gnome.system.proxy.http", "host", host).CombinedOutput()
		if err != nil {
			return err
		}
		_, err = system.Command("gsettings", "set", "org.gnome.system.proxy.http", "port", fmt.Sprint(port)).CombinedOutput()
		if err != nil {
			return err
		}
		_, err = system.Command("gsettings", "set", "org.gnome.system.proxy.https", "host", host).CombinedOutput()
		if err != nil {
			return err
		}
		_, err = system.Command("gsettings", "set", "org.gnome.system.proxy.https", "port", fmt.Sprint(port)).CombinedOutput()
		return err
	default:
		return fmt.Errorf("system proxy on %s: use manual setup", runtime.GOOS)
	}
}

func clearSystemProxy() error {
	return setSystemProxy("", 0)
}

func setWindowsProxy(host string, port int) error {
	if host == "" {
		_, err := system.Command("reg", "add",
			`HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
			"/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f").CombinedOutput()
		refreshWindowsProxySettings()
		return err
	}
	proxy := fmt.Sprintf("%s:%d", host, port)
	if _, err := system.Command("reg", "add",
		`HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		"/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f").CombinedOutput(); err != nil {
		return err
	}
	if _, err := system.Command("reg", "add",
		`HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		"/v", "ProxyOverride", "/t", "REG_SZ", "/d", "<local>;127.0.0.1;localhost;httptoolkitpro.vercel.app;*.vercel.app;*.vercel.com", "/f").CombinedOutput(); err != nil {
		return err
	}
	_, err := system.Command("reg", "add",
		`HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		"/v", "ProxyServer", "/t", "REG_SZ", "/d", proxy, "/f").CombinedOutput()
	refreshWindowsProxySettings()
	return err
}

func refreshWindowsProxySettings() {
	refreshCmd := `
$sig = '[DllImport("wininet.dll", SetLastError = true)] public static extern bool InternetSetOption(IntPtr hInit, int opt, IntPtr buf, int bufLen);'
$type = Add-Type -MemberDefinition $sig -Name "WinInet" -Namespace "Win32" -PassThru
$type::InternetSetOption([IntPtr]::Zero, 39, [IntPtr]::Zero, 0)
$type::InternetSetOption([IntPtr]::Zero, 37, [IntPtr]::Zero, 0)
`
	_, _ = system.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", refreshCmd).CombinedOutput()
}
