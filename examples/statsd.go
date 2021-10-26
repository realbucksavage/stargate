package examples

import (
	"log"
	"net/http"

	"github.com/realbucksavage/stargate"
	"github.com/realbucksavage/stargate/pkg/listers"
	"github.com/realbucksavage/stargate/pkg/middleware"
)

func main() {
	l := listers.Static{
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
