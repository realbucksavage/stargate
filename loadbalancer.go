package stargate

type LoadBalancer interface {
	NextServer() *DownstreamServer
	Length() int
}

type LoadBalancerMaker func(*Context, []string) (LoadBalancer, error)
