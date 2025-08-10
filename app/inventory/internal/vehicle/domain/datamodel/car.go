package datamodel

import (
	"vehicle-sharing-go/pkg/domain/datamodel"
)

type Car struct {
	VinNumber string `gorm:"type:varchar(255);unique"`
	Color     string `gorm:"type:varchar(255)"`
	*datamodel.AggregateRoot
}
