package domain

import (
	"time"

	"github.com/google/uuid"
)

type AggregateRoot struct {
	*BaseEntity
	recordedEvents []*Event
}

func NewAggregateRoot(id uuid.UUID, nowFun func() time.Time) *AggregateRoot {
	return &AggregateRoot{BaseEntity: &BaseEntity{id: id, createdAt: nowFun(), updatedAt: nowFun()}}
}

func (a *AggregateRoot) RecordedEvents() []*Event {
	return a.recordedEvents
}

func (a *AggregateRoot) RecordEvent(eventID uuid.UUID, eventType, aggregateType string, payload any, timestamp time.Time) {
	a.recordedEvents = append(
		a.recordedEvents,
		NewEvent(
			eventID,
			a.id,
			aggregateType,
			eventType,
			payload,
			timestamp,
		),
	)
}

func (a *AggregateRoot) ToDTO() *AgRootDTO {
	return &AgRootDTO{
		ID:        a.id,
		CreatedAt: a.createdAt,
		UpdatedAt: a.updatedAt,
	}
}

type AgRootDTO struct {
	ID             uuid.UUID `gorm:"<-:create;type:varchar(36)"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RecordedEvents []*EventDTO `gorm:"-"`
}

func (dto AgRootDTO) ToAggRoot() *AggregateRoot {
	var recordedEvents []*Event
	for _, evtDTO := range dto.RecordedEvents {
		recordedEvents = append(recordedEvents, evtDTO.ToEvent())
	}

	return &AggregateRoot{
		BaseEntity:     &BaseEntity{id: dto.ID, createdAt: dto.CreatedAt, updatedAt: dto.UpdatedAt},
		recordedEvents: recordedEvents,
	}
}
