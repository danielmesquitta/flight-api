package auth

import (
	"context"
	"testing"

	"github.com/danielmesquitta/flight-api/internal/config"
	"github.com/danielmesquitta/flight-api/internal/pkg/jwtutil"
	"github.com/danielmesquitta/flight-api/internal/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestLoginUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		args    LoginUseCaseInput
		wantErr bool
	}{
		{
			name: "signs in",
			args: LoginUseCaseInput{
				Email:    "johndoe@email.com",
				Password: "P@ssw0rd",
			},
			wantErr: false,
		},
		{
			name: "fails with invalid email",
			args: LoginUseCaseInput{
				Email:    "invalidemail.com",
				Password: "P@ssw0rd",
			},
			wantErr: true,
		},
		{
			name: "fails without email",
			args: LoginUseCaseInput{
				Email:    "",
				Password: "P@ssw0rd",
			},
			wantErr: true,
		},
		{
			name: "fails without password",
			args: LoginUseCaseInput{
				Email:    "",
				Password: "P@ssw0rd",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newLoginUseCase()

			got, err := l.Execute(context.Background(), tt.args)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, got)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.NotEmpty(t, got.AccessToken)
		})
	}
}

func newLoginUseCase() *LoginUseCase {
	v := validator.New()
	e := config.LoadConfig(v)
	return &LoginUseCase{
		v: v,
		j: jwtutil.NewJWT(e),
	}
}
