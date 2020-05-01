package stargate

type LoadBalancer interface {
	NextServer() *DownstreamServer
	Length() int
}

type LoadBalancerMaker func([]string, DirectorFunc) (LoadBalancer, error)
