# Go Proxy Admin API (port 45456)

Replaces Mockttp PluggableAdmin for browser `proxyAdminClient.ts`.

## HTTP

| Method | Path | Body | Response |
|--------|------|------|----------|
| GET | `/metadata` | — | `{ running, httpPort, webrtcEnabled }` |
| POST | `/session/start` | `{ http2?, portRange?, cors? }` | `{ httpProxyPort, socksProxyPort, dockerTunnelPort, webrtcEnabled }` |
| POST | `/session/stop` | — | `{ ok: true }` |
| PUT | `/rules/http` | `{ rules: RequestRuleData[] }` | `{ ok: true }` |
| PUT | `/rules/ws` | `{ rules: WebSocketRuleData[] }` | `{ ok: true }` |
| PUT | `/rules/rtc` | `{ rules: RTCRule[] }` | `{ ok: true }` |
| GET | `/session/:id` | — | `{ sessionId, running }` |
| POST | `/session/:id` | `{ query: "..." }` (GraphQL) or `{ event, data }` | GraphQL result or `{ ok: true }` |

### WebRTC / MockRTC

WebRTC mocking is handled natively by the Go process using
[`pion/webrtc`](https://github.com/pion/webrtc) (`internal/rtc/native/`).
No Node.js runtime or sidecar subprocess is required.

- `PUT /rules/rtc` feeds rules into the native Pion rule engine
  (echo, send, close steps; has-data-channel / channel-label / message
  matchers).
- `POST /session/:id` with a GraphQL body (`{"query": "..."}`) dispatches
  to the native GraphQL handler implementing the webextension's MockRTC
  protocol (createOffer, answerOffer, completeOffer, getSeenMessages).

## WebSocket `GET /events`

Client connects; server pushes:

```json
{ "event": "request-initiated", "data": { ... } }
```

Event names match `MOCKTTP_HTTP_EVENTS` in `trafficHttpHandlers.ts`.

RTC events on same socket:

```json
{ "plugin": "webrtc", "event": "peer-connected", "data": { ... } }
```

## CORS

Same origin allowlist as REST API; `Access-Control-Allow-Private-Network: true`.
