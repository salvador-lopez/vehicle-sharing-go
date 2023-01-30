package model

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `gorm:"<-:create;type:varchar(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Car struct {
	VIN   string `gorm:"type:varchar(255);unique"`
	Color string `gorm:"type:varchar(255)"`
	Base
}
