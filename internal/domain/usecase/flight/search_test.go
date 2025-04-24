package flight

import (
	"context"
	"testing"
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/entity"
	"github.com/danielmesquitta/flight-api/internal/pkg/validator"
	"github.com/danielmesquitta/flight-api/internal/provider/cache/mockcache"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi/mockflightapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchFlightUseCase_Execute(t *testing.T) {
	type fields struct {
		v validator.Validator
		c *mockcache.MockCache
		f []flightapi.FlightAPI
	}
	type Test struct {
		name    string
		fields  fields
		args    SearchFlightUseCaseInput
		want    *SearchFlightUseCaseOutput
		wantErr bool
	}
	tests := []Test{
		func() Test {
			c := mockcache.NewMockCache(t)
			c.EXPECT().Scan(context.Background(), mock.Anything, mock.Anything).Return(false, nil)
			c.EXPECT().
				Set(context.Background(), mock.Anything, mock.Anything, mock.Anything).
				Return(nil)

			flights := []entity.Flight{
				{
					ID:           "123",
					Origin:       "LAX",
					Destination:  "JFK",
					Price:        100,
					Duration:     int64(time.Hour) * 2,
					FlightNumber: "TX 123",
					DepartureAt:  time.Now(),
					ArrivalAt:    time.Now().Add(time.Hour * 2),
					IsCheapest:   false,
					IsFastest:    false,
				},
			}

			f := mockflightapi.NewMockFlightAPI(t)
			f.EXPECT().
				SearchFlights(context.Background(), "LAX", "JFK", mock.Anything).
				Return(flights, nil)

			wantFlights := make([]entity.Flight, len(flights))
			copy(wantFlights, flights)

			wantFlights[0].IsCheapest = true
			wantFlights[0].IsFastest = true

			return Test{
				name: "searches for flights",
				fields: fields{
					v: validator.New(),
					c: c,
					f: []flightapi.FlightAPI{
						f,
					},
				},
				args: SearchFlightUseCaseInput{
					Origin:      "LAX",
					Destination: "JFK",
					Date:        time.Now(),
					SortBy:      "price",
					SortOrder:   "asc",
				},
				want: &SearchFlightUseCaseOutput{
					Data: wantFlights,
				},
				wantErr: false,
			}
		}(),
		func() Test {
			c := mockcache.NewMockCache(t)
			c.EXPECT().Scan(context.Background(), mock.Anything, mock.Anything).Return(false, nil)
			c.EXPECT().
				Set(context.Background(), mock.Anything, mock.Anything, mock.Anything).
				Return(nil)

			flight1 := entity.Flight{
				ID:           "456",
				Origin:       "LAX",
				Destination:  "JFK",
				Price:        150,
				Duration:     int64(time.Hour) * 2,
				FlightNumber: "TX 456",
				DepartureAt:  time.Now(),
				ArrivalAt:    time.Now().Add(time.Hour * 3),
				IsCheapest:   false,
				IsFastest:    false,
			}

			flight2 := entity.Flight{
				ID:           "123",
				Origin:       "LAX",
				Destination:  "JFK",
				Price:        100,
				Duration:     int64(time.Hour) * 3,
				FlightNumber: "TX 123",
				DepartureAt:  time.Now(),
				ArrivalAt:    time.Now().Add(time.Hour * 2),
				IsCheapest:   false,
				IsFastest:    false,
			}

			flights := []entity.Flight{
				flight1,
				flight2,
			}

			f := mockflightapi.NewMockFlightAPI(t)
			f.EXPECT().
				SearchFlights(context.Background(), "LAX", "JFK", mock.Anything).
				Return(flights, nil)

			flight2.IsCheapest = true
			flight1.IsFastest = true

			wantFlights := []entity.Flight{
				flight2,
				flight1,
			}

			return Test{
				name: "order flights by price",
				fields: fields{
					v: validator.New(),
					c: c,
					f: []flightapi.FlightAPI{
						f,
					},
				},
				args: SearchFlightUseCaseInput{
					Origin:      "LAX",
					Destination: "JFK",
					Date:        time.Now(),
					SortBy:      "price",
					SortOrder:   "asc",
				},
				want: &SearchFlightUseCaseOutput{
					Data: wantFlights,
				},
				wantErr: false,
			}
		}(),
		func() Test {
			c := mockcache.NewMockCache(t)
			c.EXPECT().Scan(context.Background(), mock.Anything, mock.Anything).Return(false, nil)
			c.EXPECT().
				Set(context.Background(), mock.Anything, mock.Anything, mock.Anything).
				Return(nil)

			flight1 := entity.Flight{
				ID:           "456",
				Origin:       "LAX",
				Destination:  "JFK",
				Price:        150,
				Duration:     int64(time.Hour) * 2,
				FlightNumber: "TX 456",
				DepartureAt:  time.Now(),
				ArrivalAt:    time.Now().Add(time.Hour * 3),
				IsCheapest:   false,
				IsFastest:    false,
			}

			flight2 := entity.Flight{
				ID:           "123",
				Origin:       "LAX",
				Destination:  "JFK",
				Price:        100,
				Duration:     int64(time.Hour) * 3,
				FlightNumber: "TX 123",
				DepartureAt:  time.Now(),
				ArrivalAt:    time.Now().Add(time.Hour * 2),
				IsCheapest:   false,
				IsFastest:    false,
			}

			flights := []entity.Flight{
				flight1,
				flight2,
			}

			f := mockflightapi.NewMockFlightAPI(t)
			f.EXPECT().
				SearchFlights(context.Background(), "LAX", "JFK", mock.Anything).
				Return(flights, nil)

			flight1.IsFastest = true
			flight2.IsCheapest = true

			wantFlights := []entity.Flight{
				flight1,
				flight2,
			}

			return Test{
				name: "order flights by duration",
				fields: fields{
					v: validator.New(),
					c: c,
					f: []flightapi.FlightAPI{
						f,
					},
				},
				args: SearchFlightUseCaseInput{
					Origin:      "LAX",
					Destination: "JFK",
					Date:        time.Now(),
					SortBy:      "duration",
					SortOrder:   "asc",
				},
				want: &SearchFlightUseCaseOutput{
					Data: wantFlights,
				},
				wantErr: false,
			}
		}(),
		func() Test {
			c := mockcache.NewMockCache(t)
			c.EXPECT().Scan(context.Background(), mock.Anything, mock.Anything).Return(false, nil)
			c.EXPECT().
				Set(context.Background(), mock.Anything, mock.Anything, mock.Anything).
				Return(nil)

			flight1 := entity.Flight{
				ID:           "456",
				Origin:       "LAX",
				Destination:  "JFK",
				Price:        150,
				Duration:     int64(time.Hour) * 2,
				FlightNumber: "TX 456",
				DepartureAt:  time.Now().Add(time.Hour * 1),
				ArrivalAt:    time.Now().Add(time.Hour * 3),
				IsCheapest:   false,
				IsFastest:    false,
			}

			flight2 := entity.Flight{
				ID:           "123",
				Origin:       "LAX",
				Destination:  "JFK",
				Price:        100,
				Duration:     int64(time.Hour) * 3,
				FlightNumber: "TX 123",
				DepartureAt:  time.Now(),
				ArrivalAt:    time.Now().Add(time.Hour * 2),
				IsCheapest:   false,
				IsFastest:    false,
			}

			flights := []entity.Flight{
				flight1,
				flight2,
			}

			f := mockflightapi.NewMockFlightAPI(t)
			f.EXPECT().
				SearchFlights(context.Background(), "LAX", "JFK", mock.Anything).
				Return(flights, nil)

			flight1.IsFastest = true
			flight2.IsCheapest = true

			wantFlights := []entity.Flight{
				flight2,
				flight1,
			}

			return Test{
				name: "order flights by departure",
				fields: fields{
					v: validator.New(),
					c: c,
					f: []flightapi.FlightAPI{
						f,
					},
				},
				args: SearchFlightUseCaseInput{
					Origin:      "LAX",
					Destination: "JFK",
					Date:        time.Now(),
					SortBy:      "departure",
					SortOrder:   "asc",
				},
				want: &SearchFlightUseCaseOutput{
					Data: wantFlights,
				},
				wantErr: false,
			}
		}(),
		func() Test {
			c := mockcache.NewMockCache(t)
			c.EXPECT().Scan(context.Background(), mock.Anything, mock.Anything).Return(true, nil)
			f := mockflightapi.NewMockFlightAPI(t)

			return Test{
				name: "searches for cached flights",
				fields: fields{
					v: validator.New(),
					c: c,
					f: []flightapi.FlightAPI{
						f,
					},
				},
				args: SearchFlightUseCaseInput{
					Origin:      "LAX",
					Destination: "JFK",
					Date:        time.Now(),
					SortBy:      "price",
					SortOrder:   "asc",
				},
				want:    &SearchFlightUseCaseOutput{},
				wantErr: false,
			}
		}(),
		func() Test {
			c := mockcache.NewMockCache(t)
			c.EXPECT().Scan(context.Background(), mock.Anything, mock.Anything).Return(false, nil)

			flights := []entity.Flight{}

			f := mockflightapi.NewMockFlightAPI(t)
			f.EXPECT().
				SearchFlights(context.Background(), "LAX", "JFK", mock.Anything).
				Return(flights, nil)

			return Test{
				name: "fails to find flights",
				fields: fields{
					v: validator.New(),
					c: c,
					f: []flightapi.FlightAPI{
						f,
					},
				},
				args: SearchFlightUseCaseInput{
					Origin:      "LAX",
					Destination: "JFK",
					Date:        time.Now(),
					SortBy:      "price",
					SortOrder:   "asc",
				},
				want:    nil,
				wantErr: true,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SearchFlightUseCase{
				v: tt.fields.v,
				c: tt.fields.c,
				f: tt.fields.f,
			}

			got, err := s.Execute(context.Background(), tt.args)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, got)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, len(tt.want.Data), len(got.Data))
			assert.Equal(t, tt.want.Data, got.Data)
		})
	}
}
