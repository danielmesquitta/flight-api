package container

import (
	"log"
	"strings"
)

type testLogger struct{}

func (t testLogger) Printf(format string, v ...interface{}) {
	if strings.HasPrefix(format, "🐳") ||
		strings.HasPrefix(format, "✅") ||
		strings.HasPrefix(format, "⏳") ||
		strings.HasPrefix(format, "🚫") ||
		strings.HasPrefix(format, "🔔") {
		return
	}

	log.Printf(format, v...)
}

func newLogger() *testLogger {
	return &testLogger{}
}
