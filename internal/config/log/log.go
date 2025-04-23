package log

import (
	"log/slog"
	"os"

	"github.com/danielmesquitta/flight-api/internal/config/env"
)

func SetDefaultLogger(
	e *env.Env,
) {
	switch e.Environment {
	case env.EnvironmentProduction, env.EnvironmentStaging:
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	default:
	}
}
