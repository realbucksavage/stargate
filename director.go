package stargate

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type DirectorFunc func(*url.URL) func(*http.Request)

func defaultDirector(route string) DirectorFunc {
	return func(origin *url.URL) func(*http.Request) {
		return func(r *http.Request) {
			r.URL.Path = strings.Replace(r.URL.Path, route, "", 1)
			r.URL.RawQuery = combineQueries(origin, r.URL)

			r.Header.Add("X-Forwarded-For", r.Host)
			r.Header.Add("X-Origin-Host", origin.Host)
			r.URL.Scheme = origin.Scheme

			r.URL.Host = origin.Host
		}
	}
}

func combineQueries(src *url.URL, dst *url.URL) string {
	q := ""

	if src.RawQuery != "" {
		q = src.RawQuery
	}

	if dst.RawQuery != "" {
		if q == "" {
			q = dst.RawQuery
		} else {
			q = fmt.Sprintf("%s&%s", q, dst.RawQuery)
		}
	}

	return q
}
