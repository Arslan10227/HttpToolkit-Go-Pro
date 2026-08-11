# API contracts (Phase 0)

Golden fixtures and schemas for Go server parity with Node httptoolkit-server.

## REST API (port 45457)

See `rest-api.openapi.yaml` for route definitions matching
`httptoolkit-server/src/api/rest-api.ts`.

## Proxy admin API (port 45456)

See `admin-api.md` for the Go-native admin protocol consumed by
`src/services/proxyAdminClient.ts`.

## Event shapes

- `mockttp-events.json` — HTTP/WebSocket traffic event names and sample payloads
- `mockrtc-events.json` — WebRTC event names and sample payloads

## Rule JSON

Rules sent as `RequestRuleData[]` / `WebSocketRuleData[]` JSON (built by
`src/utils/mockttpRuleBuilder.ts`). Go interprets the serialized mockttp schema
in `internal/proxy/rules/`.

## Contract tests

Run: `go test ./contracts/...` from `httptoolkit-go/`.
