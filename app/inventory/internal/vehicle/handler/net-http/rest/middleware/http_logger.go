package middleware

import (
	"log"
	"net/http"
	"time"
)

func LogRequest(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Printf("Started %s %s", r.Method, r.URL.Path)

		lrw := &LoggingResponseWriter{ResponseWriter: w}

		next.ServeHTTP(lrw, r)

		logger.Printf("Completed %s %s with status %d in %v", r.Method, r.URL.Path, lrw.statusCode, time.Since(start))
	})
}

// LoggingResponseWriter wraps the http.ResponseWriter to capture the status code
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code when the response is written.
func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}
