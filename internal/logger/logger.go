// Package logger provides the application-wide logging API used across
// httptoolkit-go. Internally it is backed by the stdlib log/slog structured
// logger; the package keeps its own simple Info/Debug/Error(map[string]any)
// API so the ~20 call sites across the codebase don't need to change.
package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
)

var (
	std      *slog.Logger
	logFile  *os.File
	levelVar = new(slog.LevelVar) // starts at slog.LevelInfo (verbose=false)
)

func init() {
	levelVar.Set(slog.LevelInfo)

	dir := "logs"
	if execPath, err := os.Executable(); err == nil {
		dir = filepath.Join(filepath.Dir(execPath), "logs")
	}
	_ = os.MkdirAll(dir, 0o755)

	var handler slog.Handler
	if f, err := os.OpenFile(filepath.Join(dir, "httptoolkit.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
		logFile = f
		handler = slog.NewJSONHandler(&syncWriter{f}, &slog.HandlerOptions{Level: levelVar})
	} else {
		// Fallback to stderr if file opening fails.
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: levelVar})
	}
	std = slog.New(handler)
}

// syncWriter fsyncs the underlying file after every write, matching the
// previous hand-rolled logger's per-entry durability (important for
// post-crash diagnostics on a desktop app).
type syncWriter struct{ f *os.File }

func (w *syncWriter) Write(p []byte) (int, error) {
	n, err := w.f.Write(p)
	_ = w.f.Sync()
	return n, err
}

// SetVerbose toggles verbose (debug-level) logging.
func SetVerbose(on bool) {
	if on {
		levelVar.Set(slog.LevelDebug)
	} else {
		levelVar.Set(slog.LevelInfo)
	}
}

// Info logs an informational message.
func Info(msg string, fields map[string]any) {
	log(slog.LevelInfo, msg, fields)
}

// Debug logs a verbose debug message (only emitted when SetVerbose(true)).
func Debug(msg string, fields map[string]any) {
	log(slog.LevelDebug, msg, fields)
}

// Error logs an error message.
func Error(err error, fields map[string]any) {
	if err == nil {
		return
	}
	log(slog.LevelError, err.Error(), fields)
}

func log(level slog.Level, msg string, fields map[string]any) {
	if std == nil || !std.Enabled(context.Background(), level) {
		return
	}
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	std.Log(context.Background(), level, msg, args...)
}

// GetLogPath returns the absolute path to the log file.
func GetLogPath() string {
	if logFile == nil {
		return ""
	}
	dir := "logs"
	if execPath, err := os.Executable(); err == nil {
		dir = filepath.Join(filepath.Dir(execPath), "logs")
	}
	return filepath.Join(dir, "httptoolkit.log")
}
