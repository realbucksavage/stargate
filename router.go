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

type downstreamRoute struct {
	pathPrefix string
	handler    http.Handler
}

// Router implements http.Handler and handles all requests that are to be reverse-proxied.
type Router struct {
	lister            ServiceLister
	routes            []downstreamRoute
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
		routes: make([]downstreamRoute, 0),
		mut:    sync.RWMutex{},
	}

	for _, opt := range options {
		opt(router)
	}

	if router.loadBalancerMaker == nil {
		router.loadBalancerMaker = RoundRobin
	}

	if err := router.Reload(); err != nil {
		return nil, nil
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
	newRoutes := make([]downstreamRoute, 0)

	for route, svc := range routes {

		if _, ok := mappedRoutes[route]; ok {
			return errors.Errorf("route %q is already mapped", route)
		}

		lb, err := r.loadBalancerMaker(svc, defaultDirector(route))
		if err != nil {
			return errors.Wrapf(err, "cannot create load balancer to downstream service %q", svc)
		}

		handler := r.createHandler(lb, r.middlewareFuncs...)
		newRoutes = append(newRoutes, downstreamRoute{pathPrefix: route, handler: handler})
		mappedRoutes[route] = struct{}{}

		Log.Debug("Route initialized - %s -> %s", route, svc)
	}

	sort.SliceStable(newRoutes, func(i, j int) bool { return newRoutes[i].pathPrefix > newRoutes[j].pathPrefix })

	r.mut.Lock()
	r.routes = newRoutes
	r.mut.Unlock()

	return nil
}

// ServeHTTP satisfies http.Handler
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
	rw.Write([]byte(`<h1>Page Not Found</h1><br><small>Stargate Router`))
}

func (r *Router) createHandler(lb LoadBalancer, mwf ...MiddlewareFunc) http.Handler {

	handler := serve(lb)
	for _, mw := range mwf {
		handler = mw(handler)
	}

	return handler
}
