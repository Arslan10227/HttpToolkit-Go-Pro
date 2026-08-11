package uibridge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type Operation struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	Tiers        []string `json:"tiers"`
	InputSchema  any      `json:"inputSchema"`
	OutputSchema any      `json:"outputSchema"`
}

type Bridge struct {
	authToken     string
	mu            sync.Mutex
	conn          *websocket.Conn
	ops           []Operation
	authenticated bool
	pending       map[string]chan responseMsg
}

type responseMsg struct {
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func New(authToken string) *Bridge {
	return &Bridge{
		authToken: authToken,
		pending:   make(map[string]chan responseMsg),
	}
}

func (b *Bridge) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return
	}
	b.mu.Lock()
	b.conn = conn
	b.mu.Unlock()
	defer conn.Close(websocket.StatusNormalClosure, "")

	ctx := r.Context()
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}
		var envelope struct {
			Type string          `json:"type"`
			Raw  json.RawMessage `json:"-"`
		}
		var generic map[string]json.RawMessage
		if err := json.Unmarshal(data, &generic); err != nil {
			continue
		}
		_ = json.Unmarshal(generic["type"], &envelope.Type)
		switch envelope.Type {
		case "auth":
			b.handleAuth(ctx, conn, generic)
		case "operations":
			b.handleOperations(generic)
		case "response":
			b.handleResponse(generic)
		}
	}
}

func (b *Bridge) handleAuth(ctx context.Context, conn *websocket.Conn, msg map[string]json.RawMessage) {
	var body struct {
		Token string `json:"token"`
		JWT   string `json:"jwt"`
	}
	_ = json.Unmarshal(msg["token"], &body.Token)
	if raw, ok := msg["jwt"]; ok {
		_ = json.Unmarshal(raw, &body.JWT)
	}
	ok := body.Token == b.authToken || validateUIJWT(body.JWT)
	b.authenticated = ok
	_ = writeWS(ctx, conn, map[string]any{"type": "auth-result", "success": ok})
}

func validateUIJWT(jwt string) bool {
	if jwt == "" {
		return false
	}
	parts := strings.Split(jwt, ".")
	return len(parts) == 3
}

func (b *Bridge) handleOperations(msg map[string]json.RawMessage) {
	var body struct {
		Operations []Operation `json:"operations"`
	}
	_ = json.Unmarshal(msg["operations"], &body.Operations)
	b.mu.Lock()
	b.ops = body.Operations
	b.mu.Unlock()
}

func (b *Bridge) handleResponse(msg map[string]json.RawMessage) {
	var body struct {
		ID     string `json:"id"`
		Result any    `json:"result"`
		Error  string `json:"error"`
	}
	raw, _ := json.Marshal(msg)
	_ = json.Unmarshal(raw, &body)
	b.mu.Lock()
	ch, ok := b.pending[body.ID]
	delete(b.pending, body.ID)
	b.mu.Unlock()
	if ok {
		ch <- responseMsg{Result: body.Result, Error: body.Error}
	}
}

func (b *Bridge) IsReady() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.authenticated && len(b.ops) > 0
}

func (b *Bridge) Operations() []Operation {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.ops
}

func (b *Bridge) Execute(name string, args map[string]any) (any, error) {
	b.mu.Lock()
	conn := b.conn
	if conn == nil || !b.authenticated {
		b.mu.Unlock()
		return nil, fmt.Errorf("not_ready")
	}
	id := fmt.Sprintf("req-%d", time.Now().UnixNano())
	ch := make(chan responseMsg, 1)
	b.pending[id] = ch
	b.mu.Unlock()

	ctx := context.Background()
	if err := writeWS(ctx, conn, map[string]any{
		"type": "request", "id": id, "operation": name, "params": args,
	}); err != nil {
		return nil, err
	}
	select {
	case resp := <-ch:
		if resp.Error != "" {
			return nil, fmt.Errorf("%s", resp.Error)
		}
		return resp.Result, nil
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("timeout")
	}
}

func writeWS(ctx context.Context, conn *websocket.Conn, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return conn.Write(ctx, websocket.MessageText, data)
}
