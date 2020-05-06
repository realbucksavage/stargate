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
	loggerName = "stargate.requests"
)

var (
	defaultWriter = os.Stdout
	defaultLevel  = log.INFO
)

// LoggerConfig facilitates configuration of the request log writer. As of now, only the logging level and its output
// are configurable.
// TODO: Add more config options like "formatter"
type LoggerConfig struct {
	// Out is the Writer instance the logger will write to.
	Out io.Writer

	// Level is logging level for the requests logger.
	Level log.Level
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader makes loggingResponseWriter implement http.ResponseWriter.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware creates the middleware with the logger set to default config.
func LoggingMiddleware() stargate.Middleware {
	return LoggerWithConfig(LoggerConfig{})
}

// LoggedWithOutput creates the middleware with the logger set to write to the passed in io.Writer
func LoggerWithOutput(w io.Writer) stargate.Middleware {
	return LoggerWithConfig(LoggerConfig{Out: w})
}

// LoggerWithConfig takes in an entire LoggerConfig struct and creates the middleware with passed in configuration.
func LoggerWithConfig(conf LoggerConfig) stargate.Middleware {
	if conf.Out == nil {
		conf.Out = defaultWriter
	}

	if conf.Level == 0 {
		conf.Level = defaultLevel
	}

	l := log.MustGetLogger(loggerName)
	log.SetLevel(conf.Level, loggerName)

	return func(ctx *stargate.Context, next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			lrw := &loggingResponseWriter{w, http.StatusOK}
			next.ServeHTTP(lrw, r)

			l.Infof("[%s | %d] %s\t\t(%dms)",
				r.Method,
				lrw.status,
				r.RequestURI,
				time.Now().Sub(start).Milliseconds())
		}
	}
}
