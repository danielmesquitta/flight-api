package auth

import (
	"context"
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/danielmesquitta/flight-api/internal/pkg/jwtutil"
)

type LoginUseCase struct {
	j *jwtutil.JWT
}

func NewLoginUseCase(
	j *jwtutil.JWT,
) *LoginUseCase {
	return &LoginUseCase{
		j: j,
	}
}

type LoginUseCaseInput struct {
	Email    string `json:"username" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUseCaseOutput struct {
	AccessToken string `json:"access_token"`
}

func (l *LoginUseCase) Execute(
	ctx context.Context,
	in LoginUseCaseInput,
) (*LoginUseCaseOutput, error) {
	in7Days := time.Now().Add(time.Hour * 24 * 7)
	tokenClaims := jwtutil.UserClaims{
		Issuer:    in.Email,
		IssuedAt:  time.Now(),
		ExpiresAt: in7Days,
	}

	accessToken, err := l.j.NewToken(tokenClaims, jwtutil.TokenTypeAccess)
	if err != nil {
		return nil, errs.New(err)
	}

	return &LoginUseCaseOutput{AccessToken: accessToken}, nil
}
