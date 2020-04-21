package stargate

type LoadBalancer interface {
	InitRoutes(svc []string)
	NextServer() *DownstreamServer
}

type LoadBalancerMaker func() LoadBalancer
