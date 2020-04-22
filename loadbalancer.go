package stargate

type LoadBalancer interface {
	NextServer() *DownstreamServer
}

type LoadBalancerMaker func([]string) (LoadBalancer, error)
