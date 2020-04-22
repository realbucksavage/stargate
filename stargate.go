package stargate

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Proxy struct {
	mux *mux.Router
}

func (s Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s Proxy) UseMiddleware(mw func(http.Handler) http.Handler) {
	s.mux.Use(mw)
}

func NewProxy(l ServiceLister, loadBalancerMaker LoadBalancerMaker) (Proxy, error) {
	r := mux.NewRouter()

	routes := l.ListAll()
	for route, svc := range routes {
		lb, err := loadBalancerMaker(svc)
		if err != nil {
			return Proxy{}, err
		}

		r.HandleFunc(route, serve(lb))
	}

	return Proxy{mux: r}, nil
}

func serve(lb LoadBalancer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		server := lb.NextServer()
		if server == nil {
			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(http.StatusServiceUnavailable)

			_, err := w.Write([]byte(`<h1>503 Service Unavailable</h1>"`))
			if err != nil {
				log.Printf("Unable to write response to client: %v\n", err)
			}
			return
		}
		server.Backend.ServeHTTP(w, r)
	}
}
