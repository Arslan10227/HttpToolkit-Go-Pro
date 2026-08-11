package rules

import (
	"encoding/json"
	"log"
	"net/url"
	"regexp"
	"strings"
)

type requestCtx struct {
	method  string
	url     string
	headers map[string]string
	body    []byte
}

func (e *Engine) MatchFirst(method, rawURL string, headers map[string]string, body []byte) int {
	ctx := requestCtx{
		method:  strings.ToUpper(method),
		url:     rawURL,
		headers: normalizeHeaders(headers),
		body:    body,
	}
	for i, rule := range e.httpRules {
		if allMatchersMatch(rule.Matchers, ctx) {
			return i
		}
	}
	return -1
}

func allMatchersMatch(matchers []json.RawMessage, ctx requestCtx) bool {
	if len(matchers) == 0 {
		return true
	}
	for _, m := range matchers {
		if !matcherMatches(m, ctx) {
			return false
		}
	}
	return true
}

func matcherMatches(raw json.RawMessage, ctx requestCtx) bool {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return false
	}
	typ, _ := obj["type"].(string)
	switch typ {
	case "wildcard", "":
		return true
	case "method":
		want, _ := obj["method"].(string)
		return strings.EqualFold(want, ctx.method)
	case "host":
		want, _ := obj["host"].(string)
		return hostMatches(ctx, want)
	case "hostname":
		want, _ := obj["hostname"].(string)
		return hostnameMatches(ctx, want)
	case "protocol":
		want, _ := obj["protocol"].(string)
		u, err := url.Parse(ctx.url)
		if err != nil {
			return false
		}
		return strings.EqualFold(u.Scheme, want)
	case "simple-path", "path":
		want, _ := obj["path"].(string)
		return matchPath(ctx, want)
	case "flexible-path":
		want, _ := obj["path"].(string)
		u, err := url.Parse(ctx.url)
		if err != nil {
			return false
		}
		return strings.Contains(u.Path, want)
	case "regex-path":
		return regexFieldMatches(obj, "regexSource", ctx.url)
	case "regex-url":
		return regexFieldMatches(obj, "regexSource", ctx.url)
	case "header":
		return headerMatcherMatches(obj, ctx.headers)
	case "query":
		return queryMatcherMatches(obj, ctx.url)
	case "exact-query":
		return exactQueryMatcherMatches(obj, ctx.url)
	case "raw-body-includes", "raw-body":
		content, _ := obj["content"].(string)
		if content == "" {
			content, _ = obj["value"].(string)
		}
		return strings.Contains(string(ctx.body), content)
	case "json-body", "json-body-flexible", "json-body-matching":
		return jsonBodyMatches(obj, ctx.body)
	default:
		log.Printf("[rules] unknown matcher type %q — no match", typ)
		return false
	}
}

func jsonBodyMatches(obj map[string]any, body []byte) bool {
	wantRaw, ok := obj["value"]
	if !ok {
		wantRaw = obj["body"]
	}
	if wantRaw == nil {
		return len(body) == 0
	}
	want, err := json.Marshal(wantRaw)
	if err != nil {
		return false
	}
	var got any
	if err := json.Unmarshal(body, &got); err != nil {
		return false
	}
	var expected any
	if err := json.Unmarshal(want, &expected); err != nil {
		return false
	}
	gotBytes, _ := json.Marshal(got)
	expBytes, _ := json.Marshal(expected)
	return string(gotBytes) == string(expBytes)
}

func regexFieldMatches(obj map[string]any, field, target string) bool {
	pat, _ := obj[field].(string)
	if pat == "" {
		if r, ok := obj["regex"].(string); ok {
			pat = r
		}
	}
	if pat == "" {
		return false
	}
	re, err := regexp.Compile(pat)
	if err != nil {
		return false
	}
	u, _ := url.Parse(target)
	if u != nil && field == "regexSource" && obj["type"] == "regex-path" {
		return re.MatchString(u.Path)
	}
	return re.MatchString(target)
}

func headerMatcherMatches(obj map[string]any, headers map[string]string) bool {
	if h, ok := obj["headers"].(map[string]any); ok {
		for k, v := range h {
			if !headerValueMatches(headers, k, toString(v)) {
				return false
			}
		}
		return true
	}
	key, _ := obj["key"].(string)
	val, _ := obj["value"].(string)
	if key != "" {
		return headerValueMatches(headers, key, val)
	}
	return false
}

func headerValueMatches(headers map[string]string, key, want string) bool {
	got, ok := headers[strings.ToLower(key)]
	if !ok {
		return false
	}
	if strings.Contains(want, "*") {
		pat := "^" + regexp.QuoteMeta(want) + "$"
		pat = strings.ReplaceAll(pat, "\\*", ".*")
		re, err := regexp.Compile(pat)
		return err == nil && re.MatchString(got)
	}
	return got == want
}

func queryMatcherMatches(obj map[string]any, rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	q := u.Query()
	if qm, ok := obj["query"].(map[string]any); ok {
		for k, v := range qm {
			vals := q[k]
			if len(vals) == 0 {
				return false
			}
			if !strings.Contains(toString(v), "*") && vals[0] != toString(v) {
				return false
			}
		}
		return true
	}
	name, _ := obj["name"].(string)
	val, _ := obj["value"].(string)
	return q.Get(name) == val
}

func exactQueryMatcherMatches(obj map[string]any, rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	q := u.Query()
	if qm, ok := obj["query"].(map[string]any); ok {
		if len(qm) != len(q) {
			return false
		}
		for k, v := range qm {
			if q.Get(k) != toString(v) {
				return false
			}
		}
		return true
	}
	return false
}

func hostMatches(ctx requestCtx, host string) bool {
	u, err := url.Parse(ctx.url)
	if err != nil {
		return false
	}
	h := u.Host
	if strings.Contains(h, ":") {
		h = strings.Split(h, ":")[0]
	}
	return strings.EqualFold(h, host) || headerContainsHost(ctx.headers, host)
}

func hostnameMatches(ctx requestCtx, host string) bool {
	return hostMatches(ctx, host)
}

func headerContainsHost(headers map[string]string, host string) bool {
	got, ok := headers["host"]
	if !ok {
		return false
	}
	gotHost := got
	if strings.Contains(gotHost, ":") {
		gotHost = strings.Split(gotHost, ":")[0]
	}
	return strings.EqualFold(gotHost, host)
}

func normalizeHeaders(h map[string]string) map[string]string {
	out := make(map[string]string, len(h))
	for k, v := range h {
		out[strings.ToLower(k)] = v
	}
	return out
}

type StepAction struct {
	Kind        string
	StatusCode  int
	Body        string
	Headers     map[string]string
	DelayMs     int
	FilePath    string
	RuleID      string
	Passthrough PassthroughOptions
	Webhook     WebhookAction
	StreamFile  string
	IsCallback  bool
}

func (e *Engine) ActionForRule(index int) StepAction {
	if index < 0 || index >= len(e.httpRules) {
		return StepAction{Kind: "passthrough"}
	}
	rule := e.httpRules[index]
	action := parseSteps(rule.Steps)
	action.RuleID = rule.ID
	action.Webhook = extractWebhookFromSteps(rule.Steps)
	if action.Kind == "passthrough" {
		action.Passthrough = extractPassthroughFromSteps(rule.Steps)
	}
	return action
}

func (e *Engine) RuleID(index int) string {
	if index < 0 || index >= len(e.httpRules) {
		return ""
	}
	return e.httpRules[index].ID
}

func parseSteps(steps []json.RawMessage) StepAction {
	action := StepAction{Kind: "passthrough"}
	for _, s := range steps {
		var obj map[string]any
		if err := json.Unmarshal(s, &obj); err != nil {
			continue
		}
		typ, _ := obj["type"].(string)
		switch typ {
		case "simple":
			status := 200
			if sc, ok := obj["status"].(float64); ok {
				status = int(sc)
			}
			body, _ := obj["data"].(string)
			headers := map[string]string{}
			if h, ok := obj["headers"].(map[string]any); ok {
				for k, v := range h {
					headers[k] = toString(v)
				}
			}
			return StepAction{Kind: "fixed-response", StatusCode: status, Body: body, Headers: headers}
		case "file":
			status := 200
			if sc, ok := obj["status"].(float64); ok {
				status = int(sc)
			}
			fp, _ := obj["filePath"].(string)
			return StepAction{Kind: "file", StatusCode: status, FilePath: fp}
		case "close-connection":
			return StepAction{Kind: "close-connection"}
		case "reset-connection":
			return StepAction{Kind: "reset-connection"}
		case "timeout":
			ms := 30000
			if v, ok := obj["delayMs"].(float64); ok {
				ms = int(v)
			}
			if v, ok := obj["timeoutMs"].(float64); ok {
				ms = int(v)
			}
			return StepAction{Kind: "timeout", DelayMs: ms}
		case "delay":
			ms := 0
			if v, ok := obj["delayMs"].(float64); ok {
				ms = int(v)
			}
			action.DelayMs += ms
		case "passthrough", "forward", "wait-for-body":
			action.Kind = "passthrough"
			action.Passthrough = parsePassthroughOptions(obj)
		case "callback":
			action.Kind = "passthrough"
			action.IsCallback = true
			action.Passthrough = parsePassthroughOptions(obj)
		case "breakpoint":
			action.Kind = "breakpoint"
			action.Passthrough = parsePassthroughOptions(obj)
		case "webhook":
			action.Kind = "passthrough"
			action.Webhook = parseWebhookStep(obj)
			action.Passthrough = parsePassthroughOptions(obj)
		case "stream":
			action.Kind = "stream"
			if fp, _ := obj["filePath"].(string); fp != "" {
				action.StreamFile = fp
			}
			action.Passthrough = parsePassthroughOptions(obj)
		case "redirect":
			status := 302
			if sc, ok := obj["status"].(float64); ok {
				status = int(sc)
			}
			loc, _ := obj["location"].(string)
			if loc == "" {
				loc, _ = obj["url"].(string)
			}
			return StepAction{Kind: "redirect", StatusCode: status, Body: loc}
		case "abort":
			return StepAction{Kind: "abort"}
		case "ws-passthrough", "pass-through-ws":
			action.Kind = "ws-passthrough"
			action.Passthrough = parsePassthroughOptions(obj)
		case "ws-echo", "echo-ws":
			action.Kind = "ws-echo"
		}
	}
	return action
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func matchPath(ctx requestCtx, want string) bool {
	u, err := url.Parse(ctx.url)
	if err != nil {
		return false
	}
	reqPath := u.Path
	if reqPath == "" {
		reqPath = "/"
	}

	// Case 1: want is absolute URL
	if strings.HasPrefix(want, "http://") || strings.HasPrefix(want, "https://") {
		wantURL, err := url.Parse(want)
		if err != nil {
			return false
		}
		wantPath := wantURL.Path
		if wantPath == "" {
			wantPath = "/"
		}
		// Match host and path
		wantHost := wantURL.Host
		reqHost := u.Host
		if !strings.Contains(wantHost, ":") && strings.Contains(reqHost, ":") {
			reqHost = strings.Split(reqHost, ":")[0]
		}
		return strings.EqualFold(reqHost, wantHost) && (reqPath == wantPath || strings.HasPrefix(reqPath, wantPath))
	}

	// Case 2: want is relative path starting with /
	if strings.HasPrefix(want, "/") {
		return reqPath == want || strings.HasPrefix(reqPath, want)
	}

	// Case 3: want is host/path format (Mockttp flexible path matching behavior)
	parts := strings.SplitN(want, "/", 2)
	wantHost := parts[0]
	wantPath := "/"
	if len(parts) > 1 {
		wantPath = "/" + parts[1]
	}
	reqHost := u.Host
	if !strings.Contains(wantHost, ":") && strings.Contains(reqHost, ":") {
		reqHost = strings.Split(reqHost, ":")[0]
	}
	return strings.EqualFold(reqHost, wantHost) && (reqPath == wantPath || strings.HasPrefix(reqPath, wantPath))
}
