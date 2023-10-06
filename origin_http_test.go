package stargate

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPHealthCheck(t *testing.T) {

	statusCode := http.StatusOK
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(statusCode)
	})

	server := httptest.NewServer(handler)
	origin, err := NewOriginServer(&RouteOptions{
		Address: makeRouteOption(server, "http").Address,
		HealthCheck: &HealthCheckOptions{
			Interval: 2 * time.Second,
		},
	}, defaultDirector(""))
	if err != nil {
		t.Fatalf("cannot create testing origin server: %v", err)
	}
	defer server.Close()

	t.Log("testing hc with a healthy server")
	if !origin.Healthy() {
		t.Fatal("got unexpected unhealthy origin on a valid server")
	}

	statusCode = http.StatusInsufficientStorage
	time.Sleep(3 * time.Second)
	t.Log("testing hc with a closed server")
	if !origin.Healthy() {
		t.Fatal("got unexpected unhealthy origin on an immediately closed server")
	}

	time.Sleep(6 * time.Second)
	t.Log("testing hc with a closed server after delay")
	if origin.Healthy() {
		t.Fatalf("got unexpected healthy origin after a considerable downtime")
	}

	statusCode = http.StatusOK
	time.Sleep(2 * time.Second)
	if origin.Healthy() {
		t.Fatalf("got unexpected healthy origin after an immediate startup")
	}

	time.Sleep(6 * time.Second)
	if !origin.Healthy() {
		t.Fatalf("got unexpected unhealthy origin after a prolonged uptime")
	}
}
