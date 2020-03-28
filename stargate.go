package stargate

import "net/http"

func (s StargateProxy) HandlerFunc() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		server := s.LoadBalancer.NextServer()
		server.Backend.ServeHTTP(w, r)
	}
}
