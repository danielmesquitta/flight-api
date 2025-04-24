package duffelapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/entity"
	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/itlightning/dateparse"
)

type SearchFlightsResponse struct {
	Data SearchFlightsData `json:"data"`
}

type SearchFlightsData struct {
	Offers []SearchFlightsOffer `json:"offers"`
}

type SearchFlightsOffer struct {
	TotalAmount string                    `json:"total_amount"`
	Slices      []SearchFlightsOfferSlice `json:"slices"`
}

type SearchFlightsOfferSlice struct {
	Segments []SearchFlightsSegment `json:"segments"`
	Duration string                 `json:"duration"`
}

type SearchFlightsSegment struct {
	DepartingAt                  string                        `json:"departing_at"`
	ArrivingAt                   string                        `json:"arriving_at"`
	MarketingCarrierFlightNumber string                        `json:"marketing_carrier_flight_number"`
	MarketingCarrier             SearchFlightsMarketingCarrier `json:"marketing_carrier"`
}

type SearchFlightsMarketingCarrier struct {
	IataCode string `json:"iata_code"`
}

func (d *DuffelAPI) SearchFlights(
	ctx context.Context,
	origin, destination string,
	date time.Time,
) ([]entity.Flight, error) {
	reqBody, err := d.buildRequestBody(origin, destination, date)
	if err != nil {
		return nil, errs.New(err)
	}

	res, err := d.c.R().
		SetContext(ctx).
		SetBody(reqBody).
		Post("/air/offer_requests")
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

	flights := []entity.Flight{}
	for _, offer := range data.Data.Offers {
		if len(offer.Slices) == 0 {
			continue
		}

		firstOffer := offer.Slices[0]

		if len(firstOffer.Segments) == 0 {
			continue
		}

		firstSegment := firstOffer.Segments[0]
		departureAt, err := dateparse.ParseAny(firstSegment.DepartingAt)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"failed to parse departure date",
				"error",
				err,
			)
			continue
		}

		lastSegment := firstOffer.Segments[len(firstOffer.Segments)-1]
		arrivalAt, err := dateparse.ParseAny(lastSegment.ArrivingAt)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"failed to parse arrival date",
				"error",
				err,
			)
			continue
		}

		duration, err := d.parseDuration(firstOffer.Duration)
		if err != nil {
			slog.ErrorContext(ctx, "failed to parse duration", "error", err)
			continue
		}

		price, err := strconv.ParseFloat(offer.TotalAmount, 64)
		if err != nil {
			return nil, errs.New(err)
		}

		flightNumber := firstSegment.MarketingCarrierFlightNumber
		if firstSegment.MarketingCarrier.IataCode != "" {
			flightNumber = firstSegment.MarketingCarrier.IataCode + " " + flightNumber
		}

		id := fmt.Sprintf(
			"duffel-%s",
			strings.ReplaceAll(strings.ToLower(flightNumber), " ", "-"),
		)

		flights = append(flights, entity.Flight{
			ID:           id,
			FlightNumber: flightNumber,
			Origin:       origin,
			Destination:  destination,
			DepartureAt:  departureAt,
			ArrivalAt:    arrivalAt,
			Duration:     int64(duration),
			Price:        int64(price * 100),
			IsCheapest:   false,
			IsFastest:    false,
		})
	}

	return flights, nil
}

func (d *DuffelAPI) buildRequestBody(
	origin, destination string,
	date time.Time,
) ([]byte, error) {
	type Slice struct {
		Origin        string `json:"origin"`
		Destination   string `json:"destination"`
		DepartureDate string `json:"departure_date"`
	}
	type Passenger struct {
		Type string `json:"type"`
	}
	type Data struct {
		Slices     []Slice     `json:"slices"`
		Passengers []Passenger `json:"passengers"`
	}
	type RequestBody struct {
		Data Data `json:"data"`
	}

	reqBody := RequestBody{
		Data: Data{
			Slices: []Slice{
				{
					Origin:        origin,
					Destination:   destination,
					DepartureDate: date.Format(time.DateOnly),
				},
			},
			Passengers: []Passenger{
				{
					Type: "adult",
				},
			},
		},
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errs.New(err)
	}

	return reqBodyJSON, nil
}

func (d *DuffelAPI) parseDuration(duration string) (time.Duration, error) {
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
