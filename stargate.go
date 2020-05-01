package stargate

import (
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

type Proxy struct {
	mux           *mux.Router
	lister        ServiceLister
	balancerMaker LoadBalancerMaker
	ctx           *Context
	middleware    []Middleware
	balancers     map[string]LoadBalancer

	mutex sync.Mutex
}

type Middleware func(*Context, http.Handler) http.HandlerFunc

func (s *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	s.mux.ServeHTTP(w, r)
	s.mutex.Unlock()
}

func (s *Proxy) Reload() {

	rtr := mux.NewRouter()
	routes := s.lister.ListAll()
	for route, svc := range routes {

		lb, err := s.balancerMaker(svc, defaultDirector(s.ctx, route))
		if err != nil {
			Logger.Errorf("Cannot create a loadBalancer for route %s : %v", route, err)
			continue
		}
		s.balancers[route] = lb

		handler := createHandler(s.ctx, lb, s.middleware)
		rtr.NewRoute().Path(route).Handler(handler)

		Logger.Infof("Route created - %s", route)
	}

	s.mutex.Lock()
	s.mux = rtr
	s.mutex.Unlock()
}

func NewProxy(ctx *Context, l ServiceLister, loadBalancerMaker LoadBalancerMaker, mwf ...Middleware) (Proxy, error) {
	r := mux.NewRouter()

	routes := l.ListAll()
	bm := map[string]LoadBalancer{}
	for route, svc := range routes {
		lb, err := loadBalancerMaker(svc, defaultDirector(ctx, route))
		if err != nil {
			Logger.Errorf("Cannot create a loadBalancer for route %s : %v", route, err)
			return Proxy{}, err
		}
		bm[route] = lb

		handler := createHandler(ctx, lb, mwf)
		r.HandleFunc(route, handler)
	}

	return Proxy{
		mux:           r,
		lister:        l,
		ctx:           ctx,
		middleware:    mwf,
		balancerMaker: loadBalancerMaker,
		balancers:     bm,
	}, nil
}

func createHandler(ctx *Context, lb LoadBalancer, mwf []Middleware) http.HandlerFunc {
	handler := http.HandlerFunc(serve(lb))
	for _, m := range mwf {
		handler = m(ctx, handler)
	}
	return handler
}
