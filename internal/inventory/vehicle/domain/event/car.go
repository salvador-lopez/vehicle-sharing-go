package event

import "time"

type CarCreatedPayload struct {
	VinNumber string
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
