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
	backend := httptest.NewServer(http.HandlerFunc(testServerHandler))
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
	if resp != "ok" {
		t.Errorf(`Expected "ok" but got "%s"`, resp)
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
