# Stargate

A lightweight, extendable, blazing fast reverse proxy load balancer.

## Specs

Running basic program

```go
func main() {

	l := listers.StaticLister{
		Routes: map[string][]string{
			"/":    {"http://localhost:8081", "http://localhost:8082"},
			"/api": {"http://localhost:8083"},
		},
	}
	sg := stargate.NewProxy(l, func() stargate.LoadBalancer { return &balancers.RoundRobin{} })
	http.Handle("/", sg)
	log.Fatal(http.ListenAndServe(":7000", nil))
}
```

### Backends

```go
type DownstreamServer struct {
    Backend httputil.ReverseProxy
    HealthURL string
    Alive bool
}

func (s BackendServer) IsAlive() bool {
    // Check for server health
    return true
}
```

### Backend list providers

```go
type ServiceLister interface {
	List(string) []string
	ListALl() map[string][]string
}
```

#### Static

#### Eureka Discovery

### Load Balancers

```go
type LoadBalancer interface {
	InitRoutes(svc []string)
	NextServer() *DownstreamServer
}
```

#### Round-Robin

### Middlewares

### SSL Termination