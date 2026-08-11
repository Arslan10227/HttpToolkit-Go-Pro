package contracts_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestMockttpEventsFixture(t *testing.T) {
	path := filepath.Join(fixtureDir(), "mockttp-events.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		EventNames []string `json:"eventNames"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	if len(doc.EventNames) < 10 {
		t.Fatalf("expected at least 10 mockttp events, got %d", len(doc.EventNames))
	}
}

func TestMockRTCEventsFixture(t *testing.T) {
	path := filepath.Join(fixtureDir(), "mockrtc-events.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		EventNames []string `json:"eventNames"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	if len(doc.EventNames) != 10 {
		t.Fatalf("expected 10 mockrtc events, got %d", len(doc.EventNames))
	}
}

func fixtureDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}
