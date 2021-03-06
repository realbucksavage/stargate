package stargate

import (
	"net/http"
	"net/http/httputil"
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

	_, err := http.Get(d.BaseURL)
	if err != nil {
		Logger.Errorf("Alive-check failed for server %s : %v", d.BaseURL, err)
		return false
	}

	// TODO : Ignore status check for now.
	d.lastAlive = time.Now()
	return true

}
