package stargate

type LoadBalancer interface {
	NextServer() *DownstreamServer
}
