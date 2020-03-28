package balancers

import (
	"github.com/realbucksavage/stargate"
	"net/http"
	"net/http/httputil"
	"net/url"
)

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

func NewRoundRobinBalancer(servers []string) *RoundRobin {
	rb := RoundRobin{}
	for _, s := range servers {
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

		rb.servers = append(rb.servers, &localServer)
	}
	rb.latest = -1

	return &rb
}
