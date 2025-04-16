package nethttp

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
	"vehicle-sharing-go/app/inventory/internal/vehicle/database/gorm"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/net-http/rest"
)

func HandleHTTPServer(
	ctx context.Context,
	addr *url.URL,
	carQueryService *gorm.CarProjectionRepository,
	wg *sync.WaitGroup,
	errc chan<- error,
	logger *log.Logger,
	_ bool,
) {
	mux := http.NewServeMux()

	carHandler := rest.NewCarHandler(carQueryService)
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
		server.Shutdown(shutdownCtx)
	}()
}
