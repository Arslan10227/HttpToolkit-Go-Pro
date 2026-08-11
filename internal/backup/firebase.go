package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	firebase "firebase.google.com/go"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/logger"
	"google.golang.org/api/option"
)

// Client backs up user data to Firebase Cloud Firestore using the Firebase
// Admin SDK. It is intentionally isolated so that if the SDK cannot initialize
// (e.g. no service account) the rest of the server still works.
type Client struct {
	projectID string
	app       *firebase.App
	machineID string
}

// NewFirebase creates a Firebase backup client for the configured project.
// It returns nil if Firebase is disabled (no project ID) or if the SDK cannot
// initialize without surfacing a fatal error.
func NewFirebase(cfg *config.Config) *Client {
	if cfg.FirebaseProjectID == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conf := &firebase.Config{
		ProjectID: cfg.FirebaseProjectID,
	}

	var opts []option.ClientOption
	// Prefer an explicit service account key in the config directory, then
	// fall back to GOOGLE_APPLICATION_CREDENTIALS / gcloud default credentials.
	credPath := filepath.Join(cfg.ConfigDir, "firebase-service-account.json")
	if _, err := os.Stat(credPath); err == nil {
		opts = append(opts, option.WithCredentialsFile(credPath))
	}

	// The SDK will use GOOGLE_APPLICATION_CREDENTIALS, gcloud default
	// credentials, or credentials loaded by the running environment.
	app, err := firebase.NewApp(ctx, conf, opts...)
	if err != nil {
		logger.Error(fmt.Errorf("firebase init: %w", err), map[string]any{
			"projectId": cfg.FirebaseProjectID,
		})
		return nil
	}

	return &Client{
		projectID: cfg.FirebaseProjectID,
		app:       app,
		machineID: machineID(cfg),
	}
}

// machineID returns a stable, per-machine identifier. It prefers a saved
// machine id file in the config directory, and falls back to the hostname.
func machineID(cfg *config.Config) string {
	idPath := cfg.ConfigDir + "/machine-id"
	if data, err := os.ReadFile(idPath); err == nil && len(data) > 0 {
		return string(data)
	}
	host, _ := os.Hostname()
	if host == "" {
		host = "unknown"
	}
	_ = os.WriteFile(idPath, []byte(host), 0o600)
	return host
}

// SaveSettings stores arbitrary JSON-able settings under the current machine
// document in the "userSettings" collection. The caller must ensure the
// Firebase service account has write access to Firestore.
func (c *Client) SaveSettings(ctx context.Context, data any) error {
	if c == nil {
		return nil
	}

	db, err := c.app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("firestore client: %w", err)
	}
	defer db.Close()

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}

	m := make(map[string]any)
	if err := json.Unmarshal(payload, &m); err != nil {
		return fmt.Errorf("unmarshal settings: %w", err)
	}
	m["_machineId"] = c.machineID
	m["_updatedAt"] = time.Now().UTC().Format(time.RFC3339)

	_, err = db.Collection("machines").Doc(c.machineID).Set(ctx, m)
	if err != nil {
		return fmt.Errorf("firestore set: %w", err)
	}

	logger.Info("Firebase settings backup saved", map[string]any{
		"projectId": c.projectID,
		"machineId": c.machineID,
	})
	return nil
}

// RestoreSettings attempts to read previously backed-up settings into the
// provided map. It returns an error if the document does not exist or cannot
// be read.
func (c *Client) RestoreSettings(ctx context.Context, out map[string]any) error {
	if c == nil {
		return fmt.Errorf("firebase backup not initialized")
	}

	db, err := c.app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("firestore client: %w", err)
	}
	defer db.Close()

	doc, err := db.Collection("machines").Doc(c.machineID).Get(ctx)
	if err != nil {
		return fmt.Errorf("firestore get: %w", err)
	}

	for k, v := range doc.Data() {
		if k == "_machineId" || k == "_updatedAt" {
			continue
		}
		out[k] = v
	}
	return nil
}

// DefaultCredentialsPath returns the path the SDK is currently using for
// Application Default Credentials, or an empty string if none is set.
func DefaultCredentialsPath() string {
	if p := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); p != "" {
		return p
	}
	return ""
}

// Ensure Option is imported so the firebase.NewApp call can optionally be
// extended with explicit credentials later.
var _ = option.WithCredentialsFile
