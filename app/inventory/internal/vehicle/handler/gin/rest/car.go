package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"vehicle-sharing-go/app/inventory/internal/vehicle/projection"
	"vehicle-sharing-go/pkg/handler/rest"
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
