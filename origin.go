package stargate

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/pkg/errors"
)

var errUnknownScheme = errors.New("unknown scheme")

// OriginServer is an abstraction for Stargate that represents the server to be reverse proxied.
// NewDownstreamServer returns an appropriate implementation of this interface.
type OriginServer interface {
	io.Closer
	http.Handler
	Address() string
	Healthy() bool
	startHealthCheck(options *HealthCheckOptions)
}

// NewOriginServer returns a DownstreamServer implementation backed by HTTP or WebSockets, depending
// on the RouteOptions' address. The said address must have http, https, ws, or wss protocol. Anything
// else passed to this function will make it return an "unknown scheme" error.
func NewOriginServer(routeOptions *RouteOptions, director DirectorFunc) (OriginServer, error) {
	origin, err := url.Parse(routeOptions.Address)
	if err != nil {
		return nil, err
	}

	directorFunc := director(origin)

	var server OriginServer
	switch origin.Scheme {
	case "http", "https":
		server = &httpOriginServer{
			url:     routeOptions.Address,
			backend: &httputil.ReverseProxy{Director: directorFunc},
			alive:   false,
		}
	case "ws", "wss":
		server = &websocketOriginServer{
			url:      routeOptions.Address,
			director: directorFunc,
		}
	default:
		return nil, errors.Wrap(errUnknownScheme, origin.Scheme)
	}

	if routeOptions.HealthCheck != nil {
		go server.startHealthCheck(routeOptions.HealthCheck)
	}

	return server, nil
}
