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
// Implemented by httpDownstream and websocketDownstream.
type DownstreamServer interface {
	http.Handler
	Address() string
	Healthy(ctx context.Context) error
}

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
