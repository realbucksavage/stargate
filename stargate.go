package stargate

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Proxy struct {
	mux *mux.Router
}

type Middleware func(*Context, http.Handler) http.HandlerFunc

func (s Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func NewProxy(ctx *Context, l ServiceLister, loadBalancerMaker LoadBalancerMaker, mwf ...Middleware) (Proxy, error) {
	r := mux.NewRouter()

	routes := l.ListAll()
	for route, svc := range routes {
		lb, err := loadBalancerMaker(svc, defaultDirector(ctx))
		if err != nil {
			return Proxy{}, err
		}

		handler := http.HandlerFunc(serve(lb))
		for _, m := range mwf {
			handler = m(ctx, handler)
		}

		r.HandleFunc(route, handler)
	}

	return Proxy{r}, nil
}

func serve(lb LoadBalancer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var server *DownstreamServer

		if lb.Length() > 0 {
			serverCount := 0
			for sv := lb.NextServer(); serverCount < lb.Length(); sv = lb.NextServer() {
				if sv.IsAlive() {
					server = sv
					break
				}
				Logger.Debugf("Backend %s is not alive. Skipped.", sv.BaseURL)
				serverCount++
			}
		}

		if server == nil {
			Logger.Errorf("No alive server available for route %s", r.URL)

			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(http.StatusServiceUnavailable)

			_, err := w.Write([]byte(`<h1>503 Service Unavailable</h1>"`))
			if err != nil {
				Logger.Errorf("Unable to write response to client: %v\n", err)
			}
			return
		}

		Logger.Debugf("Resolved backend %s", server.BaseURL)
		server.Backend.ServeHTTP(w, r)
	}
}
