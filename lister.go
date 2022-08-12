package stargate

// ServiceLister provides all available routes and their downstream services
type ServiceLister interface {
	List(string) ([]string, error)
	ListAll() (map[string][]string, error)
}
