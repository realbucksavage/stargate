package middleware

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/realbucksavage/stargate"

	log "github.com/op/go-logging"
)

const (
	loggerName = "stargate-requests"
)

var (
	defaultWriter = os.Stdout
	defaultLevel  = log.INFO
)

type (
	LoggerConfig struct {
		Out   io.Writer
		Level log.Level
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
	return LoggerWithConfig(LoggerConfig{})
}

func LoggerWithOutput(w io.Writer) stargate.MiddlewareFunc {
	return LoggerWithConfig(LoggerConfig{Out: w})
}

func LoggerWithConfig(conf LoggerConfig) stargate.MiddlewareFunc {
	if conf.Out == nil {
		conf.Out = defaultWriter
	}

	if conf.Level == 0 {
		conf.Level = defaultLevel
	}

	l := log.MustGetLogger(loggerName)
	log.SetLevel(conf.Level, loggerName)

	return func(ctx *stargate.Context) func(http.Handler) http.Handler {
		return func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				start := time.Now()

				lrw := &loggingResponseWriter{w, http.StatusOK}
				handler.ServeHTTP(lrw, r)

				l.Infof("[%s | %d] %s\t\t(%dms))",
					r.Method,
					lrw.status,
					r.RequestURI,
					time.Now().Sub(start).Milliseconds())

			})
		}
	}
}