package flight

import (
	"context"
	"log/slog"
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/entity"
	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/danielmesquitta/flight-api/internal/pkg/validator"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi"
	"golang.org/x/sync/errgroup"
)

type SearchFlightUseCase struct {
	v validator.Validator
	f []flightapi.FlightAPI
}

func NewSearchFlightUseCase(
	v validator.Validator,
	f []flightapi.FlightAPI,
) *SearchFlightUseCase {
	return &SearchFlightUseCase{
		v: v,
		f: f,
	}
}

type SearchFlightUseCaseInput struct {
	Origin      string `json:"origin"      validate:"required,length=3"`
	Destination string `json:"destination" validate:"required,length=3"`
	Date        string `json:"date"        validate:"required,date=2006-01-02"`
}

type SearchFlightUseCaseOutput struct {
	CheapestFlight  *entity.Flight  `json:"cheapest_flight,omitzero"`
	FastestFlight   *entity.Flight  `json:"fastest_flight,omitzero"`
	PriceComparison []entity.Flight `json:"price_comparison,omitzero"`
}

func (s *SearchFlightUseCase) Execute(
	ctx context.Context,
	in SearchFlightUseCaseInput,
) (*SearchFlightUseCaseOutput, error) {
	if err := s.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	parsedDate, err := time.Parse(time.DateOnly, in.Date)
	if err != nil {
		return nil, errs.New(err)
	}

	allFlights := []entity.Flight{}
	g := errgroup.Group{}
	for _, api := range s.f {
		g.Go(func() error {
			flights, err := api.GetFlightDetails(
				ctx,
				in.Origin,
				in.Destination,
				parsedDate,
			)
			if err != nil {
				return err
			}
			allFlights = append(allFlights, flights...)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		slog.ErrorContext(
			ctx,
			"one or more api calls failed while searching for allFlights",
			"error", err,
		)
	}

	if len(allFlights) == 0 {
		return nil, errs.ErrFlightSearchNotFound
	}

	var cheapest, fastest *entity.Flight
	for _, flight := range allFlights {
		if cheapest == nil || flight.Price < cheapest.Price {
			cheapest = &flight
		}

		duration := flight.ArrivalAt.Sub(flight.DepartureAt)
		cmpDuration := fastest.ArrivalAt.Sub(fastest.DepartureAt)
		if fastest == nil || duration < cmpDuration {
			fastest = &flight
		}
	}

	out := SearchFlightUseCaseOutput{
		CheapestFlight:  cheapest,
		FastestFlight:   fastest,
		PriceComparison: allFlights,
	}

	return &out, nil
}
