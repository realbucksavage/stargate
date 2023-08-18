package stargate

import (
	"net/http/httptest"
	"testing"
)

func TestRoundRobin(t *testing.T) {

	maxServers := 3
	backends := make([]*httptest.Server, maxServers)
	backendAddreses := map[int]string{}

	for i := 1; i <= maxServers; i++ {
		b := httptest.NewServer(namedHandler("x"))
		backends[i-1] = b

		backendAddreses[i] = "http://" + b.Listener.Addr().String()
	}

	lister := newLister(backends, "http")
	routes, err := lister.List("/")
	if err != nil {
		t.Fatalf("cannot list routes: %v", err)
	}

	lb, err := RoundRobin(routes, defaultDirector("/"))
	if err != nil {
		t.Fatalf("cannot create round robin balancer: %v", err)
	}

	for i, j := 1, 1; i < 10; i++ {

		sv := lb.NextServer()
		if addr := sv.Address(); addr != backendAddreses[j] {
			t.Fatalf("expected address for iteration %d to be %q, got %q", i, backendAddreses[j], addr)
		}

		j++
		if j > maxServers {
			j = 1
		}
	}
}
