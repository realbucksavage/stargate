package stargate

type LoadBalancer interface {
	NextServer() *DownstreamServer
	Length() int
	UpdateRoutes([]string, DirectorFunc) error
}

type LoadBalancerMaker func([]string, DirectorFunc) (LoadBalancer, error)
