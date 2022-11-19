package stargate

import "net/http"

func serve(lb LoadBalancer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var server DownstreamServer

		if lb.Length() > 0 {
			serverCount := 0
			for sv := lb.NextServer(); serverCount < lb.Length(); sv = lb.NextServer() {
				if err := sv.Healthy(r.Context()); err == nil {
					server = sv
					break
				} else {
					Log.Debug("backend %q is not alive: %v", sv.Address(), err)
					serverCount++
				}
			}
		}

		if server == nil {
			Log.Error("No alive server available for downstreamRoute %s", r.URL)

			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(http.StatusServiceUnavailable)

			_, err := w.Write([]byte(`<h1>503 Service Unavailable</h1>"`))
			if err != nil {
				Log.Error("Unable to write response to client: %v\n", err)
			}
			return
		}

		Log.Debug("Resolved backend %s", server.Address())
		server.ServeHTTP(w, r)
	})
}
