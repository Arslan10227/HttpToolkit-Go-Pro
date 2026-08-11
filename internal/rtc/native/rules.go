package native

import (
	"encoding/json"
	"strings"
	"sync"
)

// ActionKind is the effect a matched data-channel rule applies.
type ActionKind int

const (
	ActionNone ActionKind = iota
	ActionEcho
	ActionSend
	ActionClose
)

// DataChannelAction is the resolved behaviour to apply for one matched rule.
type DataChannelAction struct {
	Kind        ActionKind
	Text        string
	Binary      []byte
	LastMessage []byte // populated by the caller for ActionEcho
}

// RTCRule mirrors the (loosely-typed) rule JSON accepted by `PUT /rules/rtc`,
// scoped for v1 to data-channel matching/steps — see port.md's "Native
// WebRTC migration plan" for why media-track rules are out of scope for now.
type RTCRule struct {
	Matchers []json.RawMessage `json:"matchers"`
	Steps    []json.RawMessage `json:"steps"`
}

// RuleEngine evaluates PUT /rules/rtc payloads against data-channel activity.
type RuleEngine struct {
	mu    sync.RWMutex
	rules []RTCRule
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{}
}

func (e *RuleEngine) SetRules(rules []RTCRule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = rules
}

// MatchOpen returns actions to apply as soon as a data channel with the
// given label opens (e.g. a "send" step with no message-matcher).
func (e *RuleEngine) MatchOpen(label string) []DataChannelAction {
	return e.match(label, nil, false)
}

// MatchMessage returns actions to apply after a message is received on the
// given channel (e.g. "echo" or a matched "send" reply).
func (e *RuleEngine) MatchMessage(label string, data []byte) []DataChannelAction {
	return e.match(label, data, true)
}

func (e *RuleEngine) match(label string, data []byte, hasMessage bool) []DataChannelAction {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var actions []DataChannelAction
	for _, rule := range e.rules {
		if !matchersApply(rule.Matchers, label, data, hasMessage) {
			continue
		}
		for _, raw := range rule.Steps {
			var obj map[string]any
			if err := json.Unmarshal(raw, &obj); err != nil {
				continue
			}
			switch typ, _ := obj["type"].(string); typ {
			case "echo":
				action := DataChannelAction{Kind: ActionEcho}
				if hasMessage {
					action.LastMessage = data
				}
				actions = append(actions, action)
			case "send":
				action := DataChannelAction{Kind: ActionSend}
				if s, ok := obj["value"].(string); ok {
					action.Text = s
				} else if s, ok := obj["message"].(string); ok {
					action.Text = s
				}
				actions = append(actions, action)
			case "close":
				actions = append(actions, DataChannelAction{Kind: ActionClose})
			}
		}
	}
	return actions
}

// matchersApply checks the (small, v1) set of data-channel matchers this
// engine supports: "has-data-channel" (by label) and "message-content"
// (substring match against a received message).
func matchersApply(matchers []json.RawMessage, label string, data []byte, hasMessage bool) bool {
	if len(matchers) == 0 {
		return true
	}
	for _, raw := range matchers {
		var obj map[string]any
		if err := json.Unmarshal(raw, &obj); err != nil {
			return false
		}
		typ, _ := obj["type"].(string)
		switch typ {
		case "has-data-channel", "channel-label":
			want, _ := obj["label"].(string)
			if want != "" && want != label {
				return false
			}
		case "message-content", "message-includes":
			if !hasMessage {
				return false
			}
			want, _ := obj["value"].(string)
			if want != "" && !strings.Contains(string(data), want) {
				return false
			}
		case "wildcard", "":
			// matches anything
		default:
			// Unknown matcher type: fail closed (don't misapply unrelated
			// rules), matching the same conservative behaviour as the HTTP
			// rule engine's `default` case in internal/proxy/rules/engine.go.
			return false
		}
	}
	return true
}
