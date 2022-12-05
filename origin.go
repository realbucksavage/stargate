package stargate

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/pkg/errors"
)

var errUnknownScheme = errors.New("unknown scheme")

// OriginServer is an abstraction for Stargate that represents the server to be reverse proxied.
// NewDownstreamServer returns an appropriate implementation of this interface.
type OriginServer interface {
	http.Handler
	Address() string
	Healthy(ctx context.Context) error
}

// NewOriginServer returns a DownstreamServer implementation backed by http or WebSockets, depending
// on the protocol of the passed address. The address to be passed must have http, https, ws, or wss
// protocol. Anything else passed to this function will make it return an "unknown scheme" error.
func NewOriginServer(address string, director DirectorFunc) (OriginServer, error) {
	origin, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	directorFunc := director(origin)

	scheme := origin.Scheme
	if scheme == "http" || scheme == "https" {
		return &httpOriginServer{
			url:     address,
			backend: &httputil.ReverseProxy{Director: directorFunc},
			alive:   false,
		}, nil
	}

	if scheme == "ws" || scheme == "wss" {
		return &websocketOriginServer{
			url:      address,
			director: directorFunc,
		}, nil
	}

	return nil, errors.Wrap(errUnknownScheme, scheme)
}
