package handler

import (
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
// @Success 200 {object} dto.FlightSearchResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/flights/search [get]
func (h *FlightHandler) Search(c *fiber.Ctx) error {
	origin := c.Query(QueryParamOrigin)
	destination := c.Query(QueryParamDestination)

	date, err := parseDateQueryParam(c, QueryParamDate)
	if err != nil {
		return errs.New(err)
	}

	in := flight.SearchFlightUseCaseInput{
		Origin:      origin,
		Destination: destination,
		Date:        date,
	}

	out, err := h.sfuc.Execute(c.UserContext(), in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.FlightSearchResponse{
		SearchFlightUseCaseOutput: out,
	})
}
