package stargate

import "testing"

func TestHealthCounter(t *testing.T) {

	checker := newHealthCounter(3)
	if !checker.ok() {
		t.Fatal("checker should be initially healthy")
	}

	checker.countUnhealthy()
	if !checker.ok() {
		t.Fatalf("checker should be ok until 3 failures")
	}

	checker.countUnhealthy()
	checker.countUnhealthy()
	if checker.ok() {
		t.Fatalf("checker should be unhealthy after 3 failures")
	}

	checker.countHealthy()
	if checker.ok() {
		t.Fatalf("checker should not be ok until 3 healthy pings")
	}

	checker.countHealthy()
	checker.countHealthy()
	if !checker.ok() {
		t.Fatalf("checker sholuld be ok after 3 healthy pings")
	}
}
