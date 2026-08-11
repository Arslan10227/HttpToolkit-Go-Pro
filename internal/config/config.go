package config

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	DefaultServerPort = 45457
	DefaultAdminPort  = 45456
	DefaultProxyPort  = 8000
	CACommonName      = "HttpToolkit Pro"
)

// ServerVersion returns the current server version from version.yaml.
// Use this instead of a hardcoded constant so version changes only need
// editing version.yaml.
func ServerVersion() string { return AppVersion() }

type Config struct {
	ServerPort        int
	AdminPort         int
	AuthToken         string
	ConfigDir         string
	AssetsDir         string
	DevMode           bool
	FirebaseProjectID string
	UpstashRedisURL   string
	GoogleClientID    string
	GoogleClientSecret string
}

func LoadDefault() (*Config, error) {
	loadDotEnv()
	token := os.Getenv("HTK_SERVER_TOKEN")
	if token == "" {
		token = randomToken()
	}
	configDir, err := defaultConfigDir()
	if err != nil {
		return nil, err
	}
	assetsDir := os.Getenv("HTK_ASSETS_DIR")
	if assetsDir == "" {
		assetsDir = findAssetsDir()
	}
	dev := os.Getenv("HTK_DEV") != "0"
	firebaseProject := os.Getenv("FIREBASE_PROJECT_ID")
	if firebaseProject == "" {
		firebaseProject = "httptoolkitpro"
	}
	return &Config{
		ServerPort:         envInt("HTK_SERVER_PORT", DefaultServerPort),
		AdminPort:          envInt("HTK_ADMIN_PORT", DefaultAdminPort),
		AuthToken:          token,
		ConfigDir:          configDir,
		AssetsDir:          assetsDir,
		DevMode:            dev,
		FirebaseProjectID:  firebaseProject,
		UpstashRedisURL:    envOrString("UPSTASH_REDIS_URL", defaultUpstashRedisURL),
		GoogleClientID:     envOrString("GOOGLE_CLIENT_ID", defaultGoogleClientID),
		GoogleClientSecret: envOrString("GOOGLE_CLIENT_SECRET", defaultGoogleClientSecret),
	}, nil
}

func LoadFromFlags() (*Config, error) {
	loadDotEnv()
	serverPort := flag.Int("server-port", envInt("HTK_SERVER_PORT", DefaultServerPort), "REST API port")
	adminPort := flag.Int("admin-port", envInt("HTK_ADMIN_PORT", DefaultAdminPort), "Proxy admin port")
	token := flag.String("token", os.Getenv("HTK_SERVER_TOKEN"), "Bearer auth token")
	dev := flag.Bool("dev", os.Getenv("HTK_DEV") == "1", "Development mode (relaxed CORS)")
	flag.Parse()

	if *token == "" {
		*token = randomToken()
	}

	configDir, err := defaultConfigDir()
	if err != nil {
		return nil, err
	}

	assetsDir := os.Getenv("HTK_ASSETS_DIR")
	if assetsDir == "" {
		assetsDir = findAssetsDir()
	}

	firebaseProject := os.Getenv("FIREBASE_PROJECT_ID")
	if firebaseProject == "" {
		firebaseProject = "httptoolkitpro"
	}
	return &Config{
		ServerPort:         *serverPort,
		AdminPort:          *adminPort,
		AuthToken:          *token,
		ConfigDir:          configDir,
		AssetsDir:          assetsDir,
		DevMode:            os.Getenv("HTK_DEV") != "0" && (*dev || os.Getenv("HTK_DEV") == "1"),
		FirebaseProjectID:  firebaseProject,
		UpstashRedisURL:    envOrString("UPSTASH_REDIS_URL", defaultUpstashRedisURL),
		GoogleClientID:     envOrString("GOOGLE_CLIENT_ID", defaultGoogleClientID),
		GoogleClientSecret: envOrString("GOOGLE_CLIENT_SECRET", defaultGoogleClientSecret),
	}, nil
}

func defaultConfigDir() (string, error) {
	if d := os.Getenv("HTK_CONFIG_DIR"); d != "" {
		return d, nil
	}
	var base string
	switch runtime.GOOS {
	case "windows":
		base = os.Getenv("LOCALAPPDATA")
		if base == "" {
			base = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
	default:
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			base = xdg
		} else {
			home, _ := os.UserHomeDir()
			base = filepath.Join(home, ".config")
		}
	}
	dir := filepath.Join(base, "httptoolkit-pro-go")
	return dir, os.MkdirAll(dir, 0o700)
}

func findAssetsDir() string {
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		// Check for standalone/Wails distribution resource directory:
		candidate := filepath.Join(dir, "resources", "httptoolkit-server")
		if st, err := os.Stat(filepath.Join(candidate, "overrides")); err == nil && st.IsDir() {
			return candidate
		}
		// Standard assets folder
		candidate = filepath.Join(dir, "assets")
		if st, err := os.Stat(candidate); err == nil && st.IsDir() {
			return candidate
		}
	}
	// Dev: repo httptoolkit-go/assets
	if wd, err := os.Getwd(); err == nil {
		for _, rel := range []string{"assets", "httptoolkit-go/assets", "../httptoolkit-go/assets"} {
			p := filepath.Join(wd, rel)
			if st, err := os.Stat(p); err == nil && st.IsDir() {
				return p
			}
		}
	}
	return "assets"
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		var n int
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return n
		}
	}
	return def
}

func envOrString(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func randomToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// loadDotEnv reads a .env file next to the executable and in the working
// directory. It only sets values that are not already present in the process
// environment, so real env vars always win.
func loadDotEnv() {
	var candidates []string
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), ".env"))
	}
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(wd, ".env"))
	}

	for _, path := range candidates {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		defer f.Close()

		scan := bufio.NewScanner(f)
		for scan.Scan() {
			line := strings.TrimSpace(scan.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			key, value, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			// Remove matching quotes if present.
			if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
				value = value[1 : len(value)-1]
			}
			if os.Getenv(key) == "" {
				_ = os.Setenv(key, value)
			}
		}
	}
}
