package flight

import (
	"context"
	"log/slog"
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/entity"
	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/danielmesquitta/flight-api/internal/pkg/validator"
	"github.com/danielmesquitta/flight-api/internal/provider/cache"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi"
	"golang.org/x/sync/errgroup"
)

type SearchFlightUseCase struct {
	v validator.Validator
	c cache.Cache
	f []flightapi.FlightAPI
}

func NewSearchFlightUseCase(
	v validator.Validator,
	c cache.Cache,
	f []flightapi.FlightAPI,
) *SearchFlightUseCase {
	return &SearchFlightUseCase{
		v: v,
		c: c,
		f: f,
	}
}

type SearchFlightUseCaseInput struct {
	Origin      string    `json:"origin"      validate:"required,len=3"`
	Destination string    `json:"destination" validate:"required,len=3"`
	Date        time.Time `json:"date"        validate:"required"`
}

type SearchFlightUseCaseOutput struct {
	CheapestFlight *entity.Flight  `json:"cheapest_flight,omitzero"`
	FastestFlight  *entity.Flight  `json:"fastest_flight,omitzero"`
	Flights        []entity.Flight `json:"flights,omitzero"`
}

func (s *SearchFlightUseCase) Execute(
	ctx context.Context,
	in SearchFlightUseCaseInput,
) (*SearchFlightUseCaseOutput, error) {
	if err := s.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	cacheKey := in.Origin + "_" + in.Destination + "_" + in.Date.Format(time.DateOnly)

	out := &SearchFlightUseCaseOutput{}
	ok, err := s.c.Scan(ctx, cacheKey, out)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to scan cache for search flight use case",
			"error", err,
		)
	}
	if ok {
		return out, nil
	}

	allFlights := []entity.Flight{}
	g := errgroup.Group{}
	for _, api := range s.f {
		g.Go(func() error {
			flights, err := api.SearchFlights(
				ctx,
				in.Origin,
				in.Destination,
				in.Date,
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
	for i := range allFlights {
		flight := &allFlights[i]

		if cheapest == nil || flight.Price < cheapest.Price {
			cheapest = flight
		}

		duration := flight.ArrivalAt.Sub(flight.DepartureAt)
		if fastest == nil || duration < fastest.ArrivalAt.Sub(fastest.DepartureAt) {
			fastest = flight
		}
	}

	out = &SearchFlightUseCaseOutput{
		CheapestFlight: cheapest,
		FastestFlight:  fastest,
		Flights:        allFlights,
	}

	if err := s.c.Set(ctx, cacheKey, out, time.Second*30); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to set cache for search flight use case",
			"error", err,
		)
	}

	return out, nil
}
