package mitm

import (
	"bytes"
	"testing"
)

func TestDecodeWSFramePayloadText(t *testing.T) {
	raw := []byte{0x81, 0x82, 0x12, 0x34, 0x56, 0x78, 'h' ^ 0x12, 'i' ^ 0x34}
	payload, binary := decodeWSFramePayload(raw)
	if binary {
		t.Fatal("expected text frame")
	}
	if !bytes.Equal(payload, []byte("hi")) {
		t.Fatalf("got %q", payload)
	}
}

func TestEncodeWSFrameRoundTrip(t *testing.T) {
	payload := []byte("hello")
	frame := encodeWSFrame(payload, false)
	got, binary := decodeWSFramePayload(frame)
	if binary {
		t.Fatal("expected text frame")
	}
	if !bytes.Equal(payload, got) {
		t.Fatalf("got %q", got)
	}
}

func TestComputeWebSocketAccept(t *testing.T) {
	key := "dGhlIHNhbXBsZSBub25jZQ=="
	accept := computeWebSocketAccept(key)
	if accept == "" || accept == key {
		t.Fatalf("unexpected accept %q", accept)
	}
	if computeWebSocketAccept(key) != accept {
		t.Fatal("accept should be deterministic")
	}
}
