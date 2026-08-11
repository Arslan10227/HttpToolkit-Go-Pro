package rules

import (
	"encoding/json"
)

// PassthroughOptions mirrors mockttp PassThroughStepConnectionOptions JSON.
type PassthroughOptions struct {
	IgnoreHostHttpsErrors []string
	ProxyConfig           any
	LookupServers         []string
	ClientCertHostMap     map[string]ClientCert
	TransformRequest      *TransformRequestOptions
	TransformResponse     *TransformResponseOptions
}

type TransformRequestOptions struct {
	ReplaceHost       *ReplaceHostOptions
	SetProtocol       string
	ReplaceHeaders    map[string]string
	UpdateHeaders     map[string]string
	ReplaceBody       *string
	ReplaceBodyFile   *string
	UpdateJsonBody    map[string]any
	PatchJsonBody     []any
	MatchReplaceBody  []MatchReplacePair
	MatchReplaceHost  []MatchReplacePair
	MatchReplacePath  []MatchReplacePair
	MatchReplaceQuery []MatchReplacePair
}

type TransformResponseOptions struct {
	ReplaceStatus    int
	ReplaceHeaders   map[string]string
	UpdateHeaders    map[string]string
	ReplaceBody      *string
	ReplaceBodyFile  *string
	UpdateJsonBody   map[string]any
	PatchJsonBody    []any
	MatchReplaceBody []MatchReplacePair
}

type ReplaceHostOptions struct {
	TargetHost       string
	UpdateHostHeader bool
}

type MatchReplacePair struct {
	Match   string
	Replace string
}

type ClientCert struct {
	PFX        string
	Passphrase string
}

func parsePassthroughOptions(obj map[string]any) PassthroughOptions {
	opts := PassthroughOptions{
		IgnoreHostHttpsErrors: stringSliceAny(obj["ignoreHostHttpsErrors"]),
		ProxyConfig:           obj["proxyConfig"],
	}
	if lo, ok := obj["lookupOptions"].(map[string]any); ok {
		opts.LookupServers = stringSliceAny(lo["servers"])
	}
	if m, ok := obj["clientCertificateHostMap"].(map[string]any); ok {
		opts.ClientCertHostMap = parseClientCertMap(m)
	}
	if tr, ok := obj["transformRequest"].(map[string]any); ok {
		opts.TransformRequest = parseTransformRequest(tr)
	}
	if tr, ok := obj["transformResponse"].(map[string]any); ok {
		opts.TransformResponse = parseTransformResponse(tr)
	}
	return opts
}

func parseStringMap(v any) map[string]string {
	m, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, val := range m {
		if s, ok := val.(string); ok {
			out[k] = s
		} else if val == nil {
			out[k] = ""
		}
	}
	return out
}

func parseBodyString(v any) *string {
	if v == nil {
		return nil
	}
	if s, ok := v.(string); ok {
		return &s
	}
	if m, ok := v.(map[string]any); ok {
		if typ, ok := m["type"].(string); ok && typ == "Buffer" {
			if dataSlice, ok := m["data"].([]any); ok {
				bytes := make([]byte, len(dataSlice))
				for i, item := range dataSlice {
					if f, ok := item.(float64); ok {
						bytes[i] = byte(f)
					}
				}
				s := string(bytes)
				return &s
			}
		}
	}
	return nil
}

func parseMatchReplacePairs(v any) []MatchReplacePair {
	var pairs []MatchReplacePair
	slice, ok := v.([]any)
	if !ok {
		return nil
	}
	for _, item := range slice {
		if pairSlice, ok := item.([]any); ok && len(pairSlice) == 2 {
			var matchStr string
			if matchObj, ok := pairSlice[0].(map[string]any); ok {
				matchStr, _ = matchObj["source"].(string)
			} else if str, ok := pairSlice[0].(string); ok {
				matchStr = str
			}
			replaceStr, _ := pairSlice[1].(string)
			if matchStr != "" {
				pairs = append(pairs, MatchReplacePair{Match: matchStr, Replace: replaceStr})
			}
			continue
		}
		if pairMap, ok := item.(map[string]any); ok {
			matchStr, _ := pairMap["match"].(string)
			replaceStr, _ := pairMap["replace"].(string)
			if matchStr == "" {
				matchStr, _ = pairMap["key"].(string)
			}
			if matchStr != "" {
				pairs = append(pairs, MatchReplacePair{Match: matchStr, Replace: replaceStr})
			}
		}
	}
	return pairs
}

func parseTransformRequest(m map[string]any) *TransformRequestOptions {
	opts := &TransformRequestOptions{}
	if rh, ok := m["replaceHost"].(map[string]any); ok {
		th, _ := rh["targetHost"].(string)
		uhh := true
		if val, ok := rh["updateHostHeader"].(bool); ok {
			uhh = val
		}
		opts.ReplaceHost = &ReplaceHostOptions{
			TargetHost:       th,
			UpdateHostHeader: uhh,
		}
	}
	if sp, ok := m["setProtocol"].(string); ok {
		opts.SetProtocol = sp
	}
	if rh, ok := m["replaceHeaders"]; ok {
		opts.ReplaceHeaders = parseStringMap(rh)
	}
	if uh, ok := m["updateHeaders"]; ok {
		opts.UpdateHeaders = parseStringMap(uh)
	}
	opts.ReplaceBody = parseBodyString(m["replaceBody"])
	if rbf, ok := m["replaceBodyFromFile"].(string); ok && rbf != "" {
		opts.ReplaceBodyFile = &rbf
	}
	if ujb, ok := m["updateJsonBody"].(map[string]any); ok {
		opts.UpdateJsonBody = ujb
	}
	if pjb, ok := m["patchJsonBody"].([]any); ok {
		opts.PatchJsonBody = pjb
	}
	opts.MatchReplaceBody = parseMatchReplacePairs(m["matchReplaceBody"])
	if mrh, ok := m["matchReplaceHost"].(map[string]any); ok {
		opts.MatchReplaceHost = parseMatchReplacePairs(mrh["replacements"])
	}
	opts.MatchReplacePath = parseMatchReplacePairs(m["matchReplacePath"])
	opts.MatchReplaceQuery = parseMatchReplacePairs(m["matchReplaceQuery"])
	return opts
}

func parseTransformResponse(m map[string]any) *TransformResponseOptions {
	opts := &TransformResponseOptions{}
	if rs, ok := m["replaceStatus"].(float64); ok {
		opts.ReplaceStatus = int(rs)
	}
	if rh, ok := m["replaceHeaders"]; ok {
		opts.ReplaceHeaders = parseStringMap(rh)
	}
	if uh, ok := m["updateHeaders"]; ok {
		opts.UpdateHeaders = parseStringMap(uh)
	}
	opts.ReplaceBody = parseBodyString(m["replaceBody"])
	if rbf, ok := m["replaceBodyFromFile"].(string); ok && rbf != "" {
		opts.ReplaceBodyFile = &rbf
	}
	if ujb, ok := m["updateJsonBody"].(map[string]any); ok {
		opts.UpdateJsonBody = ujb
	}
	if pjb, ok := m["patchJsonBody"].([]any); ok {
		opts.PatchJsonBody = pjb
	}
	opts.MatchReplaceBody = parseMatchReplacePairs(m["matchReplaceBody"])
	return opts
}

func stringSliceAny(v any) []string {
	switch s := v.(type) {
	case []string:
		return s
	case []any:
		out := make([]string, 0, len(s))
		for _, item := range s {
			if str, ok := item.(string); ok {
				out = append(out, str)
			}
		}
		return out
	default:
		return nil
	}
}

func parseClientCertMap(m map[string]any) map[string]ClientCert {
	out := make(map[string]ClientCert, len(m))
	for host, v := range m {
		cm, ok := v.(map[string]any)
		if !ok {
			continue
		}
		pfx, _ := cm["pfx"].(string)
		pass, _ := cm["passphrase"].(string)
		out[host] = ClientCert{PFX: pfx, Passphrase: pass}
	}
	return out
}

func extractPassthroughFromSteps(steps []json.RawMessage) PassthroughOptions {
	for _, s := range steps {
		var obj map[string]any
		if err := json.Unmarshal(s, &obj); err != nil {
			continue
		}
		typ, _ := obj["type"].(string)
		if typ == "passthrough" || typ == "forward" {
			return parsePassthroughOptions(obj)
		}
	}
	return PassthroughOptions{}
}
