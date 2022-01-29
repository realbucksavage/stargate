package examples

import (
	"github.com/realbucksavage/stargate"
	"github.com/realbucksavage/stargate/pkg/listers"
	"github.com/realbucksavage/stargate/pkg/middleware"

	"log"
	"net/http"
)

func main() {
	l := listers.Static{
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

func someMiddleware(next http.Handler) http.HandlerFunc {
	count := 0
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		count++
	}
}
