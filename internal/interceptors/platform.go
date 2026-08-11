package interceptors

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	certmgr "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	dockersvc "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/docker"
)

type firefoxFresh struct {
	base  *stubInterceptor
	certs *certmgr.Manager
}

func (f *firefoxFresh) IsActivable() (bool, error) {
	if _, err := findFirefoxBinary(); err != nil {
		return false, nil
	}
	return firefoxCertutilAvailable(f.base.cfg.AssetsDir), nil
}

func (f *firefoxFresh) IsActive(p int) (bool, error) { return f.base.IsActive(p) }

func (f *firefoxFresh) Activate(proxyPort int, _ map[string]any) (map[string]any, error) {
	firefoxPath, err := findFirefoxBinary()
	if err != nil {
		return nil, err
	}
	profile := firefoxProfilePath(f.base.cfg, firefoxPath)
	certPath := filepath.Join(f.base.cfg.ConfigDir, "ca.pem")
	if err := ensureFirefoxProfile(f.base.cfg, firefoxPath, profile, certPath); err != nil {
		return nil, fmt.Errorf("firefox profile: %w", err)
	}
	if err := writeFirefoxUserJS(profile, proxyPort); err != nil {
		return nil, err
	}

	startURL := "about:blank"
	var checkSrv *amiusingServer
	if certFilesExist(f.base.cfg.ConfigDir) {
		var srvErr error
		checkSrv, srvErr = startAmiusingServer(f.certs)
		if srvErr == nil {
			startURL = checkSrv.URL()
		}
	}

	args := []string{"-profile", profile, "-no-remote", startURL}
	if isFlatpak(firefoxPath) {
		args = append([]string{"run", "org.mozilla.firefox"}, args...)
	}
	cmd := exec.Command(firefoxPath, args...)
	if runtime.GOOS == "windows" {
		cmd.Env = append(os.Environ(), "MOZ_NO_REMOTE=1")
	}
	if err := cmd.Start(); err != nil {
		if checkSrv != nil {
			checkSrv.Stop()
		}
		return nil, err
	}

	if checkSrv != nil {
		if err := checkSrv.WaitSuccess(25 * time.Second); err != nil {
			_ = cmd.Process.Kill()
			checkSrv.Stop()
			return nil, fmt.Errorf("Firefox certificate trust check failed: %w", err)
		}
		go func() {
			time.Sleep(3 * time.Second)
			checkSrv.Stop()
		}()
	}

	return map[string]any{"pid": cmd.Process.Pid, "profile": profile}, nil
}

func (f *firefoxFresh) Deactivate(_ int, _ map[string]any) error { return nil }
func (f *firefoxFresh) Metadata(string) (any, error) {
	return map[string]any{"browser": "firefox"}, nil
}

type safariFresh struct {
	base        *stubInterceptor
	systemProxy *systemProxy
}

func (s *safariFresh) IsActivable() (bool, error)   { return runtime.GOOS == "darwin", nil }
func (s *safariFresh) IsActive(p int) (bool, error) { return s.base.IsActive(p) }

func (s *safariFresh) Activate(proxyPort int, options map[string]any) (map[string]any, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("Safari only on macOS")
	}
	if s.systemProxy == nil {
		s.systemProxy = &systemProxy{base: s.base}
	}
	if _, err := s.systemProxy.Activate(proxyPort, options); err != nil {
		return nil, err
	}
	if err := exec.Command("open", "-a", "Safari").Start(); err != nil {
		_ = s.systemProxy.Deactivate(proxyPort, options)
		return nil, err
	}
	return map[string]any{"ok": true}, nil
}

func (s *safariFresh) Deactivate(proxyPort int, options map[string]any) error {
	if s.systemProxy != nil {
		return s.systemProxy.Deactivate(proxyPort, options)
	}
	return nil
}

func (s *safariFresh) Metadata(string) (any, error) { return map[string]any{"browser": "safari"}, nil }

type terminalInterceptor struct {
	base    *stubInterceptor
	fresh   bool
	mu      sync.Mutex
	servers map[int]*existingTerminalServer
	procs   map[int]*exec.Cmd
}

type existingTerminalServer struct {
	httpServer *http.Server
	active     bool
	port       int
}

func (t *terminalInterceptor) IsActivable() (bool, error) { return true, nil }
func (t *terminalInterceptor) IsActive(p int) (bool, error) {
	if t.fresh {
		return t.base.IsActive(p)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if srv, ok := t.servers[p]; ok {
		return srv.active, nil
	}
	return false, nil
}

func (t *terminalInterceptor) Activate(proxyPort int, _ map[string]any) (map[string]any, error) {
	env := buildTerminalEnv(t.base.cfg, proxyPort)
	if t.fresh {
		cmd, err := launchFreshTerminal(t.base.cfg, env, proxyPort)
		if err != nil {
			return nil, err
		}
		t.mu.Lock()
		if t.procs == nil {
			t.procs = make(map[int]*exec.Cmd)
		}
		t.procs[proxyPort] = cmd
		t.mu.Unlock()
		return map[string]any{"pid": cmd.Process.Pid, "proxyEnv": strings.Join(env, " ")}, nil
	}

	t.mu.Lock()
	if t.servers == nil {
		t.servers = make(map[int]*existingTerminalServer)
	}
	server, exists := t.servers[proxyPort]
	t.mu.Unlock()

	if exists {
		t.mu.Lock()
		server.active = false
		t.base.mu.Lock()
		t.base.active[proxyPort] = false
		t.base.mu.Unlock()
		port := server.port
		t.mu.Unlock()
		return map[string]any{
			"port":     port,
			"commands": getShellCommands(port),
		}, nil
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	serverPort := listener.Addr().(*net.TCPAddr).Port

	mux := http.NewServeMux()
	serverState := &existingTerminalServer{
		port: serverPort,
	}

	mux.HandleFunc("/setup", func(w http.ResponseWriter, r *http.Request) {
		envVars := getSetupEnvVars(t.base.cfg, proxyPort, "posix")
		callbackURL := fmt.Sprintf("http://localhost:%d/success", serverPort)
		script := getBashScript(callbackURL, envVars)
		w.Header().Set("Content-Type", "text/x-shellscript")
		_, _ = w.Write([]byte(script))
	})

	mux.HandleFunc("/gb-setup", func(w http.ResponseWriter, r *http.Request) {
		envVars := getSetupEnvVars(t.base.cfg, proxyPort, "posix")
		overridePath := filepath.Join(t.base.cfg.AssetsDir, "overrides")
		binPath := filepath.Join(overridePath, "path")
		posixBinPath := toPosixPath(binPath)
		envVars["PATH"] = posixBinPath + ":" + "$PATH"

		callbackURL := fmt.Sprintf("http://localhost:%d/success", serverPort)
		script := getBashScript(callbackURL, envVars)
		w.Header().Set("Content-Type", "text/x-shellscript")
		_, _ = w.Write([]byte(script))
	})

	mux.HandleFunc("/fish-setup", func(w http.ResponseWriter, r *http.Request) {
		envVars := getSetupEnvVars(t.base.cfg, proxyPort, "posix")
		callbackURL := fmt.Sprintf("http://localhost:%d/success", serverPort)
		script := getFishScript(callbackURL, envVars)
		w.Header().Set("Content-Type", "application/x-fish")
		_, _ = w.Write([]byte(script))
	})

	mux.HandleFunc("/ps-setup", func(w http.ResponseWriter, r *http.Request) {
		envVars := getSetupEnvVars(t.base.cfg, proxyPort, "powershell")
		callbackURL := fmt.Sprintf("http://localhost:%d/success", serverPort)
		script := getPowerShellScript(callbackURL, envVars)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(script))
	})

	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		t.mu.Lock()
		serverState.active = true
		t.base.mu.Lock()
		t.base.active[proxyPort] = true
		t.base.mu.Unlock()
		t.mu.Unlock()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	httpSrv := &http.Server{Handler: mux}
	serverState.httpServer = httpSrv

	go func() {
		_ = httpSrv.Serve(listener)
	}()

	t.mu.Lock()
	t.servers[proxyPort] = serverState
	t.mu.Unlock()

	return map[string]any{
		"port":     serverPort,
		"commands": getShellCommands(serverPort),
	}, nil
}

func (t *terminalInterceptor) Deactivate(proxyPort int, _ map[string]any) error {
	if t.fresh {
		t.mu.Lock()
		cmd, ok := t.procs[proxyPort]
		if ok {
			_ = cmd.Process.Kill()
			delete(t.procs, proxyPort)
		}
		t.mu.Unlock()
		return nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if proxyPort == 0 {
		for p, srv := range t.servers {
			if srv.httpServer != nil {
				_ = srv.httpServer.Close()
			}
			delete(t.servers, p)
		}
		return nil
	}
	if srv, ok := t.servers[proxyPort]; ok {
		if srv.httpServer != nil {
			_ = srv.httpServer.Close()
		}
		delete(t.servers, proxyPort)
	}
	return nil
}

func (t *terminalInterceptor) Metadata(string) (any, error) {
	return map[string]any{"fresh": t.fresh}, nil
}

func getSetupEnvVars(cfg *config.Config, proxyPort int, mode string) map[string]string {
	overridePath := filepath.Join(cfg.AssetsDir, "overrides")
	certPath := filepath.Join(cfg.ConfigDir, "ca.pem")
	proxyURL := fmt.Sprintf("http://127.0.0.1:%d", proxyPort)

	binPath := filepath.Join(overridePath, "path")
	rubyPath := filepath.Join(overridePath, "gems")
	pythonPath := filepath.Join(overridePath, "pythonpath")
	phpPath := filepath.Join(overridePath, "php")
	nodeScript := filepath.ToSlash(filepath.Join(overridePath, "js", "prepend-node.js"))
	javaAgent := javaAgentPath(cfg)

	var pathSep, pathInherit, rubyInherit, pythonInherit, nodeOptionsInherit, javaInherit, phpInherit string

	if mode == "powershell" {
		pathSep = ";"
		pathInherit = "$env:PATH"
		rubyInherit = "$env:RUBYLIB"
		pythonInherit = "$env:PYTHONPATH"
		nodeOptionsInherit = "$env:NODE_OPTIONS"
		javaInherit = "$env:JAVA_TOOL_OPTIONS"
		phpInherit = "$env:PHP_INI_SCAN_DIR"
	} else {
		pathSep = ":"
		if runtime.GOOS == "windows" {
			pathSep = ";"
		}
		pathInherit = "$PATH"
		rubyInherit = "$RUBYLIB"
		pythonInherit = "$PYTHONPATH"
		nodeOptionsInherit = "$NODE_OPTIONS"
		javaInherit = "$JAVA_TOOL_OPTIONS"
		phpInherit = "$PHP_INI_SCAN_DIR"
	}

	env := map[string]string{
		"HTTP_PROXY":                           proxyURL,
		"HTTPS_PROXY":                          proxyURL,
		"http_proxy":                           proxyURL,
		"https_proxy":                          proxyURL,
		"WS_PROXY":                             proxyURL,
		"WSS_PROXY":                            proxyURL,
		"GLOBAL_AGENT_HTTP_PROXY":              proxyURL,
		"CGI_HTTP_PROXY":                       proxyURL,
		"npm_config_proxy":                     proxyURL,
		"npm_config_https_proxy":               proxyURL,
		"npm_config_scripts_prepend_node_path": "false",
		"SSL_CERT_FILE":                        certPath,
		"NODE_EXTRA_CA_CERTS":                  certPath,
		"DENO_CERT":                            certPath,
		"PERL_LWP_SSL_CA_FILE":                 certPath,
		"GIT_SSL_CAINFO":                       certPath,
		"CARGO_HTTP_CAINFO":                    certPath,
		"CURL_CA_BUNDLE":                       certPath,
		"AWS_CA_BUNDLE":                        certPath,
		"HTTP_TOOLKIT_ACTIVE":                  "true",
		"HTTP_TOOLKIT_OVERRIDE_PATH":           overridePath,
		"DOCKER_BUILDKIT":                      "0",
	}

	env["PATH"] = binPath + pathSep + pathInherit
	env["RUBYLIB"] = rubyPath + ":" + rubyInherit
	env["PYTHONPATH"] = pythonPath + ":" + pythonInherit

	nodePrepend := "--require " + quoteIfNeeded(nodeScript)
	if mode == "powershell" {
		nodePrepend = "--require " + strings.ReplaceAll(nodeScript, `"`, "`\"")
	}
	env["NODE_OPTIONS"] = nodeOptionsInherit + " " + nodePrepend

	javaAgentOpt := fmt.Sprintf(`-javaagent:"%s"=127.0.0.1|%d|%s`, javaAgent, proxyPort, certPath)
	env["JAVA_TOOL_OPTIONS"] = javaInherit + " " + javaAgentOpt
	env["PHP_INI_SCAN_DIR"] = phpInherit + pathSep + phpPath

	if dockersvc.IsAvailable() {
		env["DOCKER_HOST"] = dockersvc.ProxyHostEnv(proxyPort)
	}

	return env
}

func getBashScript(callbackURL string, env map[string]string) string {
	var sb strings.Builder
	for k, v := range env {
		escaped := strings.ReplaceAll(v, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `"`, `\"`)
		sb.WriteString(fmt.Sprintf("    export %s=\"%s\"\n", k, escaped))
	}
	sb.WriteString(`
    if command -v winpty >/dev/null 2>&1; then
        alias php=php
        alias node=node
    fi

    if command -v curl >/dev/null 2>&1; then
        (curl --noproxy '*' -X POST "` + callbackURL + `" >/dev/null 2>&1 &) &> /dev/null
    fi

    echo 'HTTP Toolkit interception enabled'
`)
	return sb.String()
}

func getFishScript(callbackURL string, env map[string]string) string {
	var sb strings.Builder
	for k, v := range env {
		escaped := strings.ReplaceAll(v, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `"`, `\"`)
		sb.WriteString(fmt.Sprintf("    set -x %s \"%s\"\n", k, escaped))
	}
	sb.WriteString(`
    if command -v winpty >/dev/null 2>&1
        alias php=php
        alias node=node
    end

    if command -v curl >/dev/null 2>&1
        curl --noproxy '*' -X POST "` + callbackURL + `" >/dev/null 2>&1 &
    end

    echo 'HTTP Toolkit interception enabled'
`)
	return sb.String()
}

func getPowerShellScript(callbackURL string, env map[string]string) string {
	var sb strings.Builder
	sb.WriteString("    $HTTPTOOLKIT_envVars = Get-ChildItem Env:\n\n")
	for k, v := range env {
		escaped := strings.ReplaceAll(v, `"`, "`\"")
		sb.WriteString(fmt.Sprintf("    $Env:%s = \"%s\"\n", k, escaped))
	}
	sb.WriteString(`
    function Stop-Intercepting {
        $currentEnvVars = Get-ChildItem Env:
        foreach ($envVar in $currentEnvVars) {
            [System.Environment]::SetEnvironmentVariable($envVar.Name, $null)
        }
        foreach ($var in $HTTPTOOLKIT_envVars) {
            [System.Environment]::SetEnvironmentVariable($var.Name, $var.Value)
        }
        $PSDefaultParameterValues.Remove("invoke-webrequest:proxy")
        $PSDefaultParameterValues.Remove("invoke-webrequest:SkipCertificateCheck")
        Write-Host 'HTTP Toolkit interception disabled'
    }

    $PSDefaultParameterValues["invoke-webrequest:proxy"] = $Env:HTTP_PROXY
    $PSDefaultParameterValues["invoke-webrequest:SkipCertificateCheck"] = $True

    Start-Job -ScriptBlock { Invoke-WebRequest "` + callbackURL + `" -NoProxy -Method 'POST' } | out-null

    Write-Host "HTTP Toolkit interception enabled` + "`" + `nTo stop intercepting type " -NoNewline
    Write-Host "Stop-Intercepting" -ForegroundColor Red
`)
	return sb.String()
}

func getShellCommands(port int) map[string]any {
	if runtime.GOOS == "windows" {
		return map[string]any{
			"Powershell": map[string]any{
				"description": "Powershell",
				"command":     fmt.Sprintf("Invoke-Expression (Invoke-WebRequest http://localhost:%d/ps-setup).Content", port),
			},
			"Git Bash": map[string]any{
				"description": "Git Bash",
				"command":     fmt.Sprintf(`eval "$(curl -sS localhost:%d/gb-setup)"`, port),
			},
		}
	}
	return map[string]any{
		"Bash": map[string]any{
			"description": "Bash-compatible",
			"command":     fmt.Sprintf(`eval "$(curl -sS localhost:%d/setup)"`, port),
		},
		"Fish": map[string]any{
			"description": "Fish",
			"command":     fmt.Sprintf(`curl -sS localhost:%d/fish-setup | source`, port),
		},
		"Powershell": map[string]any{
			"description": "Powershell",
			"command":     fmt.Sprintf("Invoke-Expression (Invoke-WebRequest http://localhost:%d/ps-setup).Content", port),
		},
	}
}

func toPosixPath(winPath string) string {
	cleaned := filepath.ToSlash(winPath)
	if len(cleaned) >= 2 && cleaned[1] == ':' {
		drive := strings.ToLower(string(cleaned[0]))
		return "/" + drive + cleaned[2:]
	}
	return cleaned
}

func launchFreshTerminal(cfg *config.Config, env []string, proxyPort int) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		psSetup := fmt.Sprintf("Invoke-Expression (Invoke-WebRequest http://localhost:%d/ps-setup).Content", proxyPort)
		if wt, err := exec.LookPath("wt"); err == nil {
			// Windows Terminal: open a new tab with PowerShell
			cmd = exec.Command(wt, "new-tab", "--title", "HTTP Toolkit", "powershell.exe", "-NoExit", "-Command", psSetup)
		} else if pwsh, err := exec.LookPath("pwsh"); err == nil {
			// PowerShell 7
			cmd = exec.Command("cmd.exe", "/c", "start", pwsh, "-NoExit", "-Command", psSetup)
		} else if _, err := exec.LookPath("powershell"); err == nil {
			// Built-in Windows PowerShell
			cmd = exec.Command("cmd.exe", "/c", "start", "powershell.exe", "-NoExit", "-Command", psSetup)
		} else {
			gitBashPath := "C:\\Program Files\\Git\\git-bash.exe"
			if _, err := os.Stat(gitBashPath); err == nil {
				cmd = exec.Command(gitBashPath)
			} else if bp, err := exec.LookPath("git-bash"); err == nil {
				cmd = exec.Command(bp)
			} else {
				cmd = exec.Command("cmd.exe", "/c", "start", "cmd.exe", "/k",
					fmt.Sprintf("curl -sS localhost:%d/setup | sh", proxyPort))
			}
		}
	case "darwin":
		if _, err := os.Stat("/Applications/iTerm.app"); err == nil {
			cmd = exec.Command("osascript", "-e",
				fmt.Sprintf(`tell application "iTerm" to create window with default profile command "eval \"$(curl -sS localhost:%d/setup)\""`, proxyPort))
		} else {
			cmd = exec.Command("osascript", "-e",
				fmt.Sprintf(`tell application "Terminal" to do script "eval \"$(curl -sS localhost:%d/setup)\""`, proxyPort))
		}
	default:
		terminals := []string{"x-terminal-emulator", "gnome-terminal", "konsole", "xfce4-terminal", "kitty", "xterm"}
		found := "xterm"
		for _, t := range terminals {
			if _, err := exec.LookPath(t); err == nil {
				found = t
				break
			}
		}
		cmd = exec.Command(found, "-e", fmt.Sprintf(`sh -c 'eval "$(curl -sS localhost:%d/setup)"; exec sh'`, proxyPort))
	}

	cmd.Env = append(os.Environ(), env...)
	cmd.Dir = os.Getenv("HOME")
	if cmd.Dir == "" {
		cmd.Dir = os.Getenv("USERPROFILE")
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}
