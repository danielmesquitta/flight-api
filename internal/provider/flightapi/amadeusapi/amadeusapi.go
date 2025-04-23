package amadeusapi

import (
	"resty.dev/v3"

	"github.com/danielmesquitta/flight-api/internal/config/env"
)

type AmadeusAPI struct {
	e *env.Env
	c *resty.Client
}

func NewAmadeusAPI(env *env.Env) *AmadeusAPI {
	client := resty.New().
		SetBaseURL("https://test.api.amadeus.com")

	return &AmadeusAPI{
		e: env,
		c: client,
	}
}
