# Stargate

A lightweight, extendable, and blazing fast API Gateway.

## Getting started

Stargate's concept is to take in a table of routes and downstream services and create a load balancer that reverse
proxies to them. This table is created by a `stargate.ServiceLister` instance and is passed to `stargate.NewProxy`
function.

Check the [basic example](https://github.com/realbucksavage/stargate/blob/master/examples/basic.go) that implements a
`stargate.ServiceLister` to create a static table of routes and uses round-robin appproch to load balance the request.

### Using dynamic route tables.

If the `stargate.ServiceLister`'s implementation updates the route table, the `stargate.Proxy` instance can be told
to update the routing by calling the `Reload()` method.

Check the [reloading routes example](https://github.com/realbucksavage/stargate/blob/master/examples/reload.go).

### Middleware

A middleware is a function that is defined like this

```go
type Middleware func(*Context, http.Handler) http.HandlerFunc
```

Any middleware to be applied must be passed to the `NewProxy` function like shown.

```go
sg, err := stargate.NewProxy(lister, balancers.RoundRobin, middleware1, middleware2)
```

Check the [middleware example](https://github.com/realbucksavage/stargate/blob/master/examples/middleware.go). There's
already one [middleware implemented in the `middleware` package](https://github.com/realbucksavage/stargate/blob/master/middleware/logger.go)
that logs http responses and execution time.

## Open TODOs

- Improve logging
- Improve documentation
- Write tests
- WebSockets

### `ServiceLister` implementations

- Eureka
- Consuul

### `LoadBalancer` implementations

- Priority round robin

### Test with

- HTTP2
- See how Multipart will work
