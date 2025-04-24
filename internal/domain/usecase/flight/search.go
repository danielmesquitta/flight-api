package flight

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/entity"
	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/danielmesquitta/flight-api/internal/pkg/validator"
	"github.com/danielmesquitta/flight-api/internal/provider/cache"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi"
	"golang.org/x/sync/errgroup"
)

type SearchFlightsUseCase struct {
	v validator.Validator
	c cache.Cache
	f []flightapi.FlightAPI
}

func NewSearchFlightsUseCase(
	v validator.Validator,
	c cache.Cache,
	f []flightapi.FlightAPI,
) *SearchFlightsUseCase {
	return &SearchFlightsUseCase{
		v: v,
		c: c,
		f: f,
	}
}

type SearchFlightsUseCaseInput struct {
	Origin      string    `json:"origin"      validate:"required,len=3"`
	Destination string    `json:"destination" validate:"required,len=3"`
	Date        time.Time `json:"date"        validate:"required"`
	SortBy      string    `json:"sort_by"     validate:"omitempty,oneof=price duration departure"`
	SortOrder   string    `json:"sort_order"  validate:"omitempty,oneof=asc desc"`
}

type SearchFlightsUseCaseOutput struct {
	Data []entity.Flight `json:"data,omitzero"`
}

func (s *SearchFlightsUseCase) Execute(
	ctx context.Context,
	in SearchFlightsUseCaseInput,
) (*SearchFlightsUseCaseOutput, error) {
	if err := s.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	in.SortBy = cmp.Or(in.SortBy, "price")
	in.SortOrder = cmp.Or(in.SortOrder, "asc")

	cacheKey := fmt.Sprintf(
		"%s_%s_%s_%s_%s",
		in.Origin,
		in.Destination,
		in.Date.Format(time.DateOnly),
		in.SortBy,
		in.SortOrder,
	)

	out := &SearchFlightsUseCaseOutput{}
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
		return nil, errs.ErrSearchFlightsNotFound
	}

	s.setFastestAndCheapest(allFlights)

	s.sortFlights(allFlights, in.SortBy, in.SortOrder)

	out = &SearchFlightsUseCaseOutput{
		Data: allFlights,
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

func (s *SearchFlightsUseCase) setFastestAndCheapest(
	flights []entity.Flight,
) {
	var cheapest, fastest *entity.Flight
	for i := range flights {
		flight := &flights[i]

		if cheapest == nil || flight.Price < cheapest.Price {
			cheapest = flight
		}

		if fastest == nil || flight.Duration < fastest.Duration {
			fastest = flight
		}
	}

	if cheapest != nil {
		cheapest.IsCheapest = true
	}
	if fastest != nil {
		fastest.IsFastest = true
	}
}

func (s *SearchFlightsUseCase) sortFlights(
	flights []entity.Flight,
	sortBy, sortOrder string,
) {
	sort.Slice(flights, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "duration":
			less = flights[i].Duration < flights[j].Duration

		case "departure":
			less = flights[i].DepartureAt.Before(flights[j].DepartureAt)

		default: // "price"
			less = flights[i].Price < flights[j].Price
		}
		if sortOrder == "desc" {
			return !less
		}
		return less
	})
}
