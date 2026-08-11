# httptoolkit-go — Node → Go migration status

This document tracks progress porting **httptoolkit-node** (Node/Mockttp PluggableAdmin) to **httptoolkit-go** (Go MITM proxy + Wails desktop shell). The hosted UI at `https://httptoolkitpro.vercel.app` talks to the local backend via REST and a Go-native proxy admin client.

**Reference:** Node server lives in `httptoolkit-node/`. UI admin client: `webui/src/services/proxyAdminClient.ts`.

**Last verified:** `go build ./...` and `go test ./...` pass in `httptoolkit-go/` (full suite, ~45s — previously `go test ./...` would hang for 10+ minutes and fail on a Windows machine without `frida`/full NSS assets; see "Bugs found and fixed" below).

---

## Architecture

| Component | Default port | Role |
|-----------|--------------|------|
| REST API | **45457** | Interceptors, certs, snippets, MCP status, `/ui-operations` WebSocket bridge |
| Proxy admin | **45456** | Session start/stop, HTTP/WS/RTC rules, `/events` WebSocket (replaces Mockttp PluggableAdmin) |
| MITM proxy | **8000+** (dynamic) | HTTP/HTTPS/WebSocket interception, SOCKS5 (optional) |
| MockRTC | — | Native Go (Pion) in-process WebRTC mocking — no Node sidecar |
| Wails desktop | — | Embeds Go server; webview loads hosted UI |

```
┌─────────────────┐     REST :45457      ┌──────────────────┐
│  Vercel UI      │◄────────────────────►│  htk-go REST     │
│  (React)        │     Admin :45456     │  + interceptors  │
└────────┬────────┘◄────────────────────►└────────┬─────────┘
         │ proxyAdminClient.ts                     │
         │ WS /events                              │ MITM
         └────────────────────────────────────────►│ :8000+
                                                   │
         ┌─────────────────────────────────────────┤
         │ Native MockRTC (Pion, in-process)       │
         │ internal/rtc/native/                    │
         └─────────────────────────────────────────┘
```

### Binaries

| Command | Path | Purpose |
|---------|------|---------|
| `htk-server` | `cmd/htk-server/` | Standalone Go backend |
| `htk-mcp` | `cmd/htk-mcp/` | MCP stdio server (`htk-mcp stdio`) |
| `htk-ctl` | `cmd/htk-ctl/` | Control pipe client for UI operations |
| Wails app | `main.go`, `app.go` | Desktop shell |

---

## Migration approach

1. **Replace Mockttp/PluggableAdmin** with a Go MITM server and a simplified admin HTTP/WS API documented in `contracts/admin-api.md`.
2. **Keep UI JSON contracts** — rules, events, and interceptor REST shapes match what the React app already sends.
3. **Native WebRTC (Pion)** — in-process Go WebRTC mocking (`internal/rtc/native/`) replaces the former Node `mockrtc` sidecar. No Node subprocess is spawned.
4. **Bundle assets** — overrides, NSS certutil, webextension build, Frida scripts copied via `npm run copy:go-assets`.

---

## Completed work

### Core server & config

- [x] Config loading (`HTK_SERVER_PORT`, `HTK_ADMIN_PORT`, `HTK_SERVER_TOKEN`, `HTK_ASSETS_DIR`, `HTK_CONFIG_DIR`, `HTK_DEV`)
- [x] CA generation, SPKI fingerprint, PEM/P12 export
- [x] System cert install (Windows registry, macOS `security`, Linux trust stores)
- [x] Java cert detect/install (`keytool`)
- [x] Session manager for `/config` (proxy port, DNS, rule parameter keys, HTTP/2, WebRTC, SOCKS, docker tunnel)
- [x] CORS + private-network headers for browser UI
- [x] Bearer auth on REST (auto-generated token)
- [x] Graceful shutdown on SIGINT/SIGTERM and `POST /shutdown`
- [x] Wails integration (`internal/server/app.go`)

### REST API (port 45457)

Routes match `httptoolkit-node/src/api/rest-api.ts` unless noted.

| Route | Status | Notes |
|-------|--------|-------|
| `GET /health`, `/ready` | ✅ | |
| `GET /version` | ✅ | Returns `1.0.0-go` |
| `GET /config` | ✅ | Active session metadata |
| `GET /config/network-interfaces` | ✅ | |
| `GET/POST /auth/client-token` | ✅ | |
| `GET/POST /auth/desktop-session` | ✅ | Desktop OAuth token handoff |
| `GET /certificate/*` | ✅ | status, export (pem/crt/cer/p12), install |
| `GET/POST /java/*` | ✅ | versions, cert status/install |
| `GET /interceptors` | ✅ | All registered interceptors |
| `GET /interceptors/:id/metadata` | ✅ | |
| `GET /interceptors/:id/metadata/:subId` | ✅ | devices, containers, JVM targets, Frida targets |
| `POST /interceptors/:id/activate/:port` | ✅ | Per-interceptor impl |
| `POST /interceptors/:id/deactivate/:port` | ✅ | |
| `POST /client/send` | ✅ | HAR-shaped request replay |
| `POST /snippets/generate` | ✅ | Broad target/client matrix (see Snippets) |
| `GET /mcp/status`, `/mcp/tools` | ✅ | |
| `POST /shutdown` | ✅ | |
| `POST /update` | ✅ stub-by-design | Returns `{ success: true }`, no auto-update — **matches Node** (`apiModel.updateServer()` is a no-op there too; no updater exists in either backend) |
| `POST /rules/bulk-create` | ✅ 501-by-design | **Corrected:** Node's `src/api/rest-api.ts` also returns the exact same 501 + message for this route — this is intentional 1:1 parity (UI never calls it; uses admin rule APIs), not a Go gap |
| `POST /webhooks/capture` | ✅ 501-by-design | **Corrected:** Node also stubs this route with the same 501 message — intentional parity, not a gap |
| `GET /` (GraphQL) | ⚠️ removed | Returns 405; UI uses admin client instead |
| `WS /ui-operations` | ✅ | UI operation bridge + MCP ctl pipe |

### Proxy admin API (port 45456)

Documented in `contracts/admin-api.md`. Consumed by `proxyAdminClient.ts`.

| Route | Status | Notes |
|-------|--------|-------|
| `GET /metadata` | ✅ | `{ running, httpPort, webrtcEnabled }` |
| `POST /session/start` | ✅ | HTTP/2, SOCKS, port range, DNS, docker proxy, native RTC |
| `POST /session/stop` | ✅ | Cleans proxy, native RTC, webextension config |
| `PUT /rules/http` | ✅ | |
| `PUT /rules/ws` | ✅ | |
| `PUT /rules/rtc` | ✅ | Forwards to native RTC rule engine |
| `GET/POST /session/:id` | ✅ | Session info; native GraphQL handler; JSON event relay |
| `WS /events` | ✅ | HTTP + RTC events on one socket |

### MITM proxy (`internal/proxy/mitm/`)

- [x] HTTP/HTTPS CONNECT MITM with dynamic CA
- [x] HTTP/2 upstream (optional, session flag)
- [x] SOCKS5 proxy (optional, default on at session start)
- [x] TLS passthrough + raw passthrough events
- [x] Request/response event stream (Mockttp-compatible names)
- [x] WebSocket upgrade, passthrough, echo (`ws_proxy.go`)
- [x] Rule engine integration (match → act)
- [x] SDP sniffing for WebRTC offer detection (`maybeNotifySDP`)
- [x] Passthrough options parsing (ignore cert errors, lookup servers, client certs — partial application)
- [x] Docker tunnel SOCKS routing via session rule parameters

**HTTP events emitted:** `request-initiated`, `request`, `response`, `abort`, `tls-client-error`, `tls-passthrough-opened/closed`, `raw-passthrough-opened/closed`, `client-error`, `rule-event`, all `websocket-*` events (see `contracts/mockttp-events.json`).

### Rule engine (`internal/proxy/rules/`)

**Matchers implemented:**

> **Corrected 2026-08 audit:** cross-checked against `webui/src/utils/mockttpRuleBuilder.ts`
> (`mapUiMatcher()`), which is the only place the UI ever constructs matcher JSON. The UI
> only ever emits `wildcard`, `method`, `flexible-path`, `regex-path`, `regex-url`, `host`,
> `hostname`, `header`, `query`, `raw-body-includes`, `json-body` — **all already implemented**
> in `internal/proxy/rules/engine.go`, with several extras (`simple-path`, `exact-query`,
> `raw-body`, `json-body-flexible`) implemented beyond what the UI currently needs.
> **Matcher parity with the UI is 100%.** The "other mockttp matchers" row below refers only
> to matcher types that exist in the upstream `mockttp` library's TypeScript API surface but
> that this fork's UI never sends — not an active gap.

| Matcher | Status |
|---------|--------|
| `wildcard` | ✅ |
| `method` | ✅ |
| `host`, `hostname` | ✅ |
| `protocol` | ✅ |
| `simple-path`, `path`, `flexible-path` | ✅ |
| `regex-path`, `regex-url` | ✅ |
| `header` | ✅ |
| `query`, `exact-query` | ✅ |
| `raw-body`, `raw-body-includes` | ✅ |
| `json-body`, `json-body-flexible`, `json-body-matching` | ✅ (exact JSON equality) |
| Other mockttp matchers never emitted by this fork's UI (`cookie`, `port`, etc.) | ❌ logged as unknown, no match — not reachable from the UI today |

**Steps implemented:**

> **Corrected 2026-08 audit:** the table below previously listed `redirect`, `abort`,
> `forward-to-host`/transforms, and `req-res-transformer` as unported. That was **stale** —
> all four are implemented in `internal/proxy/mitm/server.go` (`case "redirect"`,
> `case "abort"`, `applyRequestTransforms`/`applyResponseTransforms`) and
> `internal/proxy/rules/passthrough.go` (`parseTransformRequest`/`parseTransformResponse`).
> Cross-checked against `webui/src/utils/ruleStepMapper.ts` (the only place the UI ever
> constructs step JSON) — every step type the UI can emit maps 1:1 to a handled case in
> the Go engine. **Rule step parity with the UI is effectively 100%** as of this audit;
> no further rule-step work is required unless the UI's step vocabulary changes.

| Step | Status | Notes |
|------|--------|-------|
| `simple` (fixed response) | ✅ | |
| `file` | ✅ | |
| `close-connection` | ✅ | |
| `reset-connection` | ✅ | |
| `timeout` | ✅ | |
| `delay` | ✅ | Accumulates before passthrough |
| `passthrough`, `forward`, `wait-for-body` | ✅ | |
| `callback` | ⚠️ | Emits `rule-event`; no live pause/resume. **Note:** the UI (`ruleStepMapper.ts`) never actually emits a `callback` step type today — only `breakpoint`, which intentionally passthroughs — so this has no observable UI impact currently. |
| `webhook` | ✅ | POST on match (`webhook.go`) |
| `stream` | ✅ | Chunked file response |
| `ws-passthrough`, `ws-echo` | ✅ | |
| `redirect` | ✅ | Corrected: implemented via `writeFixedResponse` + `Location` header in `mitm/server.go` |
| `abort` | ✅ | Corrected: implemented (`case "abort"` resets the connection and emits an `abort` event) |
| `forward-to-host` / `transformRequest` | ✅ | Corrected: UI maps to `passthrough` + `transformRequest.replaceHost`; Go applies it in `applyRequestTransforms` |
| `req-res-transformer` | ✅ | Corrected: UI maps to `passthrough` + `transformRequest`/`transformResponse`; Go applies both in `server.go` |
| `breakpoint` | ⚠️ | UI intentionally maps to passthrough (no live callback bridge yet) — this is a genuine gap, not a doc error |

### Docker interception (`internal/docker/`)

- [x] Docker API proxy (Unix socket / Windows named pipe)
- [x] Container create/start transform (inject env, labels, volumes)
- [x] Compose label transform
- [x] Dockerfile injection (`dockerfile.go`)
- [x] Build context tar injection (`build_inject.go`)
- [x] Build output step renumbering (`build_output.go`)
- [x] Network alias monitor for `*.httptoolkit.localhost`
- [x] Local SOCKS tunnel for container egress (`tunnel.go`)
- [x] `docker-attach` interceptor (restart + inject)
- [x] `DOCKER_HOST` override in terminal env
- [x] Session rule parameter `docker-tunnel-proxy-{port}`

**Not ported:** upstream `ghcr.io/httptoolkit/docker-socks-tunnel` container; Go uses in-process SOCKS instead.

### Interceptors (`internal/interceptors/`)

All IDs registered (parity with Node `buildInterceptors`). Implementation depth varies.

| Interceptor ID | Status | Implementation |
|----------------|--------|----------------|
| `fresh-chrome`, `fresh-chromium`, `fresh-edge`, `fresh-brave`, `fresh-opera` (+ beta/dev/canary variants) | ✅ | Launch with SPKI pin, temp profile, webextension |
| `existing-chrome`, `existing-chromium`, `existing-arc` | ✅ | Attach to running browser |
| `fresh-firefox`, `fresh-firefox-dev`, `fresh-firefox-nightly` | ✅ | NSS certutil, profile, `user.js` proxy prefs, cert-check page |
| `fresh-safari` | ⚠️ | macOS only; basic `networksetup` + `open Safari` |
| `system-proxy` | ✅ | Windows / macOS / Linux (gsettings) set + clear |
| `fresh-terminal`, `existing-terminal` | ✅ | Env overrides; fresh spawns shell |
| `attach-jvm` | ⚠️ | `jattach` + bundled `java-agent.jar`; basic `jps` metadata |
| `android-adb` | ⚠️ | Proxy via `settings put global http_proxy`; basic cert push |
| `android-frida`, `ios-frida` | ⚠️ | Spawn `frida` with bundled scripts; simplified vs Node |
| `electron` | ⚠️ | `--require` prepend script + proxy env |
| `docker-attach` | ✅ | Full inject/restart flow |

**Chromium webextension** (`internal/webextension/`):

- [x] Copy build to temp dir
- [x] Per-proxy config `127_0_0_1.{port}` with MockRTC admin URL
- [x] Assets from upstream v1.2.0 via `scripts/copy-go-assets.mjs`

**Firefox NSS** (`firefox_nss.go`):

- [x] Bundled `certutil` per OS in `assets/nss/{win32,linux,darwin}`
- [x] Profile under `{configDir}/firefox-profile`
- [x] Cert-check auxiliary HTTPS server (`certcheck.go`)

**Terminal env** (`terminal_env.go`):

- [x] Node parity for proxy vars, `NODE_OPTIONS`, `JAVA_TOOL_OPTIONS`, `SSL_CERT_FILE`, `PHP_INI_SCAN_DIR`, PATH overrides, `DOCKER_HOST`
- [x] Inherit filtering for desktop-specific vars

### WebRTC / MockRTC (`internal/rtc/`)

- [x] Event normalization + validation (`events.go`, `contracts/mockrtc-events.json`)
- [x] SDP pattern detection fix (`sdp_detect.go`)
- [x] MITM hooks: SDP in request/response bodies → `peer-connected`
- [x] Admin JSON POST relay on `/session/:id`
- [x] **Native Go WebRTC (Pion)** (`internal/rtc/native/`)
  - In-process `pion/webrtc/v4` PeerConnections — no Node subprocess
  - GraphQL handler for webextension's MockRTC protocol
  - `PUT /rules/rtc` → native rule engine (echo/send/close steps)
  - `stopSession` closes all native peer connections
- [x] Webextension `adminBaseUrl` points at Go admin server

**RTC events (when native RTC + extension active):** `peer-connected`, `peer-disconnected`, `external-peer-attached`, data-channel and media-track events per `mockrtc-events.json`.

### MCP

- [x] REST `/mcp/status`, `/mcp/tools`
- [x] Stdio JSON-RPC (`cmd/htk-mcp`, `internal/mcp/stdio.go`): `initialize`, `tools/list`, `tools/call`, `notifications/tools/list_changed`
- [x] Control pipe server (`internal/mcp/ctl/`) bridged to `/ui-operations`

### Snippets (`internal/snippets/generate.go`)

Targets: shell, http, javascript/node, python, go, java, php, ruby, csharp, swift, kotlin, c, objc, r, clojure, ocaml, crystal, rust, powershell.

Clients include: curl, httpie, wget, axios, fetch, jQuery, okhttp, requests, RestSharp, Faraday, net/http, libcurl, and others.

### Certificates & DNS

- [x] Internal DNS server for `*.httptoolkit.localhost` (`internal/proxy/dns/`)
- [x] SPKI pinning for Chromium interceptors

### Contracts & tests

- [x] `contracts/mockttp-events.json`, `mockrtc-events.json`
- [x] `contracts/admin-api.md`
- [x] Contract tests (`go test ./contracts/...`)
- [x] Unit tests across docker, rules, mitm, interceptors, snippets, rtc, webextension, session

---

## Additional verification pass (2026-08)

- **Passthrough options full-application** — added dedicated tests
  (`internal/proxy/mitm/passthrough_test.go`) for the previously-untested paths: DNS
  `lookupServers` (custom `DialContext` installed only when set),
  `proxyConfig`/`noProxy` (string form, `[]any` form, and the `<local>` loopback token), and
  client-certificate host-map selection. All confirmed correctly applied — no code changes
  needed, just closed a test-coverage gap.
- **End-to-end rule-engine verification** — added `TestMITMFixedResponseRule`
  (`internal/proxy/mitm/server_test.go`): starts a real MITM proxy, sets a `simple`
  fixed-response HTTP rule via the engine, sends a real HTTP request through the live proxy,
  and asserts the mocked response (not the real upstream's) is what comes back — closing the
  gap where only default passthrough was previously covered end-to-end.
- **Firefox cert-check wait** — re-verified against code: `firefoxFresh.Activate`
  (`internal/interceptors/platform.go`) already blocks on `checkSrv.WaitSuccess(25s)` and
  fails activation on error/timeout. `port.md`'s P2 backlog previously listed this as
  outstanding; **corrected** — see the P2 table above.
- **Android ADB reverse tunnel** — re-verified: `internal/interceptors/mobile.go` already
  implements `adb reverse tcp:{local} tcp:{remote}` with periodic re-establishment and
  cleanup on deactivate. **Corrected** in the P2 table above. QR-connect flow and full cert
  install parity with `adb-commands.ts` were **not** re-verified in this pass and should
  still be treated as open until checked.

## Modernization pass (2026-08)

- **REST API + admin routing** (`internal/api/server.go`, `internal/proxy/admin/server.go`):
  migrated every single-method route to Go 1.22+ `net/http.ServeMux` method-qualified
  patterns (`"POST /shutdown"`, `"PUT /rules/http"`, etc.), removing ~10 hand-rolled
  `if r.Method != http.MethodX { ... }` checks. `internal/proxy/admin` (no catch-all `/`
  pattern) now gets real automatic `405 Method Not Allowed` responses from the stdlib mux
  for free (verified: `PUT`-only `/rules/http` correctly 405s a `GET`). `internal/api`
  still has a catch-all `/` handler (`handleRoot`) which masks that specific stdlib
  behavior per a known Go limitation ([golang/go#65648](https://github.com/golang/go/issues/65648)) —
  documented in a code comment; not a functional risk since no client sends the wrong
  method to these routes. `handleDesktopSession` (previously a single handler switching on
  `r.Method`) was split into `handleDesktopSessionGet`/`handleDesktopSessionPost` registered
  as separate method-qualified patterns.
- **Logging** (`internal/logger/logger.go`): reimplemented on stdlib `log/slog` (a
  `slog.JSONHandler` + `slog.LevelVar` for the existing `SetVerbose` toggle), replacing the
  hand-rolled JSON-marshal-and-write logger. The public API (`Info`/`Debug`/`Error`/
  `SetVerbose`/`GetLogPath`) is unchanged, so none of the ~20 call sites across the codebase
  needed to change. Per-entry `fsync` durability (useful for post-crash log diagnostics on
  a desktop app) is preserved via a small wrapping `io.Writer`.
- **Not done / deliberately deferred:**
  - `graphql-go/graphql` dependency: **kept** (not removed as originally considered) — it
    turned out to be required for the native WebRTC GraphQL schema (see above), not just
    the legacy REST `/graphql` mirror.
  - Wails v2 → v3: **not attempted** — v3 is still alpha/beta; v2 remains the correct choice
    until v3 stabilizes.
  - `coder/websocket` vs `gorilla/websocket`: no change — `gorilla/websocket` is only an
    indirect (Wails-internal) dependency, not a duplicate direct usage.

## Bugs found and fixed (2026-08 audit)

While re-verifying this document against actual code and running `go test ./...` to
completion (previously never run to completion — it was hanging), two real production
bugs were found and fixed:

1. **Interceptor probe commands could hang forever.** `internal/interceptors/mobile.go`'s
   `execCommand` calls for `frida-ps`, `adb`, `jps`, `powershell` process probes had no
   timeout. `Registry.List()` (backing `GET /interceptors`, called on every UI refresh of
   the interceptors panel) fans out `Metadata()` calls across a `sync.WaitGroup` — a single
   hung probe (e.g. `frida-ps` waiting on a USB device enumeration that never resolves)
   blocked the *entire* interceptors list response indefinitely, freezing that part of the
   UI. Fixed by adding `probeTimeout` (5s) bounded execution helpers
   (`execRunTimeout`/`execOutputTimeout`/`execCombinedOutputTimeout`) and rewiring every
   probe call site to use them. `go test ./internal/interceptors/...` went from a 10-minute
   timeout/failure to ~5-10s.
2. **Bundled Firefox NSS `certutil.exe` (and `modutil.exe`, `pk12util.exe`, `shlibsign.exe`,
   `signmar.exe`) were missing from `assets/nss/win32/`** (and from
   `httptoolkit-node/nss/win32/` too), because the repo's root `.gitignore` had a blanket
   `*.exe` rule that silently stripped them from version control. This breaks Firefox
   certificate interception out of the box on a fresh Windows clone/checkout for **both**
   backends. Fixed by restoring the binaries from `Original/httptoolkit-server-main/nss/win32/`
   and adding `.gitignore` negation rules (`!httptoolkit-go/assets/nss/**/*.exe` etc.) so this
   can't silently regress again.

## Remaining work

Prioritized backlog for full Node parity.

### P1 — High impact for daily use

> **2026-09 update:** the rule-step items were already implemented, and a self-hosted
> `amiusingServer` has replaced the external `https://amiusing.httptoolkit.tech`
> dependency. The new server is used by both Chromium and Firefox interceptors. It
> also disables QUIC (`--disable-quic`) so Chrome cannot bypass the HTTPS proxy.

| Item | Description | Reference (Node) |
|------|-------------|------------------|
| **Native Go WebRTC (replace Node sidecar)** | ✅ **Complete** — Node sidecar deleted; native Pion handler is the only backend. See "Native WebRTC migration" status below. | `mockrtc` PluggableAdmin event stream |
| **Fresh-Chrome TLS interception** | ✅ **Fixed** — self-hosted `amiusingServer` replaces the external `amiusing.httptoolkit.tech` page and the broken TLS cert-check server. QUIC is disabled so Chrome cannot bypass the HTTPS proxy, and `--ignore-certificate-errors-spki-list` handles trust. | `hide-warning-server.ts` |
| ~~Rule step: redirect~~ | ✅ Already implemented — see corrected table above | |
| ~~Rule step: abort~~ | ✅ Already implemented — see corrected table above | |
| ~~Request/response transforms~~ | ✅ Already implemented — see corrected table above | |
| **Passthrough options (full)** | `ignoreHostHttpsErrors`, `lookupServers`, `clientCertificateHostMap` are parsed **and applied** in `transport.go`'s `buildPassthroughTransport`. `proxyConfig` parsing exists (`resolveProxyURL`) — verify it covers every shape the UI can send (string, array, callback-shaped) with a dedicated test matrix. | `passthrough.go` + `transport.go` |

### P1b — Close confirmation, settings, and cloud hooks

> **2026-09 update:** app-level settings are now persisted on disk; a native Wails close-confirmation dialog is wired to `confirmBeforeClose`; and local backup/notification/cloud REST hooks are available for the UI.

| Item | Description |
|------|-------------|
| **Close confirmation** | ✅ **Wails `OnBeforeClose` shows a native "Are you sure?" dialog when `confirmBeforeClose` is enabled.** Toggle is available through `POST /settings` and persisted in `settings.json`. |
| **Settings persistence** | ✅ **`/settings` GET/POST stores app-level settings on disk (`ConfigDir/settings.json`).** Cloud account settings remain in the UI/Supabase layer. |
| **Backup/restore hooks** | ✅ **`/backup/export` and `/backup/import` endpoints** export/import settings and are ready to be extended for rules, snippets, and session state. |
| **Notifications** | ✅ **`/notifications/register` and `/notifications/status` endpoints** accept and report push/notification preferences. |
| **Cloud sync hooks** | ✅ **`/cloud/status` and `/cloud/sync` endpoints** report and trigger sync status for cloud account integration. |

### P2 — Interceptor & platform depth

| Item | Description |
|------|-------------|
| **Breakpoint / callback pause** | Live pause/resume bridge via UI operations; today breakpoints map to passthrough. |
| **Terminal (PowerShell)** | Dedicated PowerShell startup path; inherit-mode parity with `existing-terminal-interceptor.ts`. |
| **Terminal startup scripts** | Full `fresh-terminal-interceptor.ts` platform scripts (macOS Terminal.app, Windows Terminal, etc.). |
| **Firefox** | ~~Wait for cert-check success before reporting activate~~ ✅ **already implemented** — now uses the same self-hosted `amiusingServer` as Chromium. Remaining: snap/flatpak profile paths; richer `prefs.js` from Node. |
| **Android ADB** | ~~adb reverse tunnel~~ ✅ **already implemented** (`mobile.go`: periodic `adb reverse tcp:{local} tcp:{remote}` re-establishment + cleanup on deactivate). ~~Cert status/install REST endpoints~~ ✅ **implemented** — `GET /interceptors/android-adb/certificate/status` and `POST /interceptors/android-adb/certificate/install` now match the Node REST API, backed by `getDeviceCertificateStatus`/`installDeviceCertificate` methods on the `androidAdb` interceptor. Remaining: QR connect flow. |
| **Frida Android/iOS** | Structured target list, host IDs, app spawn/kill lifecycle vs raw `frida-ps` output. |
| **JVM attach** | Fallback attach APIs, deactivate detach, `interceptedByProxy` tracking. |
| **Safari** | Full fresh-safari flow (cert trust, automation). |
| **Electron** | Preload injection parity, app picker metadata. |
| **OAuth callback page** | ✅ **Implemented** — `GET /auth/callback` now serves the same HTML page as the Node backend, reading OAuth tokens from the URL fragment and POSTing them to `/auth/desktop-session`. |

### Native WebRTC migration — status: **implemented (v1), behind a feature flag**

> **2026-08 implementation update:** contrary to the original plan below (which assumed a
> simple event-relay contract), inspecting the actual bundled webextension
> (`assets/overrides/webextension/build/background.js`) showed the browser-side MockRTC
> client speaks a **GraphQL** protocol against `/session/:id`
> (`createOffer`/`createExternalOffer`/`answerOffer`/`answerExternalOffer`/`completeOffer`
> mutations + `getSeenMessages` query — see `MockRTCRemotePeer`/`RemoteSessionApi` in that
> bundle). The implementation below matches that real wire protocol, not a simplified one.

**What's implemented** (`internal/rtc/native/`):
- `peer.go` — `Handler` managing real `pion/webrtc/v4` `PeerConnection`s per mock session:
  `AnswerOffer` (browser-initiated), `CreateOffer`/`CompleteOffer` (mock-initiated),
  `GetSeenMessages`. Publishes `peer-connected`/`peer-disconnected`/`data-channel-opened`/
  `data-channel-message-received`/`data-channel-message-sent`/`data-channel-closed`/
  `media-track-opened` onto the existing `events.Bus.PublishRTC` via `rtc.NormalizeEvent`
  — **zero event-contract changes**, so `webui/src/services/proxyAdminClient.ts` and
  `contracts/mockrtc-events.json` need no changes.
- `graphql.go` — a GraphQL schema/resolver set (reusing the existing
  `github.com/graphql-go/graphql` dependency already used by `internal/api/gql` — **not** a
  second GraphQL library) implementing the exact mutations/query the webextension bundle
  calls, backed by the `Handler` above.
- `rules.go` — a v1 data-channel rule engine for `PUT /rules/rtc` (`echo`, `send`, `close`
  steps matched by `has-data-channel`/`channel-label` and `message-content`/
  `message-includes` matchers). Media-track transcoding/recording rules are **out of scope
  for v1** (matches the plan's original scoping rationale).
- Wired into `internal/proxy/admin/server.go` — the native handler is always
  present when WebRTC is enabled (no `HTK_RTC_BACKEND` env var, no Node
  subprocess). `handleMockRTCSession`'s POST branch detects GraphQL-shaped
  bodies (`{"query": ...}`) and dispatches to `native.GraphQLHandler`;
  `handleRTCRules` feeds `PUT /rules/rtc` payloads into the native rule
  engine; `stopSession` closes all native peer connections.
- **Tests** (`internal/rtc/native/peer_test.go`) validate both flows end-to-end using a real
  second `pion/webrtc` PeerConnection standing in for a browser tab: full offer/answer
  negotiation to `PeerConnectionStateConnected`, a real data-channel message round-trip,
  rule-driven echo, and `getSeenMessages` — not just unit-level mocks. `go test
  ./internal/rtc/native/...` passes in ~5s.
- **Sidecar removed (2026-08):** `internal/rtc/sidecar.go`,
  `webui/scripts/mockrtc-sidecar.mjs`, `HTK_RTC_BACKEND`, `HTK_REPO_ROOT`,
  and `HTK_SIDECAR_PORT` have all been deleted. The native Pion handler is
  the only WebRTC backend.

**Remaining limitations:**
1. **No live browser validation yet.** The tests above simulate the browser side with a
   second pion `PeerConnection`, which validates the actual SDP/ICE/DTLS/SCTP wire protocol
   correctly, but a real Chrome/Firefox tab going through the actual MITM proxy + real
   webextension bundle has not been exercised end-to-end.
2. **ICE/STUN/TURN config** — `peerConnectionConfig()` currently returns an empty
   `webrtc.Configuration{}` (no STUN servers), relying on host-candidate-only negotiation
   since both sides are on localhost/LAN. Revisit if mocking WebRTC across a NAT boundary
   is ever needed.
3. **`createExternalOffer`/`answerExternalOffer`** are currently implemented identically to
   `createOffer`/`answerOffer` (mockrtc distinguishes "external" peers connecting to a
   real upstream vs the mock itself — that distinction isn't modeled yet).

Risk: Pion opens UDP ports for ICE — needs the same kind of firewall documentation already
written for the TCP MITM port in `docs/internal/status.md`.

### P3 — Nice to have / legacy

| Item | Description |
|------|-------------|
| **Native MockRTC (Pion)** | ✅ Complete — in-process WebRTC mocking, no Node dependency. |
| **GraphQL REST on /** | Legacy Mockttp GraphQL endpoint; UI no longer needs it. |
| **`/rules/bulk-create`** | REST bulk rule import (501 today; UI uses admin PUT). |
| **`/webhooks/capture`** | Inbound webhook capture inbox for rule debugging. |
| **`POST /update`** | Auto-update channel (Electron/updater). |
| **Integration tests** | Live server tests: start session, set rules, assert WS events, native RTC. |
| **OpenAPI fixture** | `rest-api.openapi.yaml` referenced in contracts README but not committed. |
| **docker-socks-tunnel image** | Optional container-based tunnel vs local SOCKS. |

### Phase 3 — MITM hardening (test coverage + bug fixes)

> **2026-09 update:** comprehensive test coverage was added for the three
> previously-untested MITM subsystems (SOCKS5, HTTP/2, WebSocket E2E), and
> several production bugs were found and fixed in the process.

**SOCKS5 (`socks_test.go`):**
- 6 tests covering handshake with IPv4, domain, and IPv6 address types;
  unsupported command/address-type rejection; wrong-version rejection.
- **Bug fixed:** `handshake()` used `io.ReadAtLeast(conn, buf, 2)` which can
  over-read past the 2-byte greeting, consuming method bytes and causing a
  deadlock. Replaced with `bufio.NewReader` + `io.ReadFull` for precise reads.
- **Bug fixed:** IPv6 target was formatted as `::1:443` (ambiguous) instead of
  `[::1]:443`. Now uses `net.JoinHostPort`.

**HTTP/2 (`http2_test.go`):**
- 3 tests: ALPN negotiation (h2 when enabled, http/1.1 fallback when disabled),
  and request interception through the HTTP/2 path with a fixed-response rule.
- **Bug fixed:** `writeFixedResponse` always wrote raw `HTTP/1.1 ...` bytes to
  the `io.Writer`, corrupting HTTP/2 framing when the writer was an
  `http.ResponseWriter` from `http2.Server.ServeConn`. Now detects
  `http.ResponseWriter` and uses `WriteHeader`/`Write` instead.

**WebSocket E2E (`websocket_e2e_test.go`):**
- 4 tests: echo mode round-trip, passthrough relay to upstream, non-upgrade
  request handling, and `isWebSocketUpgrade` detection (6 sub-cases).
- **Bug fixed:** `encodeWSFrame` masked server-to-client frames with a hardcoded
  mask, violating RFC 6455 §5.3 (servers MUST NOT mask). This caused clients to
  receive garbled payloads. Now sends unmasked frames.
- **Bug fixed:** `handleHTTP` didn't populate `req.URL.Host` from the Host
  header for non-proxy-style requests, causing WebSocket upstream connections
  to fail with an empty host. Now sets `req.URL.Host = req.Host`.
- **Cleanup:** Removed stray `fmt.Printf` debug statement in `handleHTTP`.

**Route parity (`server.go`, `handlers.go`):**
- Added `GET /auth/callback` — OAuth callback HTML page matching Node's
  `api-server.ts` endpoint.
- Added `GET /interceptors/android-adb/certificate/status` and
  `POST /interceptors/android-adb/certificate/install` — backed by new
  `getDeviceCertificateStatus`/`installDeviceCertificate` methods on the
  `androidAdb` interceptor, matching Node's REST API.

---

## Parity summary

| Area | Estimate | Notes |
|------|----------|-------|
| REST API | **~99%** | `update`/`bulk-create`/`webhooks` stubs match Node's own behavior 1:1 (verified); `auth/callback` and `android-adb/certificate/*` routes now implemented |
| Proxy admin | **~90%** | Core session + rules + events complete |
| MITM / HTTP traffic | **~97%** | Corrected: transform/redirect/abort steps and passthrough options are implemented and applied; HTTP/2 response writing fixed; WebSocket echo/passthrough E2E tested |
| Rule matchers | **100% of UI-reachable matchers** | Corrected: every matcher type `mockttpRuleBuilder.ts` can emit is implemented; extra mockttp matcher types unreachable from this fork's UI remain unimplemented but are dead code paths |
| Rule steps | **100% of UI-reachable steps** | Corrected: every step type `ruleStepMapper.ts` can emit is implemented, including `redirect`/`abort`/`forward-to-host`/`req-res-transformer`. Only true gap: live breakpoint/callback pause-resume (UI intentionally avoids sending it today) |
| Interceptors | **~75%** | All IDs present; mobile/Frida/JVM/Safari shallow — genuine remaining depth work (see P2) |
| Docker | **~90%** | Local SOCKS vs container tunnel |
| WebRTC / MockRTC | **~80%** | Native Go (Pion) — no Node sidecar. Offer/answer + data-channel round-trip validated by automated tests; **not yet validated against a real browser**, and `createExternalOffer`/media-track rules are simplified. |
| MCP | **~85%** | Stdio + ctl; full tool surface depends on UI bridge ops |
| Snippets | **~90%** | Most UI catalog targets covered |
| Desktop (Wails) | **~80%** | Server solid; installer/packaging ongoing |

**Bottom line:** this backlog was significantly overstated before the 2026-08 audit. The
proxy/rule engine core is essentially at parity with what the UI can exercise. The Node
WebRTC sidecar has been removed and replaced with native Pion. Real remaining work is
concentrated in: (1) deeper mobile/Frida/JVM interceptor implementations, and (2) live
browser validation of the native WebRTC path. The two runtime bugs fixed in this pass
(probe timeouts, missing NSS binaries) suggest an integration/E2E test that actually
exercises `go test ./...` to completion should run in CI to catch regressions like these
going forward.

---

## Key file map

| Area | Path |
|------|------|
| Entry | `cmd/htk-server/main.go` |
| App wiring | `internal/server/app.go` |
| Config | `internal/config/config.go` |
| REST API | `internal/api/` |
| Proxy admin | `internal/proxy/admin/` |
| MITM + WS | `internal/proxy/mitm/` |
| Rules | `internal/proxy/rules/` |
| Events bus | `internal/proxy/events/` |
| Docker | `internal/docker/` |
| RTC (native) | `internal/rtc/`, `internal/rtc/native/` |
| WebExtension | `internal/webextension/` |
| Interceptors | `internal/interceptors/` |
| Snippets | `internal/snippets/` |
| MCP | `internal/mcp/`, `cmd/htk-mcp/` |
| UI bridge | `internal/uibridge/` |
| Contracts | `contracts/` |
| Assets | `assets/` (overrides, nss, webextension) |
| UI admin client | `../webui/src/services/proxyAdminClient.ts` |
| Node reference | `../httptoolkit-node/src/` |

---

## Build, run, test

All UI development and build scripts must run with the working directory set to `webui/`:

```powershell
cd webui
npm run copy:go-assets    # overrides + nss + webextension into httptoolkit-go/assets
npm run go:server         # REST :45457, admin :45456
npm run go:build          # bin/htk-server
npm run dev:go            # Vite UI + Go server
npm run wails:dev         # Wails desktop + Go + hosted UI
```

From `httptoolkit-go/`:

```powershell
go build ./...
go test ./...
go test ./contracts/...
```

### Environment variables

| Variable | Purpose |
|----------|---------|
| `HTK_SERVER_TOKEN` | Bearer auth (auto-generated if unset) |
| `HTK_SERVER_PORT` | REST port (default 45457) |
| `HTK_ADMIN_PORT` | Admin port (default 45456) |
| `HTK_ASSETS_DIR` | Bundled overrides path |
| `HTK_CONFIG_DIR` | CA + Firefox profile storage |
| `HTK_DEV` | Relaxed CORS (`0` = strict) |

### WebRTC

WebRTC mocking is handled natively by the Go process using
[`pion/webrtc`](https://github.com/pion/webrtc) (`internal/rtc/native/`).
No Node.js runtime or sidecar subprocess is required. The webextension's
MockRTC client speaks GraphQL directly to the Go admin server's
`/session/:id` endpoint.

---

## UI compatibility notes

- **`proxyAdminClient.ts`** is the Go admin client (not PluggableAdmin GraphQL for HTTP rules).
- **Rule JSON** is built by `webui/src/utils/mockttpRuleBuilder.ts` and sent as `PUT /rules/http|ws|rtc`.
- **Breakpoints** in UI intentionally map to passthrough until callback bridge exists.
- **Event handlers** in `webui/src/services/trafficHttpHandlers.ts` expect Mockttp event names; Go emits the same names where implemented.

---

## Recommended next steps

1. **Live browser WebRTC validation** — exercise the native Pion handler with a real Chrome tab + webextension to validate end-to-end.
2. **Rule gaps** — implement `redirect`, `abort`, and passthrough transforms in `engine.go` / `mitm/server.go`.
3. **Integration test** — start server, activate Chromium, verify one HTTP rule and one RTC event on WS.
4. **Android / Frida depth** — port critical paths from `httptoolkit-node/src/interceptors/android/` and `frida/`.
5. **Commit `rest-api.openapi.yaml`** — or remove reference from `contracts/README.md`.

---

## Changelog (migration milestones)

| Milestone | Summary |
|-----------|---------|
| Phase 0 | Contracts, admin API design, `proxyAdminClient.ts` |
| Core proxy | MITM, events WS, basic rules, cert manager |
| WebSocket | `ws_proxy.go` — upgrade, passthrough, echo |
| Docker | API proxy, build inject, attach interceptor, tunnel SOCKS |
| Interceptors | Chromium, Firefox NSS, terminal, system proxy, ADB, Frida stubs |
| WebExtension | Temp install, per-session MockRTC config |
| MockRTC (native) | Pion WebRTC + GraphQL handler + RTC rules engine |
| MCP | stdio server + ctl pipe |
| Snippets | Expanded language/client matrix |
| Rule steps | webhook, callback event, stream, wait-for-body |
| Polish | system-proxy deactivate all platforms, terminal env parity, cert-check server |

---

*Update this file when closing migration gaps or changing architecture decisions.*
