package handler

import (
	"github.com/danielmesquitta/flight-api/internal/app/server/dto"
	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/danielmesquitta/flight-api/internal/domain/usecase/auth"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	luc *auth.LoginUseCase
}

func NewAuthHandler(
	luc *auth.LoginUseCase,
) *AuthHandler {
	return &AuthHandler{
		luc: luc,
	}
}

// @Summary Login
// @Description Use e-mail and password to login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Request body"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := dto.LoginRequest{}
	if err := c.BodyParser(&req); err != nil {
		return errs.New(err)
	}

	out, err := h.luc.Execute(c.UserContext(), *req.LoginUseCaseInput)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.LoginResponse{
		LoginUseCaseOutput: out,
	})
}
