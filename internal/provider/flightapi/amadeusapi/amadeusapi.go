package amadeusapi

import (
	"context"
	"time"

	"resty.dev/v3"

	"github.com/danielmesquitta/flight-api/internal/config/env"
	"github.com/danielmesquitta/flight-api/internal/domain/entity"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi"
)

type AmadeusAPI struct {
	e *env.Env
	c *resty.Client

	expiresAt time.Time
}

func NewAmadeusAPI(env *env.Env) *AmadeusAPI {
	client := resty.New().
		SetBaseURL("https://test.api.amadeus.com")

	return &AmadeusAPI{
		e:         env,
		c:         client,
		expiresAt: time.Time{},
	}
}

func (a *AmadeusAPI) GetFlightDetails(
	ctx context.Context,
	origin, destination string,
	date time.Time,
) ([]entity.Flight, error) {
	if err := a.refreshAccessToken(ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

var _ flightapi.FlightAPI = (*AmadeusAPI)(nil)
