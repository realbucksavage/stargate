package stargate

import (
	"net/http/httptest"
	"testing"
)

func TestWithLoadBalancer(t *testing.T) {

	lister := newLister([]*httptest.Server{}, "http")

	tests := []struct {
		name  string
		maker LoadBalancerMaker
		value string
	}{
		{"no balancer maker should default to RoundRobin", nil, "Stargate/RoundRobin"},
		{"check RoundRobin", RoundRobin, "Stargate/RoundRobin"},
	}

	for _, tc := range tests {
		t.Logf("running test %q", tc.name)

		var (
			router *Router
			err    error
		)

		if tc.maker != nil {
			router, err = NewRouter(lister, WithLoadBalancer(tc.maker))
		} else {
			router, err = NewRouter(lister)
		}

		if err != nil {
			t.Fatalf("cannot create stargate proxy: %v", err)
		}

		rs, err := lister.List("/")
		if err != nil {
			t.Fatalf("cannot query routes: %v", err)
		}

		lb, err := router.loadBalancerMaker(rs, defaultDirector("/"))
		if err != nil {
			t.Fatalf("cannot create lb: %v", err)
		}

		if lb.Name() != tc.value {
			t.Fatalf("expected %q load balancer, got %q", tc.value, lb.Name())
		}
	}
}
