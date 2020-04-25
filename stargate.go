package stargate

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Proxy struct {
	mux *mux.Router
	ctx *Context
}

type MiddlewareFunc func(*Context) func(http.Handler) http.Handler

func (s Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s Proxy) UseMiddleware(mw MiddlewareFunc) {
	s.mux.Use(mw(s.ctx))
}

func NewProxy(l ServiceLister, loadBalancerMaker LoadBalancerMaker) (Proxy, error) {
	r := mux.NewRouter()
	ctx := &Context{}

	routes := l.ListAll()
	for route, svc := range routes {
		lb, err := loadBalancerMaker(ctx, svc)
		if err != nil {
			return Proxy{}, err
		}

		r.HandleFunc(route, serve(lb))
	}

	return Proxy{mux: r, ctx: ctx}, nil
}

func serve(lb LoadBalancer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var server *DownstreamServer

		serverCount := 0
		for sv := lb.NextServer(); serverCount < lb.Length(); sv = lb.NextServer() {
			if sv != nil && sv.IsAlive() {
				server = sv
				break
			}
			serverCount++
		}

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
