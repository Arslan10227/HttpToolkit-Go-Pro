package events

import (
	"encoding/json"
	"sync"
)

type Bus struct {
	mu   sync.RWMutex
	subs map[chan []byte]struct{}
}

func NewBus() *Bus {
	return &Bus{subs: make(map[chan []byte]struct{})}
}

func (b *Bus) Subscribe() chan []byte {
	ch := make(chan []byte, 256)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *Bus) Unsubscribe(ch chan []byte) {
	b.mu.Lock()
	delete(b.subs, ch)
	b.mu.Unlock()
	close(ch)
}

func (b *Bus) PublishHTTP(event string, data any) {
	b.publish("http", event, data)
}

func (b *Bus) PublishRTC(event string, data any) {
	b.publish("webrtc", event, data)
}

func (b *Bus) publish(plugin, event string, data any) {
	msg, err := json.Marshal(map[string]any{
		"plugin": plugin,
		"event":  event,
		"data":   data,
	})
	if err != nil {
		return
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.subs {
		select {
		case ch <- msg:
		default:
		}
	}
}
