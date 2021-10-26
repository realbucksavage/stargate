package stargate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type httptestLister struct {
	routes map[string][]string
}

func (h httptestLister) List(s string) []string {
	return h.routes[s]
}

func (h httptestLister) ListAll() (map[string][]string, error) {
	return h.routes, nil
}

func newLister(servers []*httptest.Server) ServiceLister {
	sv := make([]string, 0)

	for _, s := range servers {
		sv = append(sv, toUrl(s))
	}

	return httptestLister{routes: map[string][]string{
		"/": sv,
	}}
}

func TestRoundRobin(t *testing.T) {
	backends := make([]*httptest.Server, 0)
	maxServers := 3

	for i := 1; i <= maxServers; i++ {
		n := fmt.Sprintf("server_%d", i)
		backends = append(backends, httptest.NewServer(namedHandler(n)))
		t.Logf("Named server %s ready", n)
	}

	ls := newLister(backends)
	sg, err := NewProxy(ls, RoundRobin)
	if err != nil {
		t.Errorf("Cannot create stargate proxy : %v", err)
	}

	server := httptest.NewServer(&sg)
	defer server.Close()

	t.Logf("Stargate ready at %s", server.Listener.Addr().String())

	client := &http.Client{}
	for i, j := 1, 1; i <= 10; i++ {
		get, err := client.Get(toUrl(server))
		if err != nil {
			t.Error(err)
		}

		b, _ := ioutil.ReadAll(get.Body)
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

	for _, s := range backends {
		s.Close()
	}
}
