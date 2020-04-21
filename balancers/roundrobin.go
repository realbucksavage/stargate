package balancers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/realbucksavage/stargate"
)

// RoundRobin implements the round-robin load balancing algorithm
type RoundRobin struct {
	servers []*stargate.DownstreamServer
	latest  int
}

func (r *RoundRobin) InitRoutes(svc []string) {
	r.servers = []*stargate.DownstreamServer{}

	for _, s := range svc {
		var localServer stargate.DownstreamServer

		origin, _ := url.Parse(s)
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
}

func (r *RoundRobin) NextServer() *stargate.DownstreamServer {
	if len(r.servers) == 0 {
		return nil
	}

	i := (r.latest + 1) % len(r.servers)
	r.latest = i

	return r.servers[i]
}
