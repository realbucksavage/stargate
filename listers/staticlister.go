package listers

type StaticLister struct {
	Routes map[string][]string
}

func (s StaticLister) List(route string) []string {
	return s.Routes[route]
}

func (s StaticLister) ListALl() map[string][]string {
	return s.Routes
}
