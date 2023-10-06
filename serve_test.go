package stargate

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServe(t *testing.T) {
	okCode := "ok"
	backend := httptest.NewServer(namedHandler(okCode))
	defer backend.Close()

	origin, err := NewOriginServer(makeRouteOption(backend, "http"), defaultDirector("/"))
	if err != nil {
		t.Fatal(err)
	}

	backends := []OriginServer{origin}
	roundRobin, err := RoundRobin(backends)
	if err != nil {
		t.Errorf("Cannot create roundRobin LB : %v", err)
	}

	server := httptest.NewServer(serve(roundRobin))
	defer server.Close()

	client := &http.Client{}
	get, err := client.Get(makeRouteOption(server, "http").Address)
	if err != nil {
		t.Errorf("Cannot execute GET request : %v", err)
	}

	b, err := io.ReadAll(get.Body)
	if err != nil {
		t.Fatalf("cannot read response from server: %v", err)
	}

	resp := string(b)
	if resp != okCode {
		t.Errorf(`Expected "%s" but got "%s"`, okCode, resp)
	}
}

func namedHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := io.WriteString(w, name); err != nil {
			panic(err)
		}
	}
}

func makeRouteOption(s *httptest.Server, protocol string) *RouteOptions {
	return &RouteOptions{Address: fmt.Sprintf("%s://%s", protocol, s.Listener.Addr().String())}
}
