package stargate

type roundRobinBalancer struct {
	servers []DownstreamServer
	latest  int
}

// NextServer returns the next server that should serve the request.
func (r *roundRobinBalancer) NextServer() DownstreamServer {
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

// RoundRobin creates new instance of LoadBalancer that implements the Round-Robin load balancing algorithm.
func RoundRobin(svc []string, director DirectorFunc) (LoadBalancer, error) {
	r := &roundRobinBalancer{}
	if err := r.createRoutes(svc, director); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *roundRobinBalancer) createRoutes(svc []string, director DirectorFunc) error {
	r.servers = []DownstreamServer{}

	for _, s := range svc {
		localServer, err := NewDownstreamServer(s, director)
		if err != nil {
			return err
		}

		r.servers = append(r.servers, localServer)
	}
	r.latest = -1

	return nil
}
