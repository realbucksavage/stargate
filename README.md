# Stargate

A lightweight, extendable, blazing fast reverse proxy load balancer.

## Specs

```go
sg := stargate.RoundRobinLoadBalancer{
    BackendListProvider: stargate.StaticBackendList{
        Servers: []string{"http://backend-1", "http://backend-2"},
    }
}

http.HandlerFunc("/", func (w http.ResponseWriter, r *http.Request) {
    sg.Serve(w, r)
})
log.Fatal(http.ListenAndServe(":8080", nil))
```

### Backends

```go
type DownstreamServer struct {
    Route string
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
type BackendListProvider interface {
    List() []BackendServer
}
```

#### Static

#### Eureka Discovery

### Load Balancers

```go
type LoadBalancer interface {
    NextServer() BackendServer
    Serve(w http.ResponseWriter, r *http.Request)
}
```

#### Round-Robin

### Middlewares

### SSL Termination