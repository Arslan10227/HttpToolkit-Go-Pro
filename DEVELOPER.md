# Developer Guide

This document is for developers who want to understand, extend, or contribute to HttpToolkit Go Pro.

---

## Development Environment Setup

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| **Go** | 1.26+ | Build the Go backend |
| **Wails CLI** | v2 | Build the desktop app (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`) |
| **WebView2** | Latest | Windows desktop runtime (pre-installed on Windows 11) |
| **Node.js** | v20+ | Only needed if you want to rebuild the web UI assets |
| **Git** | 2.40+ | Version control |

### First-Time Setup

```bash
# Clone
git clone https://github.com/Arslan10227/HttpToolkit-Go-Pro.git
cd HttpToolkit-Go-Pro

# Create the secret defaults file from template
cp internal/config/env_defaults.go.example internal/config/env_defaults.go
# Edit env_defaults.go with your real OAuth credentials (or keep placeholders for dev)

# Verify the build compiles
go build ./...

# Run tests
go test ./...

# Run the standalone server
go run ./cmd/htk-server -v

# Run the Wails desktop app in dev mode (hot reload)
wails dev
```

### Verbose Logging

Pass `-v` or `--verbose` to enable verbose logging:

```bash
./htk-server -v
# or
wails dev -v
```

Log files are written to `logs/httptoolkit.log` (or the platform's config directory in production).

---

## Module Layout

```
cmd/                    Binary entry points
  htk-server/           Standalone MITM server (REST + admin + proxy)
  htk-mcp/              MCP stdio server for AI assistants
  htk-ctl/              UI operations control pipe client

internal/               Private packages (not importable externally)
  api/                  REST API (port 45457)
    gql/                GraphQL provider for admin queries
    server.go           Route registration & middleware
    handlers.go         HTTP handler functions
    settings.go         Settings REST endpoints
    mcp.go              MCP status/tools endpoints
    client_send.go      Request replay (HAR-shaped)
    breakpoint_test.go  Breakpoint API tests

  auth/                 Google OAuth flow
    google.go           OAuth token exchange & verification

  backup/               Cloud sync backends
    backup.go           Backup interface
    firebase.go         Firebase Firestore sync
    upstash.go          Upstash Redis sync

  cert/                 CA certificate management
    manager.go          CA generation, SPKI, PEM/P12 export
    system.go           System cert store install (Win/Mac/Linux)
    manager_test.go     Cert manager tests

  config/               Configuration
    config.go           Config struct, env loading, .env parser
    appmeta.go          App identity from version.yaml (go:embed)
    version.yaml        ← Single source of truth for name/version/title
    env_defaults.go     Secret defaults (gitignored)
    env_defaults.go.example  Template for env_defaults.go

  docker/               Docker interception
    session.go          Docker container attach/detach
    network.go          Network bridge setup
    commands.go         Docker CLI wrappers

  interceptors/         All interceptor implementations
    registry.go         Interceptor registry & lifecycle
    build.go            Interceptor factory (registers all IDs)
    stub.go             Stub interceptor (shared base)
    chromium.go         Chrome/Chromium/Edge/Brave/Opera launch
    firefox_nss.go      Firefox + NSS cert trust
    mobile.go           Android ADB / iOS Frida
    terminal_env.go     Terminal env-var injection
    system_proxy.go     OS system proxy
    platform.go         Platform-specific helpers
    webext.go           Browser extension interceptor
    amiusing.go         "amiusing" detection endpoint

  logger/               Structured logging (leveled, file + console)

  mcp/                  Model Context Protocol server
    status.go           MCP status endpoint
    stdio.go            MCP stdio JSON-RPC server
    ctl/                UI operations control pipe

  origins/              CORS origin allowlist

  proxy/
    admin/              Proxy admin API (port 45456)
      server.go         Admin route registration
      mockrtc.go        WebRTC rule management
      rtc_rules.go      RTC rule types
      logger.go         Request/event logging
    mitm/               MITM proxy engine
      server.go         Proxy server core
      transport.go      Upstream transport (HTTP/HTTPS)
      ws_proxy.go       WebSocket proxying
      amiusing.go       "amiusing" injection
      passthrough_events.go  TLS passthrough event handling

  rtc/                  Native WebRTC (Pion)
    manager.go          WebRTC session manager
    native/peer.go      Pion peer connection
    sdp_test.go         SDP handling tests

  server/               Server orchestration
    app.go              App struct — ties REST + admin + proxy together

  session/              Proxy session lifecycle

  settings/             Persistent settings store (JSON)

  snippets/             Code snippet generation (cURL, fetch, etc.)

  system/               OS-level operations
    Windows registry (deep links, HAR association)
    macOS / Linux helpers

  uibridge/             UI operation WebSocket bridge

  webextension/         Browser extension protocol support

assets/                 Embedded assets (go:embed)
  overrides/            Java agent, Frida scripts, web extension
  nss/                  NSS certutil binaries (Win/Mac/Linux)
  assets/               Pre-built web UI (HTML/JS/CSS)
  *.png                 Logos and icons
  index.html            Wails WebView entry point

contracts/              API contracts & schemas
  admin-api.md          Admin API documentation
  mockrtc-events.json   WebRTC event schema
  mockttp-events.json   HTTP event schema
  contract_test.go      Contract validation tests

frontend/wailsjs/       Wails-generated JS bindings for ShellApp

app.go                  Wails ShellApp (bound methods for frontend)
main.go                 Wails entry point (embeds assets, starts server)
wails.json              Wails build configuration
go.mod                  Go module: github.com/Arslan10227/HttpToolkit-Go-Pro
version.yaml            ← Edit this to change app name/version/title
```

---

## How-To Guides

### Add a New Interceptor

1. **Create the interceptor file** in `internal/interceptors/`:

```go
// internal/interceptors/my_interceptor.go
package interceptors

import (
    "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
    "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
)

type myInterceptor struct {
    *stubInterceptor
}

func newMyInterceptor(cfg *config.Config, spki string, certs *cert.Manager) *myInterceptor {
    return &myInterceptor{
        stubInterceptor: newStub(cfg, "my-interceptor", spki, certs),
    }
}

// Override methods as needed:
// - IsActivable() — check if the interceptor can run on this OS
// - Activate() — start intercepting (set proxy, launch process, etc.)
// - Deactivate() — stop intercepting and clean up
// - Metadata() — return interceptor-specific metadata
```

2. **Register it** in `internal/interceptors/build.go`:

```go
func buildAll(cfg *config.Config, spki string, certs *cert.Manager) []Interceptor {
    ids := []string{
        // ... existing IDs ...
        "my-interceptor",
    }
    // ...
}
```

3. **Write tests** in `internal/interceptors/my_interceptor_test.go`.

### Add a New REST Route

1. **Add the handler** in `internal/api/handlers.go` or a new file in `internal/api/`:

```go
func (s *Server) handleMyRoute(w http.ResponseWriter, r *http.Request) {
    // Check auth, method, etc.
    // Process request
    writeJSON(w, map[string]any{"ok": true})
}
```

2. **Register the route** in `internal/api/server.go`:

```go
mux.HandleFunc("/my-route", s.handleMyRoute)
```

3. **Add a contract test** in `contracts/contract_test.go` if the route is part of the public API.

### Add a New Proxy Rule Type

1. **Define the rule struct** in `internal/proxy/admin/` (e.g., `rules.go`).

2. **Implement matching logic** in `internal/proxy/mitm/server.go` or `transport.go`.

3. **Add the rule endpoint** in `internal/proxy/admin/server.go` if it needs a new HTTP route.

4. **Write tests** in `internal/proxy/mitm/` — follow the pattern of existing `*_test.go` files.

### Change App Name / Version / Title

Edit **one file**: `internal/config/version.yaml`

```yaml
name: "My Custom Name"
version: "2.0.0"
title: "My Custom Title"
```

The Go code embeds this file at build time via `//go:embed` in `internal/config/appmeta.go`. No other code changes needed. Rebuild to apply.

---

## Debugging

### Verbose Mode

```bash
./htk-server -v       # Enable verbose logging
./htk-server --verbose
```

### Log File Location

- **Dev:** `logs/httptoolkit.log` (in the working directory)
- **Production:** `<config_dir>/httptoolkit-pro-go/httptoolkit.log`
  - Windows: `%LOCALAPPDATA%\httptoolkit-pro-go\`
  - macOS/Linux: `~/.config/httptoolkit-pro-go/`

### Common Issues

| Issue | Fix |
|-------|-----|
| `config: env_defaults.go not found` | Run `cp internal/config/env_defaults.go.example internal/config/env_defaults.go` |
| `wails: command not found` | Run `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |
| WebView2 not found | Install [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) |
| Port already in use | Set `HTK_SERVER_PORT` and `HTK_ADMIN_PORT` env vars |
| Certificate trust errors | Use the "Install Certificate" button in Settings, or manually trust the CA |

---

## Testing

### Unit Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/proxy/mitm/...

# With coverage
go test -cover ./...

# With race detector
go test -race ./...
```

### Contract Tests

Contract tests validate that the Go API matches the expected JSON shapes from the Node backend:

```bash
go test ./contracts/
```

### Integration Tests

Integration tests start a real proxy server and test end-to-end flows:

```bash
go test -v ./internal/proxy/mitm/... -run Integration
```

---

## CI Pipeline

The CI workflow (`.github/workflows/go.yml`) runs on every push and PR:

### Jobs

| Job | OS | What it does |
|-----|----|-------------|
| `build` | Ubuntu, Windows, macOS | Builds `htk-server`, `htk-mcp`, `htk-ctl`; runs `go test ./...` |
| `build-wails-windows` | Windows | Builds the Wails desktop app (`HttpToolkit-Pro.exe`) |

### Secrets

CI generates `env_defaults.go` from GitHub secrets at build time. If secrets are missing (e.g., fork PRs), placeholder values are used and the build still compiles.

See [README.md](./README.md#required-github-secrets) for the list of required secrets.

### Artifacts

Built binaries are uploaded as GitHub Actions artifacts and can be downloaded from the Actions run page.

---

## Releasing

1. **Update `version.yaml`** with the new version number:

```yaml
version: "1.1.0"
```

2. **Commit and push** to `main`:

```bash
git add internal/config/version.yaml
git commit -m "Release v1.1.0"
git push origin main
```

3. **Tag the release**:

```bash
git tag v1.1.0
git push origin v1.1.0
```

4. **Create a GitHub Release** from the tag, attaching the CI-built artifacts.

5. CI will automatically build and upload binaries as artifacts. Download them from the Actions run and attach to the release.

---

## Code Style

- Follow standard `gofmt` / `goimports` formatting
- Use `golangci-lint` if available: `golangci-lint run ./...`
- Package names are lowercase, single-word (e.g., `config`, `cert`, `proxy`)
- Internal packages stay in `internal/` — no external imports of internal code
- Error handling: return errors, don't panic (except in `recover()` guards)
- Logging: use `logger.Info()` / `logger.Error()` with structured `map[string]any` context
- Tests: table-driven where possible, use `t.Run` for subtests
- No unused code — run `go mod tidy` before committing

### Import Order

```go
import (
    // Standard library
    "fmt"
    "os"

    // Third-party
    "github.com/wailsapp/wails/v2"

    // Local (this module)
    "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
)
```

---

## Architecture Notes

### Server Lifecycle

1. `main.go` loads config (`config.LoadDefault()`)
2. Extracts embedded overrides (java-agent.jar, frida scripts) to config dir
3. Starts Wails WebView2 shell
4. In `OnStartup`, launches the Go server in a goroutine (`server.Run()`)
5. `server.Run()` starts:
   - REST API on port 45457
   - Admin API on port 45456
   - MITM proxy on dynamic port 8000+ (when session starts)
6. On shutdown, gracefully stops all services

### Config Resolution Order

1. Real environment variables (highest priority)
2. `.env` file (next to executable or in working directory)
3. `env_defaults.go` compiled-in defaults (lowest priority)

### Embedded Assets

The Wails desktop app embeds all assets via `//go:embed all:assets`:
- Web UI (HTML/JS/CSS) — served by Wails' asset server
- `overrides/` — java-agent.jar, frida scripts, web extension
- `nss/` — NSS certutil binaries for Firefox cert trust
- Logos and icons

At runtime, `extractEmbeddedOverrides()` copies the overrides directory to the user's config directory so interceptors can access the files.

---

## License

AGPL-3.0 — see [version.yaml](./internal/config/version.yaml) and the main [README.md](./README.md#license).
