package serpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/entity"
	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/itlightning/dateparse"
)

type SearchFlightsResponse struct {
	BestFlights  []Flight `json:"best_flights"`
	OtherFlights []Flight `json:"other_flights"`
}

type Flight struct {
	Flights       []FlightElement `json:"flights"`
	TotalDuration int64           `json:"total_duration"`
	Price         float64         `json:"price"`
}

type FlightElement struct {
	FlightNumber     string  `json:"flight_number"`
	DepartureAirport Airport `json:"departure_airport"`
	ArrivalAirport   Airport `json:"arrival_airport"`
}

type Airport struct {
	Time string `json:"time"`
}

func (a *SerpAPI) SearchFlights(
	ctx context.Context,
	origin, destination string,
	date time.Time,
) ([]entity.Flight, error) {
	res, err := a.c.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"departure_id":  origin,
			"arrival_id":    destination,
			"outbound_date": date.Format(time.DateOnly),
			// Assuming a 7-day return, since return_date is required
			"return_date": date.AddDate(0, 0, 7).
				Format(time.DateOnly),
		}).
		Get("/")
	if err != nil {
		return nil, errs.New(err)
	}
	body := res.Bytes()
	if res.IsError() {
		return nil, errs.New(string(body))
	}

	data := SearchFlightsResponse{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, errs.New(err)
	}

	allFlights := append(data.BestFlights, data.OtherFlights...)

	flights := make([]entity.Flight, 0, len(allFlights))
	for _, flight := range allFlights {
		firstSegment := flight.Flights[0]
		departureAt, err := dateparse.ParseAny(
			firstSegment.DepartureAirport.Time,
		)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"failed to parse departure date",
				"error",
				err,
			)
			continue
		}

		lastSegment := flight.Flights[len(flight.Flights)-1]
		arrivalAt, err := dateparse.ParseAny(lastSegment.ArrivalAirport.Time)
		if err != nil {
			slog.ErrorContext(ctx, "failed to parse arrival date", "error", err)
			continue
		}

		duration := time.Duration(flight.TotalDuration) * time.Minute

		id := fmt.Sprintf(
			"serp-%s",
			strings.ReplaceAll(
				strings.ToLower(firstSegment.FlightNumber),
				" ",
				"-",
			),
		)

		flightData := entity.Flight{
			ID:           id,
			FlightNumber: firstSegment.FlightNumber,
			Origin:       origin,
			Destination:  destination,
			DepartureAt:  departureAt,
			ArrivalAt:    arrivalAt,
			Duration:     int64(duration),
			Price:        int64(flight.Price * 100),
		}

		flights = append(flights, flightData)
	}

	return flights, nil
}
