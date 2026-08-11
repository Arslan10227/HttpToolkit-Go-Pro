package admin

import "github.com/Arslan10227/HttpToolkit-Go-Pro/internal/logger"

// Info logs an informational message.
func Info(msg string, fields map[string]any) {
	logger.Info(msg, fields)
}

// Debug logs a debug message.
func Debug(msg string, fields map[string]any) {
	logger.Debug(msg, fields)
}

// Error logs an error message.
func Error(err error, fields map[string]any) {
	logger.Error(err, fields)
}
