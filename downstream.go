package stargate

import (
	"context"
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
func (d DownstreamServer) IsAlive(ctx context.Context) bool {
	if time.Since(d.lastAlive).Seconds() < 30.0 {
		return true
	}

	u, err := url.Parse(d.BaseURL)
	if err != nil {
		Log.Error("invalid URL %s: %v", d.BaseURL, err)
		return false
	}

	if u.Scheme == "" {
		Log.Warn("no scheme specified in %s, assuming http", d.BaseURL)
		u.Scheme = "http"
	}

	Log.Debug("checking if %q is up...", u)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		Log.Error("cannot create a new GET request to %q: %v", u, err)
		return false
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
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
