package rest

import (
	"context"
	"github.com/gin-gonic/gin"
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

func (h *CarHandler) Get(c *gin.Context) {
	carID := uuid.MustParse(c.Param("id"))
	carProjection, _ := h.queryService.Find(c, carID)

	c.JSON(http.StatusOK, carProjection)
}
