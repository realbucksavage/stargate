package stargate

type LoadBalancer interface {
	NextServer() *DownstreamServer
}

type LoadBalancerMaker func(*Context, []string) (LoadBalancer, error)
