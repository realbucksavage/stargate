package stargate

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/pkg/errors"
)

// DownstreamServer is a backend service to connect downstream
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

	return nil, errors.Errorf("unknown scheme %q (only http, https, ws, and wss are supported)", scheme)
}
