package listers

import "github.com/realbucksavage/stargate/v1"

type Static struct {
	Routes map[string][]string
}

func (s Static) List(route string) ([]string, error) {
	return s.Routes[route], nil
}

func (s Static) ListAll() (map[string][]string, error) {
	return s.Routes, nil
}

var _ stargate.ServiceLister = (*Static)(nil)
