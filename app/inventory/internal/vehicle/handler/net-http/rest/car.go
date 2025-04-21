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

// Get godoc
// @Summary      Get a car by ID
// @Description  Returns a car resource with decoded VIN data
// @Tags         car
// @Produce      json
// @Param        id   path      string  true  "Car UUID"
// @Success      200  {object}  projection.Car
// @Failure      400  {object}  errorResponse  "bad request"
// @Failure      404  {object}  errorResponse  "not found"
// @Failure      500  {object}  errorResponse  "internal error"
// @Router       /cars/{id} [get]
func (h *CarHandler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	carID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newBadRequest(err))
		return
	}

	carProjection, err := h.queryService.Find(ctx, carID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(newInternalError())
		return
	}

	if carProjection == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(newNotFound(carID))
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(carProjection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(newInternalError())
		return
	}
}
