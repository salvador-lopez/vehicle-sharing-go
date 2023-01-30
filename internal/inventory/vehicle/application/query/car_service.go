package query

import "github.com/google/uuid"

//go:generate mockgen -destination=mock/car_service_mock.go -package=mock . CarService
type CarService interface {
	Find(id uuid.UUID) Car
}

type Car struct {
	ID            uuid.UUID
	VIN           string
	Country       string
	Manufacturer  string
	Brand         string
	EngineSize    string
	EngineType    string
	Model         string
	Year          string
	AssemblyPlant string
	SN            string
}
