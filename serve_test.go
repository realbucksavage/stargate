package stargate

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServe(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer backend.Close()

	sl := StaticLister{
		Routes: map[string][]string{
			"/": {toUrl(backend)},
		},
	}
	sg, err := NewProxy(sl, RoundRobin)
	if err != nil {
		panic(err)
	}

	server := httptest.NewServer(&sg)
	defer server.Close()

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
		},
	}

	get, err := cli.Get(toUrl(server))
	if err != nil {
		panic(err)
	}

	all, _ := ioutil.ReadAll(get.Body)
	if string(all) != "ok" {
		t.Errorf(`Expected "ok" but got "%s"`, string(all))
	}
}

func testServerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, "ok"); err != nil {
		panic(err)
	}
}

func toUrl(s *httptest.Server) string {
	return fmt.Sprintf("http://%s", s.Listener.Addr().String())
}
