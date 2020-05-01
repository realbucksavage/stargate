# Stargate

A lightweight, extendable, and blazing fast API Gateway.

## Basics

```go
package main

import (
	"log"
	"net/http"

	"github.com/realbucksavage/stargate"
	"github.com/realbucksavage/stargate/balancers"
	"github.com/realbucksavage/stargate/listers"
)

func main() {

	l := listers.StaticLister{
		Routes: map[string][]string{
			"/downstream_1": {"http://localhost:8081", "http://localhost:8082"},
			"/downstream_2": {"http://localhost:8083"},
		},
	}
	sg, err := stargate.NewProxy(l, balancers.RoundRobin)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", &sg)
	log.Fatal(http.ListenAndServe(":7000", nil))
}
```

### Middleware

A middleware is a function that is defined like this

```go
type Middleware func(*Context, http.Handler) http.HandlerFunc
```

You can create your middleware and pass them to `NewProxy` like this:

```go
package main

import (
	"net/http"

	"github.com/realbucksavage/stargate"
	"github.com/realbucksavage/stargate/balancers"
)

func main() {
	// declare a context and a lister
	sg, err := stargate.NewProxy(lister, balancers.RoundRobin, myMiddleware)
}

func myMiddleware(ctx *stargate.Context, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Some code here
		next.ServeHTTP(w, r)
	})
}
```

There is a middleware that logs all requests in the `middlware` package

```go
package main

import (
	"github.com/realbucksavage/stargate"
	"github.com/realbucksavage/stargate/balancers"
	mw "github.com/realbucksavage/stargate/middleware"
)

func main() {
	stargate.NewProxy(lister, balancers.RoundRobin, myMiddleware, mw.LoggingMiddleware())
}
```

## Open TODOs

- Improve logging
- Write tests
- Implement a Eureka Client
- Priority RoundRobin implementation
- HTTP2
- WebSockets
- See how Multipart will work
