package backup

import (
	"context"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
)

// Backuper is the interface the settings manager uses to persist a copy of
// settings to a cloud backend.
type Backuper interface {
	SaveSettings(ctx context.Context, data any) error
	RestoreSettings(ctx context.Context, out map[string]any) error
}

// New creates a cloud backup client. Upstash Redis is preferred when a Redis
// URL is configured; otherwise it falls back to Firebase if a Firebase project
// is configured.
func New(cfg *config.Config) Backuper {
	if cfg.UpstashRedisURL != "" {
		return NewUpstash(cfg)
	}
	return NewFirebase(cfg)
}
