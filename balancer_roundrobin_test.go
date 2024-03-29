package stargate

import (
	"net/http/httptest"
	"testing"
)

func TestRoundRobin(t *testing.T) {

	maxServers := 3
	backends := make([]*httptest.Server, maxServers)
	backendAddresses := map[int]string{}

	for i := 1; i <= maxServers; i++ {
		b := httptest.NewServer(namedHandler("x"))
		backends[i-1] = b

		backendAddresses[i] = "http://" + b.Listener.Addr().String()
	}

	lister := newLister(backends, "http")
	routes, err := lister.List("/")
	if err != nil {
		t.Fatalf("cannot list routes: %v", err)
	}

	servers := make([]OriginServer, 0)
	for _, options := range routes {
		server, err := NewOriginServer(options, defaultDirector("/"))
		if err != nil {
			t.Fatal(err)
		}

		servers = append(servers, server)
	}

	lb, err := RoundRobin(servers)
	if err != nil {
		t.Fatalf("cannot create round robin balancer: %v", err)
	}

	for i, j := 1, 1; i < 10; i++ {

		sv := lb.NextServer()
		if addr := sv.Address(); addr != backendAddresses[j] {
			t.Fatalf("expected address for iteration %d to be %q, got %q", i, backendAddresses[j], addr)
		}

		j++
		if j > maxServers {
			j = 1
		}
	}
}
