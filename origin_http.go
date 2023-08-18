package stargate

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type httpOriginServer struct {
	url       string
	backend   *httputil.ReverseProxy
	alive     bool
	lastAlive time.Time
}

func (h *httpOriginServer) Address() string {
	return h.url
}

func (h *httpOriginServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.backend.ServeHTTP(rw, r)
}

func (h *httpOriginServer) Healthy(ctx context.Context) error {
	if time.Since(h.lastAlive).Seconds() < 30.0 {
		return nil
	}

	u, err := url.Parse(h.url)
	if err != nil {
		return err
	}

	Log.Debug("checking if %q is up...", u)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			Log.Error("cannot close response body from [%s] after alive-check: %v", h.url, err)
		}
	}()

	// TODO : Ignore status check for now.
	h.lastAlive = time.Now()
	return nil
}
