package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/logger"
	"github.com/redis/go-redis/v9"
)

// UpstashClient backs up settings to a Redis-compatible Upstash database.
type UpstashClient struct {
	client    *redis.Client
	machineID string
	key       string
}

// NewUpstash creates a Redis backup client from a URL such as the one Upstash
// provides (rediss://default:<token>@<host>:<port>).
func NewUpstash(cfg *config.Config) *UpstashClient {
	if cfg.UpstashRedisURL == "" {
		return nil
	}

	opt, err := redis.ParseURL(cfg.UpstashRedisURL)
	if err != nil {
		logger.Error(fmt.Errorf("upstash redis parse: %w", err), nil)
		return nil
	}

	c := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.Ping(ctx).Err(); err != nil {
		logger.Error(fmt.Errorf("upstash redis ping: %w", err), nil)
		_ = c.Close()
		return nil
	}

	machineID := machineID(cfg)
	return &UpstashClient{
		client:    c,
		machineID: machineID,
		key:       "machines:" + machineID + ":settings",
	}
}

func (c *UpstashClient) SaveSettings(ctx context.Context, data any) error {
	if c == nil || c.client == nil {
		return nil
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}

	if err := c.client.Set(ctx, c.key, string(payload), 0).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	logger.Info("Upstash settings backup saved", map[string]any{
		"machineId": c.machineID,
		"key":       c.key,
	})
	return nil
}

func (c *UpstashClient) RestoreSettings(ctx context.Context, out map[string]any) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("upstash backup not initialized")
	}

	val, err := c.client.Get(ctx, c.key).Result()
	if err != nil {
		return fmt.Errorf("redis get: %w", err)
	}

	if err := json.Unmarshal([]byte(val), &out); err != nil {
		return fmt.Errorf("unmarshal settings: %w", err)
	}
	return nil
}

// userKey returns a Redis key scoped to the signed-in user.
func userKey(email string) string {
	return "users:" + email + ":config"
}

// SaveUserConfig persists a full user config snapshot for the given email.
func (c *UpstashClient) SaveUserConfig(ctx context.Context, email string, data any) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("upstash backup not initialized")
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal user config: %w", err)
	}

	m := make(map[string]any)
	if err := json.Unmarshal(payload, &m); err != nil {
		return fmt.Errorf("unmarshal user config: %w", err)
	}
	m["_savedAt"] = time.Now().UTC().Format(time.RFC3339)

	out, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshal wrapped config: %w", err)
	}

	if err := c.client.Set(ctx, userKey(email), string(out), 0).Err(); err != nil {
		return fmt.Errorf("redis set user config: %w", err)
	}

	logger.Info("Upstash user config backup saved", map[string]any{
		"email": email,
		"key":   userKey(email),
	})
	return nil
}

// RestoreUserConfig returns the stored user config snapshot and its saved-at time.
func (c *UpstashClient) RestoreUserConfig(ctx context.Context, email string) (map[string]any, error) {
	if c == nil || c.client == nil {
		return nil, fmt.Errorf("upstash backup not initialized")
	}

	val, err := c.client.Get(ctx, userKey(email)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redis get user config: %w", err)
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(val), &m); err != nil {
		return nil, fmt.Errorf("unmarshal user config: %w", err)
	}
	return m, nil
}

// UserSyncStatus returns the last saved timestamp for the user's config.
func (c *UpstashClient) UserSyncStatus(ctx context.Context, email string) (map[string]any, error) {
	if c == nil || c.client == nil {
		return nil, fmt.Errorf("upstash backup not initialized")
	}

	val, err := c.client.Get(ctx, userKey(email)).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get user sync status: %w", err)
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(val), &m); err != nil {
		return nil, fmt.Errorf("unmarshal user sync status: %w", err)
	}
	return map[string]any{"lastSyncAt": m["_savedAt"], "lastSyncError": nil}, nil
}
