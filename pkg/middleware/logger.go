package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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
func LoggingMiddleware() mux.MiddlewareFunc {
	return LoggerWithConfig(LoggerConfig{log.New(os.Stdout, loggerName, log.LstdFlags)})
}

// LoggerWithConfig creates a stargate.Middleware that logs on the LoggerConfig's Logger instance
func LoggerWithConfig(conf LoggerConfig) mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

			lrw := &loggingResponseWriter{rw, http.StatusOK}
			defer func(begin time.Time) {
				conf.Logger.Printf(
					"[%s | %d] %s\t(%v)",
					r.Method,
					lrw.status,
					r.RequestURI,
					time.Since(begin),
				)
			}(time.Now())
      
			next.ServeHTTP(lrw, r)
		})
	}
}
