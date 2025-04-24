package flightapi

import (
	"context"
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/entity"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi/amadeusapi"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi/duffelapi"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi/serpapi"
)

type FlightAPI interface {
	SearchFlights(
		ctx context.Context,
		origin, destination string,
		date time.Time,
	) ([]entity.Flight, error)
}

func NewFlightAPIs(
	a *amadeusapi.AmadeusAPI,
	s *serpapi.SerpAPI,
	d *duffelapi.DuffelAPI,
) []FlightAPI {
	return []FlightAPI{
		a,
		s,
		d,
	}
}
