package admin

import (
	"encoding/json"
	"net/http"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/rtc/native"
)

func (s *Server) handleRTCRules(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Rules []json.RawMessage `json:"rules"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if s.nativeRTC != nil {
		nativeRules := make([]native.RTCRule, 0, len(body.Rules))
		for _, raw := range body.Rules {
			var rule struct {
				Matchers []json.RawMessage `json:"matchers"`
				Steps    []json.RawMessage `json:"steps"`
			}
			if err := json.Unmarshal(raw, &rule); err != nil {
				continue
			}
			nativeRules = append(nativeRules, native.RTCRule{Matchers: rule.Matchers, Steps: rule.Steps})
		}
		s.nativeRTC.SetRules(nativeRules)
	}
	writeJSON(w, map[string]any{"ok": true})
}
