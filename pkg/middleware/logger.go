package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/realbucksavage/stargate"
)

const (
	loggerName = "[stargate.requests] "
)

// LoggerConfig facilitates configuration of the request log writer.
type LoggerConfig struct {
	// Logger the underlying *log.Logger that will receive logging info
	Logger *log.Logger
}

// LoggingMiddleware creates the middleware with the logger set to default config.
func LoggingMiddleware() stargate.Middleware {
	return LoggerWithConfig(LoggerConfig{log.New(os.Stdout, loggerName, log.LstdFlags)})
}

// LoggerWithConfig creates a stargate.Middleware that logs on the LoggerConfig's Logger instance
func LoggerWithConfig(conf LoggerConfig) stargate.Middleware {

	return func(next http.Handler) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lrw := &loggingResponseWriter{rw, http.StatusOK}
			next.ServeHTTP(lrw, r)

			conf.Logger.Printf("[%s | %d] %s\t\t(%v)",
				r.Method,
				lrw.status,
				r.RequestURI,
				time.Since(start))
		}
	}
}
