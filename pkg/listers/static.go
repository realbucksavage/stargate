package listers

type Static struct {
	Routes map[string][]string
}

func (s Static) List(route string) []string {
	return s.Routes[route]
}

func (s Static) ListAll() (map[string][]string, error) {
	return s.Routes, nil
}
