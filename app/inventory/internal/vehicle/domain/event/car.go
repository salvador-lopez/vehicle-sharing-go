package event

import "time"

type CarCreatedPayload struct {
	VinNumber string `gorm:"type:varchar(255);unique"`
	Color     string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
