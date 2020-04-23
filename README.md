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

	mw "github.com/realbucksavage/stargate/middleware"
)

func main() {

	l := listers.StaticLister{
		Routes: map[string][]string{
			"/downstream_1":    {"http://localhost:8081", "http://localhost:8082"},
			"/downstream_2":    {"http://localhost:8083"},
		},
	}
	sg, err := stargate.NewProxy(l, balancers.MakeRoundRobin)
	if err != nil {
		log.Fatal(err)
	}

	sg.UseMiddleware(mw.LoggingMiddleware())

	http.Handle("/", sg)
	log.Fatal(http.ListenAndServe(":7000", nil))
}
```

### Custom middleware

```go
func headerAddingMiddleware(ctx *stargate.Context) func(http.Handler) http.Handler {

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx.AddHeader("Test", "Abcd1234")
			h.ServeHTTP(w, r)
		})
	}
}
```

## Open TODOs

- Make `LoadBalancer` aware of changes in downstream server list
- Implement a Eureka Client