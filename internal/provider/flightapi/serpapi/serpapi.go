package serpapi

import (
	"resty.dev/v3"

	"github.com/danielmesquitta/flight-api/internal/config/env"
)

type SerpAPI struct {
	e *env.Env
	c *resty.Client
}

func NewSerpAPI(e *env.Env) *SerpAPI {
	c := resty.New().
		SetBaseURL("https://serpapi.com/search").
		SetQueryParams(map[string]string{
			"api_key": e.SerpAPIKey,
			"engine":  "google_flights",
		})

	return &SerpAPI{
		e: e,
		c: c,
	}
}
