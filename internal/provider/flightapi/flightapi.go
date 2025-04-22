package flightapi

import (
	"context"
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/entity"
)

type FlightAPI interface {
	GetFlightDetails(
		ctx context.Context,
		origin, destination string,
		date time.Time,
	) ([]entity.Flight, error)
}
