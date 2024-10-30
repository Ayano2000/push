package middleware

import (
	"context"
	"github.com/Ayano2000/push/internal/types"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

const loggerContextKey types.LoggerContextKey = "logger"

func LogRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get the logger from the default context logger
		log := zerolog.DefaultContextLogger.With().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_ip", r.RemoteAddr).
			Logger()

		// Store logger in request context
		ctx := context.WithValue(r.Context(), loggerContextKey, &log)
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request details after completion
		log.Info().
			Dur("duration_ms", time.Since(start)).
			Msg("Request completed")
	}
}
