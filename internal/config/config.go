package config

import (
	"github.com/danielmesquitta/flight-api/internal/config/env"
	"github.com/danielmesquitta/flight-api/internal/config/log"
	"github.com/danielmesquitta/flight-api/internal/config/time"
	"github.com/danielmesquitta/flight-api/internal/pkg/validator"
)

func LoadConfig(v validator.Validator) *env.Env {
	e := env.NewEnv(v)
	log.SetDefaultLogger(e)
	time.SetServerTimeZone()

	return e
}
