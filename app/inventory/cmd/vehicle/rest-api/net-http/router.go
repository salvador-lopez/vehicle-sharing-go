package nethttp

import (
	"context"
	"log"
	"net/http"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/net-http/rest"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/net-http/rest/middleware"
)

func registerHandlers(
	ctx context.Context,
	mux *http.ServeMux,
	carHandler *rest.CarHandler,
	logger *log.Logger, // Pass logger as a parameter
) {
	mux.Handle("/api/inventory/vehicles/cars", middleware.LogRequest(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			carHandler.Create(ctx, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/inventory/vehicles/cars/{id}", middleware.LogRequest(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			carHandler.Get(ctx, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))
}
