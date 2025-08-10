package domain

import (
	"vehicle-sharing-go/pkg/domain/datamodel"
	"vehicle-sharing-go/pkg/domain/event"
)

func ToDataModel(a *AggregateRoot) *datamodel.AggregateRoot {
	return &datamodel.AggregateRoot{
		ID:        a.id,
		CreatedAt: a.createdAt,
		UpdatedAt: a.updatedAt,
	}
}

func FromModel(aggRootModel *datamodel.AggregateRoot) *AggregateRoot {
	var recordedEvents []*event.Event
	for _, evt := range aggRootModel.RecordedEvents {
		recordedEvents = append(recordedEvents, evt)
	}

	return &AggregateRoot{
		BaseEntity:     &BaseEntity{id: aggRootModel.ID, createdAt: aggRootModel.CreatedAt, updatedAt: aggRootModel.UpdatedAt},
		recordedEvents: recordedEvents,
	}
}
