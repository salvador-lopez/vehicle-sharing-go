package gin

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
	"vehicle-sharing-go/app/inventory/internal/vehicle/command"
	"vehicle-sharing-go/app/inventory/internal/vehicle/database/gorm"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/gin/rest"
)

func HandleHTTPServer(
	ctx context.Context,
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

	// Shutdown hook
	go func() {
		<-ctx.Done()
		logger.Println("shutting down gin server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Println("error shutting down gin server.")
			return
		}
		logger.Println("gin server shutdown gracefully.")
	}()
}
