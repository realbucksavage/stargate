package stargate

import (
	"net/http"
	"net/http/httputil"
	"time"
)

// DownstreamServer is a backend service to connect downstream
// TODO : Add health checking
type DownstreamServer struct {
	BaseURL string
	Backend *httputil.ReverseProxy
	Alive   bool

	lastAlive time.Time
}

// ServiceLister provides all available routes and their downstream services
type ServiceLister interface {
	List(string) []string
	ListAll() map[string][]string
}

// IsAlive performs a healthcheck on the server and returns true if the server responds back
func (d DownstreamServer) IsAlive() bool {
	if time.Since(d.lastAlive).Seconds() < 30.0 {
		return true
	}
	_, err := http.Get(d.BaseURL)
	if err != nil {
		return false
	}

	// TODO : Ignore status check for now.
	d.lastAlive = time.Now()
	return true

}
