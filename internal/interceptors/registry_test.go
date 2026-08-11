package interceptors

import (
	"testing"

	certmgr "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
)

func TestRegistryOrderAndList(t *testing.T) {
	cfg := &config.Config{ConfigDir: t.TempDir()}
	certs, err := certmgr.NewManager(cfg)
	if err != nil {
		t.Fatalf("failed to create cert manager: %v", err)
	}
	r := NewRegistry(cfg, "", certs)

	// Verify that r.order is populated and matches buildAll length.
	all := buildAll(cfg, "", certs)
	if len(r.order) != len(all) {
		t.Fatalf("expected order length %d, got %d", len(all), len(r.order))
	}

	// Verify order matches buildAll output exactly.
	for i, id := range r.order {
		if id != all[i].ID() {
			t.Fatalf("expected order element %d to be %s, got %s", i, all[i].ID(), id)
		}
	}

	// Verify that r.List() returns items in the same order as r.order.
	list := r.List(8080)
	if len(list) != len(r.order) {
		t.Fatalf("expected list length %d, got %d", len(r.order), len(list))
	}
	for i, item := range list {
		if item["id"] != r.order[i] {
			t.Errorf("mismatched list item at index %d: expected %s, got %v", i, r.order[i], item["id"])
		}
	}
}
