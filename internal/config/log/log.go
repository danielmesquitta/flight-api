package log

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/danielmesquitta/flight-api/internal/config/env"
	"github.com/fatih/color"
)

type prettyHandler struct {
	slog.Handler
	l *log.Logger
}

func newPrettyHandler(
	out io.Writer,
) *prettyHandler {
	return &prettyHandler{
		Handler: slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
		l: log.New(out, "", 0),
	}
}

func (h *prettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()

	var colorFunc func(format string, a ...interface{}) string
	switch r.Level {
	case slog.LevelDebug:
		colorFunc = color.CyanString
	case slog.LevelInfo:
		colorFunc = color.BlueString
	case slog.LevelWarn:
		colorFunc = color.YellowString
	case slog.LevelError:
		colorFunc = color.RedString
	}

	fields := make(map[string]any, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	level = colorFunc(level)
	timeStr := colorFunc(r.Time.Format(time.RFC3339))
	msg := r.Message

	attrs := string(b)
	if attrs == "{}" {
		attrs = ""
	}

	h.l.Printf("%s %s %s %s", timeStr, level, msg, attrs)

	return nil
}

func SetDefaultLogger(
	e *env.Env,
) {
	var logger *slog.Logger

	switch e.Environment {
	case env.EnvironmentProduction, env.EnvironmentStaging:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	default:
		logger = slog.New(newPrettyHandler(os.Stdout))
	}

	slog.SetDefault(logger)
}
