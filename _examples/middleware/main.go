package main

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
	sg, err := stargate.NewRouter(
		l,
		stargate.WithMiddleware(someMiddleware, middleware.LoggingMiddleware()),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", sg))
}

func someMiddleware(next http.Handler) http.Handler {
	count := 0
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("count is not %d", count)
		next.ServeHTTP(w, r)
		count++
	})
}
