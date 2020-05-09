package middleware

import (
	"fmt"
	"github.com/realbucksavage/stargate"
	"net/http"
	"time"

	"github.com/quipo/statsd"
)

const (
	formatStatusCode = "http_status_%d"
	formatTime       = "response_time"
)

func StatsdMiddleware(address, prefix string) stargate.Middleware {

	client := statsd.NewStatsdClient(address, prefix)
	if err := client.CreateSocket(); err != nil {
		stargate.Logger.Warningf("Cannot start statsd client for %s: %s", address, err)
	}

	return func(ctx *stargate.Context, next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()

			lrw := &loggingResponseWriter{w, http.StatusOK}
			next.ServeHTTP(lrw, r)

			if err := client.Timing(formatTime, time.Since(t).Milliseconds()); err != nil {
				stargate.Logger.Warningf("Time stat not recorded: %v", err)
			}

			if err := client.Incr(fmt.Sprintf(formatStatusCode, lrw.status), int64(1)); err != nil {
				stargate.Logger.Warningf("Incr stat not recorded: %s", err)
			}
		}
	}
}
