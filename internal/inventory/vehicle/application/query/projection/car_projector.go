package projection

import (
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -destination=mock/car_projector_mock.go -package=mock . CarProjector
type CarProjector interface {
	Project(car *Car) error
}

type Car struct {
	ID            uuid.UUID `gorm:"<-:create;type:varchar(36)"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	VIN           string `gorm:"type:varchar(255);unique"`
	Country       string `gorm:"type:varchar(255)"`
	Manufacturer  string `gorm:"type:varchar(255)"`
	Brand         string `gorm:"type:varchar(255)"`
	EngineSize    string `gorm:"type:varchar(255)"`
	FuelType      string `gorm:"type:varchar(255)"`
	Model         string `gorm:"type:varchar(255)"`
	Year          string `gorm:"type:varchar(255)"`
	AssemblyPlant string `gorm:"type:varchar(255)"`
	SN            string `gorm:"type:varchar(255)"`
}
