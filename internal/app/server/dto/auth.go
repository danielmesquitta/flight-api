package dto

import (
	"github.com/danielmesquitta/flight-api/internal/domain/usecase/auth"
)

type LoginResponse struct {
	*auth.LoginUseCaseOutput
}

type LoginRequest struct {
	*auth.LoginUseCaseInput
}
