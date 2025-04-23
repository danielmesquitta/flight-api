package mockflightapi

import "github.com/danielmesquitta/flight-api/internal/provider/flightapi"

func NewMockFlightAPIs(
	m *MockFlightAPI,
) []flightapi.FlightAPI {
	return []flightapi.FlightAPI{
		m,
	}
}
