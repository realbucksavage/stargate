package stargate

// ServiceLister provides all available routes and their downstream services
// ServiceLister provides all available routes and their downstream services
type ServiceLister interface {
	List(string) []string
	ListAll() (map[string][]string, error)
}

type StaticLister struct {
	Routes map[string][]string
}

func (s StaticLister) List(route string) []string {
	return s.Routes[route]
}

func (s StaticLister) ListAll() (map[string][]string, error) {
	return s.Routes, nil
}
