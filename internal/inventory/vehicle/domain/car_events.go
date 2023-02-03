package domain

import "time"

type CarCreatedEventPayload struct {
	vin       *VIN
	color     string
	createdAt time.Time
	updatedAt time.Time
}

type CarCreatedEventPayloadDTO struct {
	VinNumber string
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (dto CarCreatedEventPayloadDTO) ToPayload() any {
	return &CarCreatedEventPayload{
		vin:       &VIN{dto.VinNumber},
		color:     dto.Color,
		createdAt: dto.CreatedAt,
		updatedAt: dto.UpdatedAt,
	}
}
