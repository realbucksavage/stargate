# Stargate

A lightweight, extendable, and blazing fast API Gateway

## Specs

Running basic program

```go
func main() {

	l := listers.StaticLister{
		Routes: map[string][]string{
			"/":    {"http://localhost:8081", "http://localhost:8082"},
			"/api": {"https://jsonplaceholder.typicode.com/posts"},
		},
	}
	sg, err := stargate.NewProxy(l, balancers.MakeRoundRobin)
	if err != nil {
		log.Fatal(err)
	}

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
	ListAll() map[string][]string
}
```

#### Static

#### Eureka Discovery

### Load Balancers

```go
type LoadBalancer interface {
	NextServer() *DownstreamServer
}
```

#### Round-Robin

### Middlewares

### SSL Termination