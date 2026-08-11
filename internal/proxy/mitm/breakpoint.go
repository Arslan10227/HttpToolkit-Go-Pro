package mitm

import (
	"net/http"
	"sync"
	"time"
)

type BreakpointActionType string

const (
	ActionResume  BreakpointActionType = "resume"
	ActionRespond BreakpointActionType = "respond"
	ActionAbort   BreakpointActionType = "abort"
)

type BreakpointAction struct {
	Type        BreakpointActionType
	Method      string
	URL         string
	Headers     map[string]string
	Body        []byte
	StatusCode  int
	RespBody    []byte
	RespHeaders map[string]string
}

type PausedRequest struct {
	ID           string
	RuleID       string
	Req          *http.Request
	OriginalBody []byte
	ActionChan   chan BreakpointAction
	CreatedAt    time.Time
}

type BreakpointManager struct {
	mu     sync.RWMutex
	paused map[string]*PausedRequest
}

func NewBreakpointManager() *BreakpointManager {
	return &BreakpointManager{
		paused: make(map[string]*PausedRequest),
	}
}

func (bm *BreakpointManager) Pause(id string, ruleID string, req *http.Request, body []byte) (*PausedRequest, chan BreakpointAction) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	ch := make(chan BreakpointAction, 1)
	pr := &PausedRequest{
		ID:           id,
		RuleID:       ruleID,
		Req:          req,
		OriginalBody: body,
		ActionChan:   ch,
		CreatedAt:    time.Now(),
	}
	bm.paused[id] = pr
	return pr, ch
}

func (bm *BreakpointManager) Remove(id string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	delete(bm.paused, id)
}

func (bm *BreakpointManager) Get(id string) *PausedRequest {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return bm.paused[id]
}
