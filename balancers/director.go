package balancers

import (
	"net/http"
	"net/url"
)

func directorFunc(origin *url.URL) func(*http.Request) {
	return func(r *http.Request) {
		r.Header.Add("X-Forwarded-For", r.Host)
		r.Header.Add("X-Origin-Host", origin.Host)
		r.URL.Scheme = "http"

		r.URL.Host = origin.Host
	}
}
