package stargate

import (
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

// Proxy implements the http.Handler interface and handles all requests that are to be reverse proxied. Proxy wraps a
// *mux.Router for route-matching.
type Proxy struct {
	mux           *mux.Router
	lister        ServiceLister
	balancerMaker LoadBalancerMaker
	ctx           *Context
	middleware    []Middleware

	mutex sync.Mutex
}

// Middleware receives a http.Handler and returns an http.HandlerFunc. The returned http.HandlerFunc is a closure that
// can execute some code before and after a request is served. The closure must call ServeHTTP(...) on received
// http.Handler to continue execution of the chain.
type Middleware func(*Context, http.Handler) http.HandlerFunc

func (s *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	s.mux.ServeHTTP(w, r)
	s.mutex.Unlock()
}

// Reload queries the ServiceLister used to create the Proxy instance re-initializes the underlying *mux.Router. This
// method should be called after the ServiceLister has updated its routes.
func (s *Proxy) Reload() {

	rtr := mux.NewRouter()
	routes := s.lister.ListAll()
	for route, svc := range routes {

		lb, err := s.balancerMaker(svc, defaultDirector(s.ctx, route))
		if err != nil {
			Logger.Errorf("Cannot create a loadBalancer for route %s : %v", route, err)
			continue
		}

		handler := createHandler(s.ctx, lb, s.middleware)
		rtr.NewRoute().Path(route).Handler(handler)

		Logger.Infof("Route updated -\t%s", route)
	}

	s.mutex.Lock()
	s.mux = rtr
	s.mutex.Unlock()
}

// NewProxy takes in a ServiceLister, LoadBalancerMaker, and a chain of Middleware and creates a functional Proxy
// instance. The keys of the map returned from ServiceLister.ListAll() are used as the base path of the routes added
// for them. The values of the map returned from ServiceLister.ListAll() is used to create a LoadBalancer for each route.
func NewProxy(l ServiceLister, loadBalancerMaker LoadBalancerMaker, mwf ...Middleware) (Proxy, error) {
	r := mux.NewRouter()
	ctx := &Context{}

	routes := l.ListAll()
	for route, svc := range routes {
		lb, err := loadBalancerMaker(svc, defaultDirector(ctx, route))
		if err != nil {
			Logger.Errorf("Cannot create a loadBalancer for route %s : %v", route, err)
			return Proxy{}, err
		}

		handler := createHandler(ctx, lb, mwf)
		r.HandleFunc(route, handler)

		Logger.Infof("Route initialized -\t%s", route)
	}

	return Proxy{
		mux:           r,
		lister:        l,
		ctx:           ctx,
		middleware:    mwf,
		balancerMaker: loadBalancerMaker,
	}, nil
}

func createHandler(ctx *Context, lb LoadBalancer, mwf []Middleware) http.HandlerFunc {
	handler := http.HandlerFunc(serve(lb))
	for _, m := range mwf {
		handler = m(ctx, handler)
	}
	return handler
}
