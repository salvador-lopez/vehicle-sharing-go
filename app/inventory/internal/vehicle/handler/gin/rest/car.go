package rest

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"goa.design/goa/v3/dsl"
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

func (h *CarHandler) Get(c *gin.Context) {
	carID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, rest.NewBadRequest(err))
		return
	}
	carProjection, err := h.queryService.Find(c, carID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, rest.NewInternalError())
		return
	}

	if carProjection == nil {
		c.JSON(http.StatusNotFound, rest.NewNotFound(carID))
		return
	}

	c.JSON(http.StatusOK, carProjection)
}

func (h *CarHandler) Create(c *gin.Context) {
	var createCarCommand command.CreateCar
	if err := c.ShouldBindJSON(&createCarCommand); err != nil {
		c.JSON(dsl.StatusBadRequest, rest.NewBadRequest(err))
		return
	}

	err := h.commandHandler.Handle(c, &createCarCommand)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			c.JSON(http.StatusConflict, rest.NewDomainConflict(err))
			return
		}

		c.JSON(http.StatusInternalServerError, rest.NewInternalError())
		return
	}

	c.String(http.StatusCreated, "")
}
