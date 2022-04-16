package stargate

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// ServiceLister provides all available routes and their downstream services
type ServiceLister interface {
	List(string) ([]string, error)
	ListAll() (map[string][]string, error)
}

// Proxy implements http.Handler and handles all requests that are to be reverse-proxied. Proxy wraps a *mux.Router for
// route-matching.
type Proxy struct {
	router        *mux.Router
	lister        ServiceLister
	balancerMaker LoadBalancerMaker
	middleware    []mux.MiddlewareFunc

	mutex sync.RWMutex
}

func (s *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.router.ServeHTTP(w, r)
}

// Reload queries the ServiceLister used to create the Proxy instance, and re-initializes the underlying *mux.Router.
// This method is used to implement hot-reloading and can be called after the ServiceLister has updated its routes.
func (s *Proxy) Reload() error {

	rtr := mux.NewRouter()
	routes, err := s.lister.ListAll()
	if err != nil {
		Log.Error("Cannot query lister for routes : %v", err)
		return err
	}

	for route, svc := range routes {
		lb, err := s.balancerMaker(svc, defaultDirector(route))
		if err != nil {
			Log.Error("Cannot create a loadBalancer for route %s : %v", route, err)
			return err
		}

		handler := createHandler(lb, s.middleware...)
		rtr.PathPrefix(route).Handler(handler)

		Log.Debug("Route updated -\t%s", route)
	}

	// Wait for requests to finish before swapping
	s.mutex.Lock()
	s.router = rtr
	s.mutex.Unlock()

	return nil
}

// NewProxy takes in a ServiceLister, LoadBalancerMaker, a chain of Middleware and creates a functional Proxy
// instance. The keys of the map returned from ServiceLister.ListAll() are used as the base path of the routes added
// for them. The values of the map returned from ServiceLister.ListAll() are used to create a LoadBalancer for each route.
func NewProxy(l ServiceLister, loadBalancerMaker LoadBalancerMaker, mwf ...mux.MiddlewareFunc) (Proxy, error) {
	r := mux.NewRouter()

	routes, err := l.ListAll()
	if err != nil {
		Log.Error("Cannot query lister for routes : %v", err)
		return Proxy{}, err
	}

	for route, svc := range routes {
		lb, err := loadBalancerMaker(svc, defaultDirector(route))
		if err != nil {
			Log.Error("Cannot create a loadBalancer for route %s : %v", route, err)
			return Proxy{}, err
		}

		handler := createHandler(lb, mwf...)
		r.PathPrefix(route).Handler(handler)

		Log.Info("Route initialized -\t%s", route)
	}

	return Proxy{
		router:        r,
		lister:        l,
		middleware:    mwf,
		balancerMaker: loadBalancerMaker,
	}, nil
}

func createHandler(lb LoadBalancer, mwf ...mux.MiddlewareFunc) http.Handler {
	handler := serve(lb)
	for _, m := range mwf {
		handler = m(handler)
	}
	return handler
}
