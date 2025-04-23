package amadeusapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/entity"
	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/itlightning/dateparse"
)

type SearchFlightsResponse struct {
	Data []SearchFlightsResponseData `json:"data"`
}

type SearchFlightsResponseData struct {
	ID          string                           `json:"id"`
	Itineraries []SearchFlightsResponseItinerary `json:"itineraries"`
	Price       SearchFlightsResponseDataPrice   `json:"price"`
}

type SearchFlightsResponseItinerary struct {
	Duration string                         `json:"duration"`
	Segments []SearchFlightsResponseSegment `json:"segments"`
}

type SearchFlightsResponseSegment struct {
	Departure SearchFlightsResponseArrival `json:"departure"`
	Arrival   SearchFlightsResponseArrival `json:"arrival"`
}

type SearchFlightsResponseArrival struct {
	At string `json:"at"`
}

type SearchFlightsResponseDataPrice struct {
	GrandTotal string `json:"grandTotal"`
}

func (a *AmadeusAPI) SearchFlights(
	ctx context.Context,
	origin, destination string,
	date time.Time,
) ([]entity.Flight, error) {
	if err := a.refreshAccessToken(ctx); err != nil {
		return nil, err
	}

	res, err := a.c.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"originLocationCode":      origin,
			"destinationLocationCode": destination,
			"departureDate":           date.Format(time.DateOnly),
			"adults":                  "1",
		}).
		Get("/v2/shopping/flight-offers")
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

	flights := make([]entity.Flight, 0, len(data.Data))
	for _, flight := range data.Data {
		itinerary := flight.Itineraries[0]

		firstSegment := itinerary.Segments[0]
		departureAt, err := dateparse.ParseAny(firstSegment.Departure.At)
		if err != nil {
			slog.ErrorContext(ctx, "failed to parse departure date", "error", err)
			continue
		}

		lastSegment := itinerary.Segments[len(itinerary.Segments)-1]
		arrivalAt, err := dateparse.ParseAny(lastSegment.Arrival.At)
		if err != nil {
			slog.ErrorContext(ctx, "failed to parse arrival date", "error", err)
			continue
		}

		duration, err := a.parseDuration(itinerary.Duration)
		if err != nil {
			slog.ErrorContext(ctx, "failed to parse duration", "error", err)
			continue
		}

		price, err := strconv.ParseFloat(flight.Price.GrandTotal, 64)
		if err != nil {
			return nil, errs.New(err)
		}

		flightData := entity.Flight{
			Origin:      origin,
			Destination: destination,
			DepartureAt: departureAt,
			ArrivalAt:   arrivalAt,
			Duration:    duration,
			Price:       int64(price * 100),
		}

		flights = append(flights, flightData)
	}

	return flights, nil
}

func (a *AmadeusAPI) parseDuration(duration string) (time.Duration, error) {
	split := strings.Split(duration, "PT")
	if len(split) != 2 {
		return 0, errs.New("invalid duration format")
	}

	durationStr := strings.ToLower(split[1])

	parsedDuration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, errs.New(err)
	}

	return parsedDuration, nil
}
