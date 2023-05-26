package stargate

import (
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// MiddlewareFunc is a function that takes an http.Handler and returns another http.Handler. The returned http.Handler
// is a closure that can call the passed in http.Handler to move the HTTP call forward. Optionally the returned closure
// can do some extra processing - like authentication - with http.ResponseWriter and http.Request it receives.
type MiddlewareFunc func(next http.Handler) http.Handler

var errNoLister = errors.New("a lister is required")

type originRoute struct {
	pathPrefix string
	handler    http.Handler
	servers    []OriginServer
}

// Router implements http.Handler and handles all requests that are to be reverse-proxied.
type Router struct {
	lister            ServiceLister
	routes            []originRoute
	loadBalancerMaker LoadBalancerMaker
	middlewareFuncs   []MiddlewareFunc

	mut sync.RWMutex
}

// NewRouter creates a Router instance out of the downstream services supplied by ServiceLister parameter.
func NewRouter(lister ServiceLister, options ...RouterOption) (*Router, error) {

	if lister == nil {
		return nil, errNoLister
	}

	router := &Router{
		lister: lister,
		routes: make([]originRoute, 0),
		mut:    sync.RWMutex{},
	}

	for _, opt := range options {
		opt(router)
	}

	if router.loadBalancerMaker == nil {
		router.loadBalancerMaker = RoundRobin
	}

	if err := router.Reload(); err != nil {
		return nil, errors.Wrap(err, "cannot initialize routes")
	}

	return router, nil
}

// Reload queries the ServiceLister used with NewRouter and creates the internal routing table used by ServeHTTP.
func (r *Router) Reload() error {

	routes, err := r.lister.ListAll()
	if err != nil {
		return err
	}

	mappedRoutes := make(map[string]struct{})
	newRoutes := make([]originRoute, 0)

	for route, routeOptions := range routes {

		if _, ok := mappedRoutes[route]; ok {
			return errors.Errorf("route %q is already mapped", route)
		}

		servers := make([]OriginServer, 0)
		for _, routeOption := range routeOptions {
			sv, err := NewOriginServer(routeOption, defaultDirector(route))
			if err != nil {
				return errors.Wrapf(err, "cannot create an origin server pointing to %q for route %q", route, routeOption.Address)
			}

			servers = append(servers, sv)
		}

		lb, err := r.loadBalancerMaker(servers)
		if err != nil {
			return errors.Wrapf(err, "cannot create load balancer for route %q", route)
		}

		handler := r.createHandler(lb, r.middlewareFuncs...)
		newRoutes = append(newRoutes, originRoute{pathPrefix: route, handler: handler, servers: servers})
		mappedRoutes[route] = struct{}{}

		Log.Debug("Route initialized - %s -> %s", route, routeOptions)
	}

	sort.SliceStable(newRoutes, func(i, j int) bool { return newRoutes[i].pathPrefix > newRoutes[j].pathPrefix })

	r.mut.Lock()
	defer closeServers(r.routes)
	r.routes = newRoutes
	r.mut.Unlock()

	return nil
}

// ServeHTTP satisfies http.Handler. It prioritizes full URL matches from the internal routing table, and tries until /
// is reached. For example, to serve a request to https://somehost.com/some/test/url, ServeHTTP tries to look for URLs
// in the routing table in this order ; /some/test/url -> /some/test -> /some -> /
//
// The downstream service pertaining to the first matched URL is picked and the request is reverse proxied to that.
func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	r.mut.RLock()
	defer r.mut.RUnlock()

	path := req.URL.Path
	for _, route := range r.routes {
		if strings.HasPrefix(path, route.pathPrefix) {
			route.handler.ServeHTTP(rw, req)
			return
		}
	}

	rw.WriteHeader(http.StatusNotFound)
	rw.Header().Set("content-type", "text/html")
	if _, err := rw.Write([]byte(`<h1>Page Not Found</h1><br><small>Stargate Router`)); err != nil {
		Log.Error("cannot write not found response to client: %v", err)
	}
}

func (r *Router) createHandler(lb LoadBalancer, mwf ...MiddlewareFunc) http.Handler {

	handler := serve(lb)
	for _, mw := range mwf {
		handler = mw(handler)
	}

	return handler
}

func closeServers(routes []originRoute) {
	for i := range routes {
		for _, server := range routes[i].servers {
			if err := server.Close(); err != nil {
				Log.Warn("cannot close origin server %q: %v", server.Address(), err)
			}
		}
	}
}
