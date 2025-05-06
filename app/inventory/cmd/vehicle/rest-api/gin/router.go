package gin

import (
	"github.com/gin-gonic/gin"
	"net/url"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/gin/rest"
)

func registerHandlers(
	r *gin.Engine,
	addr *url.URL,
	carHandler *rest.CarHandler,
) {
	api := r.Group(addr.Path)
	api.POST("/cars", carHandler.Create)
	api.GET("cars/:id", carHandler.Get)

}
