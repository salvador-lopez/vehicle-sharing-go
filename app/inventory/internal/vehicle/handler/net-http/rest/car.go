package rest

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"vehicle-sharing-go/app/inventory/internal/vehicle/projection"
)

//go:generate mockgen -destination=mock/find_car_query_service_mock.go -package=mock . FindCarQueryService
type FindCarQueryService interface {
	Find(ctx context.Context, id uuid.UUID) (*projection.Car, error)
}

type CarHandler struct {
	queryService FindCarQueryService
}

func NewCarHandler(qs FindCarQueryService) *CarHandler {
	return &CarHandler{queryService: qs}
}

func (h *CarHandler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	carID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	carProjection, err := h.queryService.Find(ctx, carID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if carProjection == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(carProjection)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
