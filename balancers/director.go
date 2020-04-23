package balancers

import (
	"github.com/realbucksavage/stargate"
	"net/http"
	"net/url"
)

func directorFunc(ctx *stargate.Context, origin *url.URL) func(*http.Request) {
	return func(r *http.Request) {
		r.Header.Add("X-Forwarded-For", r.Host)
		r.Header.Add("X-Origin-Host", origin.Host)
		r.URL.Scheme = "http"

		for h, v := range ctx.GetHeaders() {
			r.Header.Add(h, v)
		}

		r.URL.Host = origin.Host
	}
}
