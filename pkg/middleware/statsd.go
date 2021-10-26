package middleware

import (
	"fmt"
	"net/http"
	"time"

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
func StatsdMiddleware(address, prefix string) stargate.Middleware {

	client := statsd.NewStatsdClient(address, prefix)
	if err := client.CreateSocket(); err != nil {
		stargate.Logger.Warningf("Cannot start statsd client for %s: %s", address, err)
	}

	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()

			lrw := &loggingResponseWriter{w, http.StatusOK}
			next.ServeHTTP(lrw, r)

			if err := client.Timing(tagResponseTime, time.Since(t).Milliseconds()); err != nil {
				stargate.Logger.Warningf(formatStatError, "Time", err)
			}

			if err := client.Incr(fmt.Sprintf(formatStatusCode, lrw.status), int64(1)); err != nil {
				stargate.Logger.Warningf(formatStatError, "Incr", err)
			}
		}
	}
}
