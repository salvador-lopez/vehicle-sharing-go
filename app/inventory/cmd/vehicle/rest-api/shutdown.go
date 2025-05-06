package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

type shutdownHookFunc func(server *http.Server, name string)

func registerShutdownHook(ctx context.Context, logger *log.Logger) shutdownHookFunc {
	return func(server *http.Server, name string) {
		go func() {
			<-ctx.Done()
			logger.Printf("shutting down %s server...", name)
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(shutdownCtx); err != nil {
				logger.Printf("error shutting down %s server: %v", name, err)
				return
			}
			logger.Printf("%s server shutdown gracefully.", name)
		}()
	}
}
