package rules

import (
	"encoding/json"
)

// RequestRuleData mirrors mockttp RequestRuleData JSON from the UI admin client.
type RequestRuleData struct {
	ID                string            `json:"id,omitempty"`
	Matchers          []json.RawMessage `json:"matchers"`
	Steps             []json.RawMessage `json:"steps"`
	CompletionChecker json.RawMessage   `json:"completionChecker,omitempty"`
}

type WebSocketRuleData struct {
	ID                string            `json:"id,omitempty"`
	Matchers          []json.RawMessage `json:"matchers"`
	Steps             []json.RawMessage `json:"steps"`
	CompletionChecker json.RawMessage   `json:"completionChecker,omitempty"`
}

type Engine struct {
	httpRules []RequestRuleData
	wsRules   []WebSocketRuleData
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) SetHTTPRules(rules []RequestRuleData) {
	e.httpRules = rules
}

func (e *Engine) SetWSRules(rules []WebSocketRuleData) {
	e.wsRules = rules
}

func (e *Engine) HTTPRules() []RequestRuleData { return e.httpRules }
func (e *Engine) WSRules() []WebSocketRuleData { return e.wsRules }

func (e *Engine) MatchWS(urlStr string, headers map[string]string) int {
	ctx := requestCtx{method: "GET", url: urlStr, headers: normalizeHeaders(headers)}
	for i, rule := range e.wsRules {
		if allMatchersMatch(rule.Matchers, ctx) {
			return i
		}
	}
	return -1
}

func (e *Engine) WSAction(index int) StepAction {
	if index < 0 || index >= len(e.wsRules) {
		return StepAction{Kind: "passthrough"}
	}
	action := parseSteps(e.wsRules[index].Steps)
	action.RuleID = e.wsRules[index].ID
	return action
}
