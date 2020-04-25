package balancers

import (
	"github.com/realbucksavage/stargate"
	"net/http/httputil"
	"net/url"
)

type roundRobinBalancer struct {
	servers []*stargate.DownstreamServer
	latest  int
}

// NextServer returns the next server that should serve the request.
func (r *roundRobinBalancer) NextServer() *stargate.DownstreamServer {
	if len(r.servers) == 0 {
		return nil
	}

	i := (r.latest + 1) % len(r.servers)
	r.latest = i

	return r.servers[i]
}

// RoundRobin creates new instance of LoadBalancer that implements the Round-Robin load balancing algorithm.
func RoundRobin(ctx *stargate.Context, svc []string) (stargate.LoadBalancer, error) {
	r := roundRobinBalancer{}

	r.servers = []*stargate.DownstreamServer{}

	for _, s := range svc {
		var localServer stargate.DownstreamServer

		origin, err := url.Parse(s)
		if err != nil {
			return nil, err
		}

		director := directorFunc(ctx, origin)
		localServer.Alive = localServer.IsAlive()
		localServer.Backend = &httputil.ReverseProxy{
			Director: director,
		}

		r.servers = append(r.servers, &localServer)
	}
	r.latest = -1

	return &r, nil
}
