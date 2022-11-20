package stargate

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouting(t *testing.T) {

	var sg *Router
	{
		backend1 := httptest.NewServer(namedHandler("/backends/1/"))
		defer backend1.Close()

		backend2 := httptest.NewServer(namedHandler("/backends/2/"))
		defer backend2.Close()

		backend3 := httptest.NewServer(namedHandler("/backends/"))
		defer backend3.Close()

		lister := httptestLister{
			routes: map[string][]string{
				"/backends/2": {toUrl(backend2, "http")},
				"/backends/1": {toUrl(backend1, "http")},
				"/backends":   {toUrl(backend3, "http")},
			},
		}

		var err error
		sg, err = NewRouter(lister)
		if err != nil {
			t.Fatalf("cannot create stargate server: %v", err)
		}
	}

	server := httptest.NewServer(sg)
	defer server.Close()

	table := []struct {
		input  string
		output string
	}{
		{"/backends", "/backends/"},
		{"/backends/", "/backends/"},
		{"/backends/unknown", "/backends/"},
		{"/backends/1", "/backends/1/"},
		{"/backends/1/", "/backends/1/"},
		{"/backends/1/test", "/backends/1/"},
		{"/backends/2", "/backends/2/"},
		{"/backends/2/", "/backends/2/"},
		{"/backends/2/test", "/backends/2/"},
	}

	baseURL := toUrl(server, "http")
	for _, tc := range table {
		u := fmt.Sprintf("%s%s", baseURL, tc.input)
		get, err := http.DefaultClient.Get(u)
		if err != nil {
			t.Fatalf("cannot perform get on %q: %v", u, err)
		}

		b, err := io.ReadAll(get.Body)
		if err != nil {
			t.Fatalf("cannot read response of %q: %v", u, err)
		}

		if string(b) != tc.output {
			t.Fatalf("expected response %q - got %q", tc.output, string(b))
		}
	}
}
