package internalhttp

import (
	"context"
	"log"
	"net/http"
	"github.com/google/uuid"
	"time"
)

type contextKey string

const requestIDKey contextKey = "requestID"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getRequestID(ctx context.Context) string {
	if v := ctx.Value(requestIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return "unknown"
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		log.Printf(
			"method=%s path=%s request_id=%s duration=%s",
			r.Method,
			r.URL.Path,
			getRequestID(r.Context()),
			duration,
		)
	})
}