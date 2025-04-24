package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/danielmesquitta/flight-api/internal/app/server/dto"
	"github.com/danielmesquitta/flight-api/internal/app/server/handler"
	"github.com/stretchr/testify/assert"
)

func TestSearchFlights(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description  string
		queryParams  map[string]string
		isLoggedIn   bool
		expectedCode int
	}{
		{
			description:  "fails without token",
			queryParams:  map[string]string{},
			isLoggedIn:   false,
			expectedCode: http.StatusBadRequest,
		},
		{
			description: "searches for flights",
			queryParams: map[string]string{
				handler.QueryParamOrigin:      "SYD",
				handler.QueryParamDestination: "BKK",
				handler.QueryParamDate: time.Now().
					AddDate(0, 3, 0).
					Format(time.DateOnly),
			},
			isLoggedIn:   true,
			expectedCode: http.StatusOK,
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

			var out dto.SearchFlightsResponse
			opts := []RequestOption{
				WithQueryParams(test.queryParams),
				WithResponse(&out),
			}

			if test.isLoggedIn {
				loginRes := app.Login("johndoe@email.com", "P@ssw0rd")
				opts = append(opts, WithBearerToken(loginRes.AccessToken))
			}

			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/flights/search",
				opts...,
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

			assert.Greater(t, len(out.Data), 0)
			assert.NotEmpty(t, out.Data[0].ID)
			assert.NotEmpty(t, out.Data[0].FlightNumber)
			assert.NotEmpty(t, out.Data[0].Duration)
			assert.NotEmpty(t, out.Data[0].Price)
		})
	}
}
