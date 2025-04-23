package entity

import "time"

type Flight struct {
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	DepartureAt time.Time `json:"departure_at"`
	ArrivalAt   time.Time `json:"arrival_at"`
	Duration    int64     `json:"duration"`
	Price       int64     `json:"price"`
}
