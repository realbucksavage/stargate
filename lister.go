package stargate

import "time"

type RouteOptions struct {
	Address     string
	HealthCheck *HealthCheckOptions
}

type HealthCheckOptions struct {
	Path          string
	Timeout       time.Duration
	HealthyStatus int
}

// ServiceLister provides all available routes and their downstream services
type ServiceLister interface {
	List(string) ([]*RouteOptions, error)
	ListAll() (map[string][]*RouteOptions, error)
}
