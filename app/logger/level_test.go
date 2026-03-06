package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevel(t *testing.T) {
	tc := []struct {
		current  string
		expected slog.Level
	}{{
		current:  "info",
		expected: slog.LevelInfo,
	}, {
		current:  "debug",
		expected: slog.LevelDebug,
	}, {
		current:  "warning",
		expected: slog.LevelWarn,
	}, {
		current:  "error",
		expected: slog.LevelError,
	}, {
		current:  "",
		expected: slog.LevelDebug,
	}, {
		current:  "critical",
		expected: slog.LevelDebug,
	}}

	for _, test := range tc {
		assert.Equal(t, test.expected, From(test.current), "current \"%s\", expected %v", test.current, test.expected)
	}
}
