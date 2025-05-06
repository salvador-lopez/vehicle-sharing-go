package gin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"sync"
	"vehicle-sharing-go/app/inventory/internal/vehicle/command"
	"vehicle-sharing-go/app/inventory/internal/vehicle/database/gorm"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/gin/rest"
)

func HandleHTTPServer(
	shutdownHook func(server *http.Server, name string),
	addr *url.URL,
	carQueryService *gorm.CarProjectionRepository,
	createCarCommandHandler *command.CreateCarHandler,
	wg *sync.WaitGroup,
	errc chan<- error,
	logger *log.Logger,
	debug bool,
) {
	if debug {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	carHandler := rest.NewCarHandler(createCarCommandHandler, carQueryService)
	registerHandlers(r, addr, carHandler)

	server := &http.Server{
		Addr:    addr.Host,
		Handler: r,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Printf("gin listening on %s", addr.Host)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errc <- err
		}
	}()

	shutdownHook(server, "gin")
}
