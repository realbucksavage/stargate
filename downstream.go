package stargate

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// DownstreamServer is a backend service to connect downstream
type DownstreamServer struct {
	BaseURL string
	Backend *httputil.ReverseProxy
	Alive   bool

	lastAlive time.Time
}

// IsAlive performs a healthcheck on the server and returns true if the server responds back. If a server responds to
// an initial healthcheck request, next request is made after 30 seconds.
// TODO: Make healthcheck configurable.
func (d DownstreamServer) IsAlive() bool {
	if time.Since(d.lastAlive).Seconds() < 30.0 {
		return true
	}

	u, err := url.Parse(d.BaseURL)
	if err != nil {
		Log.Error("invalid URL %s: %v", d.BaseURL, err)
		return false
	}

	if u.Scheme == "" {
		Log.Debug("no scheme specified in %s, assuming http", d.BaseURL)
		u.Scheme = "http"
	}

	res, err := http.Get(u.String())
	if err != nil {
		Log.Error("Alive-check failed for server %s : %v", d.BaseURL, err)
		return false
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			Log.Error("cannot close response body from [%s] after alive-check: %v", d.BaseURL, err)
		}
	}()

	// TODO : Ignore status check for now.
	d.lastAlive = time.Now()
	return true

}
