package stargate

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServe(t *testing.T) {
	okCode := "ok"
	backend := httptest.NewServer(namedHandler(okCode))
	defer backend.Close()

	ctx := new(Context)

	backends := []string{toUrl(backend)}
	roundRobin, err := RoundRobin(backends, defaultDirector(ctx, "/"))
	if err != nil {
		t.Errorf("Cannot create roundRobin LB : %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(serve(roundRobin)))
	defer server.Close()

	client := &http.Client{}
	get, err := client.Get(toUrl(server))
	if err != nil {
		t.Errorf("Cannot execute GET request : %v", err)
	}

	b, _ := ioutil.ReadAll(get.Body)
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
