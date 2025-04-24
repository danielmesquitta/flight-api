package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/flight-api/internal/app/server/dto"
	"github.com/danielmesquitta/flight-api/internal/domain/usecase/auth"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description  string
		body         *dto.LoginRequest
		expectedCode int
	}{
		{
			description: "signs in",
			body: &dto.LoginRequest{
				LoginUseCaseInput: &auth.LoginUseCaseInput{
					Email:    "johndoe@email.com",
					Password: "P@ssw0rd",
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			description: "fails with invalid email",
			body: &dto.LoginRequest{
				LoginUseCaseInput: &auth.LoginUseCaseInput{
					Email:    "invalidemail.com",
					Password: "P@ssw0rd",
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			description: "fails without email",
			body: &dto.LoginRequest{
				LoginUseCaseInput: &auth.LoginUseCaseInput{
					Email:    "",
					Password: "P@ssw0rd",
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			description: "fails without password",
			body: &dto.LoginRequest{
				LoginUseCaseInput: &auth.LoginUseCaseInput{
					Email:    "johndoe@email.com",
					Password: "",
				},
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			app, cleanUp := NewTestApp(t)
			defer func() {
				err := cleanUp(context.Background())
				assert.Nil(t, err)
			}()

			var actual dto.LoginResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodPost,
				"/api/v1/auth/login",
				WithBody(test.body),
				WithResponse(&actual),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
				rawBody,
			)

			if test.expectedCode != http.StatusOK {
				return
			}

			assert.NotEmpty(t, actual.AccessToken)
		})
	}
}
