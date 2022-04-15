# Stargate

A lightweight, extendable, and blazing fast gateway server library.

Stargate supports:

- Hot-reloading of routes
- Eureka service registry
- Middleware

## Getting started

Stargate's concept is to take in a table of routes and downstream services and create a load balancer that reverse
proxies to them. This table is created by a `stargate.ServiceLister` instance and is passed to `stargate.NewProxy`
function.

Check the [basic example](https://github.com/realbucksavage/stargate/blob/master/examples/basic.go) that implements a
`stargate.ServiceLister` to create a static table of routes and uses round-robin approach to load balance the request.

### Customize logging

Stargate uses `stargate.Log` variable to write its logging output. This variable is an implementation
of `stargate.Logger`. You may write your own implementation of this interface and write `stargate.Log = myOwnLogger{}`
whenever your program starts.

Check the [custom logger example](https://github.com/realbucksavage/stargate/blob/master/examples/custom_logger.go).

### Using dynamic route tables.

If the `stargate.ServiceLister`'s implementation updates the route table, the `stargate.Proxy` instance can be told to
update the routing by calling the `Reload()` method.

Check the [reloading routes example](https://github.com/realbucksavage/stargate/blob/master/examples/reload.go).

#### Eureka service discovery

`stargate.EurekaLister(string)` returns a `ServiceLister` instance that queries the specified eureka server for
registered applications. Calling the `Reload()` method on `stargate.Proxy` instance causes the Eureka lister to query
eureka server and update the routes.

```go
el := stargate.EurekaLister("http://localhost:8761/eureka")
```

Check the [eureka service discovery example](https://github.com/realbucksavage/stargate/blob/master/examples/eureka.go).

### Publishing statistics to `statsd`

This can be done by using
the [`StatsdMiddleware`](https://github.com/realbucksavage/stargate/blob/master/middleware/statsd.go#L21).

```go
mw := middleware.StatsdMiddleware("127.0.0.1:8125", "some_app.")
sg, err := stargate.NewProxy(lister, stargate.RoundRobin, mw)
```

Refer the [example](https://github.com/realbucksavage/stargate/blob/master/examples/statsd.go).

> **NOTE**: As of now, the middleware publishes only response times and HTTP response rates.
> Please use the issues board to request more stat integrations, or create a PR :smile:

### Middleware

Stargate utlizes `MiddlewareFunc` from [gorilla/mux](https://github.com/gorilla/mux), which are defined as:

```go
type MiddlewareFunc func(next http.Handler) http.Handler
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

#### `ServiceLister` implementations

- Etcd
- Consuul

#### `LoadBalancer` implementations

- Priority round robin
