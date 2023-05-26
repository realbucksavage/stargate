package listers

import "github.com/realbucksavage/stargate"

type Static struct {
	Routes map[string][]*stargate.RouteOptions
}

func (s *Static) List(route string) ([]*stargate.RouteOptions, error) {
	return s.Routes[route], nil
}

func (s *Static) ListAll() (map[string][]*stargate.RouteOptions, error) {
	return s.Routes, nil
}

var _ stargate.ServiceLister = (*Static)(nil)
