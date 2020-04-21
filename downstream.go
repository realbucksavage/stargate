package stargate

import "net/http/httputil"

// DownstreamServer is a backend service to connect downstream
// TODO : Add health checking
type DownstreamServer struct {
	Backend *httputil.ReverseProxy
	Alive   bool
}

// ServiceLister provides all available routes and their downstream services
type ServiceLister interface {
	List(string) []string
	ListAll() map[string][]string
}

// IsAlive performs a healthcheck on the server and returns true if the server responds back
// TODO: This could be a scheduled task.
func (d DownstreamServer) IsAlive() bool {
	// TODO: Implementation
	return true
}
