package stargate

// LoadBalancer is used to determine which downstream service should be invoked next to serve a request.
type LoadBalancer interface {

	// NextServer returns an instance of *DownstreamServer that should be used to serve and http request.
	NextServer() OriginServer

	// Length returns how many downstream servers are available.
	Length() int

	// Name returns a friendly name of this balancer
	Name() string
}

// LoadBalancerMaker creates a LoadBalancer from the input OriginServer slice.
type LoadBalancerMaker func([]OriginServer) (LoadBalancer, error)
