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

	backends := []string{toUrl(backend)}
	roundRobin, err := RoundRobin(backends, defaultDirector("/"))
	if err != nil {
		t.Errorf("Cannot create roundRobin LB : %v", err)
	}

	server := httptest.NewServer(serve(roundRobin))
	defer server.Close()

	client := &http.Client{}
	get, err := client.Get(toUrl(server))
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

func toUrl(s *httptest.Server) string {
	return fmt.Sprintf("http://%s", s.Listener.Addr().String())
}
