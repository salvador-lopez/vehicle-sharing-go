package model

import (
	"vehicle-sharing-go/pkg/domain/model"
)

type Car struct {
	VinNumber string `gorm:"type:varchar(255);unique"`
	Color     string `gorm:"type:varchar(255)"`
	*model.AggregateRoot
}
