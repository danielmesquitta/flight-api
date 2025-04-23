package amadeusapi

import (
	"resty.dev/v3"

	"github.com/danielmesquitta/flight-api/internal/config/env"
)

type AmadeusAPI struct {
	e *env.Env
	c *resty.Client
}

func NewAmadeusAPI(e *env.Env) *AmadeusAPI {
	c := resty.New().
		SetBaseURL("https://test.api.amadeus.com")

	return &AmadeusAPI{
		e: e,
		c: c,
	}
}
