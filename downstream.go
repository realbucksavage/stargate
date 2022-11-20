package stargate

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/pkg/errors"
)

var errUnknownScheme = errors.New("unknown scheme")

// DownstreamServer is an abstraction for Stargate that represents the server to be reverse proxied.
// NewDownstreamServer returns an appropriate implementation of this interface.
type DownstreamServer interface {
	http.Handler
	Address() string
	Healthy(ctx context.Context) error
}

// NewDownstreamServer returns a DownstreamServer implementation backed by http or WebSockets, depending
// on the protocol of the passed address. The address to be passed must have http, https, ws, or wss
// protocol. Anything else passed to this function will make it return an "unknown scheme" error.
func NewDownstreamServer(address string, director DirectorFunc) (DownstreamServer, error) {
	origin, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	scheme := origin.Scheme
	if scheme == "http" || scheme == "https" {
		return &httpDownstream{
			url:     address,
			backend: &httputil.ReverseProxy{Director: director(origin)},
			alive:   false,
		}, nil
	}

	if scheme == "ws" || scheme == "wss" {
		return &websocketDownstream{address}, nil
	}

	return nil, errors.Wrap(errUnknownScheme, scheme)
}
