package stargate

import (
	"net/http"
	"net/url"
	"strings"
)

type DirectorFunc func(*url.URL) func(*http.Request)

func defaultDirector(ctx *Context, route string) DirectorFunc {

	return func(origin *url.URL) func(*http.Request) {
		return func(r *http.Request) {
			r.URL.Path = strings.Replace(r.URL.Path, route, "/", 1)

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
