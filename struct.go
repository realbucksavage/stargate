package stargate

import "net/http/httputil"

// DownstreamServer is a backend service to connect downstream
// TODO : Add health checking
type DownstreamServer struct {
	Route   string
	Backend *httputil.ReverseProxy
	Alive   bool
}

type StargateProxy struct {
	LoadBalancer LoadBalancer
}
