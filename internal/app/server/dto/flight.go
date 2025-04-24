package dto

import "github.com/danielmesquitta/flight-api/internal/domain/usecase/flight"

type FlightSearchResponse struct {
	*flight.SearchFlightUseCaseOutput
}
