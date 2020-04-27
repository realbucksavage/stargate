package stargate

import (
	"net/http"
	"net/url"
)

type DirectorFunc func(*url.URL) func(*http.Request)

func defaultDirector(ctx *Context) DirectorFunc {

	return func(origin *url.URL) func(*http.Request) {
		return func(r *http.Request) {
			r.Header.Add("X-Forwarded-For", r.Host)
			r.Header.Add("X-Origin-Host", origin.Host)
			r.URL.Scheme = "http"

			for h, v := range ctx.headers {
				r.Header.Add(h, v)
			}

			r.URL.Host = origin.Host
		}
	}
}
