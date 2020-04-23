package middleware

import (
	"github.com/realbucksavage/stargate"
	"io"
	"net/http"
	"os"
	"time"

	logger "github.com/mgutz/logxi/v1"
)

var (
	defaultWriter = os.Stdout
)

type (
	loggerConfig struct {
		out io.Writer
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		status int
	}
)

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware() stargate.MiddlewareFunc {
	return loggerWithConfig(loggerConfig{})
}

func LoggerWithOutput(w io.Writer) stargate.MiddlewareFunc {
	return loggerWithConfig(loggerConfig{out: w})
}

func loggerWithConfig(conf loggerConfig) stargate.MiddlewareFunc {
	if conf.out == nil {
		conf.out = defaultWriter
	}

	l := logger.NewLogger(conf.out, "stargate-requests")

	return func(ctx *stargate.Context) func(http.Handler) http.Handler {
		return func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				start := time.Now()

				lrw := &loggingResponseWriter{w, http.StatusOK}
				handler.ServeHTTP(lrw, r)

				l.Info("[%d | %s] %s %d ms", lrw.status, r.Method, time.Now().Sub(start).Milliseconds())

			})
		}
	}
}
