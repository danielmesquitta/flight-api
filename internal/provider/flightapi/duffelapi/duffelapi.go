package duffelapi

import (
	"github.com/danielmesquitta/flight-api/internal/config/env"
	"resty.dev/v3"
)

type DuffelAPI struct {
	e *env.Env
	c *resty.Client
}

func NewDuffelAPI(
	e *env.Env,
) *DuffelAPI {
	c := resty.New().
		SetBaseURL("https://api.duffel.com").
		SetHeaders(map[string]string{
			"Authorization":  "Bearer " + e.DuffelAPIKey,
			"Duffel-Version": "v2",
		})

	return &DuffelAPI{
		e: e,
		c: c,
	}
}
