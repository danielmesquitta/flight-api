package handler

import (
	"log/slog"

	"github.com/danielmesquitta/flight-api/internal/app/server/dto"
	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/danielmesquitta/flight-api/internal/domain/usecase/flight"
	"github.com/gofiber/fiber/v2"
)

type FlightHandler struct {
	sfuc *flight.SearchFlightUseCase
}

func NewFlightHandler(
	sfuc *flight.SearchFlightUseCase,
) *FlightHandler {
	return &FlightHandler{
		sfuc: sfuc,
	}
}

// @Summary Flight search
// @Description Search for flights based on origin, destination, and date
// @Tags Flight
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param origin query string true "Origin airport code"
// @Param destination query string true "Destination airport code"
// @Param date query string true "Departure date (YYYY-MM-DD)"
// @Param sort_by query string false "Sort by field (price or duration)"
// @Param sort_order query string false "Sort order (asc or desc)"
// @Success 200 {object} dto.FlightSearchResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/flights/search [get]
func (h *FlightHandler) Search(c *fiber.Ctx) error {
	slog.Info("Searching for flights",
		slog.String("origin", c.Query(QueryParamOrigin)),
		slog.String("destination", c.Query(QueryParamDestination)),
		slog.String("date", c.Query(QueryParamDate)),
		slog.String("sort_by", c.Query(QueryParamSortBy)),
		slog.String("sort_order", c.Query(QueryParamSortOrder)),
	)

	origin := c.Query(QueryParamOrigin)
	destination := c.Query(QueryParamDestination)
	date, err := parseDateQueryParam(c, QueryParamDate)
	if err != nil {
		return errs.New(err)
	}
	sortBy := c.Query(QueryParamSortBy)
	sortOrder := c.Query(QueryParamSortOrder)

	in := flight.SearchFlightUseCaseInput{
		Origin:      origin,
		Destination: destination,
		Date:        date,
		SortBy:      sortBy,
		SortOrder:   sortOrder,
	}

	out, err := h.sfuc.Execute(c.UserContext(), in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.FlightSearchResponse{
		SearchFlightUseCaseOutput: out,
	})
}
