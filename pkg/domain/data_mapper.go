package domain

import "vehicle-sharing-go/pkg/domain/datamodel"

func ToDataModel(a *AggregateRoot) *datamodel.AggregateRoot {
	return &datamodel.AggregateRoot{
		ID:        a.id,
		CreatedAt: a.createdAt,
		UpdatedAt: a.updatedAt,
	}
}

func FromDataModel(dataModel *datamodel.AggregateRoot) *AggregateRoot {
	return &AggregateRoot{
		BaseEntity: &BaseEntity{id: dataModel.ID, createdAt: dataModel.CreatedAt, updatedAt: dataModel.UpdatedAt},
	}
}
