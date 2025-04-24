package auth

import (
	"context"
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/danielmesquitta/flight-api/internal/pkg/jwtutil"
	"github.com/danielmesquitta/flight-api/internal/pkg/validator"
)

type LoginUseCase struct {
	v validator.Validator
	j *jwtutil.JWT
}

func NewLoginUseCase(
	v validator.Validator,
	j *jwtutil.JWT,
) *LoginUseCase {
	return &LoginUseCase{
		v: v,
		j: j,
	}
}

type LoginUseCaseInput struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUseCaseOutput struct {
	AccessToken string `json:"access_token"`
}

func (l *LoginUseCase) Execute(
	ctx context.Context,
	in LoginUseCaseInput,
) (*LoginUseCaseOutput, error) {
	if err := l.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

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
