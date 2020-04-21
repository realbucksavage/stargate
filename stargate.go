package stargate

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type StargateProxy struct {
	mux *mux.Router
}

func (s StargateProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func NewProxy(l ServiceLister, loadBalancerMaker LoadBalancerMaker) StargateProxy {
	r := mux.NewRouter()

	routes := l.ListAll()
	for route, svc := range routes {
		lb := loadBalancerMaker()
		lb.InitRoutes(svc)
		r.HandleFunc(route, serve(lb))
	}

	return StargateProxy{mux: r}
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
