package stargate

type roundRobinBalancer struct {
	servers []OriginServer
	latest  int
}

// NextServer returns the next server that should serve the request.
func (r *roundRobinBalancer) NextServer() OriginServer {
	if len(r.servers) == 0 {
		return nil
	}

	i := (r.latest + 1) % len(r.servers)
	r.latest = i

	return r.servers[i]
}

// Length returns the number of downstream servers this load balancer can serve to.
func (r *roundRobinBalancer) Length() int {
	return len(r.servers)
}

// Name returns the name of this LoadBalancer.
func (r *roundRobinBalancer) Name() string {
	return "Stargate/RoundRobin"
}

// RoundRobin creates new instance of LoadBalancer that implements the Round-Robin load balancing algorithm.
func RoundRobin(servers []OriginServer) (LoadBalancer, error) {
	return &roundRobinBalancer{servers: servers, latest: -1}, nil
}
