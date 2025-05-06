package nethttp

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"sync"
	"vehicle-sharing-go/app/inventory/internal/vehicle/command"
	"vehicle-sharing-go/app/inventory/internal/vehicle/database/gorm"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/net-http/rest"
)

func HandleHTTPServer(
	ctx context.Context,
	shutdownHook func(server *http.Server, name string),
	addr *url.URL,
	carQueryService *gorm.CarProjectionRepository,
	createCarCommandHandler *command.CreateCarHandler,
	wg *sync.WaitGroup,
	errc chan<- error,
	logger *log.Logger,
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

	shutdownHook(server, "net-http")
}
