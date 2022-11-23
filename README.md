# Stargate

A minimal and extensible library to build gateway servers. Stargate aims to be simple while providing niche solutions
like several load balancer implementations, middleware, service discovery, etc.

Stargate supports:

- WebSockets
- Hot-reloading of routes
- Middleware

[stargatecontrb](https://github.com/realbucksavage/stargatecontrib) contains some middleware implementations that are
not in the scope of this library, but might be useful for some people.

## Getting started

Check the [basic example](./_examples/basic/main.go) that implements a
`stargate.ServiceLister` to create a static table of routes and uses round-robin approach to load balance the request.

In the same sprits, the [WebSockets example](./_examples/websockets/main.go) shows a simple WebSocket backend.

### Customize logging

Stargate uses `stargate.Log` variable to write its logging output. This variable is an implementation
of `stargate.Logger`. You may write your own implementation of this interface and write `stargate.Log = myOwnLogger{}`
whenever your program starts.

Check the [custom logger example](./_examples/logger_custom/main.go).

### Using dynamic route tables.

If the `stargate.ServiceLister`'s implementation updates the route table, the `stargate.Router` instance can be told to
update the routing by calling the `Reload()` method.

Check the [reloading routes example](./_examples/reloading_router/reload.go).

#### Eureka service discovery

Check the [eureka package in stargatecontrib](https://github.com/realbucksavage/stargatecontrib/tree/main/lister/eureka).

### Middleware

Stargate defines middleware as:

```go
type MiddlewareFunc func (http.Handler) http.Handler
```

Check the [middleware example](./_examples/middleware/main.go), that counts the number of requests served.

## Open TODOs

- Improve logging
- Improve documentation
- Write more tests
- Implementation of director functions in WebSockets reverse proxy.
- Customizable healthchecks

#### `LoadBalancer` implementations

- Priority round robin
