package stargate

import (
	"errors"
	"testing"
)

func TestNewOriginServer(t *testing.T) {

	tests := []struct {
		name    string
		address string
		err     error
	}{
		{"check for unsupported schemes", "invalid://some-address", errUnknownScheme},
		{"check that URLs without a scheme return an error", "some-address", errUnknownScheme},
		{"check that http servers are created", "http://some-address", nil},
		{"check that https servers are created", "https://some-address", nil},
		{"check that websocket servers are created", "ws://some-address", nil},
		{"check that websocket secured servers are created", "wss://some-address", nil},
	}

	for _, test := range tests {
		t.Logf("running test %q", test.name)
		_, err := NewOriginServer(test.address, defaultDirector("/"))
		if err != nil {
			if test.err == nil {
				t.Fatalf("NewDownstreamServer unexpectedly returned an error: %v", err)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("NewDownstreamServer returned an incorrect error [%v] - expected %v", err, test.err)
			}
		}
	}
}
