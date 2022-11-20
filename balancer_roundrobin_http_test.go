package stargate

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type httptestLister struct {
	routes map[string][]string
}

func (h httptestLister) List(s string) ([]string, error) {
	return h.routes[s], nil
}

func (h httptestLister) ListAll() (map[string][]string, error) {
	return h.routes, nil
}

func newLister(servers []*httptest.Server, protocol string) ServiceLister {
	sv := make([]string, 0)

	for _, s := range servers {
		sv = append(sv, toUrl(s, protocol))
	}

	return httptestLister{
		routes: map[string][]string{
			"/": sv,
		},
	}
}

func TestRoundRobinHTTP(t *testing.T) {
	maxServers := 3
	backends := make([]*httptest.Server, maxServers)

	for i := 1; i <= maxServers; i++ {
		n := fmt.Sprintf("server_%d", i)
		backends[i-1] = httptest.NewServer(namedHandler(n))
		t.Logf("Named server %s ready", n)
	}

	defer func() {
		for _, s := range backends {
			s.Close()
		}
	}()

	ls := newLister(backends, "ws")
	sg, err := NewRouter(ls)
	if err != nil {
		t.Errorf("Cannot create stargate proxy : %v", err)
	}

	server := httptest.NewServer(sg)
	defer server.Close()

	t.Logf("Stargate ready at %s", server.Listener.Addr().String())

	client := &http.Client{}
	for i, j := 1, 1; i < 10; i++ {
		get, err := client.Get(toUrl(server, "protocol"))
		if err != nil {
			t.Error(err)
		}

		b, err := io.ReadAll(get.Body)
		if err != nil {
			t.Fatalf("cannot read response from server %v", err)
		}

		resp := string(b)

		expected := fmt.Sprintf("server_%d", j)
		if resp != expected {
			t.Errorf("RoundRobin failure - got '%s' expected '%s'", resp, expected)
		}

		j++
		if j > maxServers {
			j = 1
		}

		t.Logf("> iter %d passed with response %s", i, resp)
	}
}
