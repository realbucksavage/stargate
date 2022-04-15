package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/realbucksavage/stargate"

	"github.com/quipo/statsd"
)

const (
	formatStatusCode = "http_status_%d"
	formatStatError  = "%s stat not recorded: %v"

	tagResponseTime = "response_time"
)

// StatsdMiddleware sends current response rate and response latency to the provided statd
// daemons like https://github.com/statd/statd or Amazon CloudWatch Agent.
func StatsdMiddleware(address, prefix string) mux.MiddlewareFunc {

	client := statsd.NewStatsdClient(address, prefix)
	if err := client.CreateSocket(); err != nil {
		stargate.Log.Warn("Cannot start statsd client for %s: %s", address, err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			lrw := &loggingResponseWriter{w, http.StatusOK}
			defer func(begin time.Time) {
				if err := client.Timing(tagResponseTime, time.Since(begin).Milliseconds()); err != nil {
					stargate.Log.Warn(formatStatError, "Time", err)
				}

				if err := client.Incr(fmt.Sprintf(formatStatusCode, lrw.status), int64(1)); err != nil {
					stargate.Log.Warn(formatStatError, "Incr", err)
				}
			}(time.Now())

			next.ServeHTTP(lrw, r)
		})
	}
}
