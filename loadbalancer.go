package stargate

// LoadBalancer is used to determine which downstream service should be invoked next to serve a request.
type LoadBalancer interface {

	// NextServer returns an instance of *DownstreamServer that should be used to serve and http request.
	NextServer() *DownstreamServer

	// Length returns how many downstream servers are available.
	Length() int
}

// LoadBalancerMaker takes in the addresses of downstream servers in a []string. The func(*http.Request) returned by
// DirectorFunc is used for the Director of httputil.ReverseProxy.
type LoadBalancerMaker func([]string, DirectorFunc) (LoadBalancer, error)
