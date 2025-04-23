package flightapi

import (
	"context"
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/entity"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi/amadeusapi"
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
) []FlightAPI {
	return []FlightAPI{
		a,
	}
}
