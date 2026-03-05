package logger

import "log/slog"

// From translates a log level given by the user to a respective slog.Level
func From(level string) slog.Level {
	l := slog.LevelDebug

	switch level {
	case "info":
		l = slog.LevelInfo
	case "warning":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}

	return l
}
