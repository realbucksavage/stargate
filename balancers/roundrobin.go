package balancers

import (
	"github.com/realbucksavage/stargate"
	"net/http/httputil"
	"net/url"
)

// RoundRobin implements the round-robin load balancing algorithm
type RoundRobin struct {
	servers []*stargate.DownstreamServer
	latest  int
}

// NextServer returns the next server that should serve the request.
func (r *RoundRobin) NextServer() *stargate.DownstreamServer {
	if len(r.servers) == 0 {
		return nil
	}

	i := (r.latest + 1) % len(r.servers)
	r.latest = i

	return r.servers[i]
}

// MakeRoundRobin creates new instance of RoundRobin with the passed addresses as backend servers.
func MakeRoundRobin(ctx *stargate.Context, svc []string) (stargate.LoadBalancer, error) {
	r := RoundRobin{}

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
