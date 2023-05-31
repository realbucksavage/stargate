package stargate

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type httpOriginServer struct {
	url     string
	backend *httputil.ReverseProxy

	hcMut              sync.RWMutex
	alive              bool
	healthCheckRunning bool
}

func (origin *httpOriginServer) Address() string {
	return origin.url
}

func (origin *httpOriginServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	origin.backend.ServeHTTP(rw, r)
}

func (origin *httpOriginServer) Healthy() bool {
	origin.hcMut.RLock()
	defer origin.hcMut.RUnlock()

	if !origin.healthCheckRunning {
		return true
	}

	return origin.alive
}

func (origin *httpOriginServer) Close() error {
	Log.Debug("stopping health checker to %q", origin.url)
	origin.healthCheckRunning = false
	origin.alive = false
	return nil
}

func (origin *httpOriginServer) startHealthCheck(options *HealthCheckOptions) {
	interval := options.Interval
	if interval == 0 {
		interval = DefaultHealthCheckInterval
	}

	path := strings.TrimSpace(options.Path)
	if path == "" {
		path = DefaultHealthCheckPath
	}

	okStatus := options.HealthyStatus
	if http.StatusText(okStatus) == "" {
		okStatus = DefaultHealthCheckStatus
	}

	timeout := options.Timeout
	if timeout == 0 {
		timeout = DefaultHealthCheckTimeout
	}

	Log.Debug("pinging %q every %v at %q with a timeout of %v", origin.url, interval, path, timeout)
	healthTicker := time.NewTicker(interval)
	defer healthTicker.Stop()

	origin.healthCheckRunning = true
	for origin.healthCheckRunning {
		err := origin.checkHealth(path, okStatus, timeout)

		if err != nil {
			Log.Error("%q: healthcheck failed: %v", origin.url, err)
		}

		origin.hcMut.Lock()
		origin.alive = err == nil
		origin.hcMut.Unlock()

		<-healthTicker.C
	}
}

func (origin *httpOriginServer) checkHealth(path string, status int, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	address := fmt.Sprintf("%s/%s", origin.url, path)
	req, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			Log.Error("cannot close response body of health check ping to %q: %v", address, err)
		}
	}()

	if resp.StatusCode != status {
		return errors.Errorf("invalid status %d (expected %d)", resp.StatusCode, status)
	}

	return nil
}
