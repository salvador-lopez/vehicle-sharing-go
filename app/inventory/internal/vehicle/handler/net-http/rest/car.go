package rest

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"vehicle-sharing-go/app/inventory/internal/vehicle/command"
	"vehicle-sharing-go/app/inventory/internal/vehicle/projection"
	"vehicle-sharing-go/pkg/domain"
	"vehicle-sharing-go/pkg/handler/rest"
)

//go:generate mockgen -destination=mock/create_car_command_handler_mock.go -package=mock . CreateCarCommandHandler
type CreateCarCommandHandler interface {
	Handle(ctx context.Context, cmd *command.CreateCar) error
}

//go:generate mockgen -destination=mock/find_car_query_service_mock.go -package=mock . FindCarQueryService
type FindCarQueryService interface {
	Find(ctx context.Context, id uuid.UUID) (*projection.Car, error)
}

type CarHandler struct {
	commandHandler CreateCarCommandHandler
	queryService   FindCarQueryService
}

func NewCarHandler(ch CreateCarCommandHandler, qs FindCarQueryService) *CarHandler {
	return &CarHandler{commandHandler: ch, queryService: qs}
}

func (h *CarHandler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	carID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(rest.NewBadRequest(err))
		return
	}

	carProjection, err := h.queryService.Find(ctx, carID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(rest.NewInternalError())
		return
	}

	if carProjection == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(rest.NewNotFound(carID))
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(carProjection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(rest.NewInternalError())
		return
	}
}

func (h *CarHandler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createCarCommand command.CreateCar
	err := json.NewDecoder(r.Body).Decode(&createCarCommand)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(rest.NewBadRequest(err))
		return
	}

	err = h.commandHandler.Handle(ctx, &createCarCommand)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(rest.NewDomainConflict(err))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(rest.NewInternalError())
	}
	w.WriteHeader(http.StatusCreated)
}
