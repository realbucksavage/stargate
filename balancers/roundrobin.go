package balancers

import (
	"github.com/realbucksavage/stargate"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// RoundRobin implements the round-robin load balancing algorithm
type RoundRobin struct {
	servers []*stargate.DownstreamServer
	latest  int
}

func (r *RoundRobin) NextServer() *stargate.DownstreamServer {
	if len(r.servers) == 0 {
		return nil
	}

	i := (r.latest + 1) % len(r.servers)
	r.latest = i

	return r.servers[i]
}

func MakeRoundRobin(svc []string) (stargate.LoadBalancer, error) {
	r := RoundRobin{}

	r.servers = []*stargate.DownstreamServer{}

	for _, s := range svc {
		var localServer stargate.DownstreamServer

		origin, err := url.Parse(s)
		if err != nil {
			return nil, err
		}

		director := func(r *http.Request) {
			r.Header.Add("X-Forwarded-For", r.Host)
			r.Header.Add("X-Origin-Host", origin.Host)
			r.URL.Scheme = "http"

			r.URL.Host = origin.Host
		}
		localServer.Alive = localServer.IsAlive()
		localServer.Backend = &httputil.ReverseProxy{
			Director: director,
		}

		r.servers = append(r.servers, &localServer)
	}
	r.latest = -1

	return &r, nil
}
