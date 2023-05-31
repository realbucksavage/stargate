package stargate

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPHealthCheck(t *testing.T) {

	server := httptest.NewServer(namedHandler("x"))
	origin, err := NewOriginServer(&RouteOptions{
		Address: makeRouteOption(server, "http").Address,
		HealthCheck: &HealthCheckOptions{
			Interval: 5 * time.Second,
		},
	}, defaultDirector(""))
	if err != nil {
		t.Fatalf("cannot create testing origin server: %v", err)
	}

	t.Log("testing hc with a healthy server")
	if !origin.Healthy() {
		t.Fatal("got unexpected unhealthy origin on a valid server")
	}

	server.Close()

	time.Sleep(6 * time.Second)

	t.Log("testing hc with a closed server")
	if origin.Healthy() {
		t.Fatal("got unexpected healthy origin on a closed server")
	}
}
