package log

import (
	"log/slog"
	"os"
)

func SetDefaultLogger() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}
