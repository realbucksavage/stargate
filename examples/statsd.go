package examples

import (
	"github.com/realbucksavage/stargate"
	"github.com/realbucksavage/stargate/middleware"
	"net/http"

	"log"
)

func main() {
	l := stargate.StaticLister{
		Routes: map[string][]string{
			"/": {"http://localhost:8081"},
		},
	}

	mw := middleware.StatsdMiddleware("127.0.0.1:8125", "stargate_test_")

	sg, err := stargate.NewProxy(l, stargate.RoundRobin, mw)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", &sg))
}
