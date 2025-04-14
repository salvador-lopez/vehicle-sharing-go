package rest

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"regexp"
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

	carID, err := h.extractIdFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	carProjection, err := h.queryService.Find(ctx, carID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(carProjection)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (h *CarHandler) extractIdFromPath(path string) (uuid.UUID, error) {
	re := regexp.MustCompile(`^/cars/([a-f0-9-]+)$`)

	matches := re.FindStringSubmatch(path)
	if len(matches) != 2 {
		return uuid.UUID{}, errors.New("invalid or missing UUID")
	}

	return uuid.Parse(matches[1])
}
