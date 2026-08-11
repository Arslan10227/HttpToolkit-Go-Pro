# Changelog

All notable changes to **HttpToolkit Go Pro** are documented in this file.

## [0.1.0-beta] - 2026-08-12

First public beta release of the native Go + Wails v3 desktop port of HTTP Toolkit.

### Added
- Native Go MITM proxy with HTTP/HTTPS, HTTP/2, WebSocket and WebRTC support.
- Wails v3 desktop shell with the React web UI embedded.
- Interceptors: Chrome, Chromium, Edge, Brave, Opera, Arc, Firefox, Safari, Terminal, Electron, JVM attach, Android ADB, Android Frida, iOS Frida, Docker attach, and system proxy.
- Android QR-code pairing for the official HTTP Toolkit Android app.
- Dynamic CA certificate generation and SPKI fingerprinting.
- REST + GraphQL admin APIs, MCP tool server, HAR import/export, and rule sync.
- Google OAuth sign-in and cloud sync (Firebase / Upstash Redis).
- GitHub Actions CI and automated release workflow.

### Fixed
- Android QR scan not completing because `android.httptoolkit.tech/config` was not intercepted by the proxy.
- OkHttp3 / Retrofit JVM traffic not being fully captured by the Java agent.
- UI image metadata and syntax highlighting issues in payload inspectors.
- Race conditions in interceptor activation/deactivation.

### Known Issues
- Some interceptor configuration panels are still being hardened.
- macOS and Linux desktop packages are built but not notarized/signed.

[0.1.0-beta]: https://github.com/Arslan10227/HttpToolkit-Go-Pro/releases/tag/v0.1.0-beta
