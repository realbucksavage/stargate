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

func (r *roundRobinBalancer) Length() int {
	return len(r.servers)
}

func (r *roundRobinBalancer) UpdateRoutes(svc []string, director stargate.DirectorFunc) error {
	r.servers = []*stargate.DownstreamServer{}

	for _, s := range svc {
		localServer, err := createProxy(s, director)
		if err != nil {
			return err
		}

		r.servers = append(r.servers, &localServer)
	}
	r.latest = -1

	return nil
}

// RoundRobin creates new instance of LoadBalancer that implements the Round-Robin load balancing algorithm.
func RoundRobin(svc []string, director stargate.DirectorFunc) (stargate.LoadBalancer, error) {
	r := &roundRobinBalancer{}
	if err := r.UpdateRoutes(svc, director); err != nil {
		return nil, err
	}

	return r, nil
}

func createProxy(s string, director stargate.DirectorFunc) (stargate.DownstreamServer, error) {
	var localServer stargate.DownstreamServer

	origin, err := url.Parse(s)
	if err != nil {
		return stargate.DownstreamServer{}, err
	}

	localServer.BaseURL = s
	localServer.Alive = localServer.IsAlive()
	localServer.Backend = &httputil.ReverseProxy{
		Director: director(origin),
	}
	return localServer, nil
}
