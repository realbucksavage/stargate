package stargate

import (
	"time"
)

// RouteOptions defines the configuration of a route.
type RouteOptions struct {
	// Address is the absolute address of an origin server.
	Address string

	// HealthCheck indicates that a health checker routine is to spawned, if not nil.
	HealthCheck *HealthCheckOptions
}

// HealthCheckOptions defines the behavior of the health checker routine.
type HealthCheckOptions struct {
	// Path is the relative path on the origin server that is to be used for health checking. Defaults to "/".
	Path string

	// Interval is the frequency of health checks. Defaults to 30s.
	Interval time.Duration

	// Timeout dictates how long Stargate should wait for a health check ping to finish. Defaults to 10s.
	Timeout time.Duration

	// HealthyStatus is the expected status code of a successful health check ping. Defaults to http.StatusOK.
	HealthyStatus int
}

// ServiceLister provides all available routes and their downstream services
type ServiceLister interface {
	List(string) ([]*RouteOptions, error)
	ListAll() (map[string][]*RouteOptions, error)
}
