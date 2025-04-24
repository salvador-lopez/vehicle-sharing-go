package nethttp

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
	"vehicle-sharing-go/app/inventory/internal/vehicle/command"
	"vehicle-sharing-go/app/inventory/internal/vehicle/database/gorm"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/net-http/rest"
)

//	@title			Api Inventory Vehicles
//	@swagger		3.0
//	@version		1.0
//	@description	HTTP service to interact with inventory vehicles bounded context.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @servers.url   http://localhost:8088/api/inventory/vehicles

func HandleHTTPServer(
	ctx context.Context,
	addr *url.URL,
	carQueryService *gorm.CarProjectionRepository,
	createCarCommandHandler *command.CreateCarHandler,
	wg *sync.WaitGroup,
	errc chan<- error,
	logger *log.Logger,
	_ bool,
) {
	mux := http.NewServeMux()

	carHandler := rest.NewCarHandler(createCarCommandHandler, carQueryService)
	registerHandlers(ctx, mux, carHandler, logger)

	server := &http.Server{
		Addr:    addr.Host,
		Handler: mux,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Printf("http-net listening on %s", addr.Host)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errc <- err
		}
	}()

	// Shutdown hook
	go func() {
		<-ctx.Done()
		logger.Println("shutting down net-http server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Println("error shutting down net-http server.")
			return
		}
		logger.Println("net-http server shutdown gracefully.")
	}()
}
