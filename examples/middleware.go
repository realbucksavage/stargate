package examples

import (
	"github.com/realbucksavage/stargate"
	"github.com/realbucksavage/stargate/middleware"
	"log"
	"net/http"
	"strconv"
)

func main() {
	l := stargate.StaticLister{
		Routes: map[string][]string{
			"/ds_1": {"http://app1-sv1:8080", "http://app1-sv2:8080"},
			"/ds_2": {"http://app2-sv1:8080"},
		},
	}
	sg, err := stargate.NewProxy(l, stargate.RoundRobin, someMiddleware, middleware.LoggingMiddleware())
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", &sg))
}

func someMiddleware(c *stargate.Context, next http.Handler) http.HandlerFunc {
	count := 0
	return func(w http.ResponseWriter, r *http.Request) {
		c.AddHeader("X-Hit-Count", strconv.Itoa(count))
		next.ServeHTTP(w, r)
		count++
	}
}
